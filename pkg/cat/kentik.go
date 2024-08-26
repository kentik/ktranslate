package cat

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/kt"
	model "github.com/kentik/ktranslate/pkg/util/kflow2"

	capn "zombiezen.com/go/capnproto2"
)

const (
	kentikDefaultCapnprotoDecodeLimit = 128 << 20 // 128 MiB
)

// Handler for json data, useful for testing mostly. Requires you to set content-type: application/json
func (kc *KTranslate) handleJson(cid kt.Cid, raw []byte) error {
	serBuf := make([]byte, 0)
	select {
	case jflow := <-kc.jchfChans[0]: // non blocking select on this chan.
		var base map[string]interface{}
		err := json.Unmarshal(raw, &base)
		if err != nil {
			return err
		} else {
			jflow.CustomStr = make(map[string]string)
			jflow.CustomInt = make(map[string]int32)
			jflow.CustomBigInt = make(map[string]int64)
			jflow.Provider = kt.ProviderAlert
			jflow.EventType = kt.KENTIK_EVENT_JSON

			// map any fields found into the jflow obj.
			for k, v := range base {
				switch tv := v.(type) {
				case string:
					jflow.CustomStr[k] = tv
				case int:
					jflow.CustomBigInt[k] = int64(tv)
				case uint32:
					jflow.CustomBigInt[k] = int64(tv)
				case uint64:
					jflow.CustomBigInt[k] = int64(tv)
				case int64:
					jflow.CustomBigInt[k] = int64(tv)
				case float64:
					jflow.CustomBigInt[k] = int64(tv)
				case int32:
					jflow.CustomInt[k] = tv
				case map[string]interface{}:
					for kk, sv := range tv {
						switch it := sv.(type) {
						case string:
							jflow.CustomStr[k+"_"+kk] = it
						case float64:
							jflow.CustomBigInt[k+"_"+kk] = int64(it)
						case map[string]interface{}:
							for ik, iv := range it {
								key := fmt.Sprintf("%s_%s_%s", k, kk, ik)
								switch iit := iv.(type) {
								case string:
									jflow.CustomStr[key] = iit
								case float64:
									jflow.CustomBigInt[key] = int64(iit)
								default:
									kc.log.Warnf("Unhandled json type 1: %s", sv)
								}
							}
						default:
							kc.log.Warnf("Unhandled json type 2: %s", tv)
						}
					}
				case []interface{}:
					for i, sv := range tv {
						switch it := sv.(type) {
						case map[string]interface{}:
							for ik, iv := range it {
								key := fmt.Sprintf("%s_%d_%s", k, i, ik)
								switch iit := iv.(type) {
								case string:
									jflow.CustomStr[key] = iit
								case float64:
									jflow.CustomBigInt[key] = int64(iit)
								default:
									kc.log.Warnf("Unhandled json type 3: %s", iv)
								}
							}
						}
					}
				default:
					kc.log.Warnf("Unhandled json type 4: %s", v)
				}
			}
			res, err := kc.format.To([]*kt.JCHF{jflow}, serBuf)
			jflow.Reset()
			kc.jchfChans[0] <- jflow
			if err != nil {
				return err
			}
			kc.msgsc <- res // Send  on without batching.
		}
	default: // We're out of batched flows, just drop this one.
		kc.metrics.DroppedFlows.Mark(1)
	}
	return nil
}

