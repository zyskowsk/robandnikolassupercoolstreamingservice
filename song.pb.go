// Code generated by protoc-gen-go. DO NOT EDIT.
// source: song.proto

package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SongChunk struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id                   int32    `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	RawBytes             []byte   `protobuf:"bytes,3,opt,name=raw_bytes,json=rawBytes,proto3" json:"raw_bytes,omitempty"`
	ChunkIndex           int32    `protobuf:"varint,4,opt,name=chunk_index,json=chunkIndex,proto3" json:"chunk_index,omitempty"`
	Offset               int32    `protobuf:"varint,5,opt,name=offset,proto3" json:"offset,omitempty"`
	Size                 int32    `protobuf:"varint,6,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SongChunk) Reset()         { *m = SongChunk{} }
func (m *SongChunk) String() string { return proto.CompactTextString(m) }
func (*SongChunk) ProtoMessage()    {}
func (*SongChunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_song_b96fd28c72daf125, []int{0}
}
func (m *SongChunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SongChunk.Unmarshal(m, b)
}
func (m *SongChunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SongChunk.Marshal(b, m, deterministic)
}
func (dst *SongChunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SongChunk.Merge(dst, src)
}
func (m *SongChunk) XXX_Size() int {
	return xxx_messageInfo_SongChunk.Size(m)
}
func (m *SongChunk) XXX_DiscardUnknown() {
	xxx_messageInfo_SongChunk.DiscardUnknown(m)
}

var xxx_messageInfo_SongChunk proto.InternalMessageInfo

func (m *SongChunk) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SongChunk) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SongChunk) GetRawBytes() []byte {
	if m != nil {
		return m.RawBytes
	}
	return nil
}

func (m *SongChunk) GetChunkIndex() int32 {
	if m != nil {
		return m.ChunkIndex
	}
	return 0
}

func (m *SongChunk) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *SongChunk) GetSize() int32 {
	if m != nil {
		return m.Size
	}
	return 0
}

type SongRequest struct {
	RequestId string `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	ClientId  string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// Types that are valid to be assigned to Request:
	//	*SongRequest_SongChunkRequest
	Request              isSongRequest_Request `protobuf_oneof:"request"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *SongRequest) Reset()         { *m = SongRequest{} }
func (m *SongRequest) String() string { return proto.CompactTextString(m) }
func (*SongRequest) ProtoMessage()    {}
func (*SongRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_song_b96fd28c72daf125, []int{1}
}
func (m *SongRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SongRequest.Unmarshal(m, b)
}
func (m *SongRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SongRequest.Marshal(b, m, deterministic)
}
func (dst *SongRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SongRequest.Merge(dst, src)
}
func (m *SongRequest) XXX_Size() int {
	return xxx_messageInfo_SongRequest.Size(m)
}
func (m *SongRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SongRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SongRequest proto.InternalMessageInfo

type isSongRequest_Request interface {
	isSongRequest_Request()
}

type SongRequest_SongChunkRequest struct {
	SongChunkRequest *SongChunkRequest `protobuf:"bytes,3,opt,name=song_chunk_request,json=songChunkRequest,proto3,oneof"`
}

func (*SongRequest_SongChunkRequest) isSongRequest_Request() {}

func (m *SongRequest) GetRequest() isSongRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *SongRequest) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *SongRequest) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *SongRequest) GetSongChunkRequest() *SongChunkRequest {
	if x, ok := m.GetRequest().(*SongRequest_SongChunkRequest); ok {
		return x.SongChunkRequest
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*SongRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SongRequest_OneofMarshaler, _SongRequest_OneofUnmarshaler, _SongRequest_OneofSizer, []interface{}{
		(*SongRequest_SongChunkRequest)(nil),
	}
}

func _SongRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SongRequest)
	// request
	switch x := m.Request.(type) {
	case *SongRequest_SongChunkRequest:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SongChunkRequest); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SongRequest.Request has unexpected type %T", x)
	}
	return nil
}

