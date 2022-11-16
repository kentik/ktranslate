package st

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/gnxi/utils/xpath"
	go_metrics "github.com/kentik/go-metrics"
	gnmiLib "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type KentikSTListener struct {
	logger.ContextL

	// Internal state
	internalAliases map[string]string
	cancel          context.CancelFunc
	wg              sync.WaitGroup

	config   *ktranslate.KentikSTConfig
	metrics  STListenerMetric
	apic     *api.KentikApi
	devices  map[string]*kt.Device
	jchfChan chan []*kt.JCHF
}

type STListenerMetric struct {
	Messages go_metrics.Meter
	Errors   go_metrics.Meter
}

type Worker struct {
	address  string
	tagStore *tagNode
}

type tagNode struct {
	elem     *gnmiLib.PathElem
	tagName  string
	value    *gnmiLib.TypedValue
	tagStore map[string][]*tagNode
}

type tagResults struct {
	names  []string
	values []*gnmiLib.TypedValue
}

func NewSTListener(ctx context.Context, config *ktranslate.KentikSTConfig, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF, apic *api.KentikApi) (*KentikSTListener, error) {
	ks := KentikSTListener{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "St"}, log),
		jchfChan: jchfChan,
		metrics: STListenerMetric{
			Messages: go_metrics.GetOrRegisterMeter(fmt.Sprintf("st_messages^force=true"), registry),
			Errors:   go_metrics.GetOrRegisterMeter(fmt.Sprintf("st_errors^force=true"), registry),
		},
		apic:    apic,
		devices: apic.GetDevicesAsMap(0),
	}

	go ks.run(ctx)
	return &ks, nil
}

func (c *KentikSTListener) run(ctx context.Context) error {
	var err error
	var tlscfg *tls.Config
	var request *gnmiLib.SubscribeRequest
	var ctxC context.Context
	ctxC, c.cancel = context.WithCancel(ctx)

	for i := len(c.config.Subscriptions) - 1; i >= 0; i-- {
		subscription := c.config.Subscriptions[i]
		if err = subscription.BuildFullPath(c.config); err != nil {
			return err
		}
	}
	for idx := range c.config.TagSubscriptions {
		if err = c.config.TagSubscriptions[idx].BuildFullPath(c.config); err != nil {
			return err
		}
		if len(c.config.TagSubscriptions[idx].Elements) == 0 {
			return fmt.Errorf("tag_subscription must have at least one element")
		}
	}

	// Validate configuration
	if request, err = c.newSubscribeRequest(); err != nil {
		return err
	} else if (time.Duration(c.config.RedialSec) * time.Second).Nanoseconds() <= 0 {
		return fmt.Errorf("redial duration must be positive")
	}

	// Parse TLS config
	if c.config.EnableTLS {
		if tlscfg, err = c.config.TLSConfig(); err != nil {
			return err
		}
	}

	if len(c.config.Username) > 0 {
		ctxC = metadata.AppendToOutgoingContext(ctxC, "username", c.config.Username, "password", c.config.Password)
	}

	// Invert explicit alias list and prefill subscription names
	c.internalAliases = make(map[string]string, len(c.config.Subscriptions)+len(c.config.Aliases)+len(c.config.TagSubscriptions))

	for _, s := range c.config.Subscriptions {
		if err := s.BuildAlias(c.internalAliases); err != nil {
			return err
		}
	}
	for _, s := range c.config.TagSubscriptions {
		if err := s.BuildAlias(c.internalAliases); err != nil {
			return err
		}
	}

	for alias, encodingPath := range c.config.Aliases {
		c.internalAliases[encodingPath] = alias
	}

	// Create a goroutine for each device, dial and subscribe
	c.wg.Add(len(c.config.Addresses))
	for _, addr := range c.config.Addresses {
		worker := Worker{address: addr}
		worker.tagStore = &tagNode{}
		go func(worker Worker) {
			defer c.wg.Done()
			for ctxC.Err() == nil {
				if err := c.subscribeGNMI(ctxC, &worker, tlscfg, request); err != nil && ctxC.Err() == nil {
					c.metrics.Errors.Mark(1)
					c.Errorf("Error subscribing: %v", err)
				}

				select {
				case <-ctx.Done():
				case <-time.After((time.Duration(c.config.RedialSec) * time.Second)):
				}
			}
		}(worker)
	}
	return nil
}

