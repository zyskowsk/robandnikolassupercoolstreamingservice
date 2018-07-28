package main

import (
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

var test_request_id = "1"
var test_client_id = "client1"
var test_timestamp = int64(100)
var test_song_id = int32(100)
var test_non_existing_song_id = int32(999)
var test_chunk_index = int32(789)

// Fake Address (net.Addr) used to mock TCP Connection
type FakeAddr struct {
}

func (fa FakeAddr) Network() string {
	return ""
}

func (fa FakeAddr) String() string {
	return ""
}

/*
	Fake TCP Connection
	It works by having two internal buffers:
		- read_buff = where we put requests for server to read
		- write_buff = where server puts responses for us to read
	For that reason, we have two extra metods (PutReqToReadBuffer, and ReadResFromWriteBuffer)
	Other methods are required by net.Conn interface
*/
type FakeConn struct {
	read_buff     []byte
	read_buff_len int
	write_buff    []byte
}

func NewFakeConn() FakeConn {
	return FakeConn{
		read_buff:     make([]byte, 512),
		read_buff_len: 0,
		write_buff:    make([]byte, 512),
	}
}

func (fc *FakeConn) PutReqToReadBuffer(req proto.Message) (int, error) {
	req_data, err := proto.Marshal(req)
	if err != nil {
		return -1, err
	}

	copy(fc.read_buff, req_data)
	fc.read_buff_len = len(req_data)

	return len(req_data), nil
}

func (fc *FakeConn) ReadSongServerResFromWriteBuffer(n int) (SongServerResponse, error) {
	res := &SongServerResponse{}
	err := proto.Unmarshal(fc.write_buff[:n], res)
	return *res, err
}

func (fc *FakeConn) ReadPeerServerResFromWriteBuffer(n int) (PeerServerResponse, error) {
	res := &PeerServerResponse{}
	err := proto.Unmarshal(fc.write_buff[:n], res)
	return *res, err
}

func (fc FakeConn) Read(buff []byte) (int, error) {
	copy(buff, fc.read_buff[:fc.read_buff_len])
	return fc.read_buff_len, nil
}

func (fc FakeConn) Write(b []byte) (n int, err error) {
	copy(fc.write_buff, b)
	return len(b), nil
}

func (fc FakeConn) Close() error {
	return nil
}

func (fc FakeConn) LocalAddr() net.Addr {
	return FakeAddr{}
}

func (fc FakeConn) RemoteAddr() net.Addr {
	return FakeAddr{}
}

func (c FakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (c FakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c FakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
