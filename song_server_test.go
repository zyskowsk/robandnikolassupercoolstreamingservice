package main

import (
	"github.com/golang/protobuf/proto"
	"net"
	"testing"
	"time"
)

// Some constants
var test_request_id = "1"
var test_client_id = "client1"
var test_timestamp = int64(100)
var test_song_id = int32(456)
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

func (fc *FakeConn) PutReqToReadBuffer(req SongServerRequest) (int, error) {
	req_data, err := proto.Marshal(&req)
	if err != nil {
		return -1, err
	}

	copy(fc.read_buff, req_data)
	fc.read_buff_len = len(req_data)

	return len(req_data), nil
}

func (fc *FakeConn) ReadResFromWriteBuffer(n int) (SongServerResponse, error) {
	res := &SongServerResponse{}
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

func TestSanity(t *testing.T) {
	if true != true {
		t.Error("You're insane")
	}
}

func Test_SongServer_Serve(t *testing.T) {
	ss := SongServer{
		listener: nil,
	}

	ss_res, err := ss.Serve(SongServerRequest{
		BaseRequest: &BaseRequest{
			RequestId: test_request_id,
			ClientId:  test_client_id,
			Timestamp: test_timestamp,
		},
		Request: &SongServerRequest_SongChunkRequest{
			SongChunkRequest: &SongChunkRequest{
				SongId:     test_song_id,
				ChunkIndex: test_chunk_index,
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	if ss_res.BaseResponse.RequestId != test_request_id {
		t.Error("RequestId on BaseRequest and BaseResponse doesn't match")
	}
	if ss_res.BaseResponse.ClientId != test_client_id {
		t.Error("ClientId on BaseRequest and BaseResponse doesn't match")
	}
	if ss_res.BaseResponse.Timestamp < test_timestamp {
		t.Error("Timestamps on req & res don't match")
	}

	// Maybe take a look at the *actual* response
}

func Test_SongServer_ProcessRequest(t *testing.T) {
	// Buffered channel w/ size 1 since we're running everything in the same Goroutine
	// If we run ProcessRequest in a separate Goroutine this channel can be unbuffered
	// Much headache here, it's a great thing I'm not very good at thinking about concurrency I'll learn so much
	/*
		fake_response_chan := make(chan SongServerResponse)

		go func() {
			ss := SongServer{
				listener:      nil,
				response_chan: fake_response_chan,
			}

			err := ss.ProcessRequest(FakeConn{})

			if err != nil {
				t.Error(err)
			}
		}()
	*/
	fake_response_chan := make(chan SongServerResponse, 1)
	fake_conn := NewFakeConn()

	test_req := SongServerRequest{
		BaseRequest: &BaseRequest{
			RequestId: test_request_id,
			ClientId:  test_client_id,
			Timestamp: test_timestamp,
		},
		Request: &SongServerRequest_SongChunkRequest{
			SongChunkRequest: &SongChunkRequest{
				SongId:     test_song_id,
				ChunkIndex: test_chunk_index,
			},
		},
	}

	fake_conn.PutReqToReadBuffer(test_req)

	ss := SongServer{
		listener:      nil,
		response_chan: fake_response_chan,
	}

	err := ss.ProcessRequest(fake_conn)

	if err != nil {
		t.Error(err)
	}

	t.Log("Waiting for response")

	ss_res := <-fake_response_chan

	t.Log("Response received")

	if ss_res.BaseResponse.RequestId != test_request_id {
		t.Error("RequestId on BaseRequest and BaseResponse doesn't match")
	}
	if ss_res.BaseResponse.ClientId != test_client_id {
		t.Error("ClientId on BaseRequest and BaseResponse doesn't match")
	}
	if ss_res.BaseResponse.Timestamp < test_timestamp {
		t.Error("Timestamps on req & res don't match")
	}

	// Maybe take a look at the *actual* response
}

func Test_SongServer_SendResponseToClient(t *testing.T) {
	fake_response_chan := make(chan SongServerResponse, 1)
	fake_conn := NewFakeConn()

	test_res := SongServerResponse{
		BaseResponse: &BaseResponse{
			RequestId: test_request_id,
			ClientId:  test_client_id,
			Timestamp: test_timestamp,
		},
	}

	ss := SongServer{
		listener:      nil,
		response_chan: fake_response_chan,
	}

	n, err := ss.SendResponseToClient(fake_conn, &test_res)
	t.Logf("Sent %d\n", n)
	if err != nil {
		t.Error(err)
	}

	res_from_chan := <-fake_response_chan

	if res_from_chan.BaseResponse.RequestId != test_res.BaseResponse.RequestId {
		t.Error("RequestId from Channel and RequestId from sent request don't match")
	}

	res_from_conn, err := fake_conn.ReadResFromWriteBuffer(n)
	if err != nil {
		t.Error(err)
	}

	if res_from_conn.BaseResponse.RequestId != test_res.BaseResponse.RequestId {
		t.Error("RequestId from Conn and RequestId from sent request don't match")
	}
}
