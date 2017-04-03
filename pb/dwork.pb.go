// Code generated by protoc-gen-go.
// source: pb/dwork.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	pb/dwork.proto

It has these top-level messages:
	Results
	WorkUnit
	ResultsSuccess
	WorkerID
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type Results struct {
	JobID    string `protobuf:"bytes,1,opt,name=JobID" json:"JobID,omitempty"`
	Workid   int64  `protobuf:"varint,2,opt,name=workid" json:"workid,omitempty"`
	Found    bool   `protobuf:"varint,3,opt,name=found" json:"found,omitempty"`
	Location int64  `protobuf:"varint,4,opt,name=location" json:"location,omitempty"`
}

func (m *Results) Reset()                    { *m = Results{} }
func (m *Results) String() string            { return proto.CompactTextString(m) }
func (*Results) ProtoMessage()               {}
func (*Results) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Results) GetJobID() string {
	if m != nil {
		return m.JobID
	}
	return ""
}

func (m *Results) GetWorkid() int64 {
	if m != nil {
		return m.Workid
	}
	return 0
}

func (m *Results) GetFound() bool {
	if m != nil {
		return m.Found
	}
	return false
}

func (m *Results) GetLocation() int64 {
	if m != nil {
		return m.Location
	}
	return 0
}

// The response message containing the greetings
// https://developers.google.com/protocol-buffers/docs/proto3
type WorkUnit struct {
	JobID     string `protobuf:"bytes,1,opt,name=JobID" json:"JobID,omitempty"`
	Id        int64  `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	Offset    int64  `protobuf:"varint,3,opt,name=offset" json:"offset,omitempty"`
	Size      int64  `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
	Status    int64  `protobuf:"varint,5,opt,name=status" json:"status,omitempty"`
	Timestamp int64  `protobuf:"varint,6,opt,name=timestamp" json:"timestamp,omitempty"`
	Code      string `protobuf:"bytes,7,opt,name=code" json:"code,omitempty"`
}

func (m *WorkUnit) Reset()                    { *m = WorkUnit{} }
func (m *WorkUnit) String() string            { return proto.CompactTextString(m) }
func (*WorkUnit) ProtoMessage()               {}
func (*WorkUnit) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *WorkUnit) GetJobID() string {
	if m != nil {
		return m.JobID
	}
	return ""
}

func (m *WorkUnit) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *WorkUnit) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *WorkUnit) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *WorkUnit) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *WorkUnit) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *WorkUnit) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type ResultsSuccess struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *ResultsSuccess) Reset()                    { *m = ResultsSuccess{} }
func (m *ResultsSuccess) String() string            { return proto.CompactTextString(m) }
func (*ResultsSuccess) ProtoMessage()               {}
func (*ResultsSuccess) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ResultsSuccess) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type WorkerID struct {
	UUID string `protobuf:"bytes,1,opt,name=UUID" json:"UUID,omitempty"`
}

func (m *WorkerID) Reset()                    { *m = WorkerID{} }
func (m *WorkerID) String() string            { return proto.CompactTextString(m) }
func (*WorkerID) ProtoMessage()               {}
func (*WorkerID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *WorkerID) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func init() {
	proto.RegisterType((*Results)(nil), "pb.Results")
	proto.RegisterType((*WorkUnit)(nil), "pb.WorkUnit")
	proto.RegisterType((*ResultsSuccess)(nil), "pb.ResultsSuccess")
	proto.RegisterType((*WorkerID)(nil), "pb.WorkerID")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Dwork service

type DworkClient interface {
	// Sends a greeting
	GiveWork(ctx context.Context, in *WorkerID, opts ...grpc.CallOption) (*WorkUnit, error)
	ReceiveResults(ctx context.Context, in *Results, opts ...grpc.CallOption) (*ResultsSuccess, error)
}

type dworkClient struct {
	cc *grpc.ClientConn
}

func NewDworkClient(cc *grpc.ClientConn) DworkClient {
	return &dworkClient{cc}
}

func (c *dworkClient) GiveWork(ctx context.Context, in *WorkerID, opts ...grpc.CallOption) (*WorkUnit, error) {
	out := new(WorkUnit)
	err := grpc.Invoke(ctx, "/pb.Dwork/GiveWork", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dworkClient) ReceiveResults(ctx context.Context, in *Results, opts ...grpc.CallOption) (*ResultsSuccess, error) {
	out := new(ResultsSuccess)
	err := grpc.Invoke(ctx, "/pb.Dwork/receiveResults", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Dwork service

type DworkServer interface {
	// Sends a greeting
	GiveWork(context.Context, *WorkerID) (*WorkUnit, error)
	ReceiveResults(context.Context, *Results) (*ResultsSuccess, error)
}

func RegisterDworkServer(s *grpc.Server, srv DworkServer) {
	s.RegisterService(&_Dwork_serviceDesc, srv)
}

func _Dwork_GiveWork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkerID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DworkServer).GiveWork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Dwork/GiveWork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DworkServer).GiveWork(ctx, req.(*WorkerID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dwork_ReceiveResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Results)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DworkServer).ReceiveResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Dwork/ReceiveResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DworkServer).ReceiveResults(ctx, req.(*Results))
	}
	return interceptor(ctx, in, info, handler)
}

