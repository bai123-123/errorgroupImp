package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/hello", helloServer)
	fmt.Println("http server start")
	err := srv.ListenAndServe()
	return err

}

func helloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello,world")
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	group, errorCtx := errgroup.WithContext(ctx)

	srv := &http.Server{Addr: ":9090"}

	group.Go(func() error {
		return StartHttpServer(srv)
	})

	group.Go(func() error {
		<-errorCtx.Done()
		fmt.Println("http server stop")
		return srv.Shutdown(errorCtx)
	})

	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	group.Go(func() error {
		for {
			select {
			case <-errorCtx.Done():
				return errorCtx.Err()
			case <-chanel:
				cancel()

			}
		}
		return nil
	})

	if err := group.Wait(); err!= nil{
		fmt.Println("group error: ", err)
	}
	fmt.Println("all group done!")

}