func _SongRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SongRequest)
	switch tag {
	case 3: // request.song_chunk_request
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SongChunkRequest)
		err := b.DecodeMessage(msg)
		m.Request = &SongRequest_SongChunkRequest{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SongRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SongRequest)
	// request
	switch x := m.Request.(type) {
	case *SongRequest_SongChunkRequest:
		s := proto.Size(x.SongChunkRequest)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type SongResponse struct {
	RequestId string `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	ClientId  string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// Types that are valid to be assigned to Response:
	//	*SongResponse_SongChunkResponse
	Response             isSongResponse_Response `protobuf_oneof:"response"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SongResponse) Reset()         { *m = SongResponse{} }
func (m *SongResponse) String() string { return proto.CompactTextString(m) }
func (*SongResponse) ProtoMessage()    {}
func (*SongResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_song_b96fd28c72daf125, []int{2}
}
func (m *SongResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SongResponse.Unmarshal(m, b)
}
func (m *SongResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SongResponse.Marshal(b, m, deterministic)
}
func (dst *SongResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SongResponse.Merge(dst, src)
}
func (m *SongResponse) XXX_Size() int {
	return xxx_messageInfo_SongResponse.Size(m)
}
func (m *SongResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SongResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SongResponse proto.InternalMessageInfo

type isSongResponse_Response interface {
	isSongResponse_Response()
}

type SongResponse_SongChunkResponse struct {
	SongChunkResponse *SongChunkResponse `protobuf:"bytes,3,opt,name=song_chunk_response,json=songChunkResponse,proto3,oneof"`
}

func (*SongResponse_SongChunkResponse) isSongResponse_Response() {}

func (m *SongResponse) GetResponse() isSongResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *SongResponse) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *SongResponse) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *SongResponse) GetSongChunkResponse() *SongChunkResponse {
	if x, ok := m.GetResponse().(*SongResponse_SongChunkResponse); ok {
		return x.SongChunkResponse
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*SongResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SongResponse_OneofMarshaler, _SongResponse_OneofUnmarshaler, _SongResponse_OneofSizer, []interface{}{
		(*SongResponse_SongChunkResponse)(nil),
	}
}

func _SongResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SongResponse)
	// response
	switch x := m.Response.(type) {
	case *SongResponse_SongChunkResponse:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SongChunkResponse); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SongResponse.Response has unexpected type %T", x)
	}
	return nil
}

