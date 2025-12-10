package middleware

import (
    "bytes"
    "net/http"
    "runtime"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go-api-starter/internal/pkg/response"
    "go-api-starter/internal/pkg/logger"
)

func RecoveryWithLogger(lg *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if rec := recover(); rec != nil {
                traceID := uuid.New().String()
                stack := stacktrace(3)
                lg.Logger.WithField("panic", rec).WithField("trace_id", traceID).WithField("stack", stack).Error("panic recovered")
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
