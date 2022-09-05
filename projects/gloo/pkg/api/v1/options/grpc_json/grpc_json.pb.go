// mostly copied from https://github.com/envoyproxy/envoy/blob/374dca7905fc048be74169a7655d0462606555ad/api/envoy/extensions/filters/http/grpc_json_transcoder/v3/transcoder.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.4
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/grpc_json/grpc_json.proto

package grpc_json

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// [#next-free-field: 10]
type GrpcJsonTranscoder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to DescriptorSet:
	//	*GrpcJsonTranscoder_ProtoDescriptor
	//	*GrpcJsonTranscoder_ProtoDescriptorBin
	DescriptorSet isGrpcJsonTranscoder_DescriptorSet `protobuf_oneof:"descriptor_set"`
	// A list of strings that
	// supplies the fully qualified service names (i.e. "package_name.service_name") that
	// the transcoder will translate. If the service name doesn't exist in ``proto_descriptor``,
	// Envoy will fail at startup. The ``proto_descriptor`` may contain more services than
	// the service names specified here, but they won't be translated.
	Services []string `protobuf:"bytes,2,rep,name=services,proto3" json:"services,omitempty"`
	// Control options for response JSON. These options are passed directly to
	// `JsonPrintOptions <https://developers.google.com/protocol-buffers/docs/reference/cpp/
	// google.protobuf.util.json_util#JsonPrintOptions>`_.
	PrintOptions *GrpcJsonTranscoder_PrintOptions `protobuf:"bytes,3,opt,name=print_options,json=printOptions,proto3" json:"print_options,omitempty"`
	// Set this value to true to keep the incoming request route after the outgoing headers are transformed to match the upstream gRPC service.
	// Note that you cannot set this value to true with routes for gRPC services that are not transcoded.
	// When set to false, Envoy does not match against the incoming request path.
	// For more information, see the Envoy docs <https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter#route-configs-for-transcoded-requests>.
	MatchIncomingRequestRoute bool `protobuf:"varint,5,opt,name=match_incoming_request_route,json=matchIncomingRequestRoute,proto3" json:"match_incoming_request_route,omitempty"`
	// A list of query parameters to be ignored for transcoding method mapping.
	// By default, the transcoder filter will not transcode a request if there are any
	// unknown/invalid query parameters.
	//
	// Example :
	//
	// .. code-block:: proto
	//
	//     service Bookstore {
	//       rpc GetShelf(GetShelfRequest) returns (Shelf) {
	//         option (google.api.http) = {
	//           get: "/shelves/{shelf}"
	//         };
	//       }
	//     }
	//
	//     message GetShelfRequest {
	//       int64 shelf = 1;
	//     }
	//
	//     message Shelf {}
	//
	// The request ``/shelves/100?foo=bar`` will not be mapped to ``GetShelf``` because variable
	// binding for ``foo`` is not defined. Adding ``foo`` to ``ignored_query_parameters`` will allow
	// the same request to be mapped to ``GetShelf``.
	IgnoredQueryParameters []string `protobuf:"bytes,6,rep,name=ignored_query_parameters,json=ignoredQueryParameters,proto3" json:"ignored_query_parameters,omitempty"`
	// Whether to route methods without the ``google.api.http`` option.
	//
	// Example :
	//
	// .. code-block:: proto
	//
	//     package bookstore;
	//
	//     service Bookstore {
	//       rpc GetShelf(GetShelfRequest) returns (Shelf) {}
	//     }
	//
	//     message GetShelfRequest {
	//       int64 shelf = 1;
	//     }
	//
	//     message Shelf {}
	//
	// The client could ``post`` a json body ``{"shelf": 1234}`` with the path of
	// ``/bookstore.Bookstore/GetShelfRequest`` to call ``GetShelfRequest``.
	AutoMapping bool `protobuf:"varint,7,opt,name=auto_mapping,json=autoMapping,proto3" json:"auto_mapping,omitempty"`
	// Whether to ignore query parameters that cannot be mapped to a corresponding
	// protobuf field. Use this if you cannot control the query parameters and do
	// not know them beforehand. Otherwise use ``ignored_query_parameters``.
	// Defaults to false.
	IgnoreUnknownQueryParameters bool `protobuf:"varint,8,opt,name=ignore_unknown_query_parameters,json=ignoreUnknownQueryParameters,proto3" json:"ignore_unknown_query_parameters,omitempty"`
	// Whether to convert gRPC status headers to JSON.
	// When trailer indicates a gRPC error and there was no HTTP body, take ``google.rpc.Status``
	// from the ``grpc-status-details-bin`` header and use it as JSON body.
	// If there was no such header, make ``google.rpc.Status`` out of the ``grpc-status`` and
	// ``grpc-message`` headers.
	// The error details types must be present in the ``proto_descriptor``.
	//
	// For example, if an upstream server replies with headers:
	//
	// .. code-block:: none
	//
	//     grpc-status: 5
	//     grpc-status-details-bin:
	//         CAUaMwoqdHlwZS5nb29nbGVhcGlzLmNvbS9nb29nbGUucnBjLlJlcXVlc3RJbmZvEgUKA3ItMQ
	//
	// The ``grpc-status-details-bin`` header contains a base64-encoded protobuf message
	// ``google.rpc.Status``. It will be transcoded into:
	//
	// .. code-block:: none
	//
	//     HTTP/1.1 404 Not Found
	//     content-type: application/json
	//
	//     {"code":5,"details":[{"@type":"type.googleapis.com/google.rpc.RequestInfo","requestId":"r-1"}]}
	//
	//  In order to transcode the message, the ``google.rpc.RequestInfo`` type from
	//  the ``google/rpc/error_details.proto`` should be included in the configured
	//  :ref:`proto descriptor set <config_grpc_json_generate_proto_descriptor_set>`.
	ConvertGrpcStatus bool `protobuf:"varint,9,opt,name=convert_grpc_status,json=convertGrpcStatus,proto3" json:"convert_grpc_status,omitempty"`
}

