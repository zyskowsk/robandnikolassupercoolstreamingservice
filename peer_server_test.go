package main

import (
	"testing"
)

func Test_PeerServer_Serve(t *testing.T) {
	ps := PeerServer{}

	ps_res, err := ps.Serve(PeerServerRequest{
		BaseRequest: &BaseRequest{
			RequestId: test_request_id,
			ClientId:  test_client_id,
			Timestamp: test_timestamp,
		},
		Request: &PeerServerRequest_PeerListRequest{
			PeerListRequest: &PeerListRequest{
				SongId: test_song_id,
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	if ps_res.BaseResponse.RequestId != test_request_id {
		t.Error("RequestId on BaseRequest and BaseResponse doesn't match")
	}
	if ps_res.BaseResponse.ClientId != test_client_id {
		t.Error("ClientId on BaseRequest and BaseResponse doesn't match")
	}
	if ps_res.BaseResponse.Timestamp < test_timestamp {
		t.Error("Timestamps on req & res don't match")
	}

	// Maybe take a look at the *actual* response
}

func Test_PeerServer_ProcessRequest(t *testing.T) {
	fake_response_chan := make(chan PeerServerResponse, 1)
	fake_conn := NewFakeConn()

	test_req := PeerServerRequest{
		BaseRequest: &BaseRequest{
			RequestId: test_request_id,
			ClientId:  test_client_id,
			Timestamp: test_timestamp,
		},
		Request: &PeerServerRequest_PeerListRequest{
			PeerListRequest: &PeerListRequest{
				SongId: test_song_id,
			},
		},
	}

	fake_conn.PutReqToReadBuffer(&test_req)

	ss := PeerServer{
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

func Test_PeerServer_SendResponseToClient(t *testing.T) {
	fake_response_chan := make(chan PeerServerResponse, 1)
	fake_conn := NewFakeConn()

	test_res := PeerServerResponse{
		BaseResponse: &BaseResponse{
			RequestId: test_request_id,
			ClientId:  test_client_id,
			Timestamp: test_timestamp,
		},
	}

	ss := PeerServer{
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

	res_from_conn, err := fake_conn.ReadPeerServerResFromWriteBuffer(n)
	if err != nil {
		t.Error(err)
	}

	if res_from_conn.BaseResponse.RequestId != test_res.BaseResponse.RequestId {
		t.Error("RequestId from Conn and RequestId from sent request don't match")
	}
}
