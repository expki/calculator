package config

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"
)

type algorithm uint8

const (
	algorithm_unspecified algorithm = iota
	algorithm_rsa
	algorithm_ecdsa
)

type ConfigTLSPath struct {
	// config
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
	// computed
	algorithm   algorithm       `json:"-"`
	certificate tls.Certificate `json:"-"`
	// state
	mutex sync.RWMutex `json:"-"`
}

// reloadCertificate checks if the certificate is still valid and reloads it if necessary.
func (t *ConfigTLSPath) reloadCertificate() error {
	// get certificate expiration date
	t.mutex.RLock()
	var notAfter time.Time
	if t.certificate.Leaf != nil {
		notAfter = t.certificate.Leaf.NotAfter
	} else {
		notAfter = time.Now().UTC()
	}
	t.mutex.RUnlock()
	// if certificate is still valid, return
	if notAfter.Before(time.Now().UTC().Add(-30 * 24 * time.Hour)) {
		return nil
	}
	// reload certificate
	return t.loadCertificate()
}

// loadCertificate loads the certificate from the file system.
func (t *ConfigTLSPath) loadCertificate() error {
	// Generate new
	if t.CertPath == "" || t.KeyPath == "" {
		return errors.New("could not load certificate: cert_path or key_path is empty")
	}

	// Load Filesystem
	cert, err := tls.LoadX509KeyPair(t.CertPath, t.KeyPath)
	if err != nil {
		return fmt.Errorf("could not load certificate & key: %v", err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("could not parse certificate: %v", err)
	}
	t.mutex.Lock()
	t.certificate = cert
	t.mutex.Unlock()

	// Update algorithm
	t.loadAlgorithm()

	return nil
}

// loadAlgorithm loads the algorithm of the certificate.
func (t *ConfigTLSPath) loadAlgorithm() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.certificate.Leaf == nil {
		return
	}
	switch t.certificate.Leaf.PublicKey.(type) {
	case *rsa.PublicKey:
		t.algorithm = algorithm_rsa
	case *ecdsa.PublicKey:
		t.algorithm = algorithm_ecdsa
	default:
		t.algorithm = algorithm_unspecified
	}
}