// Create a new gNMI SubscribeRequest
func (c *KentikSTListener) newSubscribeRequest() (*gnmiLib.SubscribeRequest, error) {
	// Create subscription objects
	var err error
	subscriptions := make([]*gnmiLib.Subscription, len(c.config.Subscriptions)+len(c.config.TagSubscriptions))
	for i, subscription := range c.config.TagSubscriptions {
		if subscriptions[i], err = subscription.BuildSubscription(); err != nil {
			return nil, err
		}
	}
	for i, subscription := range c.config.Subscriptions {
		if subscriptions[i+len(c.config.TagSubscriptions)], err = subscription.BuildSubscription(); err != nil {
			return nil, err
		}
	}

	// Construct subscribe request
	gnmiPath, err := parsePath(c.config.Origin, c.config.Prefix, c.config.Target)
	if err != nil {
		return nil, err
	}

	if c.config.Encoding != "proto" && c.config.Encoding != "json" && c.config.Encoding != "json_ietf" && c.config.Encoding != "bytes" {
		return nil, fmt.Errorf("unsupported encoding %s", c.config.Encoding)
	}

	return &gnmiLib.SubscribeRequest{
		Request: &gnmiLib.SubscribeRequest_Subscribe{
			Subscribe: &gnmiLib.SubscriptionList{
				Prefix:       gnmiPath,
				Mode:         gnmiLib.SubscriptionList_STREAM,
				Encoding:     gnmiLib.Encoding(gnmiLib.Encoding_value[strings.ToUpper(c.config.Encoding)]),
				Subscription: subscriptions,
				UpdatesOnly:  c.config.UpdatesOnly,
			},
		},
	}, nil
}

// SubscribeGNMI and extract telemetry data
func (c *KentikSTListener) subscribeGNMI(ctx context.Context, worker *Worker, tlscfg *tls.Config, request *gnmiLib.SubscribeRequest) error {
	var creds credentials.TransportCredentials
	if tlscfg != nil {
		creds = credentials.NewTLS(tlscfg)
	} else {
		creds = insecure.NewCredentials()
	}
	opt := grpc.WithTransportCredentials(creds)

	client, err := grpc.DialContext(ctx, worker.address, opt)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	subscribeClient, err := gnmiLib.NewGNMIClient(client).Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup subscription: %v", err)
	}

	if err = subscribeClient.Send(request); err != nil {
		// If io.EOF is returned, the stream may have ended and stream status
		// can be determined by calling Recv.
		if err != io.EOF {
			return fmt.Errorf("failed to send subscription request: %v", err)
		}
	}

	c.Debugf("Connection to gNMI device %s established", worker.address)
	defer c.Debugf("Connection to gNMI device %s closed", worker.address)
	for ctx.Err() == nil {
		var reply *gnmiLib.SubscribeResponse
		if reply, err = subscribeClient.Recv(); err != nil {
			if err != io.EOF && ctx.Err() == nil {
				return fmt.Errorf("aborted gNMI subscription: %v", err)
			}
			break
		}

		c.handleSubscribeResponse(worker, reply)
	}
	return nil
}

func (c *KentikSTListener) handleSubscribeResponse(worker *Worker, reply *gnmiLib.SubscribeResponse) {
	if response, ok := reply.Response.(*gnmiLib.SubscribeResponse_Update); ok {
		c.handleSubscribeResponseUpdate(worker, response)
	}
}

