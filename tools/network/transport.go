package network

import (
	"crypto/tls"
	"crypto/x509"
)

type TransportCerts struct {
	Cert     tls.Certificate
	CertPool *x509.CertPool
}

type TransportKey struct {
	public, private, obCert []byte
}

func (t *TransportKey) SetPublic(public []byte) *TransportKey {
	t.public = public
	return t
}
func (t *TransportKey) SetPrivate(private []byte) *TransportKey {
	t.private = private
	return t
}
func (t *TransportKey) SetOBCert(obCert []byte) *TransportKey {
	t.obCert = obCert
	return t
}
func (t *TransportKey) Build() TransportCerts {
	cert, err := tls.X509KeyPair(t.public, t.private)
	if err != nil {
		return TransportCerts{}
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(t.obCert)

	return TransportCerts{cert, caCertPool}
}
