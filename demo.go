package main

import (
	"fmt"
	"time"
)

const (
	SONG_SERVER_PORT = 4001
	PEER_SERVER_PORT = 4002
)

// No array constants in Go
var CLIENT_PORTS = [...]int32{4003, 4004}

func main() {

	/*
		This demo starts an instance of SongServer at port SONG_SERVER_PORT
		Runs that instance in a separate Goroutine

		After that it starts an instance of PeerServer at port PEER_SERVER_PORT
		Runs that instance in a separate Goroutine

		After that, it creates an instance of Client
		and requests a chunk of a song using that client in a separate Goroutine.

		At the end, it's doing basic synchronization w/ done_channel, and finally stops the servers.

	*/

	done_chan := make(chan bool)

	var s = SongServer{}
	go func() {
		s.Run(SONG_SERVER_PORT)
		done_chan <- true
	}()

	var ps = PeerServer{}
	go func() {
		ps.Run(PEER_SERVER_PORT)
		done_chan <- true
	}()

	time.Sleep(1 * time.Second)

	var c = Client{
		Port: CLIENT_PORTS[0],
	}
	go func() {
		c.RequestSongChunk(100, 1)
		peers, _ := c.RequestPeersForSongId(100)
		fmt.Println("Peers:")
		for _, p := range peers {
			fmt.Println(p)
		}
		done_chan <- true
	}()

	<-done_chan
	<-done_chan
	<-done_chan

	s.Close()
	ps.Close()
}
