package main

import (
	"os"
	"os/signal"
	"log"
	"time"
	"context"
	"syscall"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func StartServer(addr string, ctx context.Context) error {
	s := &http.Server{Addr:addr, Handler:nil,}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		s.Shutdown(ctx)
		log.Printf("%s Shutdown!", addr)
	}()

	log.Printf("Server \"%s\" start!", addr)
	return s.ListenAndServe()
}

func ListenSignal(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case <-ctx.Done():
			log.Printf("No signal coming!")
			return nil;
		case s := <-c:
			log.Printf("A signal \"%s\" coming!", s.String())
			switch s {
			case syscall.SIGINT, syscall.SIGQUIT:
				return errors.New("Exit signal is accepted!");
			default:
				continue
			}
		}
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return StartServer("127.0.0.1:8080", ctx)
	})

	g.Go(func() error {
		return StartServer("127.0.0.1:8081", ctx)
	})

	g.Go(func() error {
		return ListenSignal(ctx)		
	})

	if err := g.Wait(); err != nil {
		log.Printf("%+v", err);
	}
}