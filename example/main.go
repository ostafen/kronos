package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(time.Now(), string(data))
	})
	http.ListenAndServe(":9091", nil)
}
