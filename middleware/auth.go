package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ocfl-archive/dlza-manager-clerk/constants"
	"github.com/ocfl-archive/dlza-manager-clerk/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/oauth2"
)

// GenerateStateOauth generated a random string for aout state
func GenerateStateOauth() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func VerifyToken(ctx context.Context, keycloak models.Keycloak, verifier *oidc.IDTokenVerifier, oauth2Config oauth2.Config, domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("access_token") == nil {
			refreshedToken, err := RefreshToken(c, ctx, oauth2Config)
			if err != nil {
				c.Error(errors.Errorf("VerifyToken : RefreshToken not possible %d, err : %s", http.StatusUnauthorized, err))
				urlPath := c.Request.URL.Path
				session.Set("url_path", urlPath)
				session.Save()
				// c.Redirect(http.StatusFound, "/auth/login")
				// c.Redirect(http.StatusFound, "/auth/login?url="+urlPath)
				c.Abort()
				return
			} else {
				if refreshedToken != nil {
					// c.SetCookie("access_token", refreshedToken.AccessToken, int(time.Until(refreshedToken.Expiry).Seconds()), "/", domain, false, true)
					session.Set("access_token", refreshedToken.AccessToken)
					session.Set("refresh_token", refreshedToken.RefreshToken)
				}
			}
		}
		// rawAccessToken, errT := c.Cookie("access_token")
		// if errT != nil {
		// 	refreshedToken, err := RefreshToken(c, ctx, oauth2Config)
		// 	if err != nil {
		// 		c.Error(errors.Errorf("RefreshToken not possible VerifyToken Authorization header missing %d, err : %s", http.StatusUnauthorized, err))
		// 		urlPath := c.Request.URL.Path
		// 		session.Set("url_path", urlPath)
		// 		session.Save()
		// 		c.Redirect(http.StatusFound, "/auth/login")
		// 		// c.Redirect(http.StatusFound, "/auth/login?url="+urlPath)
		// 		c.Abort()
		// 		return
		// 	} else {
		// 		if refreshedToken != nil {
		// 			c.SetCookie("access_token", refreshedToken.AccessToken, int(time.Until(refreshedToken.Expiry).Seconds()), "/", domain, false, true)
		// 			session.Set("refresh_token", refreshedToken.RefreshToken)
		// 		}
		// 	}

		// }
		if session.Get("expiry_token") == nil {
			c.Error(errors.Errorf("expiry_token not set in session %d", http.StatusUnauthorized))
			c.Abort()
			return
		} else {

			expiryToken, err := time.Parse(time.RFC3339, session.Get("expiry_token").(string))
			if time.Until(expiryToken) < 0 {
				refreshedToken, err := RefreshToken(c, ctx, oauth2Config)
				if err != nil {
					c.Error(errors.Errorf("VerifyToken : RefreshToken not possible %d, err : %s", http.StatusUnauthorized, err))
					urlPath := c.Request.URL.Path
					session.Set("url_path", urlPath)
					session.Save()
					// c.Redirect(http.StatusFound, "/auth/login")
					// c.Redirect(http.StatusFound, "/auth/login?url="+urlPath)
					c.Abort()
					return
				} else {
					if refreshedToken != nil {
						// c.SetCookie("access_token", refreshedToken.AccessToken, int(time.Until(refreshedToken.Expiry).Seconds()), "/", domain, false, true)
						session.Set("access_token", refreshedToken.AccessToken)
						session.Set("refresh_token", refreshedToken.RefreshToken)
					}
				}
			}
			if err != nil {
				c.Error(errors.Errorf("expiry_token err %d, err %s", http.StatusUnauthorized, err))
				c.Abort()
				return
			}
		}
		rawAccessToken := session.Get("access_token").(string)
		if rawAccessToken == "" {
			c.Error(errors.Errorf("VerifyToken Authorization header missing %d", http.StatusUnauthorized))

			urlPath := c.Request.URL.Path
			session.Set("url_path", urlPath)
			session.Save()
			// c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		_, err := verifier.Verify(context.Background(), rawAccessToken)
		if err != nil {
			// c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
			c.Error(errors.Errorf("VerifyToken Invalid or malformed rawAccessToken:"+err.Error(), http.StatusUnauthorized))
			return
		}

		//get token info then user info ?
		var userClaim models.KeyCloakToken

		_, err = jwt.ParseWithClaims(rawAccessToken, &userClaim, nil)
		if err != nil && err.Error() != "no Keyfunc was provided." {
			c.Error(errors.Errorf("VerifyToken jwt.ParseWithClaims:"+err.Error(), http.StatusUnauthorized))
			return
		}

		session.Set("username", userClaim.PreferredUsername)
		session.Set("userClaim", userClaim)
		session.Save()
		c.Set("userClaim", userClaim)
		c.Set(constants.KEYCLOAK_GROUPS_CTX, userClaim.Groups)
		c.Set(constants.ADMIN_ROLE, keycloak.AdminRole)
		c.Next()
	}
}

