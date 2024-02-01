package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"emperror.dev/errors"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/constants"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	"golang.org/x/oauth2"
)

// GenerateStateOauth generated a random string for aout state
func GenerateStateOauth() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func VerifyToken(ctx context.Context, keycloack models.Keycloak, verifier *oidc.IDTokenVerifier, oauth2Config oauth2.Config, keycloak models.Keycloak, domain string) gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		rawAccessToken, errT := c.Cookie("access_token")
		if errT != nil {
			refreshedToken, err := RefreshToken(c, ctx, oauth2Config)
			if err != nil {
				c.Error(errors.Errorf("RefreshToken not possible VerifyToken Authorization header missing %d, err : %s", http.StatusUnauthorized, err))
				urlPath := c.Request.URL.Path
				session.Set("url_path", urlPath)
				session.Save()
				c.Redirect(http.StatusFound, "/auth/login?url="+urlPath)
				c.Abort()
			} else {
				if refreshedToken != nil {
					c.SetCookie("access_token", refreshedToken.AccessToken, int(time.Until(refreshedToken.Expiry).Seconds()), "/", domain, false, true)
					session.Set("refresh_token", refreshedToken.RefreshToken)
				}
			}

		}

		if rawAccessToken == "" {
			c.Error(errors.Errorf("VerifyToken Authorization header missing %d", http.StatusUnauthorized))

			urlPath := c.Request.URL.Path
			session.Set("url_path", urlPath)
			session.Save()
			c.Redirect(http.StatusFound, "/auth/login?url="+urlPath)
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
			// fmt.Println("erwrwe", err.Error(), "trq")
			fmt.Println("VerifyToken jwt.ParseWithClaims err", err)
			c.Error(errors.Errorf("VerifyToken jwt.ParseWithClaims:"+err.Error(), http.StatusUnauthorized))
			return
		}

		session.Set("username", userClaim.PreferredUsername)
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
	// fmt.Println("refreshToken : ", refreshToken, "err", errT)
	ts := oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken.(string)})
	tok, err := ts.Token()

	if err != nil {
		return nil, err
	}
	return tok, err
}
