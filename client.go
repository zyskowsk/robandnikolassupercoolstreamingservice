package main

import (
    "fmt"
    "net/http"
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var SONG_ID_GET_PARAM = "id"
var SERVER_HOST = "localhost"
var SERVER_PORT = 8080

func get_cached_song_by_id(song_id int) []byte {
	bytes := map[int][]byte{
		123: []byte("firstsongbytes"),
		345: []byte("secondsongbytes"),
	}
	fmt.Printf("song_id = %d\n", song_id)
	song_bytes := bytes[song_id]

	fmt.Printf("song_bytes = %s\n", string(song_bytes))
	if song_bytes == nil {
		return []byte{}
	}
	return song_bytes
}

func serve_cached_song_handler(w http.ResponseWriter, r *http.Request) {
	song_id, err := strconv.Atoi(r.URL.Query()[SONG_ID_GET_PARAM][0])
	if err != nil {
		os.Exit(2)
	}

	var song_bytes = get_cached_song_by_id(song_id)
	fmt.Fprintf(w, string(song_bytes))
}

func get_song(song_id int) {
	url := fmt.Sprintf("http://%s:%d?%s=%d", SERVER_HOST, SERVER_PORT, SONG_ID_GET_PARAM, song_id)
	fmt.Printf("url = %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(2)
	}

	defer response.Body.Close()
	var ips map[string]interface{}
	body_content, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body_content, &ips)
	fmt.Printf("ips = %s", ips)
}

func main() {
	get_song(123)

    http.HandleFunc("/cached-song", serve_cached_song_handler)
    log.Fatal(http.ListenAndServe(":4002", nil))
}
