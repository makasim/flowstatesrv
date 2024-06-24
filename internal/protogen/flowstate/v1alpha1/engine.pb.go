// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: flowstate/v1alpha1/engine.proto

package flowstatev1alpha1

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

type DoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateCtxs []*StateCtx `protobuf:"bytes,1,rep,name=state_ctxs,json=stateCtxs,proto3" json:"state_ctxs,omitempty"`
	Commands  []*Command  `protobuf:"bytes,2,rep,name=commands,proto3" json:"commands,omitempty"`
}

func (x *DoRequest) Reset() {
	*x = DoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoRequest) ProtoMessage() {}

func (x *DoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoRequest.ProtoReflect.Descriptor instead.
func (*DoRequest) Descriptor() ([]byte, []int) {
	return file_flowstate_v1alpha1_engine_proto_rawDescGZIP(), []int{0}
}

func (x *DoRequest) GetStateCtxs() []*StateCtx {
	if x != nil {
		return x.StateCtxs
	}
	return nil
}

func (x *DoRequest) GetCommands() []*Command {
	if x != nil {
		return x.Commands
	}
	return nil
}

type DoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateCtxs []*StateCtx      `protobuf:"bytes,1,rep,name=state_ctxs,json=stateCtxs,proto3" json:"state_ctxs,omitempty"`
	Results   []*CommandResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *DoResponse) Reset() {
	*x = DoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoResponse) ProtoMessage() {}

func (x *DoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoResponse.ProtoReflect.Descriptor instead.
func (*DoResponse) Descriptor() ([]byte, []int) {
	return file_flowstate_v1alpha1_engine_proto_rawDescGZIP(), []int{1}
}

func (x *DoResponse) GetStateCtxs() []*StateCtx {
	if x != nil {
		return x.StateCtxs
	}
	return nil
}

func (x *DoResponse) GetResults() []*CommandResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type WatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SinceRev    int64             `protobuf:"varint,1,opt,name=since_rev,json=sinceRev,proto3" json:"since_rev,omitempty"`
	SinceLatest bool              `protobuf:"varint,2,opt,name=since_latest,json=sinceLatest,proto3" json:"since_latest,omitempty"`
	Labels      map[string]string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *WatchRequest) Reset() {
	*x = WatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchRequest) ProtoMessage() {}

