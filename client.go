package main

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
)

const (
	SERVER_HOST = "127.0.0.1"
)

type SongClient struct {
	Port int32
}

func (c SongClient) Name() (string) {
	return fmt.Sprintf("Client%d", c.Port)
}

func (c *SongClient) RequestSongChunk(id int32, chunkind int32) (error) {
	fmt.Printf("%s requesting song chunk\n", c.Name())

	adr := fmt.Sprintf("%s:%d", SERVER_HOST, SERVER_PORT)
	conn, err := net.Dial(CONN_TYPE, adr)

	if err != nil {
		fmt.Println(err)
		return err
	}
	
	// learned what defer means, basically execute this line just before any `return` anywhere in the this function
	// similar to try-catch-finally in more traditional languages
	// nicer way to reduce boilerplate code 
	defer conn.Close()

	req := &SongRequest{
		RequestId: "arbitratyid",
		ClientId: c.Name(),
		Request: &SongRequest_SongChunkRequest{
			SongChunkRequest: &SongChunkRequest{
				Id: id,
				ChunkIndex: chunkind,
			},
		},
	}
	data, err := proto.Marshal(req)
	n, err := conn.Write(data)

	fmt.Printf("Sent %d bytes to SongServer at %s\n", n, adr)

	data = make([]byte, DATA_BUF_SIZE)
	n, err = conn.Read(data)
	res := &SongResponse{}
	err = proto.Unmarshal(data[:n], res)

	if err != nil {
		fmt.Println(err)
		return err
	}
	
	fmt.Printf("Received a response with length %d\n%s\n", n, res)
	return nil
}