// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: appointments.proto

// protoc --go_out=plugins=grpc:. *.proto

package appointmentsService

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Appointment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppointmentID string               `protobuf:"bytes,1,opt,name=AppointmentID,proto3" json:"AppointmentID,omitempty"`
	DoctorID      string               `protobuf:"bytes,2,opt,name=DoctorID,proto3" json:"DoctorID,omitempty"`
	PatientID     string               `protobuf:"bytes,3,opt,name=PatientID,proto3" json:"PatientID,omitempty"`
	Reason        string               `protobuf:"bytes,4,opt,name=Reason,proto3" json:"Reason,omitempty"`
	From          *timestamp.Timestamp `protobuf:"bytes,5,opt,name=From,proto3" json:"From,omitempty"`
	To            *timestamp.Timestamp `protobuf:"bytes,6,opt,name=To,proto3" json:"To,omitempty"`
}

func (x *Appointment) Reset() {
	*x = Appointment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appointments_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Appointment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Appointment) ProtoMessage() {}

func (x *Appointment) ProtoReflect() protoreflect.Message {
	mi := &file_appointments_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Appointment.ProtoReflect.Descriptor instead.
func (*Appointment) Descriptor() ([]byte, []int) {
	return file_appointments_proto_rawDescGZIP(), []int{0}
}

func (x *Appointment) GetAppointmentID() string {
	if x != nil {
		return x.AppointmentID
	}
	return ""
}

func (x *Appointment) GetDoctorID() string {
	if x != nil {
		return x.DoctorID
	}
	return ""
}

func (x *Appointment) GetPatientID() string {
	if x != nil {
		return x.PatientID
	}
	return ""
}

func (x *Appointment) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *Appointment) GetFrom() *timestamp.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Appointment) GetTo() *timestamp.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appointments_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_appointments_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_appointments_proto_rawDescGZIP(), []int{1}
}

type GetByDoctrorIDReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DoctorID string               `protobuf:"bytes,1,opt,name=DoctorID,proto3" json:"DoctorID,omitempty"`
	From     *timestamp.Timestamp `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	Till     *timestamp.Timestamp `protobuf:"bytes,3,opt,name=Till,proto3" json:"Till,omitempty"`
}

func (x *GetByDoctrorIDReq) Reset() {
	*x = GetByDoctrorIDReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appointments_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByDoctrorIDReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByDoctrorIDReq) ProtoMessage() {}

func (x *GetByDoctrorIDReq) ProtoReflect() protoreflect.Message {
	mi := &file_appointments_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByDoctrorIDReq.ProtoReflect.Descriptor instead.
func (*GetByDoctrorIDReq) Descriptor() ([]byte, []int) {
	return file_appointments_proto_rawDescGZIP(), []int{2}
}

func (x *GetByDoctrorIDReq) GetDoctorID() string {
	if x != nil {
		return x.DoctorID
	}
	return ""
}

func (x *GetByDoctrorIDReq) GetFrom() *timestamp.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *GetByDoctrorIDReq) GetTill() *timestamp.Timestamp {
	if x != nil {
		return x.Till
	}
	return nil
}

type GetByDoctorIDRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Appointments []*Appointment `protobuf:"bytes,1,rep,name=Appointments,proto3" json:"Appointments,omitempty"`
}

func (x *GetByDoctorIDRes) Reset() {
	*x = GetByDoctorIDRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appointments_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByDoctorIDRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByDoctorIDRes) ProtoMessage() {}

func (x *GetByDoctorIDRes) ProtoReflect() protoreflect.Message {
	mi := &file_appointments_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByDoctorIDRes.ProtoReflect.Descriptor instead.
func (*GetByDoctorIDRes) Descriptor() ([]byte, []int) {
	return file_appointments_proto_rawDescGZIP(), []int{3}
}

func (x *GetByDoctorIDRes) GetAppointments() []*Appointment {
	if x != nil {
		return x.Appointments
	}
	return nil
}

var File_appointments_proto protoreflect.FileDescriptor

var file_appointments_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe1, 0x01, 0x0a, 0x0b, 0x41,
	0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x41, 0x70,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x41, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x1a, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09,
	0x50, 0x61, 0x74, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x50, 0x61, 0x74, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x46, 0x72,
	0x6f, 0x6d, 0x12, 0x2a, 0x0a, 0x02, 0x54, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x54, 0x6f, 0x22, 0x07,
	0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x8f, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42,
	0x79, 0x44, 0x6f, 0x63, 0x74, 0x72, 0x6f, 0x72, 0x49, 0x44, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a,
	0x08, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x12, 0x2e, 0x0a, 0x04, 0x46, 0x72, 0x6f,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x2e, 0x0a, 0x04, 0x54, 0x69, 0x6c,
	0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69, 0x6c, 0x6c, 0x22, 0x58, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x42, 0x79, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x52, 0x65, 0x73, 0x12, 0x44, 0x0a,
	0x0c, 0x41, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0c, 0x41, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x32, 0x76, 0x0a, 0x12, 0x41, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65,
	0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x60, 0x0a, 0x0d, 0x47, 0x65, 0x74,
	0x42, 0x79, 0x44, 0x6f, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x12, 0x26, 0x2e, 0x61, 0x70, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x44, 0x6f, 0x63, 0x74, 0x72, 0x6f, 0x72, 0x49, 0x44, 0x52,
	0x65, 0x71, 0x1a, 0x25, 0x2e, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x44, 0x6f,
	0x63, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x17, 0x5a, 0x15, 0x2e,
	0x3b, 0x61, 0x70, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_appointments_proto_rawDescOnce sync.Once
	file_appointments_proto_rawDescData = file_appointments_proto_rawDesc
)

func file_appointments_proto_rawDescGZIP() []byte {
	file_appointments_proto_rawDescOnce.Do(func() {
		file_appointments_proto_rawDescData = protoimpl.X.CompressGZIP(file_appointments_proto_rawDescData)
	})
	return file_appointments_proto_rawDescData
}

var file_appointments_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_appointments_proto_goTypes = []interface{}{
	(*Appointment)(nil),         // 0: appointmentsService.Appointment
	(*Empty)(nil),               // 1: appointmentsService.Empty
	(*GetByDoctrorIDReq)(nil),   // 2: appointmentsService.GetByDoctrorIDReq
	(*GetByDoctorIDRes)(nil),    // 3: appointmentsService.GetByDoctorIDRes
	(*timestamp.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_appointments_proto_depIdxs = []int32{
	4, // 0: appointmentsService.Appointment.From:type_name -> google.protobuf.Timestamp
	4, // 1: appointmentsService.Appointment.To:type_name -> google.protobuf.Timestamp
	4, // 2: appointmentsService.GetByDoctrorIDReq.From:type_name -> google.protobuf.Timestamp
	4, // 3: appointmentsService.GetByDoctrorIDReq.Till:type_name -> google.protobuf.Timestamp
	0, // 4: appointmentsService.GetByDoctorIDRes.Appointments:type_name -> appointmentsService.Appointment
	2, // 5: appointmentsService.AppointmentService.GetByDoctorID:input_type -> appointmentsService.GetByDoctrorIDReq
	3, // 6: appointmentsService.AppointmentService.GetByDoctorID:output_type -> appointmentsService.GetByDoctorIDRes
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_appointments_proto_init() }
func file_appointments_proto_init() {
	if File_appointments_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_appointments_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Appointment); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_appointments_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_appointments_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByDoctrorIDReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_appointments_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByDoctorIDRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_appointments_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_appointments_proto_goTypes,
		DependencyIndexes: file_appointments_proto_depIdxs,
		MessageInfos:      file_appointments_proto_msgTypes,
	}.Build()
	File_appointments_proto = out.File
	file_appointments_proto_rawDesc = nil
	file_appointments_proto_goTypes = nil
	file_appointments_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AppointmentServiceClient is the client API for AppointmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AppointmentServiceClient interface {
	GetByDoctorID(ctx context.Context, in *GetByDoctrorIDReq, opts ...grpc.CallOption) (*GetByDoctorIDRes, error)
}

type appointmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppointmentServiceClient(cc grpc.ClientConnInterface) AppointmentServiceClient {
	return &appointmentServiceClient{cc}
}

func (c *appointmentServiceClient) GetByDoctorID(ctx context.Context, in *GetByDoctrorIDReq, opts ...grpc.CallOption) (*GetByDoctorIDRes, error) {
	out := new(GetByDoctorIDRes)
	err := c.cc.Invoke(ctx, "/appointmentsService.AppointmentService/GetByDoctorID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppointmentServiceServer is the server API for AppointmentService service.
type AppointmentServiceServer interface {
	GetByDoctorID(context.Context, *GetByDoctrorIDReq) (*GetByDoctorIDRes, error)
}

// UnimplementedAppointmentServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAppointmentServiceServer struct {
}

func (*UnimplementedAppointmentServiceServer) GetByDoctorID(context.Context, *GetByDoctrorIDReq) (*GetByDoctorIDRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByDoctorID not implemented")
}

func RegisterAppointmentServiceServer(s *grpc.Server, srv AppointmentServiceServer) {
	s.RegisterService(&_AppointmentService_serviceDesc, srv)
}

func _AppointmentService_GetByDoctorID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByDoctrorIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppointmentServiceServer).GetByDoctorID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appointmentsService.AppointmentService/GetByDoctorID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppointmentServiceServer).GetByDoctorID(ctx, req.(*GetByDoctrorIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _AppointmentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "appointmentsService.AppointmentService",
	HandlerType: (*AppointmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByDoctorID",
			Handler:    _AppointmentService_GetByDoctorID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "appointments.proto",
}
