package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	PEER_SERVER_HOST = "127.0.0.1"
)

type PeerServer struct {
	listener      net.Listener
	port          int32
	response_chan chan PeerServerResponse
	peer_map      map[int32][]string
}

func (s *PeerServer) Run(port int32) error {
	s.port = port
	s.response_chan = make(chan PeerServerResponse)
	s.peer_map = make(map[int32][]string)

	// Adding some dummy data for now
	s.peer_map[100] = []string{
		fmt.Sprintf("%s:%d", "127.0.0.1", CLIENT_PORTS[1]),
	}

	adr := fmt.Sprintf("%s:%d", PEER_SERVER_HOST, port)
	l, err := net.Listen(CONN_TYPE, adr)
	s.listener = l

	if err != nil {
		LoglnDebug("Error while starting PeerServer")
		LoglnDebug(err)
		return err
	}

	LogDebug("PeerServer listening on %s\n", adr)

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

func (s *PeerServer) ResponseLogger() {
	for {
		res := <-s.response_chan
		LogDebug("PeerServer sent a reponse %+v\n", res)
	}
}

func (s *PeerServer) Close() {
	LoglnDebug("PeerServer stopping")
	s.listener.Close()
}

func (s *PeerServer) SendResponseToClient(conn net.Conn, res *PeerServerResponse) (int, error) {
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
	LogDebug("PeerServer sent %d bytes to client\n", n)
	return n, nil
}

func (s *PeerServer) Serve(req PeerServerRequest) (PeerServerResponse, error) {
	switch x := req.Request.(type) {
	case *PeerServerRequest_PeerListRequest:
		peer_list_req := x.PeerListRequest
		peers := s.peer_map[peer_list_req.SongId]

		return PeerServerResponse{
			BaseResponse: &BaseResponse{
				RequestId: req.BaseRequest.RequestId,
				ClientId:  req.BaseRequest.ClientId,
				Timestamp: int64(time.Now().Unix()),
			},
			Response: &PeerServerResponse_PeerListResponse{
				PeerListResponse: &PeerListResponse{
					SongId: peer_list_req.SongId,
					Peers:  peers,
				},
			},
		}, nil
	default:
		return PeerServerResponse{}, nil
	}
}

func (s *PeerServer) ProcessRequest(conn net.Conn) error {
	defer conn.Close()

	data := make([]byte, DATA_BUF_SIZE)
	n, err := conn.Read(data)
	if err != nil {
		LoglnDebug("Error while reading bytes from connection")
		LoglnDebug(err)
		return err
	}

	LogDebug("PeerServer got %d bytes from (yet) unknown client\n", n)

	req := &PeerServerRequest{}
	err = proto.Unmarshal(data[:n], req)

	if err != nil {
		res := &PeerServerResponse{}
		LoglnDebug("Error while Unmarshaling request")
		LoglnDebug(err)
		s.SendResponseToClient(conn, res)
		return err
	}

	LoglnDebug("PeerServer received the following request")
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
