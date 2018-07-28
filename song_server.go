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
	chunks        map[int32][]SongChunk
}

func (s *SongServer) Run(port int32) error {
	s.port = port
	s.response_chan = make(chan SongServerResponse)
	adr := fmt.Sprintf("%s:%d", CONN_HOST, port)
	l, err := net.Listen(CONN_TYPE, adr)
	s.listener = l

	// Add some dummy song data for now
	s.chunks = make(map[int32][]SongChunk)
	s.chunks[100] = []SongChunk{
		SongChunk{
			Name:       "Bodak Yellow",
			Id:         100,
			RawBytes:   []byte("brr skrrrt chunk numero uno"),
			ChunkIndex: 0,
			Offset:     0,
			Size:       10,
		},
		SongChunk{
			Name:       "Bodak Yellow",
			Id:         100,
			RawBytes:   []byte("i dont dance now i make money moves"),
			ChunkIndex: 1,
			Offset:     10,
			Size:       10,
		},
		SongChunk{
			Name:       "Bodak Yellow",
			Id:         100,
			RawBytes:   []byte("and that's pretty much it"),
			ChunkIndex: 2,
			Offset:     20,
			Size:       10,
		},
	}

	if err != nil {
		LoglnDebug("Error while starting SongServer")
		LoglnDebug(err)
		return err
	}

	LogDebug("SongServer listening on %s\n", adr)

	// Spin a new goroutine that logs responses
	go s.ResponseLogger()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			LoglnDebug("Error while accepting connection, Listener closed")
			LoglnDebug(err)
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
		LogDebug("SongServer sent a reponse %+v\n", res)
	}
}

func (s *SongServer) Close() {
	LoglnDebug("SongServer stopping")
	s.listener.Close()
}

func (s *SongServer) SongChunk(req SongChunkRequest) (SongChunkResponse, error) {

	song_chunks := s.chunks[req.SongId]
	song_chunk := SongChunk{
		RawBytes: []byte{},
	}
	for _, c := range song_chunks {
		if c.ChunkIndex == req.ChunkIndex {
			song_chunk = c
			break
		}
	}

	return SongChunkResponse{
		SongChunk: &song_chunk,
	}, nil
}

func (s *SongServer) SendResponseToClient(conn net.Conn, res *SongServerResponse) (int, error) {
	s.response_chan <- *res
	data, err := proto.Marshal(res)
	if err != nil {
		LoglnDebug("Error while Marshaling response")
		LoglnDebug(err)
		return 0, err
	}

	n, err := conn.Write(data)
	if err != nil {
		LoglnDebug("Error while writing bytes to connection")
		LoglnDebug(err)
		return n, err
	}
	LogDebug("SongServer sent %d bytes to client\n", n)
	return n, nil

}

func (s *SongServer) Serve(req SongServerRequest) (SongServerResponse, error) {
	// Make sure request actually carries real request by inspecting the inner type
	switch x := req.Request.(type) {
	case *SongServerRequest_SongChunkRequest:
		song_req := x.SongChunkRequest

		LogDebug("Client %s is requesting song chunk at index %d for song %d\n", req.BaseRequest.ClientId, song_req.ChunkIndex, song_req.SongId)

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
		LoglnDebug("Error while reading bytes from connection")
		LoglnDebug(err)
		return err
	}

	LogDebug("SongServer got %d bytes from (yet) unknown client\n", n)

	req := &SongServerRequest{}
	err = proto.Unmarshal(data[:n], req)

	if err != nil {
		res := &SongServerResponse{}
		LoglnDebug("Error while Unmarshaling request")
		LoglnDebug(err)
		s.SendResponseToClient(conn, res)
		return err
	}

	LoglnDebug("SongServer received the following request")
	LoglnDebug(req)

	res, err := s.Serve(*req)

	_, err = s.SendResponseToClient(conn, &res)
	if err != nil {
		LoglnDebug("Error while sending the response")
		LoglnDebug(err)
		return err
	}

	return nil
}
