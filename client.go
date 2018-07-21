package main

import (
    "fmt"
    "net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var SONG_ID_GET_PARAM = "id"
var SERVER_HOST = "localhost"
var SERVER_PORT = 8080
var CLIENT_ENDPOINT = "song"

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

func get_uris(song_id int) []string {
	url := fmt.Sprintf("http://%s:%d?%s=%d", SERVER_HOST, SERVER_PORT, SONG_ID_GET_PARAM, song_id)
	fmt.Printf("url = %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(2)
	}

	defer response.Body.Close()
	var uris []string
	body_content, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body_content, &uris)

	return uris
}

func get_song_from_client(song_id int, client_uri string) {
	url := fmt.Sprintf("http://%s/%s?%s=%d", client_uri, CLIENT_ENDPOINT, SONG_ID_GET_PARAM, song_id)
	fmt.Printf(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(2)
	}

	defer response.Body.Close()
	body_content, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("client body_content", body_content)
}

func main() {
	port := os.Args[1:][0]
	song_id, err := strconv.Atoi(os.Args[1:][1])
	if err != nil {
		os.Exit(2)
	}

    http.HandleFunc("/" + CLIENT_ENDPOINT, serve_cached_song_handler)
    go http.ListenAndServe(":" + port, nil)

	var response string
	fmt.Scanln(&response)
	if response == "n" {
		return
	}

	uris := get_uris(song_id)
	fmt.Printf("uris = %s", uris)

	get_song_from_client(song_id, uris[0])
}
