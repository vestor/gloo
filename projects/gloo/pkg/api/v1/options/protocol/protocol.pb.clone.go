// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/protocol/protocol.proto

package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/solo-io/protoc-gen-ext/pkg/clone"
	"google.golang.org/protobuf/proto"

	github_com_golang_protobuf_ptypes_duration "github.com/golang/protobuf/ptypes/duration"

	github_com_golang_protobuf_ptypes_wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = bytes.Compare
	_ = strings.Compare
	_ = clone.Cloner(nil)
	_ = proto.Message(nil)
)

// Clone function
func (m *HttpProtocolOptions) Clone() proto.Message {
	var target *HttpProtocolOptions
	if m == nil {
		return target
	}
	target = &HttpProtocolOptions{}

	if h, ok := interface{}(m.GetIdleTimeout()).(clone.Cloner); ok {
		target.IdleTimeout = h.Clone().(*github_com_golang_protobuf_ptypes_duration.Duration)
	} else {
		target.IdleTimeout = proto.Clone(m.GetIdleTimeout()).(*github_com_golang_protobuf_ptypes_duration.Duration)
	}

	target.MaxHeadersCount = m.GetMaxHeadersCount()

	if h, ok := interface{}(m.GetMaxStreamDuration()).(clone.Cloner); ok {
		target.MaxStreamDuration = h.Clone().(*github_com_golang_protobuf_ptypes_duration.Duration)
	} else {
		target.MaxStreamDuration = proto.Clone(m.GetMaxStreamDuration()).(*github_com_golang_protobuf_ptypes_duration.Duration)
	}

	target.HeadersWithUnderscoresAction = m.GetHeadersWithUnderscoresAction()

	return target
}

// Clone function
func (m *Http1ProtocolOptions) Clone() proto.Message {
	var target *Http1ProtocolOptions
	if m == nil {
		return target
	}
	target = &Http1ProtocolOptions{}

	target.EnableTrailers = m.GetEnableTrailers()

	if h, ok := interface{}(m.GetOverrideStreamErrorOnInvalidHttpMessage()).(clone.Cloner); ok {
		target.OverrideStreamErrorOnInvalidHttpMessage = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.BoolValue)
	} else {
		target.OverrideStreamErrorOnInvalidHttpMessage = proto.Clone(m.GetOverrideStreamErrorOnInvalidHttpMessage()).(*github_com_golang_protobuf_ptypes_wrappers.BoolValue)
	}

	switch m.HeaderFormat.(type) {

	case *Http1ProtocolOptions_ProperCaseHeaderKeyFormat:

		target.HeaderFormat = &Http1ProtocolOptions_ProperCaseHeaderKeyFormat{
			ProperCaseHeaderKeyFormat: m.GetProperCaseHeaderKeyFormat(),
		}

	case *Http1ProtocolOptions_PreserveCaseHeaderKeyFormat:

		target.HeaderFormat = &Http1ProtocolOptions_PreserveCaseHeaderKeyFormat{
			PreserveCaseHeaderKeyFormat: m.GetPreserveCaseHeaderKeyFormat(),
		}

	}

	return target
}

// Clone function
func (m *Http2ProtocolOptions) Clone() proto.Message {
	var target *Http2ProtocolOptions
	if m == nil {
		return target
	}
	target = &Http2ProtocolOptions{}

	if h, ok := interface{}(m.GetMaxConcurrentStreams()).(clone.Cloner); ok {
		target.MaxConcurrentStreams = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.UInt32Value)
	} else {
		target.MaxConcurrentStreams = proto.Clone(m.GetMaxConcurrentStreams()).(*github_com_golang_protobuf_ptypes_wrappers.UInt32Value)
	}

	if h, ok := interface{}(m.GetInitialStreamWindowSize()).(clone.Cloner); ok {
		target.InitialStreamWindowSize = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.UInt32Value)
	} else {
		target.InitialStreamWindowSize = proto.Clone(m.GetInitialStreamWindowSize()).(*github_com_golang_protobuf_ptypes_wrappers.UInt32Value)
	}

	if h, ok := interface{}(m.GetInitialConnectionWindowSize()).(clone.Cloner); ok {
		target.InitialConnectionWindowSize = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.UInt32Value)
	} else {
		target.InitialConnectionWindowSize = proto.Clone(m.GetInitialConnectionWindowSize()).(*github_com_golang_protobuf_ptypes_wrappers.UInt32Value)
	}

	if h, ok := interface{}(m.GetOverrideStreamErrorOnInvalidHttpMessage()).(clone.Cloner); ok {
		target.OverrideStreamErrorOnInvalidHttpMessage = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.BoolValue)
	} else {
		target.OverrideStreamErrorOnInvalidHttpMessage = proto.Clone(m.GetOverrideStreamErrorOnInvalidHttpMessage()).(*github_com_golang_protobuf_ptypes_wrappers.BoolValue)
	}

	return target
}
