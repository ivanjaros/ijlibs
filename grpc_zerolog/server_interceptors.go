package grpc_zerolog

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/ivanjaros/ijlibs/grpc_zerolog/ctx_zerolog"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"path"
	"time"
)

func UnaryServerInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		newCtx := initCtx(ctx, logger, info.FullMethod, start)
		resp, err := handler(newCtx, req)
		code := status.Code(err)
		level := ctl(code)
		with := ctx_zerolog.Get(newCtx).Str("grpc.code", code.String()).Dur("grpc.time_ms", time.Since(start))
		if err != nil {
			with = with.Err(err)
		}
		doLog(with.Logger(), level, "finished unary call")
		return resp, err
	}
}

func StreamServerInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		newCtx := initCtx(stream.Context(), logger, info.FullMethod, start)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		err := handler(srv, wrapped)
		code := status.Code(err)
		level := ctl(code)
		with := ctx_zerolog.Get(newCtx).Str("grpc.code", code.String()).Dur("grpc.time_ms", time.Since(start))
		if err != nil {
			with = with.Err(err)
		}
		doLog(with.Logger(), level, "finished streaming call")
		return err
	}
}

func doLog(logger zerolog.Logger, level zerolog.Level, msg string) {
	switch level {
	case zerolog.DebugLevel:
		logger.Debug().Msg(msg)
	case zerolog.InfoLevel:
		logger.Info().Msg(msg)
	case zerolog.WarnLevel:
		logger.Warn().Msg(msg)
	case zerolog.ErrorLevel:
		logger.Error().Msg(msg)
	case zerolog.FatalLevel:
		logger.Fatal().Msg(msg)
	case zerolog.PanicLevel:
		logger.Panic().Msg(msg)
	}
}

func initCtx(ctx context.Context, logger zerolog.Logger, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	with := logger.With().
		Str("grpc.service", service).
		Str("grpc.method", method).
		Str("grpc.start_time", start.Format(time.RFC3339))
	if d, ok := ctx.Deadline(); ok {
		with = with.Str("grpc.request.deadline", d.Format(time.RFC3339))
	}
	return ctx_zerolog.New(ctx, with.Logger())
}

func ctl(code codes.Code) zerolog.Level {
	switch code {
	case codes.OK:
		return zerolog.DebugLevel
	case codes.Canceled:
		return zerolog.DebugLevel
	case codes.Unknown:
		return zerolog.InfoLevel
	case codes.InvalidArgument:
		return zerolog.DebugLevel
	case codes.DeadlineExceeded:
		return zerolog.InfoLevel
	case codes.NotFound:
		return zerolog.DebugLevel
	case codes.AlreadyExists:
		return zerolog.DebugLevel
	case codes.PermissionDenied:
		return zerolog.InfoLevel
	case codes.Unauthenticated:
		return zerolog.InfoLevel
	case codes.ResourceExhausted:
		return zerolog.DebugLevel
	case codes.FailedPrecondition:
		return zerolog.DebugLevel
	case codes.Aborted:
		return zerolog.DebugLevel
	case codes.OutOfRange:
		return zerolog.DebugLevel
	case codes.Unimplemented:
		return zerolog.WarnLevel
	case codes.Internal:
		return zerolog.WarnLevel
	case codes.Unavailable:
		return zerolog.WarnLevel
	case codes.DataLoss:
		return zerolog.WarnLevel
	default:
		return zerolog.InfoLevel
	}
}