func _SongResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SongResponse)
	switch tag {
	case 3: // response.song_chunk_response
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SongChunkResponse)
		err := b.DecodeMessage(msg)
		m.Response = &SongResponse_SongChunkResponse{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SongResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SongResponse)
	// response
	switch x := m.Response.(type) {
	case *SongResponse_SongChunkResponse:
		s := proto.Size(x.SongChunkResponse)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type SongChunkRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ChunkIndex           int32    `protobuf:"varint,2,opt,name=chunk_index,json=chunkIndex,proto3" json:"chunk_index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SongChunkRequest) Reset()         { *m = SongChunkRequest{} }
func (m *SongChunkRequest) String() string { return proto.CompactTextString(m) }
func (*SongChunkRequest) ProtoMessage()    {}
func (*SongChunkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_song_b96fd28c72daf125, []int{3}
}
func (m *SongChunkRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SongChunkRequest.Unmarshal(m, b)
}
func (m *SongChunkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SongChunkRequest.Marshal(b, m, deterministic)
}
func (dst *SongChunkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SongChunkRequest.Merge(dst, src)
}
func (m *SongChunkRequest) XXX_Size() int {
	return xxx_messageInfo_SongChunkRequest.Size(m)
}
func (m *SongChunkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SongChunkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SongChunkRequest proto.InternalMessageInfo

func (m *SongChunkRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SongChunkRequest) GetChunkIndex() int32 {
	if m != nil {
		return m.ChunkIndex
	}
	return 0
}

type SongChunkResponse struct {
	Success              bool       `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	SongChunk            *SongChunk `protobuf:"bytes,2,opt,name=song_chunk,json=songChunk,proto3" json:"song_chunk,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SongChunkResponse) Reset()         { *m = SongChunkResponse{} }
func (m *SongChunkResponse) String() string { return proto.CompactTextString(m) }
func (*SongChunkResponse) ProtoMessage()    {}
func (*SongChunkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_song_b96fd28c72daf125, []int{4}
}
func (m *SongChunkResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SongChunkResponse.Unmarshal(m, b)
}
func (m *SongChunkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SongChunkResponse.Marshal(b, m, deterministic)
}
func (dst *SongChunkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SongChunkResponse.Merge(dst, src)
}
func (m *SongChunkResponse) XXX_Size() int {
	return xxx_messageInfo_SongChunkResponse.Size(m)
}
func (m *SongChunkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SongChunkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SongChunkResponse proto.InternalMessageInfo

func (m *SongChunkResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *SongChunkResponse) GetSongChunk() *SongChunk {
	if m != nil {
		return m.SongChunk
	}
	return nil
}

func init() {
	proto.RegisterType((*SongChunk)(nil), "main.SongChunk")
	proto.RegisterType((*SongRequest)(nil), "main.SongRequest")
	proto.RegisterType((*SongResponse)(nil), "main.SongResponse")
	proto.RegisterType((*SongChunkRequest)(nil), "main.SongChunkRequest")
	proto.RegisterType((*SongChunkResponse)(nil), "main.SongChunkResponse")
}

func init() { proto.RegisterFile("song.proto", fileDescriptor_song_b96fd28c72daf125) }

var fileDescriptor_song_b96fd28c72daf125 = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xcd, 0x4e, 0xc2, 0x40,
	0x10, 0xc7, 0xd9, 0x0a, 0xa5, 0x9d, 0x12, 0x85, 0x35, 0xc1, 0x26, 0xc6, 0x48, 0x7a, 0xe2, 0xd4,
	0x03, 0xbe, 0x41, 0x49, 0x0c, 0xbd, 0xae, 0x37, 0x13, 0xd3, 0x94, 0x76, 0xc1, 0x8d, 0xb2, 0xc5,
	0x6e, 0x09, 0xea, 0x9b, 0x98, 0x78, 0xf1, 0x4d, 0xcd, 0xce, 0x2e, 0x24, 0x96, 0xa3, 0xb7, 0x99,
	0xff, 0x64, 0x66, 0x7e, 0xf3, 0x01, 0xa0, 0x2a, 0xb9, 0x8e, 0xb7, 0x75, 0xd5, 0x54, 0xb4, 0xbb,
	0xc9, 0x85, 0x8c, 0xbe, 0x08, 0xf8, 0x0f, 0x95, 0x5c, 0xcf, 0x9f, 0x77, 0xf2, 0x85, 0x52, 0xe8,
	0xca, 0x7c, 0xc3, 0x43, 0x32, 0x21, 0x53, 0x9f, 0xa1, 0x4d, 0xcf, 0xc1, 0x11, 0x65, 0xe8, 0x4c,
	0xc8, 0xb4, 0xc7, 0x1c, 0x51, 0xd2, 0x6b, 0xf0, 0xeb, 0x7c, 0x9f, 0x2d, 0x3f, 0x1a, 0xae, 0xc2,
	0xb3, 0x09, 0x99, 0x0e, 0x98, 0x57, 0xe7, 0xfb, 0x44, 0xfb, 0xf4, 0x16, 0x82, 0x42, 0x57, 0xca,
	0x84, 0x2c, 0xf9, 0x7b, 0xd8, 0xc5, 0x2c, 0x40, 0x29, 0xd5, 0x0a, 0x1d, 0x83, 0x5b, 0xad, 0x56,
	0x8a, 0x37, 0x61, 0x0f, 0x63, 0xd6, 0xd3, 0x9d, 0x95, 0xf8, 0xe4, 0xa1, 0x8b, 0x2a, 0xda, 0xd1,
	0x37, 0x81, 0x40, 0xb3, 0x31, 0xfe, 0xb6, 0xe3, 0xaa, 0xa1, 0x37, 0x00, 0xb5, 0x31, 0x33, 0x51,
	0x5a, 0x46, 0xdf, 0x2a, 0x29, 0x82, 0x15, 0xaf, 0x82, 0x4b, 0x8c, 0x3a, 0x18, 0xf5, 0x8c, 0x90,
	0x96, 0xf4, 0x1e, 0xa8, 0x9e, 0x3d, 0x33, 0x74, 0x36, 0x09, 0xf1, 0x83, 0xd9, 0x38, 0xd6, 0xab,
	0x88, 0x8f, 0x6b, 0xb0, 0xfd, 0x16, 0x1d, 0x36, 0x54, 0x2d, 0x2d, 0xf1, 0xa1, 0x6f, 0x93, 0xa3,
	0x1f, 0x02, 0x03, 0x83, 0xa7, 0xb6, 0x95, 0x54, 0xfc, 0x5f, 0x7c, 0x29, 0x5c, 0xfe, 0xe1, 0x33,
	0x25, 0x2d, 0xe0, 0xd5, 0x09, 0xa0, 0x09, 0x2f, 0x3a, 0x6c, 0xa4, 0xda, 0x62, 0x02, 0xe0, 0x1d,
	0xf2, 0xa3, 0x39, 0x0c, 0xdb, 0x63, 0xd9, 0x83, 0x92, 0xe3, 0x41, 0x5b, 0x37, 0x73, 0xda, 0x37,
	0x8b, 0x9e, 0x60, 0x74, 0xd2, 0x9a, 0x86, 0xd0, 0x57, 0xbb, 0xa2, 0xe0, 0x4a, 0x61, 0x29, 0x8f,
	0x1d, 0x5c, 0x1a, 0x9b, 0x37, 0x33, 0xa3, 0x60, 0xb9, 0x60, 0x76, 0xd1, 0x9e, 0xc0, 0x3f, 0x72,
	0x27, 0xee, 0x23, 0xbe, 0xe2, 0xd2, 0xc5, 0xbf, 0xbc, 0xfb, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x80,
	0x0c, 0xf9, 0xc7, 0xa5, 0x02, 0x00, 0x00,
}