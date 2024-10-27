package rest

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/oklog/ulid"

	"github.com/mfsyahrz/image_feed_api/internal/common/logger"
	"github.com/mfsyahrz/image_feed_api/internal/common/util"
	"github.com/mfsyahrz/image_feed_api/internal/interface/ioc"
)

const (
	ctxIdempotency = "idempotency_key"
	ctxRequestID   = "request_id"

	headerIdempotency = "Idempotency-Key"
	headerRequestID   = "Request-ID"
)

func generateThreadId() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	uniqueID := ulid.MustNew(ulid.Timestamp(t), entropy)
	return uniqueID.String()
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func SetupMiddleware(server *echo.Echo, container *ioc.IOC) {
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "token", "Pv", echo.HeaderContentType, "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	server.Use(RequestIDMiddleware)

	server.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.FromCtx(c.Request().Context()).Info("incoming request...")
			return h(c)
		}
	})

	server.Validator = &CustomValidator{validator: validator.New()}
}

// RequestIDMiddleware ensures assignment for a request ID in the incoming request
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Request().Header.Get(headerRequestID)
		requestID = util.Ternary(
			requestID != "",
			requestID,
			generateThreadId())

		c.Set(ctxRequestID, requestID)
		c.Response().Header().Set(headerRequestID, requestID)
		createCtxLogger(c, requestID)

		return next(c)
	}
}

func createCtxLogger(c echo.Context, requestID string) {
	log := logger.NewLogger(requestID)
	ctx := context.WithValue(c.Request().Context(), logger.CtxLoggerKey, log)
	c.SetRequest(c.Request().WithContext(ctx))
}
