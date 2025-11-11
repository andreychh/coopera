package engine

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

type shutdownEngine struct {
	origin Engine
}

func (s shutdownEngine) Start(ctx context.Context) {
	child, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.origin.Start(child)
	}()
	<-child.Done()
	// todo: use logger instead of log package
	log.Println("Received termination signal. Context cancelled.")
	wg.Wait()
	log.Println("Engine shut down gracefully.")
}

func ShutdownEngine(origin Engine) Engine {
	return shutdownEngine{origin: origin}
}
