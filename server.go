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
		fmt.Println("Error while starting SongServer")
		fmt.Println(err)
		return err
	}

	fmt.Printf("SongServer listening on %s\n", adr)

	// Spin a new goroutine that logs responses
	go s.ResponseLogger()

	for {
		conn, err := s.l.Accept()

		if err != nil {
			fmt.Println("Error while accepting connection, Listener closed")
			fmt.Println(err)
			return err
		}

		go s.Process_request(conn)
	}

	// No error
	return nil
}

func (s *SongServer) ResponseLogger() {
	for {
		res := <- s.response_chan
		fmt.Printf("SongServer sent a reponse %s\n", res)
	}
}

func (s *SongServer) Close() {
	fmt.Println("SongServer stopping")
	s.l.Close()
}

func (s *SongServer) Get_song_chunk(req SongChunkRequest) (SongChunk, error) {
	// Not implemented
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

	send_response_to_client := func (conn net.Conn, res *SongResponse) (error) {
		// Log the response
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
		fmt.Printf("SongServer sent %d bytes to client %s\n", n, res.ClientId)
		return nil
	}

	data := make([]byte, DATA_BUF_SIZE)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error while reading bytes from connection")
		fmt.Println(err)
		return err
	}

	fmt.Printf("SongServer got %d bytes from (yet) unknown client\n", n)

	req := &SongRequest{}
	err = proto.Unmarshal(data[:n], req)

	if err != nil {
		res := &SongResponse{
			RequestId: req.RequestId,
			ClientId: req.ClientId,
		}
		fmt.Println("Error while Unmarshaling request")
		fmt.Println(err)
		send_response_to_client(conn, res)
		return err
	}

	fmt.Printf("SongServer received a request from client %s\n", req.ClientId)
	fmt.Println(req)

	// We got the request now, determine it's type
	switch x := req.Request.(type) {
		case *SongRequest_SongChunkRequest: 
			song_chunk_res, err := s.Get_song_chunk(*x.SongChunkRequest)

			if err != nil {
				fmt.Println("Error while computing song chunk")
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

			send_response_to_client(conn, res)
			// No error
			return nil
	default:
		return fmt.Errorf("Request has unexpected type %T", x)
	}

	res := &SongResponse{
		RequestId: req.RequestId,
		ClientId: req.ClientId,
	}
	send_response_to_client(conn, res)
	return fmt.Errorf("Invalid request type")
}