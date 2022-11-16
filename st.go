package ktranslate

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/gnxi/utils/xpath"
	gnmiLib "github.com/openconfig/gnmi/proto/gnmi"
)

func (s *STSubscription) GetFullPath() *gnmiLib.Path {
	return s.fullPath
}

func (s *STSubscription) BuildFullPath(c *KentikSTConfig) error {
	var err error
	if s.fullPath, err = xpath.ToGNMIPath(s.Path); err != nil {
		return err
	}
	s.fullPath.Origin = s.Origin
	s.fullPath.Target = c.Target
	if c.Prefix != "" {
		prefix, err := xpath.ToGNMIPath(c.Prefix)
		if err != nil {
			return err
		}
		s.fullPath.Elem = append(prefix.Elem, s.fullPath.Elem...)
		if s.Origin == "" && c.Origin != "" {
			s.fullPath.Origin = c.Origin
		}
	}
	return nil
}

func (s *STSubscription) BuildAlias(aliases map[string]string) error {
	var err error
	var gnmiLongPath, gnmiShortPath *gnmiLib.Path

	// Build the subscription path without keys
	if gnmiLongPath, err = parsePath(s.Origin, s.Path, ""); err != nil {
		return err
	}
	if gnmiShortPath, err = parsePath("", s.Path, ""); err != nil {
		return err
	}

	longPath, _, err := handlePath(gnmiLongPath, nil, nil, "")
	if err != nil {
		return fmt.Errorf("handling long-path failed: %v", err)
	}
	shortPath, _, err := handlePath(gnmiShortPath, nil, nil, "")
	if err != nil {
		return fmt.Errorf("handling short-path failed: %v", err)
	}

	// If the user didn't provide a measurement name, use last path element
	name := s.Name
	if len(name) == 0 {
		name = path.Base(shortPath)
	}
	if len(name) > 0 {
		aliases[longPath] = name
		aliases[shortPath] = name
	}
	return nil
}

func (s *STSubscription) BuildSubscription() (*gnmiLib.Subscription, error) {
	gnmiPath, err := parsePath(s.Origin, s.Path, "")
	if err != nil {
		return nil, err
	}
	mode, ok := gnmiLib.SubscriptionMode_value[strings.ToUpper(s.SubscriptionMode)]
	if !ok {
		return nil, fmt.Errorf("invalid subscription mode %s", s.SubscriptionMode)
	}
	return &gnmiLib.Subscription{
		Path:              gnmiPath,
		Mode:              gnmiLib.SubscriptionMode(mode),
		HeartbeatInterval: uint64((time.Duration(s.HeartbeatIntervalSec) * time.Second).Nanoseconds()),
		SampleInterval:    uint64((time.Duration(s.SampleIntervalSec) * time.Second).Nanoseconds()),
		SuppressRedundant: s.SuppressRedundant,
	}, nil
}

func (c *KentikSTConfig) TLSConfig() (*tls.Config, error) {
	// This check returns a nil (aka, "use the default")
	// tls.Config if no field is set that would have an effect on
	// a TLS connection. That is, any of:
	//     * client certificate settings,
	//     * peer certificate authorities,
	//     * disabled security, or
	//     * an SNI server name.
	if c.TLSCA == "" && c.TLSKey == "" && c.TLSCert == "" && !c.InsecureSkipVerify && c.ServerName == "" {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify,
		Renegotiation:      tls.RenegotiateNever,
	}

	if c.TLSCA != "" {
		pool, err := makeCertPool([]string{c.TLSCA})
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = pool
	}

	if c.TLSCert != "" && c.TLSKey != "" {
		err := loadCertificate(tlsConfig, c.TLSCert, c.TLSKey)
		if err != nil {
			return nil, err
		}
	}

	if c.ServerName != "" {
		tlsConfig.ServerName = c.ServerName
	}

	return tlsConfig, nil
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

func makeCertPool(certFiles []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, certFile := range certFiles {
		pem, err := os.ReadFile(certFile)
		if err != nil {
			return nil, fmt.Errorf(
				"could not read certificate %q: %v", certFile, err)
		}
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf(
				"could not parse any PEM certificates %q: %v", certFile, err)
		}
	}
	return pool, nil
}

func loadCertificate(config *tls.Config, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf(
			"could not load keypair %s:%s: %v", certFile, keyFile, err)
	}

	config.Certificates = []tls.Certificate{cert}
	return nil
}
