package events

import (
	"github.com/gin-gonic/gin"
	"io"
)

type Controller interface {
	Stream(c *gin.Context)
}

type controller struct {
	h Hub
}

func NewController(h Hub) Controller {
	return &controller{h}
}

func (ctrl *controller) Stream(c *gin.Context) {
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-store")

	c.Writer.Flush()

	s := make(Subscriber)

	ctrl.h.Register(s)

	defer func() {
		ctrl.h.Unregister(s)
	}()

	c.Stream(func(w io.Writer) bool {
		if e, ok := <-s; ok {
			c.SSEvent("", e)

			return true
		}

		return false
	})
}
