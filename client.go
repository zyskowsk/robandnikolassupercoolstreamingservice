package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	SERVER_HOST = "127.0.0.1"
)

type Client struct {
	Port                int32
	RequestIdCounter    int32
	SongServer          SongServer
	play_chan           chan (SongChunk)
	need_new_chunk_chan chan (bool)
}

func Default() *Client {
	return &Client{RequestIdCounter: 0}
}

func (c Client) Name() string {
	return fmt.Sprintf("Client%d", c.Port)
}

func (c *Client) Run() {
	c.play_chan = make(chan SongChunk, 10)
	c.need_new_chunk_chan = make(chan bool)
	c.SongServer.Run(c.Port)
}

func (c *Client) Close() {
	c.SongServer.Close()
}

func (c *Client) NextRequestId() string {
	c.RequestIdCounter++
	return fmt.Sprintf("%s:%d", c.Name(), c.RequestIdCounter)
}

func (c *Client) Play() {
	for {
		chunk := <-c.play_chan
		Log("Playing chunk: %s\n", string(chunk.RawBytes))
		duration := int(chunk.Size)
		// Sleep for half the duration, wake up to notify that it's time for new chunk
		time.Sleep(time.Duration(duration/2) * time.Second)

		c.need_new_chunk_chan <- true
		// Continue sleeping
		time.Sleep(time.Duration(duration/2) * time.Second)
	}
}

func (c *Client) RequestPeersForSongId(song_id int32) ([]string, error) {
	LogDebug("%s requesting peers for song\n", c.Name())

	adr := fmt.Sprintf("%s:%d", SERVER_HOST, PEER_SERVER_PORT)
	conn, err := net.Dial(CONN_TYPE, adr)

	if err != nil {
		Logln("Error connecting to PeerServer")
		Logln(err)
		return []string{}, err
	}

	LogDebug("%s established a connection with PeerServer at %s\n", c.Name(), adr)

	defer conn.Close()

	req := &PeerServerRequest{
		BaseRequest: &BaseRequest{
			RequestId: c.NextRequestId(),
			ClientId:  c.Name(),
			Timestamp: int64(time.Now().Unix()),
		},
		Request: &PeerServerRequest_PeerListRequest{
			PeerListRequest: &PeerListRequest{
				SongId: song_id,
			},
		},
	}

	data, err := proto.Marshal(req)

	if err != nil {
		Logln("Error while Marshaling request")
		Logln(err)
		return []string{}, err
	}

	n, err := conn.Write(data)

	if err != nil {
		Logln("Error while writing bytes to connection")
		Logln(err)
		return []string{}, err
	}

	LogDebug("%s sent request %s\n", c.Name(), req)
	LogDebug("%s sent %d bytes to PeerServer at %s\n", c.Name(), n, adr)

	data = make([]byte, DATA_BUF_SIZE)
	n, err = conn.Read(data)

	if err != nil {
		Logln("Error while reading bytes from connection")
		Logln(err)
		return []string{}, err
	}

	LogDebug("%s received %d bytes from PeerServer\n", c.Name(), n)

	res := &PeerServerResponse{}
	err = proto.Unmarshal(data[:n], res)

	if err != nil {
		Logln("Error while Unmarshaling response")
		Logln(err)
		return []string{}, err
	}

	LogDebug("%s received response %s\n", c.Name(), res)

	switch x := res.Response.(type) {
	case *PeerServerResponse_PeerListResponse:
		peer_list_res := x.PeerListResponse
		return peer_list_res.Peers, nil
	default:
		return []string{}, err
	}

	// No error
	return []string{}, nil
}

func (c *Client) PlaySong(song_id int32) error {
	song_server_adr := fmt.Sprintf("%s:%d", SERVER_HOST, SONG_SERVER_PORT)
	chunkind := int32(0)

	chunk, err := c.RequestSongChunk(song_server_adr, song_id, chunkind)
	if err != nil {
		Logln("Error while playing song")
		return err
	}

	peers, err := c.RequestPeersForSongId(song_id)
	if err != nil {
		Logln("Error while fetching peers")
		return err
	}

	// Start the play Goroutine
	go c.Play()

	for len(chunk.RawBytes) > 0 {

		// Queue the chunk, and block on "need new chunk" signal
		c.play_chan <- chunk
		<-c.need_new_chunk_chan

		chunkind++

		curr_peer := peers[0]
		chunk, err = c.RequestSongChunk(curr_peer, song_id, chunkind)

		if err != nil {
			Logln("Error while fetching song chunk from Peer")
			Logln("Going to try to fetch chunk from SongServer")
			chunk, err = c.RequestSongChunk(song_server_adr, song_id, chunkind)
			if err != nil {
				Logln("Erro while fetching song chunk from SongServer")
				return err
			}

			Logln("Remove the faulty Peer")
			if len(peers) > 1 {
				peers = peers[1:]
			} else {
				Logln("No more peers")
				peers, err = c.RequestPeersForSongId(song_id)
				if err != nil {
					Logln("Error while fetching peers")
					return err
				}
			}
		}
	}

	Logln("No more chunks, song is over")

	return nil
}

func (c *Client) RequestSongChunk(adr string, id int32, chunkind int32) (SongChunk, error) {
	LogDebug("%s requesting song chunk\n", c.Name())

	// adr := fmt.Sprintf("%s:%d", SERVER_HOST, SONG_SERVER_PORT)
	conn, err := net.Dial(CONN_TYPE, adr)

	if err != nil {
		Logln("Error connecting to SongServer")
		Logln(err)
		return SongChunk{}, err
	}

	LogDebug("%s established a connection with SongServer at %s\n", c.Name(), adr)

	// Learned what defer means, basically execute this line just before any `return` anywhere in the this function
	// Similar to try-catch-finally in more traditional languages, nicer way to reduce boilerplate code
	defer conn.Close()

	req := &SongServerRequest{
		BaseRequest: &BaseRequest{
			RequestId: c.NextRequestId(),
			ClientId:  c.Name(),
			Timestamp: int64(time.Now().Unix()),
		},
		Request: &SongServerRequest_SongChunkRequest{
			SongChunkRequest: &SongChunkRequest{
				SongId:     id,
				ChunkIndex: chunkind,
			},
		},
	}

	data, err := proto.Marshal(req)

	if err != nil {
		Logln("Error while Marshaling request")
		Logln(err)
		return SongChunk{}, err
	}

	n, err := conn.Write(data)

	if err != nil {
		Logln("Error while writing bytes to connection")
		Logln(err)
		return SongChunk{}, err
	}

	LogDebug("%s sent request %s\n", c.Name(), req)
	LogDebug("%s sent %d bytes to SongServer at %s\n", c.Name(), n, adr)

	data = make([]byte, DATA_BUF_SIZE)
	n, err = conn.Read(data)

	if err != nil {
		Logln("Error while reading bytes from connection")
		Logln(err)
		return SongChunk{}, err
	}

	LogDebug("%s received %d bytes from SongServer\n", c.Name(), n)

	res := &SongServerResponse{}
	err = proto.Unmarshal(data[:n], res)

	if err != nil {
		Logln("Error while Unmarshaling response")
		Logln(err)
		return SongChunk{}, err
	}

	LogDebug("%s received response %s\n", c.Name(), res)

	switch x := res.Response.(type) {
	case *SongServerResponse_SongChunkResponse:
		song_chunk_res := x.SongChunkResponse
		return *song_chunk_res.SongChunk, nil
	default:
		return SongChunk{}, nil
	}
}