func (x *GrpcJsonTranscoder) Reset() {
	*x = GrpcJsonTranscoder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcJsonTranscoder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcJsonTranscoder) ProtoMessage() {}

func (x *GrpcJsonTranscoder) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcJsonTranscoder.ProtoReflect.Descriptor instead.
func (*GrpcJsonTranscoder) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescGZIP(), []int{0}
}

func (m *GrpcJsonTranscoder) GetDescriptorSet() isGrpcJsonTranscoder_DescriptorSet {
	if m != nil {
		return m.DescriptorSet
	}
	return nil
}

func (x *GrpcJsonTranscoder) GetProtoDescriptor() string {
	if x, ok := x.GetDescriptorSet().(*GrpcJsonTranscoder_ProtoDescriptor); ok {
		return x.ProtoDescriptor
	}
	return ""
}

func (x *GrpcJsonTranscoder) GetProtoDescriptorBin() []byte {
	if x, ok := x.GetDescriptorSet().(*GrpcJsonTranscoder_ProtoDescriptorBin); ok {
		return x.ProtoDescriptorBin
	}
	return nil
}

func (x *GrpcJsonTranscoder) GetServices() []string {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *GrpcJsonTranscoder) GetPrintOptions() *GrpcJsonTranscoder_PrintOptions {
	if x != nil {
		return x.PrintOptions
	}
	return nil
}

func (x *GrpcJsonTranscoder) GetMatchIncomingRequestRoute() bool {
	if x != nil {
		return x.MatchIncomingRequestRoute
	}
	return false
}

func (x *GrpcJsonTranscoder) GetIgnoredQueryParameters() []string {
	if x != nil {
		return x.IgnoredQueryParameters
	}
	return nil
}

func (x *GrpcJsonTranscoder) GetAutoMapping() bool {
	if x != nil {
		return x.AutoMapping
	}
	return false
}

