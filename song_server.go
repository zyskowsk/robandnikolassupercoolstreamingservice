package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	CONN_HOST     = "127.0.0.1"
	CONN_TYPE     = "tcp"
	DATA_BUF_SIZE = 4096
)

type SongServer struct {
	listener      net.Listener
	port          int32
	response_chan chan SongServerResponse
}

func (s *SongServer) Run(port int32) error {
	s.port = port
	s.response_chan = make(chan SongServerResponse)
	adr := fmt.Sprintf("%s:%d", CONN_HOST, port)
	l, err := net.Listen(CONN_TYPE, adr)
	s.listener = l

	if err != nil {
		fmt.Println("Error while starting SongServer")
		fmt.Println(err)
		return err
	}

	fmt.Printf("SongServer listening on %s\n", adr)

	// Spin a new goroutine that logs responses
	go s.ResponseLogger()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			fmt.Println("Error while accepting connection, Listener closed")
			fmt.Println(err)
			return err
		}

		go s.ProcessRequest(conn)
	}

	// No error
	return nil
}

func (s *SongServer) ResponseLogger() {
	for {
		res := <-s.response_chan
		fmt.Printf("SongServer sent a reponse %+v\n", res)
	}
}

func (s *SongServer) Close() {
	fmt.Println("SongServer stopping")
	s.listener.Close()
}

func (s *SongServer) SongChunk(req SongChunkRequest) (SongChunkResponse, error) {
	// Not implemented
	song_chunk := SongChunk{
		Name:       "SongName", // random
		Id:         req.SongId,
		RawBytes:   []byte("rawsongbytes"), // random
		ChunkIndex: req.ChunkIndex,
		Offset:     0,   // random
		Size:       100, // random
	}

	return SongChunkResponse{
		SongChunk: &song_chunk,
	}, nil
}

func (s *SongServer) SendResponseToClient(conn net.Conn, res *SongServerResponse) error {
	s.response_chan <- *res
	data, err := proto.Marshal(res)
	if err != nil {
		fmt.Println("Error while Marshaling response")
		fmt.Println(err)
		return err
	}

	n, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error while writing bytes to connection")
		fmt.Println(err)
		return err
	}
	fmt.Printf("SongServer sent %d bytes to client\n", n)
	return nil

}

func (s *SongServer) Serve(req SongServerRequest) (SongServerResponse, error) {
	// Make sure request actually carries real request by inspecting the inner type
	switch x := req.Request.(type) {
	case *SongServerRequest_SongChunkRequest:
		song_req := x.SongChunkRequest

		fmt.Printf("Client %s is requesting song chunk at index %d for song %d\n", req.BaseRequest.ClientId, song_req.ChunkIndex, song_req.SongId)

		song_res, _ := s.SongChunk(*song_req)
		res := &SongServerResponse{
			BaseResponse: &BaseResponse{
				RequestId: req.BaseRequest.RequestId,
				ClientId:  req.BaseRequest.ClientId,
				Timestamp: int64(time.Now().Unix()),
			},
			Response: &SongServerResponse_SongChunkResponse{
				SongChunkResponse: &song_res,
			},
		}

		return *res, nil
	default:
		res := &SongServerResponse{
			BaseResponse: &BaseResponse{
				RequestId: req.BaseRequest.RequestId,
				ClientId:  req.BaseRequest.ClientId,
				Timestamp: int64(time.Now().Unix()),
			},
		}
		return *res, fmt.Errorf("Request has unexpected type %T", x)
	}
}

func (s *SongServer) ProcessRequest(conn net.Conn) error {
	defer conn.Close()

	data := make([]byte, DATA_BUF_SIZE)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error while reading bytes from connection")
		fmt.Println(err)
		return err
	}

	fmt.Printf("SongServer got %d bytes from (yet) unknown client\n", n)

	req := &SongServerRequest{}
	err = proto.Unmarshal(data[:n], req)

	if err != nil {
		res := &SongServerResponse{}
		fmt.Println(n)
		fmt.Println(data[:n])
		fmt.Println("Error while Unmarshaling request")
		fmt.Println(err)
		s.SendResponseToClient(conn, res)
		return err
	}

	fmt.Println("SongServer received the following request")
	fmt.Println(req)

	res, err := s.Serve(*req)

	err = s.SendResponseToClient(conn, &res)
	if err != nil {
		fmt.Println("Error while sending the response")
		fmt.Println(err)
		return err
	}

	return nil
}