// Handle SubscribeResponse_Update message from gNMI and parse contained telemetry data
func (c *KentikSTListener) handleSubscribeResponseUpdate(worker *Worker, response *gnmiLib.SubscribeResponse_Update) {
	var prefix, prefixAliasPath string
	timestamp := time.Unix(0, response.Update.Timestamp)
	out := make([]*kt.JCHF, 0, len(response.Update.Update))
	prefixTags := make(map[string]string)

	if response.Update.Prefix != nil {
		var err error
		if prefix, prefixAliasPath, err = handlePath(response.Update.Prefix, prefixTags, c.internalAliases, ""); err != nil {
			c.Errorf("handling path %q failed: %v", response.Update.Prefix, err)
		}
	}
	prefixTags["source"], _, _ = net.SplitHostPort(worker.address)
	prefixTags["path"] = prefix

	// Process and remove tag-only updates from the response
	for i := len(response.Update.Update) - 1; i >= 0; i-- {
		update := response.Update.Update[i]
		fullPath := pathWithPrefix(response.Update.Prefix, update.Path)
		for _, tagSub := range c.config.TagSubscriptions {
			if equalPathNoKeys(fullPath, tagSub.GetFullPath()) {
				worker.storeTags(update, tagSub)
				response.Update.Update = append(response.Update.Update[:i], response.Update.Update[i+1:]...)
			}
		}
	}

	// Parse individual Update message and create measurements
	var name, lastAliasPath string
	for _, update := range response.Update.Update {
		c.metrics.Messages.Mark(1)
		in := kt.NewJCHF()
		in.CustomStr = make(map[string]string)
		in.CustomInt = make(map[string]int32)
		in.CustomBigInt = make(map[string]int64)
		in.Provider = kt.ProviderST
		in.Timestamp = timestamp.Unix()
		in.SrcAddr = prefixTags["source"]
		if dev, ok := c.devices[in.SrcAddr]; ok {
			in.DeviceName = dev.Name // Copy in any of these info we get
			in.DeviceId = dev.ID
			in.CompanyId = dev.CompanyID
			dev.SetUserTags(in.CustomStr)
		}

		fullPath := pathWithPrefix(response.Update.Prefix, update.Path)

		// Prepare tags from prefix
		tags := make(map[string]string, len(prefixTags))
		for key, val := range prefixTags {
			tags[key] = val
		}
		aliasPath, fields := c.handleTelemetryField(update, tags, prefix)

		if tagOnlyTags := worker.checkTags(fullPath); tagOnlyTags != nil {
			for k, v := range tagOnlyTags {
				if alias, ok := c.internalAliases[k]; ok {
					tags[alias] = fmt.Sprint(v)
				} else {
					tags[k] = fmt.Sprint(v)
				}
			}
		}

		// Inherent valid alias from prefix parsing
		if len(prefixAliasPath) > 0 && len(aliasPath) == 0 {
			aliasPath = prefixAliasPath
		}

		// Lookup alias if alias-path has changed
		if aliasPath != lastAliasPath {
			name = prefix
			if alias, ok := c.internalAliases[aliasPath]; ok {
				name = alias
			} else {
				c.Debugf("No measurement alias for gNMI path: %s", name)
			}
		}

		in.CustomStr["Name"] = name
		for tk, tv := range tags {
			in.CustomStr[tk] = tv
		}

		// Group metrics
		for k, v := range fields {
			key := k
			if len(aliasPath) < len(key) && len(aliasPath) != 0 {
				// This may not be an exact prefix, due to naming style
				// conversion on the key.
				key = key[len(aliasPath)+1:]
			} else if len(aliasPath) >= len(key) {
				// Otherwise use the last path element as the field key.
				key = path.Base(key)

				// If there are no elements skip the item; this would be an
				// invalid message.
				key = strings.TrimLeft(key, "/.")
				if key == "" {
					c.Errorf("invalid empty path: %q", k)
					continue
				}
			}

			in.CustomBigInt[key] = v.(int64)
		}

		out = append(out, in)
		lastAliasPath = aliasPath
	}

	c.jchfChan <- out
}

