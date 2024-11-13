package jwks

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
)

// JWKSInt defines methods for managing JSON Web Key Sets (JWKS).
// It provides domain-specific methods to retrieve and format JWKS data.
//
//go:generate moq -out mock/jwks.go -pkg mock . JWKSInt
type JWKSInt interface { //nolint // JWKSInt serves the purpose of maintaining explicit, domain-specific naming that aligns with the rest of the system.
	// JWKSGetKeyset retrieves the JWKS for the specified AWS region and user pool ID.
	JWKSGetKeyset(awsRegion, poolID string) (*JWKS, error)

	// JWKSToRSAJSONResponse formats the JWKS as an RSA JSON response.
	JWKSToRSAJSONResponse(jwks *JWKS) ([]byte, error)
}

const (
	RSAAlgorithm       = "RSA"
	RSAExponentAQAB    = "AQAB"
	RSAExponentAAEAAQ  = "AAEAAQ"
	RSADefaultExponent = 65537
)

var (
	jwksURL = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
)

// JSONKey represents a single JSON Web Key (JWK) with RSA key parameters.
type JSONKey struct {
	E   string `json:"e"`   // Exponent of the RSA public key.
	Kid string `json:"kid"` // Key ID used to match a specific key.
	Kty string `json:"kty"` // Key type (e.g., "RSA").
	N   string `json:"n"`   // Modulus of the RSA public key.
}

// JWKS represents a JSON Web Key Set (JWKS) containing multiple JSONKey entries.
type JWKS struct {
	Keys []JSONKey `json:"keys"` // List of JSON Web Keys.
}

// JWKSGetKeyset primary package method which retrives the json web key set for cognito user pool
func (j JWKS) JWKSGetKeyset(awsRegion, poolID string) (*JWKS, error) {
	resp, err := http.Get(fmt.Sprintf(jwksURL, awsRegion, poolID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
	response, err := j.JWKSToRSA(jwks)
	if err != nil {
		return nil, err
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// JWKSToRSA method returns a map of the JWKS RSA Public keys
func (j JWKS) JWKSToRSA(jwks *JWKS) (map[string]string, error) {
	var (
		response = make(map[string]string)
		err      error
	)
	for _, jwk := range jwks.Keys {
		response[jwk.Kid], err = j.JWKToRSAPublicKey(jwk)
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}

// JWKToRSAPublicKey transforms key data to PKIX, ASN.1 DER form
func (j JWKS) JWKToRSAPublicKey(jwk JSONKey) (string, error) {
	if jwk.Kty != RSAAlgorithm {
		return "", errors.New(models.JWKSUnsupportedKeyTypeDescription)
	}

	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
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
			return "", err
		}

		return base64.StdEncoding.EncodeToString(der), nil
	}
	return "", errors.New(models.JWKSExponentErrorDescription)
}

// DoGetJWKS return package interface
func (j JWKS) DoGetJWKS(_ context.Context) JWKSInt {
	return &j
}

// GetJWKSRSAKeys retrieves the JWKS RSA keys which are consumed by the authorisation middleware on startup.
func (j JWKS) GetJWKSRSAKeys(awsRegion, poolID string) (map[string]string, error) {
	jwksURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	rcClient := &dphttp.Client{
		MaxRetries: 5,
		RetryTime:  1 * time.Second,
	}

	resp, err := rcClient.Get(context.Background(), fmt.Sprintf(jwksURL, awsRegion, poolID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks JWKS
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, err
	}
	keyData := &jwks

	response, err := jwks.JWKSToRSA(keyData)
	if err != nil {
		return nil, err
	}

	return response, nil
}
