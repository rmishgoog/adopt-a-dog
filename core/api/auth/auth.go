package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/rego"
)

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type Auth struct {
	keyLookup KeyLookup
	parser    *jwt.Parser
	issuer    string
}

type KeyLookup interface {
	PublicKey(kid string) (key string, err error)
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
		return Claims{}, fmt.Errorf("key id (kid) is malformed: %w", err)
	}

	pem, err := a.keyLookup.PublicKey(kid) // Fetch the public key from the key lookup, this shall be the discovery endpoint of keycloak issuing the token.
	if err != nil {
		return Claims{}, fmt.Errorf("failed to fetch public key from the provider: %w", err)
	}

	input := map[string]any{
		"Key":   pem,
		"Token": parts[1],
		"ISS":   a.issuer,
	}

	if err := a.opaPolicyEvaluation(ctx, regoAuthentication, RuleAuthenticate, input); err != nil {
		return Claims{}, fmt.Errorf("authentication failed for the supplied token: %w", err)
	}

	fmt.Println("authenticated successfully")

	return claims, nil
}

func (a *Auth) opaPolicyEvaluation(ctx context.Context, regoScript string, rule string, input any) error {
	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", regoScript),
	).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
	}

	return nil
}
