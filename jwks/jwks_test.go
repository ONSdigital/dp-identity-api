package jwks_test

import (
	"testing"

	"github.com/ONSdigital/dp-identity-api/jwks"
	"github.com/ONSdigital/dp-identity-api/jwks/mock"

	. "github.com/smartystreets/goconvey/convey"
)

var validJWKS = &jwks.JWKS{
	Keys: []jwks.JSONKey{
		mock.KeySetOne,
		mock.KeySetTwo,
	},
}

var j = &jwks.JWKS{}

func TestJWKSToRSAJSONResponse(t *testing.T) {
	Convey("Enter a valid JWKS - check expected response", t, func() {
		response, err := j.JWKSToRSAJSONResponse(validJWKS)
		So(response, ShouldNotResemble, nil)
		So(err, ShouldEqual, nil)
	})
}

func TestJWKSToRSA(t *testing.T) {
	Convey("Enter a valid JWKS - check expected response", t, func() {
		response, err := j.JWKSToRSA(validJWKS)
		So(response, ShouldResemble, mock.JWKSData)
		So(err, ShouldEqual, nil)
	})
}

func TestConvertJwkToRsa(t *testing.T) {
	Convey("Enter a valid JWK - check expected response", t, func() {
		response, err := j.JWKToRSAPublicKey(validJWKS.Keys[0])
		So(response, ShouldNotEqual, "")
		So(err, ShouldEqual, nil)
	})

	Convey("Enter an unsupported key type - check expected error is returned", t, func() {
		validJWKS.Keys[0].Kty = "ABC"
		response, err := j.JWKToRSAPublicKey(validJWKS.Keys[0])
		So(response, ShouldEqual, "")
		So(err.Error(), ShouldEqual, "unsupported key type. Must be rsa key")
		validJWKS.Keys[0].Kty = "RSA"
	})

	Convey("Enter an invalid N value - check expected error is returned", t, func() {
		validJWKS.Keys[0].N = "!"
		response, err := j.JWKToRSAPublicKey(validJWKS.Keys[0])
		So(response, ShouldEqual, "")
		So(err.Error(), ShouldEqual, "error decoding json web key")
		validJWKS.Keys[0].N = "vBvi--N-F9MQO81xh71jIbkx81w4_sGhbztTJgIdhycV-lMzG6y3dMBWo9eRsFJuRs3MUFElmRrTVxc7EPWNQGQjUyPFW0_CnPPoGBCwgCyWtpNs5EHAkCHXsfryHb6LbJxH9LEbwOQCHR25_Bnqo_NeXSBJtvUabq3cTUgdOPc61Hskq-m19M1u7u1xu7b5DHD308Qyz3OhaEHx3cLL2za-mKxHe0VDe3sa5UfdaliTdBypFWJgNl6TsxF_G83fksgb3bVchzW45pu4dEhtNLqgXejH2-GwU8YRaAguKGW7dO_v-5uwLgDYQG9wgtAwLIMiXsFU7muig2pJEtlG2w"
	})

	Convey("Enter an unsupported exponent value - check expected error is returned", t, func() {
		validJWKS.Keys[0].E = "ABC"
		response, err := j.JWKToRSAPublicKey(validJWKS.Keys[0])
		So(response, ShouldEqual, "")
		So(err.Error(), ShouldEqual, "unexpected exponent: unable to decode JWK")
		validJWKS.Keys[0].E = "AQAB"
	})
}
