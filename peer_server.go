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
}

func (s *PeerServer) Run(port int32) error {
	s.port = port
	s.response_chan = make(chan PeerServerResponse)
	adr := fmt.Sprintf("%s:%d", PEER_SERVER_HOST, port)
	l, err := net.Listen(CONN_TYPE, adr)
	s.listener = l

	if err != nil {
		fmt.Println("Error while starting PeerServer")
		fmt.Println(err)
		return err
	}

	fmt.Printf("PeerServer listening on %s\n", adr)

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

func (s *PeerServer) ResponseLogger() {
	for {
		res := <-s.response_chan
		fmt.Printf("PeerServer sent a reponse %+v\n", res)
	}
}

func (s *PeerServer) Close() {
	fmt.Println("PeerServer stopping")
	s.listener.Close()
}

func (s *PeerServer) SendResponseToClient(conn net.Conn, res *PeerServerResponse) (int, error) {
	s.response_chan <- *res
	data, err := proto.Marshal(res)
	if err != nil {
		fmt.Println("Error while Marshaling response")
		fmt.Println(err)
		return 0, err
	}

	n, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error while writing bytes to connection")
		fmt.Println(err)
		return n, err
	}
	fmt.Printf("PeerServer sent %d bytes to client\n", n)
	return n, nil
}

func (s *PeerServer) Serve(req PeerServerRequest) (PeerServerResponse, error) {
	switch x := req.Request.(type) {
	case *PeerServerRequest_PeerListRequest:
		peer_list_req := x.PeerListRequest
		return PeerServerResponse{
			BaseResponse: &BaseResponse{
				RequestId: req.BaseRequest.RequestId,
				ClientId:  req.BaseRequest.ClientId,
				Timestamp: int64(time.Now().Unix()),
			},
			Response: &PeerServerResponse_PeerListResponse{
				PeerListResponse: &PeerListResponse{
					SongId: peer_list_req.SongId,
					Peers:  []string{"quavo", "takeoff"},
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
		fmt.Println("Error while reading bytes from connection")
		fmt.Println(err)
		return err
	}

	fmt.Printf("PeerServer got %d bytes from (yet) unknown client\n", n)

	req := &PeerServerRequest{}
	err = proto.Unmarshal(data[:n], req)

	if err != nil {
		res := &PeerServerResponse{}
		fmt.Println("Error while Unmarshaling request")
		fmt.Println(err)
		s.SendResponseToClient(conn, res)
		return err
	}

	fmt.Println("PeerServer received the following request")
	fmt.Println(req)

	res, err := s.Serve(*req)

	_, err = s.SendResponseToClient(conn, &res)
	if err != nil {
		fmt.Println("Error while sending the response")
		fmt.Println(err)
		return err
	}

	return nil
}
