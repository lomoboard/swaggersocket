package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swaggersocket"
	"github.com/go-openapi/swaggersocket/example/autogen/client/operations"
)

var (
	conn      *swaggersocket.SocketConnection
	myHandler = &handler{}
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("received -> %s\n", r.URL.Path)
	if r.URL.Path != "/ws" {
		req, _ := http.NewRequest("GET", "http://127.0.0.1/ws"+r.URL.Path, nil)
		req.Header.Set("X-Correlation-Id", conn.ID())
		if err := conn.WriteRequest(req); err != nil {
			log.Fatalf("write request: %v\n", err)
		}
		resp, err := conn.ReadResponse()
		if err != nil {
			log.Fatalf("read response: %v\n", err)
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("read body: %v\n", err)
		}
		fmt.Printf("reply <- %s\n", string(content))
	}
}

func main() {
	wsServer := swaggersocket.NewWebSocketServer(swaggersocket.SocketServerOpts{
		Addr: ":8082", KeepAlive: true,
	})
	ch, err := wsServer.EventStream()
	if err != nil {
		panic(err)
	}
	// the following loop is safe to run in a separate go-routine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case event := <-ch:
				if event.EventType == swaggersocket.ConnectionReceived {
					conn = wsServer.ConnectionFromID(event.ConnectionId)
					fmt.Printf("---- start swagger client for connection: %+v\n", conn)
					//rt := client.New("127.0.0.1", "/pets", nil)
					rt := client.New("127.0.0.1", "/api", nil)
					rt.Transport = conn
					cli := operations.New(rt, strfmt.Default)
					params := &operations.GetPetsParams{XCorrelationID: &event.ConnectionId}
					for i := 0; i < 100; i++ {
						pets, err := cli.GetPets(params)
						if err != nil {
							log.Fatal(err)
						}
						names := []string{}
						for _, p := range pets.Payload {
							names = append(names, *p.Name)
						}
						fmt.Printf("current pets: %v\n", names)
						time.Sleep(10 * time.Second)
					}
				}

				return
			case <-ctx.Done():
				return
			}
		}
	}()

	s := &http.Server{
		Addr:           ":8083",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
