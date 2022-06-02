package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller interface {
	Liveness(c *gin.Context)
}

type controller struct {
}

func NewController() Controller {
	return &controller{}
}

func (ctrl *controller) Liveness(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
