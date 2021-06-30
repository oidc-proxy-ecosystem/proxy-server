package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

type ClientCertificate interface {
	GetCertificate() []byte
	GetRsaPrivateKey() []byte
	GetSSHPublicKey() []byte
	SaveFile() error
}

type clientCertificate struct {
	rsaPrivPem   []byte
	certPem      []byte
	sshPublicKey []byte
}

func (c clientCertificate) GetCertificate() []byte {
	return c.certPem
}

func (c clientCertificate) GetRsaPrivateKey() []byte {
	return c.rsaPrivPem
}

func (c clientCertificate) GetSSHPublicKey() []byte {
	return c.sshPublicKey
}

func (c clientCertificate) SaveFile() error {
	keyFileName := "client.key"
	if err := os.WriteFile(keyFileName, c.rsaPrivPem, 0600); err != nil {
		return err
	}
	log.Println("Written " + keyFileName)
	pubFile := "client.key.pub"
	if err := os.WriteFile(pubFile, c.sshPublicKey, 0600); err != nil {
		return err
	}
	log.Println("Written " + pubFile)
	pubCertFile := "client.key-cert.pub"
	if err := os.WriteFile(pubCertFile, c.certPem, 0600); err != nil {
		return err
	}
	log.Println("Written " + pubCertFile)
	return nil
}

func newClientCertificate(filename string) (ClientCertificate, error) {
	cc := clientCertificate{}
	caName := config.NewCaName(filename)
	prvClient, cert, err := caName.GenerateClient()
	if err != nil {
		return nil, err
	}
	privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(prvClient)})
	cc.rsaPrivPem = make([]byte, len(privPem))
	copy(cc.rsaPrivPem, privPem)

	sshPub, _ := ssh.NewPublicKey(&prvClient.PublicKey)
	pubKeyBytes := ssh.MarshalAuthorizedKey(sshPub)
	cc.sshPublicKey = make([]byte, len(pubKeyBytes))
	copy(cc.sshPublicKey, pubKeyBytes)
	certBytes := ssh.MarshalAuthorizedKey(cert)
	cc.certPem = make([]byte, len(certBytes))
	copy(cc.certPem, certBytes)
	return cc, nil
}

var caGenerate = &cli.Command{
	Name:        "ca",
	Description: "Generation certificate authority",
	Subcommands: []*cli.Command{
		{
			Name: "client",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	},
	Action: func(c *cli.Context) error {
		caName := config.NewCaName(config.File.Certificate)
		return caName.GenerateCA()
	},
}