func RefreshToken(c *gin.Context, ctx context.Context, oauth2Config oauth2.Config) (*oauth2.Token, error) {

	session := sessions.Default(c)
	refreshToken := session.Get("refresh_token")

	if refreshToken == nil {
		return nil, errors.Errorf("RefreshToken, no refresh_token in session")
	}

	ts := oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken.(string)})
	tok, err := ts.Token()

	if err != nil {
		return nil, err
	}
	return tok, err
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)

		c.Set("GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ctx := context.WithValue(c, constants.Needed, "Needed to attach context")
// c.Set("keycloak_group", userClaim.Groups)
// c.Set("tenant_list", userClaim.TenantList)
// h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetAuthCodeURL(c *gin.Context) (string, error) {
	session := sessions.Default(c)

	state := GenerateStateOauth()
	session.Set("state", state)
	nonce := GenerateStateOauth()
	session.Set("nonce", nonce)
	session.Save()

	var keycloak models.Keycloak
	if c.Value("keycloak") != nil {
		keycloak = c.Value("keycloak").(models.Keycloak)
	} else {
		return "", errors.New("could't retrieve keycloak informations")
	}

	oauth2Config := GetOauth2Config(keycloak)
	return oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), nil
}

func Callback(ctx context.Context, c *gin.Context, code string) error {
	session := sessions.Default(c)
	if session.Get("state") == nil {
		c.Error(errors.Errorf("state is empty : %d", http.StatusBadRequest))
		return errors.New("state is empty")
	}
	var keycloak models.Keycloak
	if ctx.Value("keycloak") != nil {
		keycloak = ctx.Value("keycloak").(models.Keycloak)
	} else {
		return errors.New("could't retrieve keycloak informations")
	}

	oauth2Config := GetOauth2Config(keycloak)
	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		fmt.Println("oauth2Config.Exchange", err)
		fmt.Println(" code ", code)
		fmt.Println("ctx ", ctx)
		fmt.Println("gincontext ", c)
		fmt.Println("oauth2Config.RedirectURL ", oauth2Config.RedirectURL)
		fmt.Println("keycloak ", keycloak)
		return err
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return errors.New("No id_token field in oauth2 token")
	}
	verifier := GetVerifier(keycloak)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		fmt.Println("verifier error", err)
		fmt.Println(" keycloak.Callback ", keycloak.Callback)
		fmt.Println("oauth2Config.RedirectURL ", oauth2Config.RedirectURL)
		return err
	}

	nonce := session.Get("nonce").(string)
	if idToken.Nonce != nonce {
		return errors.New("nonce did not match.")
	}
	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		return err
	}

	session.Set("access_token", resp.OAuth2Token.AccessToken)
	session.Set("refresh_token", resp.OAuth2Token.RefreshToken)

	var userClaim models.KeyCloakToken
	_, err = jwt.ParseWithClaims(resp.OAuth2Token.AccessToken, &userClaim, nil)
	if err != nil && err.Error() != "no Keyfunc was provided." {
		return err
	}

	session.Set("userClaim", userClaim)
	session.Set("keycloak_group", userClaim.Groups)
	session.Set("tenant_list", userClaim.TenantList)
	expiryToken := resp.OAuth2Token.Expiry.Format(time.RFC3339)
	session.Set("expiry_token", expiryToken)
	err = session.Save()
	if err != nil {
		return err
	}
	return nil
}

func GetOauth2Config(keycloak models.Keycloak) oauth2.Config {

	// Keycloak configuration
	provider := GetProvider(keycloak)

	return oauth2.Config{
		ClientID:     keycloak.ClientId,
		ClientSecret: keycloak.ClientSecret,
		RedirectURL:  keycloak.Callback + "auth/callback",
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

}

func GetProvider(keycloak models.Keycloak) *oidc.Provider {

	// Keycloak configuration
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, keycloak.Addr+keycloak.Realm)
	if err != nil {
		panic(err)
	}

	return provider

}

func GetOidcConfig(keycloak models.Keycloak) *oidc.Config {

	return &oidc.Config{
		ClientID: keycloak.ClientId,
	}
}

