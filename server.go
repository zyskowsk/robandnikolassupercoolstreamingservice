package main

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_TYPE = "tcp"

	DATA_BUF_SIZE = 4096
)

type SongServer struct {
	l net.Listener
	port int32
	response_chan chan SongResponse
}

func (s *SongServer) Run(port int32) (error) {
	s.port = port
	s.response_chan = make(chan SongResponse)
	adr := fmt.Sprintf("%s:%d", CONN_HOST, port)
	l, err := net.Listen(CONN_TYPE, adr)
	s.l = l

	if err != nil {
		return err
	}

	fmt.Printf("SongServer listening on %s\n", adr)

	go s.ResponseLogger()

	for {
		conn, err := s.l.Accept()

		if err != nil {
			return err
		}

		fmt.Print("Got a connection\n")
		go s.Process_request(conn)
	}

	return nil
}

func (s *SongServer) ResponseLogger() {
	for {
		res := <- s.response_chan
		fmt.Printf("Sent a reponse %s\n", res)
	}
}

func (s *SongServer) Close() {
	adr := fmt.Sprintf("%s:%d", CONN_HOST, s.port)
	fmt.Printf("SongServer stopping on %s\n", adr)
	s.l.Close()
}

func (s *SongServer) Get_song_chunk(req SongChunkRequest) (SongChunk, error) {
	return SongChunk{
		Name: "SongName",	// random
		Id: req.Id,
		RawBytes: []byte("rawsongbytes"),	// random
		ChunkIndex: req.ChunkIndex,
		Offset: 0,	// random
		Size: 100,	// random
	}, nil
}

func (s *SongServer) Process_request(conn net.Conn) (error) {
	defer conn.Close()

	data := make([]byte, DATA_BUF_SIZE)
	n, err := conn.Read(data)
	fmt.Printf("Got %d bytes from client\n", n)

	req := &SongRequest{}
	err = proto.Unmarshal(data[:n], req)

	if err != nil {
		s.response_chan <- SongResponse{}
		fmt.Println(err)
		return err
	}

	fmt.Printf("Received a request with length %d\n%s\n", n, req)

	// We got the request now, determine it's type
	switch x := req.Request.(type) {
		case *SongRequest_SongChunkRequest: 
			song_chunk_res, err := s.Get_song_chunk(*x.SongChunkRequest)

			if err != nil {
				fmt.Println(err)
				break
			}

			res := &SongResponse{
				RequestId: req.RequestId,
				ClientId: req.ClientId,
				Response: &SongResponse_SongChunkResponse{
					SongChunkResponse: &SongChunkResponse{
						Success: true,
						SongChunk: &song_chunk_res,
					},
				},
			}

			s.response_chan <- *res
			data, err = proto.Marshal(res)
			n, err = conn.Write(data)
			fmt.Printf("Sent %d bytes to client", n)
			// No error
			return nil
	default:
		return fmt.Errorf("Request has unexpected type %T", x)
	}

	s.response_chan <- SongResponse{}
	return fmt.Errorf("Invalid request type")
}