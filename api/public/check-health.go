package public

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!", "from": "api/forum"})
}
