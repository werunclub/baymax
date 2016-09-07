package middleware

import (
	"io"
	"net/http/httputil"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	g_errors "github.com/go-errors/errors"

	"baymax/errors"
)

func Recovery() gin.HandlerFunc {
	return RecoveryWithWriter(gin.DefaultErrorWriter)
}

func RecoveryWithWriter(out io.Writer) gin.HandlerFunc {

	log := logrus.New()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := g_errors.Wrap(err, 3)
				reset := string([]byte{27, 91, 48, 109})

				log.Printf("[Recovery] panic recovered:\n\n%s%s\n\n%s%s", httprequest, goErr.Error(), goErr.Stack(), reset)

				c.JSON(500, errors.InternalServerError(goErr.Error()))
			}
		}()
		c.Next() // execute all the handlers
	}
}
