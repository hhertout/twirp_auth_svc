package hooks

import (
	"context"
	"log"
	"time"

	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

type LoggerContextKey string

func NewLoggingServerHooks(l *zap.Logger) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			ctx = context.WithValue(ctx, LoggerContextKey("time_taken"), time.Now())

			clientIP := "unknown"
			userAgent := "unknown"
			if md, ok := twirp.HTTPRequestHeaders(ctx); ok {
				clientIP = md.Get("X-Forwarded-For")
				if clientIP == "" {
					clientIP = md.Get("Remote-Addr")
				}
				userAgent = md.Get("User-Agent")
			}
			ctx = context.WithValue(ctx, LoggerContextKey("client_ip"), clientIP)
			ctx = context.WithValue(ctx, LoggerContextKey("user_agent"), userAgent)

			return ctx, nil
		},
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		},
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			l.Sugar().Errorf("Error code=%s; message= %s", string(twerr.Code()), twerr.Msg())
			return ctx
		},
		ResponseSent: func(ctx context.Context) {
			start, ok := ctx.Value(LoggerContextKey("time_taken")).(time.Time)
			if !ok {
				log.Println("start_time not found in context")
				return
			}
			userAgent, _ := ctx.Value(LoggerContextKey("user_agent")).(string)
			clientIP, _ := ctx.Value(LoggerContextKey("client_ip")).(string)
			method, _ := twirp.MethodName(ctx)
			status, _ := twirp.StatusCode(ctx)

			logFields := []zap.Field{
				zap.String("user_agent", userAgent),
				zap.String("client_ip", clientIP),
				zap.String("status", status),
				zap.String("method", method),
				zap.Duration("time_taken", time.Since(start)),
			}

			l.Info("GRPC Request", logFields...)
		},
	}
}