var _Dwork_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Dwork",
	HandlerType: (*DworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GiveWork",
			Handler:    _Dwork_GiveWork_Handler,
		},
		{
			MethodName: "receiveResults",
			Handler:    _Dwork_ReceiveResults_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/dwork.proto",
}

func init() { proto.RegisterFile("pb/dwork.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x91, 0xcb, 0x4a, 0xc3, 0x40,
	0x14, 0x86, 0x9b, 0xf4, 0x96, 0x1e, 0x25, 0x8b, 0x83, 0xc8, 0x50, 0x44, 0xc2, 0xac, 0x82, 0x8b,
	0x08, 0xf6, 0x15, 0x02, 0x52, 0x97, 0x23, 0xc1, 0x75, 0x2e, 0x13, 0x18, 0xda, 0x66, 0x42, 0x66,
	0x52, 0xc1, 0x07, 0xf2, 0x39, 0x65, 0x2e, 0x89, 0x6e, 0xdc, 0xfd, 0x5f, 0xf2, 0xcf, 0xf9, 0xcf,
	0x05, 0xe2, 0xbe, 0x7a, 0x6e, 0x3e, 0xe5, 0x70, 0xca, 0xfa, 0x41, 0x6a, 0x89, 0x61, 0x5f, 0x51,
	0x01, 0x5b, 0xc6, 0xd5, 0x78, 0xd6, 0x0a, 0xef, 0x60, 0xfd, 0x26, 0xab, 0x63, 0x4e, 0x82, 0x24,
	0x48, 0x77, 0xcc, 0x01, 0xde, 0xc3, 0xc6, 0x3c, 0x11, 0x0d, 0x09, 0x93, 0x20, 0x5d, 0x32, 0x4f,
	0xc6, 0xdd, 0xca, 0xb1, 0x6b, 0xc8, 0x32, 0x09, 0xd2, 0x88, 0x39, 0xc0, 0x3d, 0x44, 0x67, 0x59,
	0x97, 0x5a, 0xc8, 0x8e, 0xac, 0xac, 0x7f, 0x66, 0xfa, 0x1d, 0x40, 0xf4, 0x21, 0x87, 0x53, 0xd1,
	0x09, 0xfd, 0x4f, 0x58, 0x0c, 0xe1, 0x1c, 0x14, 0x8a, 0xc6, 0x84, 0xcb, 0xb6, 0x55, 0x5c, 0xdb,
	0x94, 0x25, 0xf3, 0x84, 0x08, 0x2b, 0x25, 0xbe, 0xb8, 0x8f, 0xb0, 0xda, 0x78, 0x95, 0x2e, 0xf5,
	0xa8, 0xc8, 0xda, 0x79, 0x1d, 0xe1, 0x03, 0xec, 0xb4, 0xb8, 0x70, 0xa5, 0xcb, 0x4b, 0x4f, 0x36,
	0xf6, 0xd7, 0xef, 0x07, 0x53, 0xa9, 0x96, 0x0d, 0x27, 0x5b, 0xdb, 0x86, 0xd5, 0xf4, 0x09, 0x62,
	0xbf, 0x93, 0xf7, 0xb1, 0xae, 0xb9, 0x52, 0x48, 0x60, 0xab, 0x9c, 0xb4, 0xfd, 0x46, 0x6c, 0x42,
	0xfa, 0xe8, 0x66, 0xe2, 0xc3, 0x31, 0x37, 0xb5, 0x8a, 0x62, 0x1e, 0xc9, 0xea, 0x97, 0x16, 0xd6,
	0xb9, 0xd9, 0x18, 0xa6, 0x10, 0xbd, 0x8a, 0x2b, 0x37, 0x66, 0xbc, 0xcd, 0xfa, 0x2a, 0x9b, 0x9e,
	0xed, 0x67, 0x32, 0x8b, 0xa1, 0x0b, 0x3c, 0x40, 0x3c, 0xf0, 0x9a, 0x8b, 0x2b, 0x9f, 0x2e, 0x73,
	0x63, 0x1c, 0x1e, 0xf6, 0xf8, 0x07, 0x7c, 0x7f, 0x74, 0x51, 0x6d, 0xec, 0x49, 0x0f, 0x3f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xf8, 0x46, 0xb7, 0xed, 0xe4, 0x01, 0x00, 0x00,
}
