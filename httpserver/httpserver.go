package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/Lameaux/core/config"
	"github.com/Lameaux/core/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const ShutdownTimeout = 5 * time.Second

func Start(c *config.AppConfig, handler *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    ":" + c.Port,
		Handler: handler,
	}

	url := fmt.Sprintf("http://0.0.0.0:%s", c.Port)
	logger.Infow("starting server", "port", c.Port, "url", url)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	return srv
}

func Shutdown(srv *http.Server, timeout time.Duration) {
	logger.Infow("shutting down API server")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("api server forced to shutdown: ", err)
	}

	logger.Infow("api server exiting")
}
