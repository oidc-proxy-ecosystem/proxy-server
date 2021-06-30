package config

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

type SamlConfig struct {
	IdpMetadata                 string `yaml:"idp_metadata"`
	ServiceProviderIssuer       string `yaml:"service_provider_issuer"`
	AssertionConsumerServiceUrl string `yaml:"assertion_consumer_service_url"`
	AudienceUri                 string `yaml:"audience_uri"`
}

func (s SamlConfig) IsFile() bool {
	u, err := url.Parse(s.IdpMetadata)
	if err != nil {
		// parseできないはファイルとする
		return true
	}
	if !u.IsAbs() {
		return true
	}
	return false
}

func (s SamlConfig) Load() (sp *saml2.SAMLServiceProvider, err error) {
	var rawMetadata []byte
	if !s.IsFile() {
		res, reErr := http.Get(s.IdpMetadata)
		if reErr != nil {
			err = reErr
			return
		}
		defer res.Body.Close()
		buf, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			err = readErr
			return
		}
		rawMetadata = buf
	} else {
		buf, readErr := os.ReadFile(s.IdpMetadata)
		if readErr != nil {
			err = readErr
			return
		}
		rawMetadata = buf
	}
	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		return
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				err = fmt.Errorf("metadata certificate(%d) must not be empty", idx)
				return
			}
			certData, decodeErr := base64.StdEncoding.DecodeString(xcert.Data)
			if decodeErr != nil {
				err = decodeErr
				return
			}

			idpCert, certErr := x509.ParseCertificate(certData)
			if certErr != nil {
				err = certErr
				return
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp = &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       s.ServiceProviderIssuer,
		AssertionConsumerServiceURL: s.AssertionConsumerServiceUrl,
		SignAuthnRequests:           true,
		AudienceURI:                 s.AudienceUri,
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}
	err = nil
	return
}

func NewSamlConfig(filename string) SamlConfig {
	var saml SamlConfig
	utils.MustReadYaml(filename, &saml)
	return saml
}