// Take flow from http requests, deserialize and pass it on to alphaChan
// Gets called from a goroutine-per-request
func (kc *KTranslate) handleFlow(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		return
	}

	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			kc.metrics.Errors.Mark(1)
			kc.log.Errorf("Error handling request: %v", err)
			fmt.Fprint(w, "BAD") // nolint: errcheck
		} else {
			fmt.Fprint(w, "GOOD") // nolint: errcheck
		}
	}()

	// Decode body in gzip format if the request header is set this way.
	body := r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		z, err := gzip.NewReader(r.Body)
		if err != nil {
			kc.log.Errorf("There was an eror when decompressing the content: %+v.", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body = z
	}
	defer body.Close()

	// check company id and other values.
	vals := r.URL.Query()
	senderId := vals.Get(HttpSenderID)
	cidBase := vals.Get(HttpCompanyID)
	cid := 0
	if cidBase != "" {
		if c, err := strconv.Atoi(cidBase); err != nil {
			kc.log.Errorf("There was an error when getting cid paramiter: %s %v.", cidBase, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			cid = c
		}
	}

	// Allocate a buffer for the expected size of the incoming data.
	var bodyBufferBytes []byte
	contentLengthString := r.Header.Get("Content-Length")
	if contentLengthString != "" {
		size, err := strconv.Atoi(contentLengthString)
		if err != nil {
			kc.log.Errorf("There was an error when getting the content length: %v.", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if size > 0 &&
			size < MaxProxyListenerBufferAlloc { // limit in case attacker sets Content-Length
			// superstitiously add extra breathing room to buffer just in case
			bodyBufferBytes = make([]byte, 0, size+(2*bytes.MinRead))
		}
	}

	// Read all data from the request (possibly gzip decoding, possibly not)
	bodyBuffer := bytes.NewBuffer(bodyBufferBytes)
	_, err = bodyBuffer.ReadFrom(body)

	if err != nil {
		kc.log.Errorf("There was an error when reading the body content: %v.", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	evt := bodyBuffer.Bytes()

	// If its http/json data, treat sperately
	if r.Header.Get("Content-Type") == "application/json" {
		err = kc.handleJson(kt.Cid(cid), evt)
		return
	}

	// If we are sending from before kentik, add offset in here.
	offset := 0
	did := 0
	deviceName := ""
	if senderId != "" && len(evt) > MSG_KEY_PREFIX && // Direct flow without enrichment.
		(evt[0] == 0x00 && evt[1] == 0x00 && evt[2] == 0x00 && evt[3] == 0x00 && evt[4] == 0x00) { // Double check with this
		offset = MSG_KEY_PREFIX
		pts := strings.Split(senderId, ":")
		if len(pts) == 3 {
			cid, _ = strconv.Atoi(strings.TrimSpace(pts[0]))
			deviceName = pts[1]
			did, _ = strconv.Atoi(strings.TrimSpace(pts[2]))
		}

	}

	// Tee any flows on to another ktrans instance if this is set up.
	if kc.tee != nil {
		kc.tee.Send(r.Context(), kt.NewOutputWithProviderAndCompanySender(evt[offset:], kt.ProviderKflow, kt.Cid(cid), kt.EventOutput, senderId))
	}

	// decompress and read (capnproto "packed" representation)
	decoder := capn.NewPackedDecoder(bytes.NewBuffer(evt[offset:]))
	decoder.MaxMessageSize = kentikDefaultCapnprotoDecodeLimit
	capnprotoMessage, err := decoder.Decode()
	if err != nil {
		return
	}

	// unpack flow messages and pass them down
	packedCHF, err := model.ReadRootPackedCHF(capnprotoMessage)
	if err != nil {
		return
	}

	messages, err := packedCHF.Msgs()
	if err != nil {
		return
	}

	var sent, dropped int64
	next := 0
	for i := 0; i < messages.Len(); i++ {
		msg := messages.At(i)
		if !msg.Big() { // Don't work on low res data
			if !msg.SampleAdj() {
				msg.SetSampleRate(msg.SampleRate() * 100) // Apply re-sample trick here.
			}
			if msg.DeviceId() == 0 && senderId != "" {
				// Fill in from the parsed senderId.
				msg.SetDeviceId(uint32(did))
			}

			// send without blocking, dropping the message if the channel buffer is full
			alpha := &Flow{CompanyId: cid, CHF: msg, DeviceName: deviceName}
			select {
			case kc.alphaChans[next] <- alpha:
				sent++
			default:
				dropped++
			}
			next++ // Round robin across processing threads.
			if next >= kc.config.ProcessingThreads {
				next = 0
			}
		}
	}
	kc.metrics.Flows.Mark(sent)
	kc.metrics.DroppedFlows.Mark(dropped)
}

func (kc *KTranslate) monitorAlphaChan(ctx context.Context, i int, seri func([]*kt.JCHF, []byte) (*kt.Output, error)) {
	cacheTicker := time.NewTicker(CacheInvalidateDuration)
	defer cacheTicker.Stop()

	sendTicker := time.NewTicker(kt.SendBatchDuration)
	defer sendTicker.Stop()

	// Set up some data structures.
	tagcache := map[uint64]string{}
	serBuf := make([]byte, 0)
	msgs := make([]*kt.JCHF, 0)
	sendBytesOn := func() {
		if len(msgs) == 0 {
			return
		}

		// Add in any extra things here.
		if kc.geo != nil || kc.asn != nil || kc.enricher != nil {
			msgs = kc.doEnrichments(ctx, msgs)
		}

		// If we have any rollups defined, send here instead of directly to the output format.
		if kc.doRollups {
			rv := make([]map[string]interface{}, len(msgs))
			for i, msg := range msgs {
				rv[i] = msg.ToMap()
			}
			for _, r := range kc.rollups {
				r.Add(rv)
			}
		}

		// Turn into a binary format here, using the passed in encoder.
		if !kc.doRollups || kc.config.RollupAndAlpha {
			// Compute and sample rate stuff here.
			keep := len(msgs)
			if kc.config.SampleRate > 1 && keep > kc.config.SampleMin {
				rand.Shuffle(len(msgs), func(i, j int) {
					msgs[i], msgs[j] = msgs[j], msgs[i]
				})
				keep = int(math.Max(float64(len(msgs))/float64(kc.config.SampleRate), 1))
				for _, msg := range msgs {
					msg.SampleRate = msg.SampleRate * uint32(kc.config.SampleRate)
				}
				kc.log.Debugf("Reduced input from %d to %d", len(msgs), keep)
			}

			// Need to only serialize flows from 1 company at a time.
			splits := map[kt.Cid][]*kt.JCHF{}
			for _, msg := range msgs[0:keep] {
				if _, ok := splits[msg.CompanyId]; !ok {
					splits[msg.CompanyId] = make([]*kt.JCHF, 0, keep)
				}
				splits[msg.CompanyId] = append(splits[msg.CompanyId], msg)
			}

			for _, cidl := range splits {
				ser, err := seri(cidl, serBuf)
				if err != nil {
					kc.log.Errorf("There was an error when converting to native: %v.", err)
				} else {
					kc.msgsc <- ser
				}
			}
		}

		for _, m := range msgs { // Give back our cache.
			m.Reset()
			kc.jchfChans[i] <- m
		}

		// match in with out.
		kc.metrics.FlowsOut.Mark(int64(len(msgs)))
		msgs = make([]*kt.JCHF, 0)
	}

	currentTime := time.Now().Unix() // Record rough time of flow sent.
	kc.log.Infof("monitorAlpha %d Online", i)
	for {
		select {
		case f := <-kc.alphaChans[i]:
			select {
			case jflow := <-kc.jchfChans[i]: // non blocking select on this chan.
				err := kc.flowToJCHF(ctx, jflow, f, currentTime, tagcache)
				if err != nil {
					kc.log.Errorf("There was an error when converting to json: %v.", err)
					jflow.Reset()
					kc.jchfChans[i] <- jflow
					continue
				}
				keep := true
				for _, f := range kc.filters {
					if !f.Filter(jflow) {
						keep = false
						break
					}
				}
				if keep {
					msgs = append(msgs, jflow) // Batch up here.
					if len(msgs) >= kc.config.MaxFlowsPerMessage {
						sendBytesOn()
					}
				} else {
					jflow.Reset()
					kc.jchfChans[i] <- jflow // Toss this guy, he doesn't meet out filter.
				}

			default: // We're out of batched flows, send what we have and re-q this one.
				sendBytesOn()
				select {
				case kc.alphaChans[i] <- f:
				default:
					kc.metrics.DroppedFlows.Mark(1)
				}
			}
		case _ = <-sendTicker.C: // Send on here.
			sendBytesOn() // Has context for everything it needs.
			currentTime = time.Now().Unix()

		case <-cacheTicker.C:
			tagcache = map[uint64]string{}

		case <-ctx.Done():
			kc.log.Infof("monitorAlpha %d Done", i)
			return
		}
	}
}
