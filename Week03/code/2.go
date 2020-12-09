package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main()  {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, GopherCon SG")
	})
	go func() {
		if err :=  http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}

// leak is a buggy function.It launches a goroutine that
// blocks receiving from a channel.Nothing will ever be
// sent on that channel and the channel is never closed so
// that goroutine will be blocked forever.
func leak() {
	ch := make(chan int)

	go func() {
		val := <- ch
		fmt.Println("we received a value:", val)
	}()
}

// search simulates a function that finds a record based
// on a search term. It takes 200ms to perform this work.
func search(term string) (string, error)  {
	time.Sleep(200 * time.Millisecond)
	return "some value", nil
}

// process is the work for the program. It finds a record
// than prints it.
func process(term string) error {
	record, err := search(term)
	if err != nil {
		return err
	}

	fmt.Println("Received:", record)
	return nil
}

// result wraps the return values from search. It allows us
// to pass both values across a single channel.
type result struct {
	record string
	err error
}

// process is the work for the program. It finds a record
// then prints it. It fails if it takes more than 100ms.
func process(term string) error  {

	// Create a context that will be canceled in 100ms.
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond)
	defer cancel()

	// Make a channel for the goroutine to report its result.
	ch := make(chan result)

	// Launch a goroutine to find the record. Create a resule
	// from the returned values to send through the channel.
	go func() {
		record, err := search(term)
		ch <- result{record, err}
	}()

	// Block waiting to either receive from the goroutine's
	// channel or for the context to be canceled.
	select {
	case <- ctx.Done():
		return errors.New("search canceled")
	case result := <-ch:
		if result.err != nil {
			return result.err
		}
		fmt.Println("Received:", result.record)
		return nil
	}
}






