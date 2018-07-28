package main

import (
	"testing"
)

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

	fake_conn.PutReqToReadBuffer(&test_req)

	ss := SongServer{
		listener:      nil,
		response_chan: fake_response_chan,
	}

	err := ss.ProcessRequest(fake_conn)

	if err != nil {
		t.Error(err)
	}

	ss_res := <-fake_response_chan

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
	if err != nil {
		t.Error(err)
	}

	res_from_chan := <-fake_response_chan

	if res_from_chan.BaseResponse.RequestId != test_res.BaseResponse.RequestId {
		t.Error("RequestId from Channel and RequestId from sent request don't match")
	}

	res_from_conn, err := fake_conn.ReadSongServerResFromWriteBuffer(n)
	if err != nil {
		t.Error(err)
	}

	if res_from_conn.BaseResponse.RequestId != test_res.BaseResponse.RequestId {
		t.Error("RequestId from Conn and RequestId from sent request don't match")
	}
}
