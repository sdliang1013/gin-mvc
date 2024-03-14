package caul

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SecureServer interface {
	Start()
	Stop(ctx context.Context) error
	Timeout() time.Duration
}

func StartSecureServer(server SecureServer) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go server.Start()
	// Wait for interrupt signal to gracefully stop the server with
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	// a timeout of 15 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), server.Timeout())
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exit.")
}
