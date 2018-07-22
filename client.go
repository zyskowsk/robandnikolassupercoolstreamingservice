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
	RequestIdCounter int32
}

func Default() *SongClient {
    return &SongClient{RequestIdCounter: 0}
}

func (c SongClient) Name() (string) {
	return fmt.Sprintf("Client%d", c.Port)
}

func (c *SongClient) NextRequestId() (string) {
	c.RequestIdCounter++
	return fmt.Sprintf("%s:%d", c.Name(), c.RequestIdCounter)
}

func (c *SongClient) RequestSongChunk(id int32, chunkind int32) (error) {
	fmt.Printf("%s requesting song chunk\n", c.Name())

	adr := fmt.Sprintf("%s:%d", SERVER_HOST, SERVER_PORT)
	conn, err := net.Dial(CONN_TYPE, adr)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s established a connection with SongServer at %s\n", c.Name(), adr)
	
	// Learned what defer means, basically execute this line just before any `return` anywhere in the this function
	// Similar to try-catch-finally in more traditional languages, nicer way to reduce boilerplate code 
	defer conn.Close()

	req := &SongRequest{
		RequestId: c.NextRequestId(),
		ClientId: c.Name(),
		Request: &SongRequest_SongChunkRequest{
			SongChunkRequest: &SongChunkRequest{
				Id: id,
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

	res := &SongResponse{}
	err = proto.Unmarshal(data[:n], res)

	if err != nil {
		fmt.Println("Error while Unmarshaling response")
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s received response %s\n", c.Name(), res)
	// assert it's actually response to our request

	// No error
	return nil
}