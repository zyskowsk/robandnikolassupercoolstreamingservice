package main

import (
	"time"
)

const (
	SONG_SERVER_PORT = 4001
	PEER_SERVER_PORT = 4002
)

// No array constants in Go
var CLIENT_PORTS = [...]int32{4003, 25012}

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

	done_chan := make(chan bool, 3)

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

	var c1 = Client{
		Port: CLIENT_PORTS[0],
	}
	var c2 = Client{
		Port: CLIENT_PORTS[1],
	}
	go func() {
		c1.Run()
	}()
	go func() {
		c2.Run()
	}()
	go func() {
		c1.PlaySong(100)
		done_chan <- true
	}()

	<-done_chan
	<-done_chan
	<-done_chan

	c2.Close()
	c1.Close()
	s.Close()
	ps.Close()
}
