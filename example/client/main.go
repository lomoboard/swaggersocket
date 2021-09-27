package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swaggersocket"
	"github.com/go-openapi/swaggersocket/example/autogen/models"
	"github.com/go-openapi/swaggersocket/example/autogen/restapi"
	"github.com/go-openapi/swaggersocket/example/autogen/restapi/operations"
)

func main() {
	// swagger api
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	pets := []*models.Pet{}
	api := operations.NewSwaggerPetstoreAPI(swaggerSpec)
	api.GetPetsHandler = operations.GetPetsHandlerFunc(func(operations.GetPetsParams) middleware.Responder {
		ok := operations.NewGetPetsOK()
		size := int64(len(pets) + 1)
		if size%2 == 0 {
			fmt.Printf("-> receve get pet request: add one cat\n")
			name := fmt.Sprintf("cat-%d", size)
			pets = append(pets, &models.Pet{ID: &size, Name: &name, Tag: "cat"})
		} else {
			fmt.Printf("-> receve get pet request: add one dog\n")
			name := fmt.Sprintf("dog-%d", size)
			pets = append(pets, &models.Pet{ID: &size, Name: &name, Tag: "dog"})
		}
		return ok.WithPayload(pets)
	})

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
	wsClient.Connection().Serve(context.Background(), api.Serve(nil))
	for {
	}
}
