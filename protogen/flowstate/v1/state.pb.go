// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: flowstate/v1/state.proto

package flowstatev1

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

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Rev int64  `protobuf:"varint,2,opt,name=rev,proto3" json:"rev,omitempty"`
	B   []byte `protobuf:"bytes,3,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_flowstate_v1_state_proto_rawDescGZIP(), []int{0}
}

func (x *Data) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Data) GetRev() int64 {
	if x != nil {
		return x.Rev
	}
	return 0
}

func (x *Data) GetB() []byte {
	if x != nil {
		return x.B
	}
	return nil
}

type DataRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Rev int64  `protobuf:"varint,2,opt,name=rev,proto3" json:"rev,omitempty"`
}

func (x *DataRef) Reset() {
	*x = DataRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1_state_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataRef) ProtoMessage() {}

func (x *DataRef) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1_state_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataRef.ProtoReflect.Descriptor instead.
func (*DataRef) Descriptor() ([]byte, []int) {
	return file_flowstate_v1_state_proto_rawDescGZIP(), []int{1}
}

func (x *DataRef) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DataRef) GetRev() int64 {
	if x != nil {
		return x.Rev
	}
	return 0
}

type StateRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Rev int64  `protobuf:"varint,2,opt,name=rev,proto3" json:"rev,omitempty"`
}

func (x *StateRef) Reset() {
	*x = StateRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1_state_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateRef) ProtoMessage() {}

func (x *StateRef) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1_state_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateRef.ProtoReflect.Descriptor instead.
func (*StateRef) Descriptor() ([]byte, []int) {
	return file_flowstate_v1_state_proto_rawDescGZIP(), []int{2}
}

func (x *StateRef) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StateRef) GetRev() int64 {
	if x != nil {
		return x.Rev
	}
	return 0
}

type StateContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Committed   *State        `protobuf:"bytes,1,opt,name=committed,proto3" json:"committed,omitempty"`
	Current     *State        `protobuf:"bytes,2,opt,name=current,proto3" json:"current,omitempty"`
	Transitions []*Transition `protobuf:"bytes,3,rep,name=transitions,proto3" json:"transitions,omitempty"`
}

func (x *StateContext) Reset() {
	*x = StateContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1_state_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateContext) ProtoMessage() {}

func (x *StateContext) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1_state_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateContext.ProtoReflect.Descriptor instead.
func (*StateContext) Descriptor() ([]byte, []int) {
	return file_flowstate_v1_state_proto_rawDescGZIP(), []int{3}
}

func (x *StateContext) GetCommitted() *State {
	if x != nil {
		return x.Committed
	}
	return nil
}

func (x *StateContext) GetCurrent() *State {
	if x != nil {
		return x.Current
	}
	return nil
}

func (x *StateContext) GetTransitions() []*Transition {
	if x != nil {
		return x.Transitions
	}
	return nil
}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                   string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Rev                  int64             `protobuf:"varint,2,opt,name=rev,proto3" json:"rev,omitempty"`
	Annotations          map[string]string `protobuf:"bytes,3,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Labels               map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CommittedAtUnixMilli int64             `protobuf:"varint,5,opt,name=committed_at_unix_milli,json=committedAtUnixMilli,proto3" json:"committed_at_unix_milli,omitempty"`
	Transition           *Transition       `protobuf:"bytes,6,opt,name=transition,proto3" json:"transition,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1_state_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1_state_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_flowstate_v1_state_proto_rawDescGZIP(), []int{4}
}

func (x *State) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *State) GetRev() int64 {
	if x != nil {
		return x.Rev
	}
	return 0
}

func (x *State) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

func (x *State) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *State) GetCommittedAtUnixMilli() int64 {
	if x != nil {
		return x.CommittedAtUnixMilli
	}
	return 0
}

func (x *State) GetTransition() *Transition {
	if x != nil {
		return x.Transition
	}
	return nil
}

type Transition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From        string            `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To          string            `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Annotations map[string]string `protobuf:"bytes,3,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Transition) Reset() {
	*x = Transition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flowstate_v1_state_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transition) ProtoMessage() {}

func (x *Transition) ProtoReflect() protoreflect.Message {
	mi := &file_flowstate_v1_state_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transition.ProtoReflect.Descriptor instead.
func (*Transition) Descriptor() ([]byte, []int) {
	return file_flowstate_v1_state_proto_rawDescGZIP(), []int{5}
}

func (x *Transition) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Transition) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Transition) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

