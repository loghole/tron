// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: strings.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type String struct {
	// string
	Str                  string   `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_0af3cde72035e609, []int{0}
}
func (m *String) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_String.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(m, src)
}
func (m *String) XXX_Size() int {
	return m.Size()
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetStr() string {
	if m != nil {
		return m.Str
	}
	return ""
}

func init() {
	proto.RegisterType((*String)(nil), "v1.String")
}

func init() { proto.RegisterFile("strings.proto", fileDescriptor_0af3cde72035e609) }

var fileDescriptor_0af3cde72035e609 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x2e, 0x29, 0xca,
	0xcc, 0x4b, 0x2f, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x33, 0x94, 0x92, 0x49,
	0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x49,
	0x2c, 0xc9, 0xcc, 0xcf, 0x83, 0xaa, 0x50, 0x92, 0xe2, 0x62, 0x0b, 0x06, 0x6b, 0x11, 0x12, 0xe0,
	0x62, 0x2e, 0x2e, 0x29, 0x92, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x8d, 0x26, 0x31,
	0x72, 0xb1, 0x43, 0x24, 0x8b, 0x85, 0xdc, 0xb8, 0xd8, 0x43, 0xf2, 0x43, 0x0b, 0x0a, 0x52, 0x8b,
	0x84, 0xb8, 0xf4, 0xca, 0x0c, 0xf5, 0x20, 0xe2, 0x52, 0x48, 0x6c, 0x25, 0xe5, 0xa6, 0xcb, 0x4f,
	0x26, 0x33, 0xc9, 0x0a, 0x49, 0x83, 0xad, 0x2a, 0x33, 0xd4, 0x87, 0xba, 0x45, 0xbf, 0x14, 0xa4,
	0x4d, 0xbf, 0xba, 0xb8, 0xa4, 0xa8, 0x56, 0xc8, 0x9e, 0x8b, 0xdd, 0x3d, 0xb5, 0xc4, 0x33, 0x2f,
	0x2d, 0x1f, 0xa7, 0x39, 0x32, 0x60, 0x73, 0xc4, 0x84, 0x44, 0xd0, 0xcd, 0xc9, 0xcc, 0x4b, 0xcb,
	0x77, 0x12, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x67,
	0x3c, 0x96, 0x63, 0x48, 0x62, 0x03, 0xfb, 0xc4, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x12,
	0x81, 0xa5, 0xfc, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StringsClient is the client API for Strings service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StringsClient interface {
	// Method to upper
	ToUpper(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	GetInfo(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
}

type stringsClient struct {
	cc *grpc.ClientConn
}

func NewStringsClient(cc *grpc.ClientConn) StringsClient {
	return &stringsClient{cc}
}

func (c *stringsClient) ToUpper(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/v1.Strings/ToUpper", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stringsClient) GetInfo(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/v1.Strings/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StringsServer is the server API for Strings service.
type StringsServer interface {
	// Method to upper
	ToUpper(context.Context, *String) (*String, error)
	GetInfo(context.Context, *String) (*String, error)
}

// UnimplementedStringsServer can be embedded to have forward compatible implementations.
type UnimplementedStringsServer struct {
}

func (*UnimplementedStringsServer) ToUpper(ctx context.Context, req *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToUpper not implemented")
}
func (*UnimplementedStringsServer) GetInfo(ctx context.Context, req *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}

func RegisterStringsServer(s *grpc.Server, srv StringsServer) {
	s.RegisterService(&_Strings_serviceDesc, srv)
}

func _Strings_ToUpper_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StringsServer).ToUpper(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.Strings/ToUpper",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StringsServer).ToUpper(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _Strings_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StringsServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.Strings/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StringsServer).GetInfo(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

var _Strings_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.Strings",
	HandlerType: (*StringsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ToUpper",
			Handler:    _Strings_ToUpper_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _Strings_GetInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "strings.proto",
}

func (m *String) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *String) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *String) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Str) > 0 {
		i -= len(m.Str)
		copy(dAtA[i:], m.Str)
		i = encodeVarintStrings(dAtA, i, uint64(len(m.Str)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintStrings(dAtA []byte, offset int, v uint64) int {
	offset -= sovStrings(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *String) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Str)
	if l > 0 {
		n += 1 + l + sovStrings(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovStrings(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStrings(x uint64) (n int) {
	return sovStrings(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *String) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStrings
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: String: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: String: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Str", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStrings
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStrings
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStrings
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Str = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStrings(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStrings
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthStrings
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipStrings(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStrings
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStrings
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStrings
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthStrings
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStrings
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStrings
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStrings        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStrings          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStrings = fmt.Errorf("proto: unexpected end of group")
)
