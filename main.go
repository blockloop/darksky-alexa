package main

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/darksky", handle)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"pong": ""}`))
	})

	log.Println("listening at :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("shutting down: %v", err)
	}
	log.Println("shutting down")
}

func handle(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, ResponseText("great job!"))
}

func ResponseText(msg string) render.M {
	return render.M{
		"response": render.M{
			"outputSpeech": render.M{
				"type": "PlainText",
				"text": msg,
			},
			"shouldEndSession": true,
		},
		"sessionAttributes": render.M{},
	}
}