func (x *GrpcJsonTranscoder) GetIgnoreUnknownQueryParameters() bool {
	if x != nil {
		return x.IgnoreUnknownQueryParameters
	}
	return false
}

func (x *GrpcJsonTranscoder) GetConvertGrpcStatus() bool {
	if x != nil {
		return x.ConvertGrpcStatus
	}
	return false
}

type isGrpcJsonTranscoder_DescriptorSet interface {
	isGrpcJsonTranscoder_DescriptorSet()
}

type GrpcJsonTranscoder_ProtoDescriptor struct {
	// Supplies the filename of
	// :ref:`the proto descriptor set <config_grpc_json_generate_proto_descriptor_set>` for the gRPC
	// services.
	ProtoDescriptor string `protobuf:"bytes,1,opt,name=proto_descriptor,json=protoDescriptor,proto3,oneof"`
}

type GrpcJsonTranscoder_ProtoDescriptorBin struct {
	// Supplies the binary content of
	// :ref:`the proto descriptor set <config_grpc_json_generate_proto_descriptor_set>` for the gRPC
	// services.
	// Note: in yaml, this must be provided as a base64 standard encoded string; yaml can't handle binary bytes
	ProtoDescriptorBin []byte `protobuf:"bytes,4,opt,name=proto_descriptor_bin,json=protoDescriptorBin,proto3,oneof"`
}

func (*GrpcJsonTranscoder_ProtoDescriptor) isGrpcJsonTranscoder_DescriptorSet() {}

func (*GrpcJsonTranscoder_ProtoDescriptorBin) isGrpcJsonTranscoder_DescriptorSet() {}

type GrpcJsonTranscoder_PrintOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Whether to add spaces, line breaks and indentation to make the JSON
	// output easy to read. Defaults to false.
	AddWhitespace bool `protobuf:"varint,1,opt,name=add_whitespace,json=addWhitespace,proto3" json:"add_whitespace,omitempty"`
	// Whether to always print primitive fields. By default primitive
	// fields with default values will be omitted in JSON output. For
	// example, an int32 field set to 0 will be omitted. Setting this flag to
	// true will override the default behavior and print primitive fields
	// regardless of their values. Defaults to false.
	AlwaysPrintPrimitiveFields bool `protobuf:"varint,2,opt,name=always_print_primitive_fields,json=alwaysPrintPrimitiveFields,proto3" json:"always_print_primitive_fields,omitempty"`
	// Whether to always print enums as ints. By default they are rendered
	// as strings. Defaults to false.
	AlwaysPrintEnumsAsInts bool `protobuf:"varint,3,opt,name=always_print_enums_as_ints,json=alwaysPrintEnumsAsInts,proto3" json:"always_print_enums_as_ints,omitempty"`
	// Whether to preserve proto field names. By default protobuf will
	// generate JSON field names using the ``json_name`` option, or lower camel case,
	// in that order. Setting this flag will preserve the original field names. Defaults to false.
	PreserveProtoFieldNames bool `protobuf:"varint,4,opt,name=preserve_proto_field_names,json=preserveProtoFieldNames,proto3" json:"preserve_proto_field_names,omitempty"`
}

func (x *GrpcJsonTranscoder_PrintOptions) Reset() {
	*x = GrpcJsonTranscoder_PrintOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcJsonTranscoder_PrintOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcJsonTranscoder_PrintOptions) ProtoMessage() {}

func (x *GrpcJsonTranscoder_PrintOptions) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcJsonTranscoder_PrintOptions.ProtoReflect.Descriptor instead.
func (*GrpcJsonTranscoder_PrintOptions) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GrpcJsonTranscoder_PrintOptions) GetAddWhitespace() bool {
	if x != nil {
		return x.AddWhitespace
	}
	return false
}

func (x *GrpcJsonTranscoder_PrintOptions) GetAlwaysPrintPrimitiveFields() bool {
	if x != nil {
		return x.AlwaysPrintPrimitiveFields
	}
	return false
}

