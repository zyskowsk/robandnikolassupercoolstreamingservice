package main

import "net"

// Interface for Server-like objects
type IRunnable interface {
	Run(port int32) error
	ProcessRequest(conn net.Conn) error
	Close()
}

type ISongServer interface {
	Serve(req SongServerRequest) (SongServerResponse, error)
	SongChunk(req SongChunkRequest) (SongChunkResponse, error)
}

type IPeerServer interface {
	Serve(req PeerServerRequest) (PeerServerResponse, error)
	ListPeers(req PeerListRequest) (PeerListResponse, error)
}
