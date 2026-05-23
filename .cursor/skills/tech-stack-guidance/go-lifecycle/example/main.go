package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/harness/go-lifecycle/lifecycle"
)

func main() {
	mgr, err := lifecycle.New(lifecycle.Config{})
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := mgr.EnsureExclusive(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	fmt.Printf("registry: %s\n", mgr.RegistryPath())
	fmt.Printf("instance key: %s\n", mgr.InstanceKey())

	if err := mgr.ListenAndServe(ctx, mux); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
