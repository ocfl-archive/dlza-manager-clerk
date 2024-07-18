package models

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ocfl-archive/dlza-manager-clerk/constants"
	"golang.org/x/exp/slices"
)

type Keycloak struct {
	Addr         string `yaml:"addr" toml:"addr"`
	Realm        string `yaml:"realm" toml:"realm"`
	Callback     string `yaml:"callback" toml:"callback"`
	ClientId     string `yaml:"clientId" toml:"clientId"`
	ClientSecret string `yaml:"clientSecret" toml:"clientSecret"`
	AdminRole    string `yaml:"admin_role" toml:"adminRole"`
}

type KeyCloakToken struct {
	jwt.RegisteredClaims
	Jti               string                 `json:"jti,omitempty"`
	Exp               int64                  `json:"exp"`
	Nbf               int64                  `json:"nbf"`
	Iat               int64                  `json:"iat"`
	Iss               string                 `json:"iss"`
	Sub               string                 `json:"sub"`
	Typ               string                 `json:"typ"`
	Azp               string                 `json:"azp,omitempty"`
	Nonce             string                 `json:"nonce,omitempty"`
	AuthTime          int64                  `json:"auth_time,omitempty"`
	SessionState      string                 `json:"session_state,omitempty"`
	Acr               string                 `json:"acr,omitempty"`
	ClientSession     string                 `json:"client_session,omitempty"`
	AllowedOrigins    []string               `json:"allowed-origins,omitempty"`
	ResourceAccess    map[string]ServiceRole `json:"resource_access,omitempty"`
	Name              string                 `json:"name"`
	PreferredUsername string                 `json:"preferred_username"`
	GivenName         string                 `json:"given_name,omitempty"`
	FamilyName        string                 `json:"family_name,omitempty"`
	Email             string                 `json:"email,omitempty"`
	RealmAccess       ServiceRole            `json:"realm_access,omitempty"`
	CustomClaims      interface{}            `json:"custom_claims,omitempty"`
	Groups            []string               `json:"groups,omitempty"`
	AtHash            string                 `json:"at_hash,omitempty"`
	EmailVerified     bool                   `json:"email_verified,omitempty"`
	Sid               string                 `json:"sid,omitempty"`
	// Aud               string                 `json:"aud,omitempty"`
	TenantList []string `json:"tenant_list,omitempty"`
}

type ServiceRole struct {
	Roles []string `json:"roles"`
}

// GetKeycloakContext retrieves keycloack info from context
func GetKeycloakContext(ctx context.Context) map[string][]string {
	var groups, accessKey []string
	if ctx.Value(constants.KEYCLOAK_GROUPS_CTX) != nil {
		groups = ctx.Value(constants.KEYCLOAK_GROUPS_CTX).([]string)
	}
	if ctx.Value(constants.KEYCLOAK_ACCESS_KEY_CTX) != nil {
		accessKey = ctx.Value(constants.KEYCLOAK_ACCESS_KEY_CTX).([]string)
	}

	return map[string][]string{
		"keycloak_groups": groups,
		"access_key":      accessKey,
	}
}

func IsAdmin(ctx context.Context) bool {
	var groups []string
	if ctx.Value(constants.KEYCLOAK_GROUPS_CTX) != nil {
		groups = ctx.Value(constants.KEYCLOAK_GROUPS_CTX).([]string)
	}

	if ctx.Value(constants.ADMIN_ROLE) == "" {
		return false
	}
	return slices.Contains(groups, ctx.Value(constants.ADMIN_ROLE).(string))
}
