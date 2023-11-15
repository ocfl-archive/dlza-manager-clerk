package models

import "github.com/golang-jwt/jwt/v4"

type Keycloak struct {
	Addr         string
	Realm        string
	Callback     string
	ClientId     string
	ClientSecret string
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
	Aud               string                 `json:"aud,omitempty"`
	TenantList        []string               `json:"tenant_list,omitempty"`
}

type ServiceRole struct {
	Roles []string `json:"roles"`
}