// HandleTelemetryField and add it to a measurement
func (c *KentikSTListener) handleTelemetryField(update *gnmiLib.Update, tags map[string]string, prefix string) (string, map[string]interface{}) {
	gpath, aliasPath, err := handlePath(update.Path, tags, c.internalAliases, prefix)
	if err != nil {
		c.Errorf("handling path %q failed: %v", update.Path, err)
	}
	fields, err := gnmiToFields(strings.Replace(gpath, "-", "_", -1), update.Val)
	if err != nil {
		c.Errorf("error parsing update value %q: %v", update.Val, err)
	}
	return aliasPath, fields
}

// Parse path to path-buffer and tag-field
func handlePath(gnmiPath *gnmiLib.Path, tags map[string]string, aliases map[string]string, prefix string) (pathBuffer string, aliasPath string, err error) {
	builder := bytes.NewBufferString(prefix)

	// Prefix with origin
	if len(gnmiPath.Origin) > 0 {
		if _, err := builder.WriteString(gnmiPath.Origin); err != nil {
			return "", "", err
		}
		if _, err := builder.WriteRune(':'); err != nil {
			return "", "", err
		}
	}

	// Parse generic keys from prefix
	for _, elem := range gnmiPath.Elem {
		if len(elem.Name) > 0 {
			if _, err := builder.WriteRune('/'); err != nil {
				return "", "", err
			}
			if _, err := builder.WriteString(elem.Name); err != nil {
				return "", "", err
			}
		}
		name := builder.String()

		if _, exists := aliases[name]; exists {
			aliasPath = name
		}
		if tags != nil {
			for key, val := range elem.Key {
				key = strings.ReplaceAll(key, "-", "_")

				// Use short-form of key if possible
				if _, exists := tags[key]; exists {
					tags[name+"/"+key] = val
				} else {
					tags[key] = val
				}
			}
		}
	}

	return builder.String(), aliasPath, nil
}

// ParsePath from XPath-like string to gNMI path structure
func parsePath(origin string, pathToParse string, target string) (*gnmiLib.Path, error) {
	gnmiPath, err := xpath.ToGNMIPath(pathToParse)
	if err != nil {
		return nil, err
	}
	gnmiPath.Origin = origin
	gnmiPath.Target = target
	return gnmiPath, err
}

// Stop listener and cleanup
func (c *KentikSTListener) Close() {
	c.cancel()
	c.wg.Wait()
}

// equalPathNoKeys checks if two gNMI paths are equal, without keys
func equalPathNoKeys(a *gnmiLib.Path, b *gnmiLib.Path) bool {
	if len(a.Elem) != len(b.Elem) {
		return false
	}
	for i := range a.Elem {
		if a.Elem[i].Name != b.Elem[i].Name {
			return false
		}
	}
	return true
}

func pathKeys(gpath *gnmiLib.Path) []*gnmiLib.PathElem {
	var newPath []*gnmiLib.PathElem
	for _, elem := range gpath.Elem {
		if elem.Key != nil {
			newPath = append(newPath, elem)
		}
	}
	return newPath
}

func pathWithPrefix(prefix *gnmiLib.Path, gpath *gnmiLib.Path) *gnmiLib.Path {
	if prefix == nil {
		return gpath
	}
	fullPath := new(gnmiLib.Path)
	fullPath.Origin = prefix.Origin
	fullPath.Target = prefix.Target
	fullPath.Elem = append(prefix.Elem, gpath.Elem...)
	return fullPath
}

func (w *Worker) storeTags(update *gnmiLib.Update, sub ktranslate.STTagSubscription) {
	updateKeys := pathKeys(update.Path)
	var foundKey bool
	for _, requiredKey := range sub.Elements {
		foundKey = false
		for _, elem := range updateKeys {
			if elem.Name == requiredKey {
				foundKey = true
			}
		}
		if !foundKey {
			return
		}
	}
	// All required keys present for this TagSubscription
	w.tagStore.insert(updateKeys, sub.Name, update.Val)
}

