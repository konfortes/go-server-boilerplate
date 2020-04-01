package main

import (
	"log"
	"net/http"

	"github.com/konfortes/go-server-utils/server"
	"github.com/konfortes/go-server-utils/utils"
)

const (
	appName = "my-app-name"
)

func main() {
	serverConfig := server.Config{
		AppName:     "my-app-name",
		Port:        utils.GetEnvOr("PORT", "3000"),
		Env:         utils.GetEnvOr("ENV", "development"),
		Handlers:    handlers(),
		WithTracing: utils.GetEnvOr("TRACING_ENABLED", "false") == "true",
	}

	srv := server.Initialize(serverConfig)

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	server.GracefulShutdown(srv)
}
