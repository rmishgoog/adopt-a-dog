package main

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

// This function converts a JWK to an RSA public key.
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

// This function fetches the JWKS from the given URL.
func fetchJWKS(url string) (*JWKS, error) {
	customTLS := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTLS}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	var jwks JWKS
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %v", err)
	}
	return &jwks, nil
}

// This function validates a JWT token using the JWKS.
func validateToken(tokenString string, jwks *JWKS) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			for _, key := range jwks.Keys {
				if key.Kid == kid {
					return getRSAPublicKey(key)
				}
			}
		}
		return nil, fmt.Errorf("unable to find appropriate key")
	})
	return token, err
}

func main() {
	url := "https://local.auth.adoptadog.com/realms/adoptadog/protocol/openid-connect/certs"
	jwks, err := fetchJWKS(url)
	if err != nil {
		fmt.Printf("error fetching JWKS: %v\n", err)
		return
	}

	// Example JWT token (replace with your actual token, this is just an example, don't use it in production & do not leave your token in the code).
	tokenString := "insert-your-token-here"

	token, err := validateToken(tokenString, jwks)
	if err != nil {
		fmt.Printf("token validation failed: %v\n", err)
		return
	}

	if token.Valid {
		fmt.Println("token is valid")
	} else {
		fmt.Println("token is invalid")
	}
}