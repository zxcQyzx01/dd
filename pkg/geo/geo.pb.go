// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v6.30.0--rc2
// source: proto/geo.proto

package geo

import (
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

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	City   string `protobuf:"bytes,1,opt,name=city,proto3" json:"city,omitempty"`
	Street string `protobuf:"bytes,2,opt,name=street,proto3" json:"street,omitempty"`
	House  string `protobuf:"bytes,3,opt,name=house,proto3" json:"house,omitempty"`
	Lat    string `protobuf:"bytes,4,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon    string `protobuf:"bytes,5,opt,name=lon,proto3" json:"lon,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_geo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_proto_geo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_proto_geo_proto_rawDescGZIP(), []int{0}
}

func (x *Address) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *Address) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

func (x *Address) GetHouse() string {
	if x != nil {
		return x.House
	}
	return ""
}

func (x *Address) GetLat() string {
	if x != nil {
		return x.Lat
	}
	return ""
}

func (x *Address) GetLon() string {
	if x != nil {
		return x.Lon
	}
	return ""
}

type SearchAddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *SearchAddressRequest) Reset() {
	*x = SearchAddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_geo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchAddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchAddressRequest) ProtoMessage() {}

func (x *SearchAddressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_geo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchAddressRequest.ProtoReflect.Descriptor instead.
func (*SearchAddressRequest) Descriptor() ([]byte, []int) {
	return file_proto_geo_proto_rawDescGZIP(), []int{1}
}

func (x *SearchAddressRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type SearchAddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses []*Address `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *SearchAddressResponse) Reset() {
	*x = SearchAddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_geo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchAddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchAddressResponse) ProtoMessage() {}

func (x *SearchAddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_geo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchAddressResponse.ProtoReflect.Descriptor instead.
func (*SearchAddressResponse) Descriptor() ([]byte, []int) {
	return file_proto_geo_proto_rawDescGZIP(), []int{2}
}

func (x *SearchAddressResponse) GetAddresses() []*Address {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type GeocodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *GeocodeRequest) Reset() {
	*x = GeocodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_geo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeocodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeocodeRequest) ProtoMessage() {}

func (x *GeocodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_geo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeocodeRequest.ProtoReflect.Descriptor instead.
func (*GeocodeRequest) Descriptor() ([]byte, []int) {
	return file_proto_geo_proto_rawDescGZIP(), []int{3}
}

func (x *GeocodeRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type GeocodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses []*Address `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *GeocodeResponse) Reset() {
	*x = GeocodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_geo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeocodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeocodeResponse) ProtoMessage() {}

func (x *GeocodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_geo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeocodeResponse.ProtoReflect.Descriptor instead.
func (*GeocodeResponse) Descriptor() ([]byte, []int) {
	return file_proto_geo_proto_rawDescGZIP(), []int{4}
}

func (x *GeocodeResponse) GetAddresses() []*Address {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_proto_geo_proto protoreflect.FileDescriptor

var file_proto_geo_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x67, 0x65, 0x6f, 0x22, 0x6f, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x63, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x68, 0x6f,
	0x75, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6c, 0x6f, 0x6e, 0x22, 0x2c, 0x0a, 0x14, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x22, 0x43, 0x0a, 0x15, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a,
	0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x65, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52,
	0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0x2a, 0x0a, 0x0e, 0x47, 0x65,
	0x6f, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x3d, 0x0a, 0x0f, 0x47, 0x65, 0x6f, 0x63, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x09, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67,
	0x65, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x65, 0x73, 0x32, 0x8a, 0x01, 0x0a, 0x0a, 0x47, 0x65, 0x6f, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x19, 0x2e, 0x67, 0x65, 0x6f, 0x2e, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x67, 0x65, 0x6f, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x07,
	0x47, 0x65, 0x6f, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x13, 0x2e, 0x67, 0x65, 0x6f, 0x2e, 0x47, 0x65,
	0x6f, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67,
	0x65, 0x6f, 0x2e, 0x47, 0x65, 0x6f, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x64, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_geo_proto_rawDescOnce sync.Once
	file_proto_geo_proto_rawDescData = file_proto_geo_proto_rawDesc
)

func file_proto_geo_proto_rawDescGZIP() []byte {
	file_proto_geo_proto_rawDescOnce.Do(func() {
		file_proto_geo_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_geo_proto_rawDescData)
	})
	return file_proto_geo_proto_rawDescData
}

var file_proto_geo_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_geo_proto_goTypes = []interface{}{
	(*Address)(nil),               // 0: geo.Address
	(*SearchAddressRequest)(nil),  // 1: geo.SearchAddressRequest
	(*SearchAddressResponse)(nil), // 2: geo.SearchAddressResponse
	(*GeocodeRequest)(nil),        // 3: geo.GeocodeRequest
	(*GeocodeResponse)(nil),       // 4: geo.GeocodeResponse
}
var file_proto_geo_proto_depIdxs = []int32{
	0, // 0: geo.SearchAddressResponse.addresses:type_name -> geo.Address
	0, // 1: geo.GeocodeResponse.addresses:type_name -> geo.Address
	1, // 2: geo.GeoService.SearchAddress:input_type -> geo.SearchAddressRequest
	3, // 3: geo.GeoService.Geocode:input_type -> geo.GeocodeRequest
	2, // 4: geo.GeoService.SearchAddress:output_type -> geo.SearchAddressResponse
	4, // 5: geo.GeoService.Geocode:output_type -> geo.GeocodeResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_geo_proto_init() }
func file_proto_geo_proto_init() {
	if File_proto_geo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_geo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_proto_geo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchAddressRequest); i {
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
		file_proto_geo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchAddressResponse); i {
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
		file_proto_geo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeocodeRequest); i {
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
		file_proto_geo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeocodeResponse); i {
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
			RawDescriptor: file_proto_geo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_geo_proto_goTypes,
		DependencyIndexes: file_proto_geo_proto_depIdxs,
		MessageInfos:      file_proto_geo_proto_msgTypes,
	}.Build()
	File_proto_geo_proto = out.File
	file_proto_geo_proto_rawDesc = nil
	file_proto_geo_proto_goTypes = nil
	file_proto_geo_proto_depIdxs = nil
}
