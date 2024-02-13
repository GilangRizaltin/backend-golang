package pkg

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Server(router *gin.Engine) *http.Server {
	var addr string = "localhost:9000"
	if os.Getenv("GO_ENV") == "DOCKER" {
		addr = ":9000"
	}
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}
	return server
}
