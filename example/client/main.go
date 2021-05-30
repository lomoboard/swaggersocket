package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-openapi/swaggersocket"
)

var myHandler = &handler{}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("client receive: %s\n", r.URL.Path)
	w.Write([]byte(r.URL.Path))
}

func main() {
	// create a websocket client
	u, err := url.Parse("ws://127.0.0.1:8082")
	if err != nil {
		panic(err)
	}

	wsClient := swaggersocket.NewWebSocketClient(swaggersocket.SocketClientOpts{
		KeepAlive: true, URL: u,
	}).WithMetaData("dummy connection metadata")
	// connect to the websocket server
	if err := wsClient.Connect(); err != nil {
		panic(err)
	}
	///////////////////
	// some code goes here
	//////////////////
	// serve the swagger api on the websocket client
	wsClient.Connection().Serve(context.Background(), myHandler)
	for {
	}
}
