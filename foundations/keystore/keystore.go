package keystore

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	jwks JWKS
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
	ks.jwks = jwks // This is set once upon the service bootstrapping.
	return nil

}

// Construct a new KeyStore upon service bootstrapping & return it with zero values.
func New() *KeyStore {
	return &KeyStore{
		jwks: JWKS{},
	}
}
