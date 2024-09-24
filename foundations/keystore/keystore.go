package keystore

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}
type JWKS struct {
	Keys []JWK `json:"keys"`
}

type KeyStore struct {
	jwks map[string]JWK
}

func (ks *KeyStore) PublicKey(discoveryURL string, skipCert bool) error {

	customTLS := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipCert}, // This is insecure, don't do this in production
	}
	client := &http.Client{Transport: customTLS}
	resp, err := client.Get(discoveryURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	var jwks JWKS
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return fmt.Errorf("failed to parse JWKS: %v", err)
	}
	keyMap := make(map[string]JWK)
	for _, key := range jwks.Keys {
		keyMap[key.Kid] = key
	}
	ks.jwks = keyMap // This is set once upon the service bootstrapping.
	return nil

}

// ValidateJWT validates the JWT token with the given kid. The token here is the JWT token string as recieved in the Authorization header.
func (ks *KeyStore) ValidateJWT(token string, kid string) error {
	// Check if the kid exists in the JWKS, kid was fetched from the token header after parsing the JWT received in the Authorization header.
	jwk, exists := ks.jwks[kid]
	if !exists {
		return fmt.Errorf("kid was not found in JWKS")
	}
	jwt, err := validateToken(token, jwk)
	if err != nil {
		return fmt.Errorf("failed to validate token: %v", err)
	}
	if !jwt.Valid {
		return fmt.Errorf("supplied token is invalid")
	}
	return nil
}

func getRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	e := int(binary.BigEndian.Uint32(append(make([]byte, 4-len(eBytes)), eBytes...)))
	n := new(big.Int).SetBytes(nBytes)

	return &rsa.PublicKey{N: n, E: e}, nil
}

func validateToken(tokenString string, jwk JWK) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getRSAPublicKey(jwk)
	})
	return token, err
}

// Construct a new KeyStore upon service bootstrapping & return it with zero values.
func New() *KeyStore {
	return &KeyStore{
		jwks: make(map[string]JWK),
	}
}