var File_flowstate_v1_state_proto protoreflect.FileDescriptor

var file_flowstate_v1_state_proto_rawDesc = []byte{
	0x0a, 0x18, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x36, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72,
	0x65, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x62,
	0x22, 0x2b, 0x0a, 0x07, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x66, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x72,
	0x65, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x65, 0x76, 0x22, 0x2c, 0x0a,
	0x08, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x76,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x65, 0x76, 0x22, 0xac, 0x01, 0x0a, 0x0c,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x31, 0x0a, 0x09,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x12,
	0x2d, 0x0a, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x3a,
	0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x96, 0x03, 0x0a, 0x05, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x72, 0x65, 0x76, 0x12, 0x46, 0x0a, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x37,
	0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x35, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6d, 0x69, 0x6c,
	0x6c, 0x69, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x55, 0x6e, 0x69, 0x78, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x12, 0x38,
	0x0a, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x3e, 0x0a, 0x10, 0x41, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0xbd, 0x01, 0x0a, 0x0a, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x4b, 0x0a, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x1a, 0x3e, 0x0a, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0xb2, 0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x65, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x6b, 0x61, 0x73, 0x69, 0x6d, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x73, 0x72, 0x76, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e,
	0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa,
	0x02, 0x0c, 0x46, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x0c, 0x46, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18,
	0x46, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x46, 0x6c, 0x6f, 0x77, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flowstate_v1_state_proto_rawDescOnce sync.Once
	file_flowstate_v1_state_proto_rawDescData = file_flowstate_v1_state_proto_rawDesc
)

func file_flowstate_v1_state_proto_rawDescGZIP() []byte {
	file_flowstate_v1_state_proto_rawDescOnce.Do(func() {
		file_flowstate_v1_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_flowstate_v1_state_proto_rawDescData)
	})
	return file_flowstate_v1_state_proto_rawDescData
}

var file_flowstate_v1_state_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_flowstate_v1_state_proto_goTypes = []any{
	(*Data)(nil),         // 0: flowstate.v1.Data
	(*DataRef)(nil),      // 1: flowstate.v1.DataRef
	(*StateRef)(nil),     // 2: flowstate.v1.StateRef
	(*StateContext)(nil), // 3: flowstate.v1.StateContext
	(*State)(nil),        // 4: flowstate.v1.State
	(*Transition)(nil),   // 5: flowstate.v1.Transition
	nil,                  // 6: flowstate.v1.State.AnnotationsEntry
	nil,                  // 7: flowstate.v1.State.LabelsEntry
	nil,                  // 8: flowstate.v1.Transition.AnnotationsEntry
}
var file_flowstate_v1_state_proto_depIdxs = []int32{
	4, // 0: flowstate.v1.StateContext.committed:type_name -> flowstate.v1.State
	4, // 1: flowstate.v1.StateContext.current:type_name -> flowstate.v1.State
	5, // 2: flowstate.v1.StateContext.transitions:type_name -> flowstate.v1.Transition
	6, // 3: flowstate.v1.State.annotations:type_name -> flowstate.v1.State.AnnotationsEntry
	7, // 4: flowstate.v1.State.labels:type_name -> flowstate.v1.State.LabelsEntry
	5, // 5: flowstate.v1.State.transition:type_name -> flowstate.v1.Transition
	8, // 6: flowstate.v1.Transition.annotations:type_name -> flowstate.v1.Transition.AnnotationsEntry
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_flowstate_v1_state_proto_init() }
func file_flowstate_v1_state_proto_init() {
	if File_flowstate_v1_state_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_flowstate_v1_state_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Data); i {
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
		file_flowstate_v1_state_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DataRef); i {
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
		file_flowstate_v1_state_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*StateRef); i {
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
		file_flowstate_v1_state_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*StateContext); i {
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
		file_flowstate_v1_state_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*State); i {
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
		file_flowstate_v1_state_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Transition); i {
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
			RawDescriptor: file_flowstate_v1_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flowstate_v1_state_proto_goTypes,
		DependencyIndexes: file_flowstate_v1_state_proto_depIdxs,
		MessageInfos:      file_flowstate_v1_state_proto_msgTypes,
	}.Build()
	File_flowstate_v1_state_proto = out.File
	file_flowstate_v1_state_proto_rawDesc = nil
	file_flowstate_v1_state_proto_goTypes = nil
	file_flowstate_v1_state_proto_depIdxs = nil
}
