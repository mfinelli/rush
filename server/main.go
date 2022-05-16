package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/mfinelli/rush/db"
)

var VERSION string = "1.0.0"

func Serve(rdb *gorm.DB) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "rush server version %s\n", VERSION)
	})

	router.GET("/ca", func(c *gin.Context) {
		t := c.DefaultQuery("t", "ed25519")

		if t == "rsa" {
			privateKey, publicKey, err := generateRSAKey()
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

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
			r, err := generateRSAHostKey(cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		} else {
			r, err := generateEd25519HostKey(cn)
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
			r, err := generateRSAUserKey(cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		} else {
			r, err := generateEd25519UserKey(cn)
			if err != nil {
				c.String(http.StatusInternalServerError, "%v\n", err)
			}

			c.JSON(http.StatusOK, r)
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
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
