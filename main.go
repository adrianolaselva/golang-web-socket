package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/ws/transactions", handler)
	error := http.ListenAndServe(":3000", nil)
	if error != nil {
		log.Fatal(error)
	}
}

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func handler(writer http.ResponseWriter, request *http.Request) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
	}

	for {

		msgType, request, err := socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		start := time.Now()
		log.Println("request: ", string(request))

		response := processMessage(request)

		err = socket.WriteMessage(msgType, response)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("response: ", string(response), time.Since(start))
	}
}


type response struct {
	Message   string `json:"message"`
}

func processMessage(request []byte) []byte {
	var payload map[string]interface{}

	if err := json.Unmarshal(request, &payload); err != nil {
		response, _ := json.Marshal(&response{
			Message: "failed to process message",
		})
		return response
	}

	response, _ := json.Marshal(&response{
		Message: "message processed",
	})

	return response
}