func (x *GrpcJsonTranscoder_PrintOptions) GetAlwaysPrintEnumsAsInts() bool {
	if x != nil {
		return x.AlwaysPrintEnumsAsInts
	}
	return false
}

func (x *GrpcJsonTranscoder_PrintOptions) GetPreserveProtoFieldNames() bool {
	if x != nil {
		return x.PreserveProtoFieldNames
	}
	return false
}

var File_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDesc = []byte{
	0x0a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x6a, 0x73, 0x6f, 0x6e,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x1e, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x2e, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x06,
	0x0a, 0x12, 0x47, 0x72, 0x70, 0x63, 0x4a, 0x73, 0x6f, 0x6e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63,
	0x6f, 0x64, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x12, 0x32, 0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x5f, 0x62, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x48,
	0x00, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x42, 0x69, 0x6e, 0x12, 0x24, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x08,
	0x01, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x64, 0x0a, 0x0d, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x3f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x2e, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x4a, 0x73, 0x6f, 0x6e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x0c, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x3f, 0x0a, 0x1c, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x6e, 0x63, 0x6f, 0x6d,
	0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x19, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x49, 0x6e,
	0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x12, 0x38, 0x0a, 0x18, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x64, 0x5f, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x16, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x64, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x21, 0x0a, 0x0c,
	0x61, 0x75, 0x74, 0x6f, 0x5f, 0x6d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x61, 0x75, 0x74, 0x6f, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x12,
	0x45, 0x0a, 0x1f, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x5f, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77,
	0x6e, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x1c, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65,
	0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72,
	0x74, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x11, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x47, 0x72, 0x70, 0x63,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0xf1, 0x01, 0x0a, 0x0c, 0x50, 0x72, 0x69, 0x6e, 0x74,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x64, 0x64, 0x5f, 0x77,
	0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0d, 0x61, 0x64, 0x64, 0x57, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x41,
	0x0a, 0x1d, 0x61, 0x6c, 0x77, 0x61, 0x79, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x5f, 0x70,
	0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x1a, 0x61, 0x6c, 0x77, 0x61, 0x79, 0x73, 0x50, 0x72, 0x69,
	0x6e, 0x74, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x12, 0x3a, 0x0a, 0x1a, 0x61, 0x6c, 0x77, 0x61, 0x79, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x6e,
	0x74, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x5f, 0x61, 0x73, 0x5f, 0x69, 0x6e, 0x74, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x16, 0x61, 0x6c, 0x77, 0x61, 0x79, 0x73, 0x50, 0x72, 0x69,
	0x6e, 0x74, 0x45, 0x6e, 0x75, 0x6d, 0x73, 0x41, 0x73, 0x49, 0x6e, 0x74, 0x73, 0x12, 0x3b, 0x0a,
	0x1a, 0x70, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x17, 0x70, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x42, 0x15, 0x0a, 0x0e, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x12, 0x03, 0xf8, 0x42,
	0x01, 0x42, 0x50, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0xc0, 0xf5, 0x04, 0x01, 0xb8, 0xf5, 0x04, 0x01, 0xd0,
	0xf5, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescData = file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDesc
)

func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescData)
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_goTypes = []interface{}{
	(*GrpcJsonTranscoder)(nil),              // 0: grpc_json.options.gloo.solo.io.GrpcJsonTranscoder
	(*GrpcJsonTranscoder_PrintOptions)(nil), // 1: grpc_json.options.gloo.solo.io.GrpcJsonTranscoder.PrintOptions
}
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_depIdxs = []int32{
	1, // 0: grpc_json.options.gloo.solo.io.GrpcJsonTranscoder.print_options:type_name -> grpc_json.options.gloo.solo.io.GrpcJsonTranscoder.PrintOptions
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcJsonTranscoder); i {
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
		file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcJsonTranscoder_PrintOptions); i {
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
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*GrpcJsonTranscoder_ProtoDescriptor)(nil),
		(*GrpcJsonTranscoder_ProtoDescriptorBin)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_rawDesc = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_grpc_json_grpc_json_proto_depIdxs = nil
}
