package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func GetTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("Failed to add ca cert to pool")
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}, nil
}
