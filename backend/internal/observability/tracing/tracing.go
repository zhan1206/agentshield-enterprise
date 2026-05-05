package tracing

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// InitTracing initializes the tracing provider
func InitTracing(serviceName, endpoint string) {
	fmt.Printf("[Tracing] Initialized for %s, endpoint: %s\n", serviceName, endpoint)
}

// Shutdown gracefully shuts down the tracing provider
func Shutdown() {
	fmt.Println("[Tracing] Shutdown complete")
}

// GinMiddleware returns a Gin middleware for request tracing
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := fmt.Sprintf("trace-%d", start.UnixNano())
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
		duration := time.Since(start)
		fmt.Printf("[Trace] %s %s %d %v trace_id=%s\n",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration, traceID)
	}
}
