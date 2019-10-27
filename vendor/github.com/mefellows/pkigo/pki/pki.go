package pki

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var CertificatePreamble = "-----BEGIN CERTIFICATE-----"
var KeyPreamble = "-----BEGIN RSA PRIVATE KEY-----"

type Pki struct {
	sync.Mutex
	ClientTlsConfig *tls.Config
	ServerTlsConfig *tls.Config
	BaseDir         string
}

var PkiConfig Pki

func (m *Pki) SetBaseDir(baseDir string) {
	m.BaseDir = baseDir
}

func (m *Pki) SetClientTLSConfig(config *tls.Config) {
	m.ClientTlsConfig = config
}

func (m *Pki) SetServerTLSConfig(config *tls.Config) {
	m.ServerTlsConfig = config
}

// General Pki Pubic Key Infrastructure functions

type PKI struct {
	Config *Config
}

type Config struct {
	Application    string
	ClientKeyPath  string
	ClientCertPath string
	ServerKeyPath  string
	ServerCertPath string
	CaKeyPath      string
	CaCertPath     string
	Insecure       bool
}

func NewWithConfig(config *Config) (*PKI, error) {
	pki := &PKI{Config: config}
	if err := pki.SetupPKI("localhost"); err != nil {
		return nil, err
	}
	return pki, nil
}

func New() (*PKI, error) {
	pki := &PKI{Config: getDefaultConfig()}
	if err := pki.SetupPKI("localhost"); err != nil {
		return nil, err
	}
	return pki, nil
}

func getDefaultConfig() *Config {
	caHomeDir := GetCADir()
	certDir := GetCertDir()
	caCertPath := filepath.Join(caHomeDir, "ca.pem")
	caKeyPath := filepath.Join(caHomeDir, "key.pem")
	certPath := filepath.Join(certDir, "cert.pem")
	keyPath := filepath.Join(certDir, "cert-key.pem")
	serverCertPath := filepath.Join(certDir, "server-cert.pem")
	serverKeyPath := filepath.Join(certDir, "server-key.pem")

	return &Config{
		ClientKeyPath:  keyPath,
		ClientCertPath: certPath,
		CaCertPath:     caCertPath,
		CaKeyPath:      caKeyPath,
		ServerCertPath: serverCertPath,
		ServerKeyPath:  serverKeyPath,
	}
}

func (p *PKI) RemovePKI() error {
	// Root CA + Certificates
	err := os.RemoveAll(filepath.Dir(p.Config.CaCertPath))
	if err != nil {
		return err
	}

	// Client certificates
	err = os.RemoveAll(filepath.Dir(p.Config.ClientCertPath))
	if err != nil {
		return err
	}

	// Server certificates
	err = os.RemoveAll(filepath.Dir(p.Config.ServerKeyPath))
	if err != nil {
		return err
	}

	return err
}

func (p *PKI) GenerateClientCertificate(hosts []string) (err error) {
	organisation := "client"
	bits := 2048

	if len(hosts) == 0 {
		hosts = []string{}
	}
	err = GenerateCertificate(hosts, p.Config.ClientCertPath, p.Config.ClientKeyPath, p.Config.CaCertPath, p.Config.CaKeyPath, organisation, bits)
	if err == nil {
		_, err = os.Stat(p.Config.ClientCertPath)
		_, err = os.Stat(p.Config.ClientKeyPath)
	}
	return
}

// Validate all components of the PKI infrastructure are properly configured
func (p *PKI) CheckSetup() error {
	var err error

	// Check directories
	if _, err = os.Stat(p.Config.CaCertPath); err == nil {
		return nil
	}

	// Check CA

	// Check server cert

	// Check client certs against CA (from conf + user?)

	// Check permissions?

	return err
}

