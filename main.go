package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/samuelastech/products-api/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	productHandler := handlers.NewProducts(logger)

	serverMux := http.NewServeMux()
	serverMux.Handle("/", productHandler)
	serverMux.Handle("/{id}", productHandler)

	server := &http.Server{
		Addr:         ":9000",
		Handler:      serverMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		server.ListenAndServe()
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	receivedSignal := <-signalChannel
	logger.Println("Received a terminate, graceful shutdown", receivedSignal)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
