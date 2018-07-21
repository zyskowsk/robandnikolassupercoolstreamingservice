package main

import (
    "fmt"
	"strings"
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
	song_bytes := bytes[song_id]
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

	fmt.Printf("Request comming, id = %d\n", song_id)

	var song_bytes = get_cached_song_by_id(song_id)
	fmt.Fprintf(w, string(song_bytes))
}

func get_uris(song_id int) []Uri {
	url := fmt.Sprintf("http://%s:%d?%s=%d", SERVER_HOST, SERVER_PORT, SONG_ID_GET_PARAM, song_id)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(2)
	}

	defer response.Body.Close()
	var uris []Uri
	body_content, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body_content, &uris)

	return uris
}

func get_song_from_client(song_id int, uri Uri) {
	url := fmt.Sprintf("http://%s:%d/%s?%s=%d", uri.Ip, uri.Host, CLIENT_ENDPOINT, SONG_ID_GET_PARAM, song_id)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(2)
	}

	defer response.Body.Close()
	body_content, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("response for song_id: %d, %s\n", song_id, body_content)
}

func main() {
	port := os.Args[1:][0]

    http.HandleFunc("/" + CLIENT_ENDPOINT, serve_cached_song_handler)
    go http.ListenAndServe(":" + port, nil)

	for {
		var response string
		fmt.Scanln(&response)
		if response == "n" {
			return
		}
		if response[:5] == "play:" {
			song_id, err := strconv.Atoi(strings.Split(response, ":")[1])
			if err != nil {
				os.Exit(2)
			}

			uris := get_uris(song_id)

			if len(uris) > 0 {
				get_song_from_client(song_id, uris[0])
			} else {
				fmt.Printf("Song is not available, song_id:%d\n", song_id)
			}
		}
	}
}