// Sets up the PKI infrastructure for client / server communications
// This involves creating directories, CAs, and client/server certs
func (p *PKI) SetupPKI(caHost string) error {
	if p.CheckSetup() == nil {
		return nil
	}
	log.Printf("Setting up PKI for '%s'...", caHost)

	bits := 2048
	if _, err := os.Stat(p.Config.CaCertPath); err == nil {
		return fmt.Errorf("CA already exists. Run --delete to remove the old CA.")
	}

	os.MkdirAll(filepath.Dir(p.Config.CaCertPath), 0700)
	if err := GenerateCACertificate(p.Config.CaCertPath, p.Config.CaKeyPath, caHost, bits); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(p.Config.CaCertPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	if _, err := os.Stat(p.Config.CaKeyPath); err != nil {
		return fmt.Errorf("Couldn't generate CA Certificate: %s", err.Error())
	}

	organisation := caHost
	hosts := []string{caHost}

	os.MkdirAll(filepath.Dir(p.Config.ServerCertPath), 0700)
	err := GenerateCertificate(hosts, p.Config.ServerCertPath, p.Config.ServerKeyPath, p.Config.CaCertPath, p.Config.CaKeyPath, organisation, bits)
	if err == nil {
		_, err = os.Stat(p.Config.ServerCertPath)
		_, err = os.Stat(p.Config.ServerKeyPath)
	}

	// Setup Client side...
	p.GenerateClientCertificate([]string{"localhost"})

	return nil
}

func (p *PKI) OutputClientKey() (string, error) {
	return OutputFileContents(p.Config.ClientKeyPath)
}

func (p *PKI) OutputClientCert() (string, error) {
	return OutputFileContents(p.Config.ClientCertPath)
}

func (p *PKI) OutputCAKey() (string, error) {
	return OutputFileContents(p.Config.CaKeyPath)
}
func (p *PKI) OutputCACert() (string, error) {
	return OutputFileContents(p.Config.CaCertPath)
}

func (p *PKI) GetClientTLSConfig() (*tls.Config, error) {

	var certificates []tls.Certificate
	if !p.Config.Insecure {
		cert, err := tls.LoadX509KeyPair(p.Config.ClientCertPath, p.Config.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, cert)
	}

	certPool, err := p.discoverCAs()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates:       certificates,
		RootCAs:            certPool,
		InsecureSkipVerify: p.Config.Insecure,
	}

	return config, err
}

func (p *PKI) GetServerTLSConfig() (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(p.Config.ServerCertPath, p.Config.ServerKeyPath)
	if err != nil {
		return nil, err
	}

	certPool, err := p.discoverCAs()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		Rand:         rand.Reader,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	if p.Config.Insecure {
		config.ClientAuth = tls.NoClientCert
	}

	return config, err
}

func (p *PKI) ImportCA(name string, certPath string) error {
	// Validate name - only alphanumeric
	var nameMatch = regexp.MustCompile(`^[a-zA-Z-_\.0-9]+$`)
	if !nameMatch.MatchString(name) {
		return errors.New("CA Name must contain only alphanumeric characters")
	}

	dstCert := filepath.Join(GetCADir(), fmt.Sprintf("%s-ca.pem", name))
	cert, err := ioutil.ReadFile(certPath)

	if err != nil {
		return err
	}

	// import Cert
	if strings.Contains(string(cert), CertificatePreamble) {
		ioutil.WriteFile(dstCert, cert, 0600)
	} else {
		return errors.New(fmt.Sprintf("Certificate provided is not valid, no %s present", CertificatePreamble))
	}

	return nil
}

// Overrides the default client certificate with a new one
func (p *PKI) ImportClientCertAndKey(certPath string, keyPath string) error {
	cert, err := ioutil.ReadFile(certPath)

	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile(keyPath)

	if err != nil {
		return err
	}

	// import cert
	if strings.Contains(string(cert), CertificatePreamble) {
		err = ioutil.WriteFile(p.Config.ClientCertPath, cert, 0600)
		if err != nil {
			return err
		}

		// import key
		if strings.Contains(string(key), KeyPreamble) {
			err = ioutil.WriteFile(p.Config.ClientKeyPath, key, 0600)
		} else {
			return errors.New(fmt.Sprintf("Key provided is not valid, no %s present", KeyPreamble))
		}

	} else {
		return errors.New(fmt.Sprintf("Certificate provided is not valid, no %s present", CertificatePreamble))
	}

	return err
}

func (p *PKI) discoverCAs() (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	var caPaths []string
	var err error

	// Read in all certs from CA dir
	readFiles, err := ioutil.ReadDir(filepath.Dir(p.Config.CaCertPath))
	if err == nil {
		caPaths = make([]string, 0)

		for _, file := range readFiles {
			if strings.HasSuffix(file.Name(), ".pem") || strings.HasSuffix(file.Name(), ".crt") {
				caPaths = append(caPaths, filepath.Join(filepath.Dir(p.Config.CaCertPath), file.Name()))
			}
		}
	}

	for _, cert := range caPaths {
		pemData, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, err
		}

		// Only add certs
		if strings.Contains(string(pemData), CertificatePreamble) {
			if ok := certPool.AppendCertsFromPEM(pemData); !ok {
				return nil, err
			}
		}
	}

	return certPool, err
}