func GetVerifier(keycloak models.Keycloak) *oidc.IDTokenVerifier {
	provider := GetProvider(keycloak)
	return provider.Verifier(GetOidcConfig(keycloak))
}

func ResetSession(c *gin.Context) error {
	session := sessions.Default(c)
	// session.Set("access_token", nil)
	// session.Set("refresh_token", nil)
	// session.Set("expiry_token", nil)
	session.Clear()
	err := session.Save()
	return err
}

func GetUser(c *gin.Context) (*models.KeyCloakToken, error) {
	session := sessions.Default(c)
	var userClaim models.KeyCloakToken

	if session.Get("userClaim") == nil {

		return nil, errors.New("no user found")
	}
	userClaim = session.Get("userClaim").(models.KeyCloakToken)
	return &userClaim, nil
}

func GraphqlVerifyToken(ctx context.Context) error {
	c, err := GinContextFromContext(ctx)
	if err != nil {
		return err
	}
	var keycloak models.Keycloak
	if ctx.Value("keycloak") != nil {
		keycloak = ctx.Value("keycloak").(models.Keycloak)
	} else {
		return errors.New("could't retrieve keycloak informations")
	}

	session := sessions.Default(c)
	oauth2Config := GetOauth2Config(keycloak)
	if session.Get("access_token") == nil {
		// refreshedToken, err := RefreshToken(c, ctx, oauth2Config)
		// if err != nil {

		// 	return err
		// } else {
		// 	if refreshedToken != nil {
		// 		session.Set("access_token", refreshedToken.AccessToken)
		// 		session.Set("refresh_token", refreshedToken.RefreshToken)
		// 	}
		// }
		return errors.New("Access denied : no access token")
	}

	if session.Get("expiry_token") == nil {
		return errors.New("expiry_token not set in session")
	} else {

		expiryToken, err := time.Parse(time.RFC3339, session.Get("expiry_token").(string))
		if err != nil {
			return err
		}
		if time.Until(expiryToken) < 0 {
			refreshedToken, err := RefreshToken(c, ctx, oauth2Config)
			if err != nil {
				return err
			} else {
				if refreshedToken != nil {
					session.Set("access_token", refreshedToken.AccessToken)
					session.Set("refresh_token", refreshedToken.RefreshToken)
				}
			}
		}

	}
	rawAccessToken := session.Get("access_token").(string)
	if rawAccessToken == "" {
		return errors.New("access_token not set in session")
	}
	verifier := GetVerifier(keycloak)
	_, err = verifier.Verify(context.Background(), rawAccessToken)
	if err != nil {
		fmt.Println("GraphqlVerifyToken verifier error", err)
		fmt.Println(" keycloak.Callback ", keycloak.Callback)
		fmt.Println("oauth2Config.RedirectURL ", oauth2Config.RedirectURL)
		return err
	}

	//get token info then user info ?
	var userClaim models.KeyCloakToken

	_, err = jwt.ParseWithClaims(rawAccessToken, &userClaim, nil)
	if err != nil && err.Error() != "no Keyfunc was provided." {
		return err
	}

	session.Set("username", userClaim.PreferredUsername)
	session.Set("userClaim", userClaim)
	session.Set("keycloak_group", userClaim.Groups)
	session.Set("tenant_list", userClaim.TenantList)
	session.Save()
	return nil
}

func TenantGroups(ctx context.Context) ([]string, []models.Tenant, error) {
	var keyCloakGroup []string
	// var tenantList []string
	var tenantList []models.Tenant
	c, err := GinContextFromContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	session := sessions.Default(c)
	if session.Get("keycloak_group") != nil {
		keyCloakGroup = session.Get("keycloak_group").([]string)
	}
	if session.Get("tenant_list") != nil {
		// tenantList = session.Get("tenant_list").([]models.Teanntstring)
		tenantList = session.Get("tenant_list").([]models.Tenant)
	}

	return keyCloakGroup, tenantList, nil
}

func GraphqlErrorWrapper(err error, ctx context.Context, httpStatus int) *gqlerror.Error {

	if strings.Contains(err.Error(), "You are not allowed to retrieve datas") {
		httpStatus = http.StatusForbidden
	} else if strings.Contains(err.Error(), "You could not retrieve more than 1000") {
		httpStatus = http.StatusBadRequest
	}
	return &gqlerror.Error{
		Err:     err,
		Path:    graphql.GetPath(ctx),
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"code": httpStatus,
		},
	}
}
