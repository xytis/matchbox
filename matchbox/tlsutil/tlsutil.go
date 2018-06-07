package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

// NewCertPool creates x509 certPool with provided CA files.
func NewCertPool(CAFiles []string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, CAFile := range CAFiles {
		pemByte, err := ioutil.ReadFile(CAFile)
		if err != nil {
			return nil, err
		}

		for {
			var block *pem.Block
			block, pemByte = pem.Decode(pemByte)
			if block == nil {
				break
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certPool.AddCert(cert)
		}
	}

	return certPool, nil
}

// NewCert generates TLS cert by using the given cert and key.
func NewCert(certfile, keyfile string) (*tls.Certificate, error) {
	cert, err := ioutil.ReadFile(certfile)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return &tlsCert, nil
}

// StaticClientCertificate Builds a client certificate responder function with static certificate
func StaticClientCertificate(cert *tls.Certificate) func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return func(unused *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return cert, nil
	}
}

// StaticServerCertificate returns static server certificate
func StaticServerCertificate(cert *tls.Certificate) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(unused *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return cert, nil
	}
}
