package api

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"google.golang.org/grpc/credentials"
)

func clientTransportCredentials(insecureSkipVerify bool, cacertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	var tlsConf tls.Config

	if clientCertFile != "" {
		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %v", err)
		}
		tlsConf.Certificates = []tls.Certificate{certificate}
	}

	if insecureSkipVerify {
		tlsConf.InsecureSkipVerify = true
	} else if cacertFile != "" {
		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(cacertFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %v", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca certs")
		}

		tlsConf.RootCAs = certPool
	}

	return credentials.NewTLS(&tlsConf), nil
}

func getAddressFromApiRoot(apiRoot string) (string, error) {
	address := apiRoot
	if strings.HasPrefix(apiRoot, "https://") {
		address = address[len("https://"):]
	}
	if strings.HasPrefix(apiRoot, "http://") {
		address = address[len("http://"):]
	}
	if !strings.HasSuffix(address, ":443") {
		address = address + ":443"
	}
	return address, nil
}
