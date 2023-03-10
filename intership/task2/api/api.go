package api

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var db sql.DB

type Music struct {
	Id       int
	Name     string
	Singer   string
	Duration int
}

func init() {
	data, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=goLANGn1nja")
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := data.Ping(); err != nil {
		log.Fatal(err)
		return
	}
}

func HandlerRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMusic(w, r)
	case http.MethodPost:
		addUser(w, r)
	case http.MethodDelete:
		deleteMusic(w, r)
	case http.MethodPut:
		updateMusic(w, r)
	}
}

func getMusic(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select * from playlist")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	defer rows.Close()

	musics := make([]Music, 0)

	for rows.Next() {
		m := Music{}
		if err := rows.Scan(&m.Id, &m.Name, &m.Singer, &m.Duration); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		musics = append(musics, m)
	}

	if err = rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	resp, err := json.Marshal(musics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var m Music

	if err := json.Unmarshal(reqBytes, &m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err = db.Exec("insert into playlist (name, singer, duraion) values ($1, $2, $3)", m.Name, m.Singer, m.Duration); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
}

func deleteMusic(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if _, err := db.Exec("delete from playlist where id=$1", id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func updateMusic(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id := r.URL.Query().Get("id")

	var m Music

	if err := json.Unmarshal(reqBytes, &m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := db.Exec("update playlist set name=$1, singer=$2, duration=$3 where id=$4", m.Name, m.Singer, m.Duration, id); err != nil {
		return 
	}
}