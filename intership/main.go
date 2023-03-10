package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/xasan4/task2/api"
)

func main() {
	data, err := sql.Open("postgres", "host=localhost port=5432 user=sql-server-music dbname=music sslmode=disable password=Qwe123456")
	if err != nil {
		log.Fatal(err)
	}

	if err := data.Ping(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/musics", api.HandlerRequest)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
