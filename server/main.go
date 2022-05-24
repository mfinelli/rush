package server

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/mfinelli/rush/db"
	"github.com/mfinelli/rush/version"
)

var DistFiles embed.FS
var Templates embed.FS

func Serve(rdb *gorm.DB) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "rush server version %s\n", version.VERSION.Version)
	})

	router.GET("/login", func(c *gin.Context) {
		t, err := Templates.ReadFile("src/login.html.tmpl")
		if err != nil {
			c.String(http.StatusInternalServerError, "%v\n", err)
		}

		tmpl, err := template.New("login").Parse(string(t))
		if err != nil {
			c.String(http.StatusInternalServerError, "%v\n", err)
		}

		data := struct {
			Lang              string
			Title             string
			ResetPasswordLink string
		}{
			Lang:              "en",
			Title:             "Login",
			ResetPasswordLink: "https://example.com",
		}

		var rendered bytes.Buffer
		if err := tmpl.Execute(&rendered, data); err != nil {
			c.String(http.StatusInternalServerError, "%v\n", err)
		}

		c.Render(http.StatusOK, render.Data{
			ContentType: "text/html; charset=utf-8",
			Data:        rendered.Bytes(),
		})
	})

	router.POST("/auth", postAuth)

	router.GET("/static/login.css", func(c *gin.Context) {
		file, err := DistFiles.ReadFile("dist/login.css")
		if err != nil {
			c.String(http.StatusInternalServerError, "%v\n", err)
		}

		// c.Header("Key", "Value")

		c.Render(http.StatusOK, render.Data{
			ContentType: "text/css",
			Data:        file,
		})
	})

	router.GET("/ca", func(c *gin.Context) {
		t := c.DefaultQuery("t", "ed25519")

		if t == "rsa" {
			privateKey, publicKey, err := generateRSAKey()
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			// TODO: error handling
			rdb.Create(&db.CACertificate{
				Type:       t,
				PublicKey:  string(ssh.MarshalAuthorizedKey(publicKey)),
				PrivateKey: string(convertRSAPrivateKeyToPem(privateKey)),
			})

			c.JSON(http.StatusOK, CACertificateResponse{
				PublicKey: string(ssh.MarshalAuthorizedKey(publicKey)),
			})
		} else {
			privateKey, publicKey, err := generateEd25519Key()
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			rdb.Create(&db.CACertificate{
				Type:       t,
				PublicKey:  string(ssh.MarshalAuthorizedKey(publicKey)),
				PrivateKey: string(convertEd25519PrivateKeyToPem(privateKey)),
			})

			c.JSON(http.StatusOK, CACertificateResponse{
				PublicKey: string(ssh.MarshalAuthorizedKey(publicKey)),
			})
		}

	})

	router.GET("/host", func(c *gin.Context) {
		t := c.DefaultQuery("t", "ed25519")
		cn := c.Query("h")

		if t == "rsa" {
			r, err := generateRSAHostKey(rdb, cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		} else {
			r, err := generateEd25519HostKey(rdb, cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		}
	})

	router.GET("/user", func(c *gin.Context) {
		t := c.DefaultQuery("t", "ed25519")
		cn := c.Query("u")

		if t == "rsa" {
			r, err := generateRSAUserKey(rdb, cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		} else {
			r, err := generateEd25519UserKey(rdb, cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		}
	})

	log.Printf("rush server version %s running on port %d", version.VERSION.Version, viper.Get("server.port"))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.Get("server.port")),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