func (node *tagNode) insert(keys []*gnmiLib.PathElem, name string, value *gnmiLib.TypedValue) {
	if len(keys) == 0 {
		node.value = value
		node.tagName = name
		return
	}
	var found *tagNode
	key := keys[0]
	keyName := key.Name
	if node.tagStore == nil {
		node.tagStore = make(map[string][]*tagNode)
	}
	if _, ok := node.tagStore[keyName]; !ok {
		node.tagStore[keyName] = make([]*tagNode, 0)
	}
	for _, node := range node.tagStore[keyName] {
		if compareKeys(node.elem.Key, key.Key) {
			found = node
			break
		}
	}
	if found == nil {
		found = &tagNode{elem: keys[0]}
		node.tagStore[keyName] = append(node.tagStore[keyName], found)
	}
	found.insert(keys[1:], name, value)
}

func (node *tagNode) retrieve(keys []*gnmiLib.PathElem, tagResults *tagResults) {
	if node.value != nil {
		tagResults.names = append(tagResults.names, node.tagName)
		tagResults.values = append(tagResults.values, node.value)
	}
	for _, key := range keys {
		if elems, ok := node.tagStore[key.Name]; ok {
			for _, node := range elems {
				if compareKeys(node.elem.Key, key.Key) {
					node.retrieve(keys, tagResults)
				}
			}
		}
	}
}

func (w *Worker) checkTags(fullPath *gnmiLib.Path) map[string]interface{} {
	results := &tagResults{}
	w.tagStore.retrieve(pathKeys(fullPath), results)
	tags := make(map[string]interface{})
	for idx := range results.names {
		vals, _ := gnmiToFields(results.names[idx], results.values[idx])
		for k, v := range vals {
			tags[k] = v
		}
	}
	return tags
}

func gnmiToFields(name string, updateVal *gnmiLib.TypedValue) (map[string]interface{}, error) {
	var value interface{}
	var jsondata []byte

	// Make sure a value is actually set
	if updateVal == nil || updateVal.Value == nil {
		return nil, nil
	}

	switch val := updateVal.Value.(type) {
	case *gnmiLib.TypedValue_AsciiVal:
		value = val.AsciiVal
	case *gnmiLib.TypedValue_BoolVal:
		value = val.BoolVal
	case *gnmiLib.TypedValue_BytesVal:
		value = val.BytesVal
	case *gnmiLib.TypedValue_DoubleVal:
		value = val.DoubleVal
	case *gnmiLib.TypedValue_DecimalVal:
		//nolint:staticcheck // to maintain backward compatibility with older gnmi specs
		value = float64(val.DecimalVal.Digits) / math.Pow(10, float64(val.DecimalVal.Precision))
	case *gnmiLib.TypedValue_FloatVal:
		//nolint:staticcheck // to maintain backward compatibility with older gnmi specs
		value = val.FloatVal
	case *gnmiLib.TypedValue_IntVal:
		value = val.IntVal
	case *gnmiLib.TypedValue_StringVal:
		value = val.StringVal
	case *gnmiLib.TypedValue_UintVal:
		value = val.UintVal
	case *gnmiLib.TypedValue_JsonIetfVal:
		jsondata = val.JsonIetfVal
	case *gnmiLib.TypedValue_JsonVal:
		jsondata = val.JsonVal
	}

	fields := make(map[string]interface{})
	if value != nil {
		fields[name] = value
	} else if jsondata != nil {
		if err := json.Unmarshal(jsondata, &value); err != nil {
			return nil, fmt.Errorf("failed to parse JSON value: %v", err)
		}
		fields[name] = value
	}
	return fields, nil
}

func compareKeys(a map[string]string, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if _, ok := b[k]; !ok {
			return false
		}
		if b[k] != v {
			return false
		}
	}
	return true
}
