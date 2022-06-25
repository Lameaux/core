package views

import (
	"github.com/gin-gonic/gin"

	"github.com/Lameaux/core/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorJSON(c *gin.Context, code int, err error) {
	logger.Errorw("api error", "error", err)
	c.JSON(code, ErrorResponse{Error: err.Error()})
}
