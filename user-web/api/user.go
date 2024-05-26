package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
