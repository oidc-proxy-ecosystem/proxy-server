package config

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/oidc-proxy-ecosystem/proxy-server/plugins/storage"
	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
	"golang.org/x/crypto/ssh"
)

type CaName struct {
	Contry             []string `yaml:"country"`
	Province           []string `yaml:"province"`
	Locality           []string `yaml:"locality"`
	OrganizationalUnit []string `yaml:"organizational_unit"`
	Organization       []string `yaml:"organization"`
	CommonName         string   `yaml:"common_name"`
	Years              int      `yaml:"years"`
	Serial             int64    `yaml:"serial"`
	Refresh            bool     `yaml:"refresh"`
}

func NewCaName(filename string) CaName {
	var config CaName
	utils.MustReadYamlExpand(filename, &config)
	if config.Years == 0 {
		config.Years = 1
	}
	return config
}

func (caName CaName) generateCertName() string {
	return "ca.crt"
}

func (caName CaName) generateKeyName() string {
	return "ca.key"
}

func (caName CaName) generatePublicName() string {
	return "ca.key.pub"
}

func (caName CaName) exists(ctx context.Context, name string) (bool, error) {
	if caName.Refresh {
		return false, nil
	}
	if v, err := storage.Client.Get(ctx, name); err != nil {
		return false, err
	} else {
		return len(v) > 0, nil
	}
}

func (caName CaName) GenerateCA() error {
	ctx := context.Background()
	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	certFileName := caName.generateCertName()
	exists, err := caName.exists(ctx, certFileName)
	if err != nil {
		return err
	}
	now := time.Now()
	caTpl := &x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(caName.Serial),
		Subject: pkix.Name{
			Country:            caName.Contry,
			Province:           caName.Province,
			Locality:           caName.Locality,
			OrganizationalUnit: caName.OrganizationalUnit,
			Organization:       caName.Organization,
			CommonName:         caName.CommonName,
		},
		NotBefore:             now.Add(-5 * time.Minute).UTC(),
		NotAfter:              now.AddDate(caName.Years, 0, 0).UTC(),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, caTpl, caTpl, prv.Public(), prv)
	if err != nil {
		return fmt.Errorf("Failed to create CA Certificate: %s", err)
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if !exists {
		if err := storage.Client.Set(ctx, certFileName, certPem); err != nil {
			return err
		}
		os.WriteFile(certFileName, certPem, 0600)
		log.Println("Written " + certFileName)
	}

	keyFileName := caName.generateKeyName()
	exists, err = caName.exists(ctx, keyFileName)
	if err != nil {
		return err
	}
	if !exists {
		prvPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(prv)})
		if err := storage.Client.Set(ctx, keyFileName, prvPem); err != nil {
			return err
		}
		os.WriteFile(keyFileName, prvPem, 0600)
		log.Println("Written " + keyFileName)
	}

	publicSSHKey, err := ssh.NewPublicKey(&prv.PublicKey)
	if err != nil {
		return err
	}
	pubFile := caName.generatePublicName()
	exists, err = caName.exists(ctx, pubFile)
	if err != nil {
		return err
	}
	if !exists {
		pubKeyBytes := ssh.MarshalAuthorizedKey(publicSSHKey)
		if err := storage.Client.Set(ctx, pubFile, pubKeyBytes); err != nil {
			return err
		}
		os.WriteFile(pubFile, pubKeyBytes, 0600)
		log.Println("Written " + pubFile)
	}
	return nil
}

func (caName CaName) ReadCertificate() (*x509.Certificate, error) {
	ctx := context.Background()
	certFileName := caName.generateCertName()
	certBuf, err := storage.Client.Get(ctx, certFileName)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certBuf)
	if block == nil {
		return nil, fmt.Errorf("invalid certificate data")
	}
	var caTpl *x509.Certificate
	if block.Type == "CERTIFICATE" {
		caTpl, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid certificate type: %s", block.Type)
	}
	return caTpl, nil
}

func (caName CaName) ReadCAPrivateKey() (*rsa.PrivateKey, error) {
	ctx := context.Background()
	keyFile := caName.generateKeyName()
	buf, err := storage.Client.Get(ctx, keyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, fmt.Errorf("invalid private key data")
	}
	var key *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	case "PRIVATE KEY":
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		key, ok = keyInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not RSA Private key")
		}
	default:
		return nil, fmt.Errorf("invalid private key type: %s", block.Type)
	}
	key.Precompute()
	return key, nil
}

func (caName CaName) ReadCASSHPublicKey() ([]byte, error) {
	return storage.Client.Get(context.Background(), caName.generatePublicName())
}

func (caName CaName) GenerateCAPublicKey() (*rsa.PublicKey, error) {
	key, err := caName.ReadCAPrivateKey()
	if err != nil {
		return nil, err
	}
	return &key.PublicKey, nil
}

func (caName CaName) SavePublicKey() error {
	pub, err := caName.GenerateCAPublicKey()
	if err != nil {
		return err
	}
	pubFileName := caName.generatePublicName()
	outFile, err := os.OpenFile(pubFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open "+pubFileName+" for writing: %s", err)
	}
	if err := pem.Encode(outFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pub)}); err != nil {
		return err
	}
	log.Println("Written " + pubFileName)
	return nil
}

func (caName CaName) GenerateClient() (*rsa.PrivateKey, *ssh.Certificate, error) {
	prvCaKey, err := caName.ReadCAPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	prvClient, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicSSHKey, err := ssh.NewPublicKey(&prvClient.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	sshCert, err := generateAndSign(prvCaKey, publicSSHKey)
	if err != nil {
		return nil, nil, err
	}
	return prvClient, sshCert, err
}

func generateCert(pub ssh.PublicKey) *ssh.Certificate {
	permissions := ssh.Permissions{
		CriticalOptions: map[string]string{},
		Extensions: map[string]string{
			"permit-pty": "",
		},
	}
	now := time.Now()

	return &ssh.Certificate{
		CertType:        ssh.UserCert,
		Permissions:     permissions,
		Key:             pub,
		ValidPrincipals: []string{"actions"},
		KeyId:           "",
		ValidBefore:     uint64(now.Add(30 * time.Minute).Unix()),
		ValidAfter:      uint64(now.Unix()),
	}
}

func generateSignerFromKey(priv *rsa.PrivateKey) (ssh.Signer, error) {
	return ssh.NewSignerFromKey(priv)
}

func generateAndSign(priv *rsa.PrivateKey, pub ssh.PublicKey) (*ssh.Certificate, error) {
	signer, err := generateSignerFromKey(priv)
	if err != nil {
		return nil, err
	}
	cert := generateCert(pub)
	return cert, cert.SignCert(rand.Reader, signer)
}