func (x *WatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchRequest.ProtoReflect.Descriptor instead.
func (*WatchRequest) Descriptor() ([]byte, []int) {
	return file_flowstate_v1alpha1_engine_proto_rawDescGZIP(), []int{2}
}

func (x *WatchRequest) GetSinceRev() int64 {
	if x != nil {
		return x.SinceRev
	}
	return 0
}

func (x *WatchRequest) GetSinceLatest() bool {
	if x != nil {
		return x.SinceLatest
	}
	return false
}

func (x *WatchRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type WatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State *State `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *WatchResponse) Reset() {
	*x = WatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchResponse) ProtoMessage() {}

func (x *WatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1alpha1_engine_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchResponse.ProtoReflect.Descriptor instead.
func (*WatchResponse) Descriptor() ([]byte, []int) {
	return file_flowstate_v1alpha1_engine_proto_rawDescGZIP(), []int{3}
}

func (x *WatchResponse) GetState() *State {
	if x != nil {
		return x.State
	}
	return nil
}

var File_flowstate_v1alpha1_engine_proto protoreflect.FileDescriptor

var file_flowstate_v1alpha1_engine_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x21, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x09, 0x44, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3b, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f,
	0x63, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x66, 0x6c, 0x6f,
	0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x74, 0x78, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x43,
	0x74, 0x78, 0x73, 0x12, 0x37, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x22, 0x86, 0x01, 0x0a,
	0x0a, 0x44, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x74, 0x78, 0x52, 0x09, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x43, 0x74, 0x78, 0x73, 0x12, 0x3b, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0xcf, 0x01, 0x0a, 0x0c, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x69, 0x6e, 0x63, 0x65, 0x5f,
	0x72, 0x65, 0x76, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x69, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x76, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x69, 0x6e, 0x63, 0x65, 0x5f, 0x6c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x69, 0x6e, 0x63, 0x65,
	0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x44, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x40, 0x0a, 0x0d, 0x57, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x32, 0xa8, 0x01, 0x0a, 0x0d, 0x45, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x02, 0x44,
	0x6f, 0x12, 0x1d, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x50, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x20, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x42, 0xe6, 0x01, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x6f,
	0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42,
	0x0b, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x56,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x6b, 0x61, 0x73,
	0x69, 0x6d, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x72, 0x76, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65,
	0x6e, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x3b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x12, 0x46,
	0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0xca, 0x02, 0x12, 0x46, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5c, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x1e, 0x46, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x46, 0x6c, 0x6f, 0x77, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flowstate_v1alpha1_engine_proto_rawDescOnce sync.Once
	file_flowstate_v1alpha1_engine_proto_rawDescData = file_flowstate_v1alpha1_engine_proto_rawDesc
)

func file_flowstate_v1alpha1_engine_proto_rawDescGZIP() []byte {
	file_flowstate_v1alpha1_engine_proto_rawDescOnce.Do(func() {
		file_flowstate_v1alpha1_engine_proto_rawDescData = protoimpl.X.CompressGZIP(file_flowstate_v1alpha1_engine_proto_rawDescData)
	})
	return file_flowstate_v1alpha1_engine_proto_rawDescData
}

var file_flowstate_v1alpha1_engine_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_flowstate_v1alpha1_engine_proto_goTypes = []interface{}{
	(*DoRequest)(nil),     // 0: flowstate.v1alpha1.DoRequest
	(*DoResponse)(nil),    // 1: flowstate.v1alpha1.DoResponse
	(*WatchRequest)(nil),  // 2: flowstate.v1alpha1.WatchRequest
	(*WatchResponse)(nil), // 3: flowstate.v1alpha1.WatchResponse
	nil,                   // 4: flowstate.v1alpha1.WatchRequest.LabelsEntry
	(*StateCtx)(nil),      // 5: flowstate.v1alpha1.StateCtx
	(*Command)(nil),       // 6: flowstate.v1alpha1.Command
	(*CommandResult)(nil), // 7: flowstate.v1alpha1.CommandResult
	(*State)(nil),         // 8: flowstate.v1alpha1.State
}
var file_flowstate_v1alpha1_engine_proto_depIdxs = []int32{
	5, // 0: flowstate.v1alpha1.DoRequest.state_ctxs:type_name -> flowstate.v1alpha1.StateCtx
	6, // 1: flowstate.v1alpha1.DoRequest.commands:type_name -> flowstate.v1alpha1.Command
	5, // 2: flowstate.v1alpha1.DoResponse.state_ctxs:type_name -> flowstate.v1alpha1.StateCtx
	7, // 3: flowstate.v1alpha1.DoResponse.results:type_name -> flowstate.v1alpha1.CommandResult
	4, // 4: flowstate.v1alpha1.WatchRequest.labels:type_name -> flowstate.v1alpha1.WatchRequest.LabelsEntry
	8, // 5: flowstate.v1alpha1.WatchResponse.state:type_name -> flowstate.v1alpha1.State
	0, // 6: flowstate.v1alpha1.EngineService.Do:input_type -> flowstate.v1alpha1.DoRequest
	2, // 7: flowstate.v1alpha1.EngineService.Watch:input_type -> flowstate.v1alpha1.WatchRequest
	1, // 8: flowstate.v1alpha1.EngineService.Do:output_type -> flowstate.v1alpha1.DoResponse
	3, // 9: flowstate.v1alpha1.EngineService.Watch:output_type -> flowstate.v1alpha1.WatchResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_flowstate_v1alpha1_engine_proto_init() }
func file_flowstate_v1alpha1_engine_proto_init() {
	if File_flowstate_v1alpha1_engine_proto != nil {
		return
	}
	file_flowstate_v1alpha1_commands_proto_init()
	file_flowstate_v1alpha1_state_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_flowstate_v1alpha1_engine_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoRequest); i {
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
		file_flowstate_v1alpha1_engine_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoResponse); i {
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
		file_flowstate_v1alpha1_engine_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchRequest); i {
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
		file_flowstate_v1alpha1_engine_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchResponse); i {
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
			RawDescriptor: file_flowstate_v1alpha1_engine_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_flowstate_v1alpha1_engine_proto_goTypes,
		DependencyIndexes: file_flowstate_v1alpha1_engine_proto_depIdxs,
		MessageInfos:      file_flowstate_v1alpha1_engine_proto_msgTypes,
	}.Build()
	File_flowstate_v1alpha1_engine_proto = out.File
	file_flowstate_v1alpha1_engine_proto_rawDesc = nil
	file_flowstate_v1alpha1_engine_proto_goTypes = nil
	file_flowstate_v1alpha1_engine_proto_depIdxs = nil
}