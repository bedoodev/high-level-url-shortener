package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

func NewRateLimiterMiddleware(limit string) func(http.Handler) http.Handler {

	// Create memory store (in-memory, process-level rate limiting)
	store := memory.NewStore()

	rate, err := limiter.NewRateFromFormatted(limit)

	if err != nil {
		panic("invalid rate format: " + err.Error())
	}

	instance := limiter.New(store, rate)

	// Middleware Logic
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the client IP
			ip := extractClientIP(r)

			// Get the context for this IP
			limiterCtx, err := instance.Get(r.Context(), ip)
			if err != nil {
				zap.L().Error("rate limiter failed", zap.Error(err))
				http.Error(w, "internal rate limiter error", http.StatusInternalServerError)
				return
			}

			// If request limit exceeded
			if limiterCtx.Reached {
				resetTime := time.Unix(limiterCtx.Reset, 0) // 0 = nanoseconds
				retryAfter := time.Until(resetTime)
				seconds := strconv.FormatInt(int64(retryAfter.Seconds()), 10) + " seconds"
				w.Header().Set("Retry-After", seconds)
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Let the request continue
			next.ServeHTTP(w, r)
		})
	}
}

func extractClientIP(r *http.Request) string {
	xForwarded := r.Header.Get("X-Forwarded-For")

	if xForwarded != "" {
		parts := strings.Split(xForwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
