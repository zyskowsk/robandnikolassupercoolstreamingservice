package main

import (
    "fmt"
    "net/http"
	"log"
	"strconv"
	"os"
	"encoding/json"
)

var SONG_ID_GET_PARAM = "id"

func get_song_hosts(song_id int) []string {
	ips := map[int][]string{
		123: []string{"127.0.0.1:4002"},
		345: []string{"127.0.0.1:4003", "456.12.21.12"},
	}
	song_ips := ips[song_id]
	if song_ips == nil {
		return []string{}
	}
	return song_ips
}

func handler(w http.ResponseWriter, r *http.Request) {
	var song_id_str = r.URL.Query()[SONG_ID_GET_PARAM][0]
	song_id, err := strconv.Atoi(song_id_str)
	if err != nil {
		os.Exit(2)
	}

	fmt.Printf("Request comming, id = %d\n", song_id)
	var song_hosts = get_song_hosts(song_id)
	response, _ := json.Marshal(song_hosts)

	fmt.Fprintf(w, string(response))
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
