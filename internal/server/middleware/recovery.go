package middleware

import (
    "bytes"
    "net/http"
    "runtime"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
    "github.com/yourname/go-backend-enterprise/internal/pkg/response"
)

// RecoveryWithLogger wraps a panic recovery and logs stacktrace with traceID.
func RecoveryWithLogger(lg *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if rec := recover(); rec != nil {
                traceID := uuid.New().String()
                stack := stacktrace(3)
                lg.Error("panic recovered",
                    logger.F("panic", rec),
                    logger.F("trace_id", traceID),
                    logger.F("stack", stack),
                    logger.F("path", c.Request.URL.Path),
                )
                response.Error(c, http.StatusInternalServerError, "internal server error (trace: "+traceID+")")
                c.Abort()
            }
        }()
        c.Next()
    }
}

func stacktrace(skip int) string {
    var buf bytes.Buffer
    for i := skip; ; i++ {
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        fn := runtime.FuncForPC(pc)
        buf.WriteString(fn.Name())
        buf.WriteString(" - ")
        buf.WriteString(file)
        buf.WriteString(":")
        buf.WriteString(strconv.Itoa(line))
        buf.WriteString("\n")
    }
    return buf.String()
}
