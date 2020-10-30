package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type LeaderboardData struct {
	Score    float64 `json:"score"`
	Position int     `json:"position"`
}

type ResponseData struct {
	Data map[int]LeaderboardData
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := `
		<h1>Gzip JSON</h1>
		<h2><a href="http://localhost:8888/normal">/normal</a><h2>
		<h2><a href="http://localhost:8888/zip">/zip</a></h2>
		`
		w.Write([]byte(resp))
	})
	mux.HandleFunc("/normal", NormalHandler)
	mux.HandleFunc("/zip", ZipHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
	})

	log.Println("Server listening on http://localhost:8888")
	http.ListenAndServe(":8888", c.Handler(mux))
}

func NormalHandler(w http.ResponseWriter, r *http.Request) {
	response, err := getData()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ZipHandler(w http.ResponseWriter, r *http.Request) {

	response, err := getData()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var gzipEnabled bool
	acceptedEncodings := r.Header.Get("Accept-Encoding")

	for _, encoding := range strings.Split(acceptedEncodings, ",") {
		if strings.ToLower(strings.TrimSpace(encoding)) == "gzip" {
			gzipEnabled = true
			w.Header().Set("Content-Encoding", "gzip")
		}
	}

	if gzipEnabled {
		w.WriteHeader(http.StatusOK)
		gz := gzip.NewWriter(w)

		defer func() {
			if err := gz.Close(); err != nil {
				log.Println(err)
			}
		}()

		err := json.NewEncoder(gz).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error processing action"}`))
			return
		}
		return
	}

}

func getData() (*ResponseData, error) {
	bs, err := ioutil.ReadFile("./leaderboarddata.json")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var res ResponseData
	if err := json.Unmarshal(bs, &res); err != nil {
		log.Println(err)
		return nil, err
	}

	return &res, nil

}
