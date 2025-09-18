package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		switch {
		case statusCode >= 500:
			log.Error().
				Str("method", method).
				Str("path", path).
				Str("ip", clientIP).
				Int("status", statusCode).
				Dur("latency", latency).
				Msg("Server error")
		case statusCode >= 400:
			log.Warn().
				Str("method", method).
				Str("path", path).
				Str("ip", clientIP).
				Int("status", statusCode).
				Dur("latency", latency).
				Msg("Client error")
		default:
			log.Info().
				Str("method", method).
				Str("path", path).
				Str("ip", clientIP).
				Int("status", statusCode).
				Dur("latency", latency).
				Msg("Request processed")
		}
	}
}