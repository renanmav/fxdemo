package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
)

func main() {
	fx.New(
		//fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		//	return &fxevent.ZapLogger{Logger: log}
		//}),
		fx.Provide(
			NewHTTPServer,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			AsRoute(NewEchoHandler),
			AsRoute(NewHelloHandler),
			zap.NewExample,
		),
		fx.Invoke(func(s *http.Server) {}),
	).Run()
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	s := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server", zap.String("addr", s.Addr))
			go s.Serve(l)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
	return s
}

type Route interface {
	http.Handler

	Pattern() string
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
