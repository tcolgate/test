// Code generated by protoc-gen-go.
// source: ping.proto
// DO NOT EDIT!

/*
Package ping is a generated protocol buffer package.

It is generated from these files:
	ping.proto

It has these top-level messages:
	PingRequest
	PingReply
*/
package ping

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type PingRequest struct {
	TimeOut int64 `protobuf:"varint,1,opt,name=timeOut" json:"timeOut,omitempty"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}

type PingReply struct {
	TimeBack int64 `protobuf:"varint,1,opt,name=timeBack" json:"timeBack,omitempty"`
}

func (m *PingReply) Reset()         { *m = PingReply{} }
func (m *PingReply) String() string { return proto.CompactTextString(m) }
func (*PingReply) ProtoMessage()    {}

// Client API for Ping service

type PingClient interface {
	// Sends a greeting
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error)
}

type pingClient struct {
	cc *grpc.ClientConn
}

func NewPingClient(cc *grpc.ClientConn) PingClient {
	return &pingClient{cc}
}

func (c *pingClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := grpc.Invoke(ctx, "/ping.Ping/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Ping service

type PingServer interface {
	// Sends a greeting
	Ping(context.Context, *PingRequest) (*PingReply, error)
}

func RegisterPingServer(s *grpc.Server, srv PingServer) {
	s.RegisterService(&_Ping_serviceDesc, srv)
}

func _Ping_Ping_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(PingRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(PingServer).Ping(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Ping_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ping.Ping",
	HandlerType: (*PingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Ping_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
