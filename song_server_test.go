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
type FakeAddr struct{}

func (fa FakeAddr) Network() string {
	return ""
}

func (fa FakeAddr) String() string {
	return ""
}

// Fake TCP Connection, with fixed Read method
type FakeConn struct{}

func (fc FakeConn) Read(buff []byte) (int, error) {
	req := SongServerRequest{
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

	req_data, err := proto.Marshal(&req)
	if err != nil {
		return -1, err
	}

	copy(buff, req_data)

	return len(req_data), nil
}

func (c FakeConn) Write(b []byte) (n int, err error) {
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

	ss := SongServer{
		listener:      nil,
		response_chan: fake_response_chan,
	}

	err := ss.ProcessRequest(FakeConn{})

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
