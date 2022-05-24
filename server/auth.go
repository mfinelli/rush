package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/tg123/go-htpasswd"
)

type authToken struct {
	jwt.RegisteredClaims
}

// TODO: either set this to a random value on startup, pull from shared
// memory cache on startup, or read from config file
var jwtSigningKey = []byte("TODO: changeme!")

func getScheme(c *gin.Context) string {
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	} else {
		if c.Request.TLS != nil {
			return "https"
		} else {
			return "http"
		}
	}
}

func buildFullHost(c *gin.Context) string {
	return fmt.Sprintf("%s://%s", getScheme(c), c.Request.Host)
}

func postAuth(c *gin.Context) {
	if viper.Get("server.auth").(string) == "htpasswd" {
		auth, err := htpasswd.New(viper.Get("server.htpasswd").(string),
			htpasswd.DefaultSystems, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "%v\n", err)
		}

		if auth.Match(c.PostForm("username"), c.PostForm("password")) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256,
				authToken{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "https://github.com/mfinelli/rush",
						Subject:   c.PostForm("username"),
						Audience:  []string{buildFullHost(c)},
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
						NotBefore: jwt.NewNumericDate(time.Now()),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ID:        "TODO",
					},
				})

			result, err := token.SignedString(jwtSigningKey)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n",
					err)
			}

			c.String(http.StatusOK, "token: %s\n", result)
		} else {
			// TODO: redirect bad password
			c.String(http.StatusOK, "failue\n")
		}
	} else {
		c.String(http.StatusInternalServerError,
			"Unsupported auth method %s\n",
			viper.Get("server.auth"))
	}
}
