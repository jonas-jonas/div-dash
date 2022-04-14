package error

import (
	"div-dash/internal/httputil"
	"div-dash/internal/logging"

	"github.com/gin-gonic/gin"
)

type (
	ErrorHandlerProvider interface {
		ErrorHandler() *Handler
	}

	handlerDependencies interface {
		logging.LoggerProvider
	}

	Handler struct {
		handlerDependencies
	}
)

func NewHandler(h handlerDependencies) *Handler {
	return &Handler{
		handlerDependencies: h,
	}
}

//
// Middleware Error Handler in server package
//
func (h *Handler) HandleErrors(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	// in this example we only will use the **last error typed as public**
	// but you could iterate over all them since c.Errors is a slice!
	errorToPrint := c.Errors.Last()
	if errorToPrint != nil {
		h.Logger().Warnf("Caught error on %s: %s", c.Request.RequestURI, errorToPrint.Error())
		httputil.AbortServerError(c)
	}
}
