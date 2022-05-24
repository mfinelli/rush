package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tg123/go-htpasswd"
)

func postAuth(c *gin.Context) {
	if viper.Get("server.auth").(string) == "htpasswd" {
		auth, err := htpasswd.New(viper.Get("server.htpasswd").(string), htpasswd.DefaultSystems, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "%v\n", err)
		}

		if auth.Match(c.PostForm("username"), c.PostForm("password")) {

			c.String(http.StatusOK, "success\n")
		} else {
			// TODO: redirect bad password
			c.String(http.StatusOK, "failue\n")
		}
	} else {
		c.String(http.StatusInternalServerError, "Unsupported auth method %s\n", viper.Get("server.auth"))
	}
}
