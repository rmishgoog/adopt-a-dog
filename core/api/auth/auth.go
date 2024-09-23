package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
)

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type Config struct {
	Log       *logger.Logger
	KeyLookup KeyLookup
	Issuer    string
}

type Auth struct {
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
}

// New creates a new Auth struct & configures it with the provided Config.
func New(cfg Config) (*Auth, error) {
	a := Auth{
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodES256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
	}
	return &a, nil
}

// Issuer returns the issuer of the token configured in the Auth struct.
func (a *Auth) Issuer() string {
	return a.issuer
}

type KeyLookup interface {
	PublicKey(discoveryURL string, skipCert bool) (key string, err error)
}

func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {

	parts := strings.Split(bearerToken, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format: Bearer {token}, token is malformed")
	}

	var claims Claims
	token, _, err := a.parser.ParseUnverified(parts[1], &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	kidRaw, exists := token.Header["kid"]
	if !exists {
		return Claims{}, fmt.Errorf("kid missing from header: %w", err)
	}

	kid, ok := kidRaw.(string)
	if !ok {
		return Claims{}, fmt.Errorf("key id (kid) %s is malformed: %w", kid, err)
	}

	return claims, nil
}
