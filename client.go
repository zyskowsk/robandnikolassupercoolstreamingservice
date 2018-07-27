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
	Port             int32
	RequestIdCounter int32
}

func Default() *Client {
	return &Client{RequestIdCounter: 0}
}

func (c Client) Name() string {
	return fmt.Sprintf("Client%d", c.Port)
}

func (c *Client) NextRequestId() string {
	c.RequestIdCounter++
	return fmt.Sprintf("%s:%d", c.Name(), c.RequestIdCounter)
}

func (c *Client) RequestPeersForSongId(song_id int32) ([]string, error) {
	fmt.Printf("%s requesting peers for song\n", c.Name())

	adr := fmt.Sprintf("%s:%d", SERVER_HOST, PEER_SERVER_PORT)
	conn, err := net.Dial(CONN_TYPE, adr)

	if err != nil {
		fmt.Println("Error connecting to PeerServer")
		fmt.Println(err)
		return []string{}, err
	}

	fmt.Printf("%s established a connection with PeerServer at %s\n", c.Name(), adr)

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
		fmt.Println("Error while Marshaling request")
		fmt.Println(err)
		return []string{}, err
	}

	n, err := conn.Write(data)

	if err != nil {
		fmt.Println("Error while writing bytes to connection")
		fmt.Println(err)
		return []string{}, err
	}

	fmt.Printf("%s sent request %s\n", c.Name(), req)
	fmt.Printf("%s sent %d bytes to PeerServer at %s\n", c.Name(), n, adr)

	data = make([]byte, DATA_BUF_SIZE)
	n, err = conn.Read(data)

	if err != nil {
		fmt.Println("Error while reading bytes from connection")
		fmt.Println(err)
		return []string{}, err
	}

	fmt.Printf("%s received %d bytes from PeerServer\n", c.Name(), n)

	res := &PeerServerResponse{}
	err = proto.Unmarshal(data[:n], res)

	if err != nil {
		fmt.Println("Error while Unmarshaling response")
		fmt.Println(err)
		return []string{}, err
	}

	fmt.Printf("%s received response %s\n", c.Name(), res)

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

func (c *Client) RequestSongChunk(id int32, chunkind int32) error {
	fmt.Printf("%s requesting song chunk\n", c.Name())

	adr := fmt.Sprintf("%s:%d", SERVER_HOST, SONG_SERVER_PORT)
	conn, err := net.Dial(CONN_TYPE, adr)

	if err != nil {
		fmt.Println("Error connecting to SongServer")
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s established a connection with SongServer at %s\n", c.Name(), adr)

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
		fmt.Println("Error while Marshaling request")
		fmt.Println(err)
		return err
	}

	n, err := conn.Write(data)

	if err != nil {
		fmt.Println("Error while writing bytes to connection")
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s sent request %s\n", c.Name(), req)
	fmt.Printf("%s sent %d bytes to SongServer at %s\n", c.Name(), n, adr)

	data = make([]byte, DATA_BUF_SIZE)
	n, err = conn.Read(data)

	if err != nil {
		fmt.Println("Error while reading bytes from connection")
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s received %d bytes from SongServer\n", c.Name(), n)

	res := &SongServerResponse{}
	err = proto.Unmarshal(data[:n], res)

	if err != nil {
		fmt.Println("Error while Unmarshaling response")
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s received response %s\n", c.Name(), res)

	// No error
	return nil
}
