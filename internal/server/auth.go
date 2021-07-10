package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func GetTLSConfig(svrCertFile, svrKeyFile, caCertFile string) (*tls.Config, error) {
	svrCert, err := tls.LoadX509KeyPair(svrCertFile, svrKeyFile)
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
		Certificates: []tls.Certificate{svrCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}, nil
}
