package middleware

import (
    "fmt"
    "net/http"
    "runtime/debug"

    "golang-api-starter/pkg/logger"
    "golang-api-starter/pkg/response"
    "github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic
                log.Errorw("Panic recovered",
                    "error", err,
                    "stack", string(debug.Stack()),
                    "path", c.Request.URL.Path,
                    "method", c.Request.Method,
                )

                // Return error response
                response.InternalServerError(
                    c,
                    "Internal server error",
                    fmt.Sprintf("Panic recovered: %v", err),
                )

                c.Abort()
            }
        }()

        c.Next()
    }
}
```

**Penjelasan:**
- Recover dari panic untuk prevent app crash
- Log panic dengan stack trace
- Return 500 error ke client
- App tetap running

**Panic Recovery Flow:**
```
Request → Handler
           │
           ├─ Something wrong
           ├─ panic("database connection lost")
           │
           ▼
      Recovery Middleware
           │
           ├─ Catch panic
           ├─ Log error + stack trace
           ├─ Return 500 to client
           │
           ▼
      App continues running ✅
```

**Without Recovery:**
```
panic → App crash ❌