package main

import (
	"time"
)

const (
	SERVER_PORT = 4001
)

// No array constants in Go
var CLIENT_PORTS = [...]int32{4002, 4003, 4004}

func main() {

	/*
		This demo starts an instance of SongServer at port SERVER_PORT
		Runs that instance in a separate Goroutine
		and then sleeps for 1 second, to give it time to spin up

		After that, it creates an instance of Client
		and requests a chunk of a song using that client in a separate Goroutine.

		At the end, it's doing basic synchronization w/ done_channel, and finally kills  the server.

	*/

	done_chan := make(chan bool)

	var s = SongServer{}
	go func() {
		s.Run(SERVER_PORT)
		done_chan <- true
	}()

	time.Sleep(1 * time.Second)

	var c = Client{
		Port: CLIENT_PORTS[0],
	}
	go func() {
		c.RequestSongChunk(100, 1)
		done_chan <- true
	}()

	<-done_chan
	<-done_chan

	s.Close()
}
