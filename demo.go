package main 

import (
	"time"
)

const (
	SERVER_PORT = 4001
)

// No array constants in Go
var CLIENT_PORTS = [...]int32 { 4002, 4003, 4004 }

func main() {
	// Start a server instance
	var s = SongServer{}
	go s.Run(SERVER_PORT)

	// Give server some time to spin up
	time.Sleep(1 * time.Second)

	// Start one client instance
	var c = SongClient{
		Port: CLIENT_PORTS[0],
	}
	go c.RequestSongChunk(100, 1)
	
	// Sleep so we make sure goroutines are executed
	time.Sleep(10 * time.Second)

	s.Close()
}