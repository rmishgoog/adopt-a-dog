package auth

import (
	_ "embed"
)

// These the current set of rules we have for auth.
const (
	RuleAuthenticate   = "auth"
	RuleAny            = "rule_any"
	RuleAdminOnly      = "rule_admin_only"
	RuleUserOnly       = "rule_user_only"
	RuleAdminOrSubject = "rule_admin_or_subject"
)

// Package name for rego code.
const (
	opaPackage string = "adoptadog.rego"
)

// Core OPA policies through embeddings.
var (
	//go:embed rego/authentication.rego
	regoAuthentication string
)
