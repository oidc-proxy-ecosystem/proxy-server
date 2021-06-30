package router

import (
	"github.com/n-creativesystem/go-fwncs"
	saml2 "github.com/russellhaering/gosaml2"
)

const spKey = "service_prorider_key"

func setServiceProvider(sp *saml2.SAMLServiceProvider) fwncs.HandlerFunc {
	return func(c fwncs.Context) {
		c.Set(spKey, sp)
	}
}
