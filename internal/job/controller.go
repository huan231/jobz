package job

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller interface {
	List(c *gin.Context)
}

type controller struct {
	s Service
}

func NewController(s Service) Controller {
	return &controller{s}
}

func (ctrl *controller) List(c *gin.Context) {
	jobs, err := ctrl.s.List(c.Request.Context())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}
