package jwks

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"
)

//go:generate moq -out mock/jwks.go -pkg mock . JWKSFuncs
type JWKSInt interface {
	JWKSGetKeyset(awsRegion, poolId string) (*JWKS, error)
	JWKSToRSAJSONResponse(jwks *JWKS) ([]byte, error)
}

const (
	RSAAlgorithm = "RSA"
	RSAExponentAQAB = "AQAB"
	RSAExponentAAEAAQ = "AAEAAQ"
	RSADefaultExponent = 65537
)

var (
	jwksURL = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
)

type JsonKey struct {
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
}

type JWKS struct {
	Keys []JsonKey `json:"keys"`
}

// JWKSGetKeyset primary package method which retrives the json web key set for cognito user pool
func (j JWKS) JWKSGetKeyset(awsRegion, poolId string) (*JWKS, error) {
	resp, err := http.Get(fmt.Sprintf(jwksURL, awsRegion, poolId))
	if err != nil {
		return nil, err
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks JWKS
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, err
	}
	return &jwks, nil
}

// JWKSToRSAJSONResponse method returns byte[] array for request response
func (j JWKS) JWKSToRSAJSONResponse(jwks *JWKS) ([]byte, error) {
	var (
		response = make(map[string]string)
		err error	
	)
	for _, jwk := range jwks.Keys {
		response[jwk.Kid], err = j.JWKToRSAPublicKey(jwk)
		if err != nil {
			return nil, err
		}
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// JWKToRSAPublicKey transforms key data to PKIX, ASN.1 DER form
func (j JWKS) JWKToRSAPublicKey(jwk JsonKey) (string, error) {
	if jwk.Kty != RSAAlgorithm {
		return "", errors.New(models.JWKSUnsupportedKeyTypeDescription)
	}

	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		log.Println(err)
		return "", errors.New(models.JWKSErrorDecodingDescription)
	}

	// Use default exponent -> 65537
	if jwk.E == RSAExponentAQAB || jwk.E == RSAExponentAAEAAQ {
		der, err := x509.MarshalPKIXPublicKey(
			&rsa.PublicKey{
				N: new(big.Int).SetBytes(nb),
				E: RSADefaultExponent,
			},
		)
		if err != nil {
			return "", errors.New(err.Error())
		}
	
		return base64.StdEncoding.EncodeToString(der), nil
	} else {
		return "", errors.New(models.JWKSExponentErrorDescription)
	}
}

// DoGetJWKS return package interface
func (j JWKS) DoGetJWKS(ctx context.Context) JWKSInt {
	return &j
}
