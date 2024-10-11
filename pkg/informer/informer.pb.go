// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: proto/informer.proto

package informer

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

// EventType represents the type of event.
type EventType int32

const (
	EventType_CREATED EventType = 0
	EventType_UPDATED EventType = 1
	EventType_DELETED EventType = 2
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "CREATED",
		1: "UPDATED",
		2: "DELETED",
	}
	EventType_value = map[string]int32{
		"CREATED": 0,
		"UPDATED": 1,
		"DELETED": 2,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_informer_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_proto_informer_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_proto_informer_proto_rawDescGZIP(), []int{0}
}

type IPInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string            `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Labels    map[string]string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Ips       []string          `protobuf:"bytes,4,rep,name=ips,proto3" json:"ips,omitempty"`
	Kind      string            `protobuf:"bytes,5,opt,name=kind,proto3" json:"kind,omitempty"`
}

func (x *IPInfo) Reset() {
	*x = IPInfo{}
	mi := &file_proto_informer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IPInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPInfo) ProtoMessage() {}

func (x *IPInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_informer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPInfo.ProtoReflect.Descriptor instead.
func (*IPInfo) Descriptor() ([]byte, []int) {
	return file_proto_informer_proto_rawDescGZIP(), []int{0}
}

func (x *IPInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IPInfo) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *IPInfo) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *IPInfo) GetIps() []string {
	if x != nil {
		return x.Ips
	}
	return nil
}

func (x *IPInfo) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

type PodInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IpInfo       *IPInfo  `protobuf:"bytes,1,opt,name=ip_info,json=ipInfo,proto3" json:"ip_info,omitempty"`
	Uid          string   `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	OwnerName    string   `protobuf:"bytes,3,opt,name=owner_name,json=ownerName,proto3" json:"owner_name,omitempty"`
	OwnerKind    string   `protobuf:"bytes,4,opt,name=owner_kind,json=ownerKind,proto3" json:"owner_kind,omitempty"`
	NodeName     string   `protobuf:"bytes,5,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	StartTimeStr string   `protobuf:"bytes,6,opt,name=start_time_str,json=startTimeStr,proto3" json:"start_time_str,omitempty"`
	ContainerIds []string `protobuf:"bytes,7,rep,name=container_ids,json=containerIds,proto3" json:"container_ids,omitempty"`
	HostName     string   `protobuf:"bytes,8,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	HostIp       string   `protobuf:"bytes,9,opt,name=host_ip,json=hostIp,proto3" json:"host_ip,omitempty"`
}

func (x *PodInfo) Reset() {
	*x = PodInfo{}
	mi := &file_proto_informer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PodInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodInfo) ProtoMessage() {}

func (x *PodInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_informer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodInfo.ProtoReflect.Descriptor instead.
func (*PodInfo) Descriptor() ([]byte, []int) {
	return file_proto_informer_proto_rawDescGZIP(), []int{1}
}

func (x *PodInfo) GetIpInfo() *IPInfo {
	if x != nil {
		return x.IpInfo
	}
	return nil
}

func (x *PodInfo) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *PodInfo) GetOwnerName() string {
	if x != nil {
		return x.OwnerName
	}
	return ""
}

func (x *PodInfo) GetOwnerKind() string {
	if x != nil {
		return x.OwnerKind
	}
	return ""
}

func (x *PodInfo) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *PodInfo) GetStartTimeStr() string {
	if x != nil {
		return x.StartTimeStr
	}
	return ""
}

func (x *PodInfo) GetContainerIds() []string {
	if x != nil {
		return x.ContainerIds
	}
	return nil
}

func (x *PodInfo) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *PodInfo) GetHostIp() string {
	if x != nil {
		return x.HostIp
	}
	return ""
}

// Event represents a single event.
type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Todo: add timestamp:
	// - avoids out-of-order events coming from different informers (local vs remote)
	// - on failure/reconnection you don't need to receive all the events again
	Type EventType `protobuf:"varint,1,opt,name=type,proto3,enum=informer.EventType" json:"type,omitempty"`
	// Types that are assignable to Resource:
	//
	//	*Event_Pod
	//	*Event_IpInfo
	Resource isEvent_Resource `protobuf_oneof:"resource"`
}

func (x *Event) Reset() {
	*x = Event{}
	mi := &file_proto_informer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_proto_informer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_proto_informer_proto_rawDescGZIP(), []int{2}
}

func (x *Event) GetType() EventType {
	if x != nil {
		return x.Type
	}
	return EventType_CREATED
}

func (m *Event) GetResource() isEvent_Resource {
	if m != nil {
		return m.Resource
	}
	return nil
}

func (x *Event) GetPod() *PodInfo {
	if x, ok := x.GetResource().(*Event_Pod); ok {
		return x.Pod
	}
	return nil
}

func (x *Event) GetIpInfo() *IPInfo {
	if x, ok := x.GetResource().(*Event_IpInfo); ok {
		return x.IpInfo
	}
	return nil
}

type isEvent_Resource interface {
	isEvent_Resource()
}

type Event_Pod struct {
	Pod *PodInfo `protobuf:"bytes,2,opt,name=pod,proto3,oneof"`
}

type Event_IpInfo struct {
	IpInfo *IPInfo `protobuf:"bytes,3,opt,name=ip_info,json=ipInfo,proto3,oneof"`
}

func (*Event_Pod) isEvent_Resource() {}

func (*Event_IpInfo) isEvent_Resource() {}

// Empty message for Subscribe RPC
type SubscribeMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubscribeMessage) Reset() {
	*x = SubscribeMessage{}
	mi := &file_proto_informer_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubscribeMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeMessage) ProtoMessage() {}

func (x *SubscribeMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_informer_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeMessage.ProtoReflect.Descriptor instead.
func (*SubscribeMessage) Descriptor() ([]byte, []int) {
	return file_proto_informer_proto_rawDescGZIP(), []int{3}
}

var File_proto_informer_proto protoreflect.FileDescriptor

var file_proto_informer_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72,
	0x22, 0xd1, 0x01, 0x0a, 0x06, 0x49, 0x50, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x34, 0x0a,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e, 0x49, 0x50, 0x49, 0x6e, 0x66, 0x6f, 0x2e,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x70, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x03, 0x69, 0x70, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0xa2, 0x02, 0x0a, 0x07, 0x50, 0x6f, 0x64, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x29, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e, 0x49, 0x50, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x06, 0x69, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6e, 0x6f, 0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x72, 0x12, 0x23,
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x73, 0x18,
	0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x49, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x70, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x70, 0x22, 0x90, 0x01, 0x0a, 0x05, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x13, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x03,
	0x70, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x69, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x65, 0x72, 0x2e, 0x50, 0x6f, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x03,
	0x70, 0x6f, 0x64, 0x12, 0x2b, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e,
	0x49, 0x50, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x06, 0x69, 0x70, 0x49, 0x6e, 0x66, 0x6f,
	0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x12, 0x0a, 0x10,
	0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2a, 0x32, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a,
	0x07, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x50,
	0x44, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x44, 0x10, 0x02, 0x32, 0x50, 0x0a, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x09, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x1a, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d,
	0x65, 0x72, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x0f, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x69, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_informer_proto_rawDescOnce sync.Once
	file_proto_informer_proto_rawDescData = file_proto_informer_proto_rawDesc
)

func file_proto_informer_proto_rawDescGZIP() []byte {
	file_proto_informer_proto_rawDescOnce.Do(func() {
		file_proto_informer_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_informer_proto_rawDescData)
	})
	return file_proto_informer_proto_rawDescData
}

var file_proto_informer_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_informer_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_informer_proto_goTypes = []any{
	(EventType)(0),           // 0: informer.EventType
	(*IPInfo)(nil),           // 1: informer.IPInfo
	(*PodInfo)(nil),          // 2: informer.PodInfo
	(*Event)(nil),            // 3: informer.Event
	(*SubscribeMessage)(nil), // 4: informer.SubscribeMessage
	nil,                      // 5: informer.IPInfo.LabelsEntry
}
var file_proto_informer_proto_depIdxs = []int32{
	5, // 0: informer.IPInfo.labels:type_name -> informer.IPInfo.LabelsEntry
	1, // 1: informer.PodInfo.ip_info:type_name -> informer.IPInfo
	0, // 2: informer.Event.type:type_name -> informer.EventType
	2, // 3: informer.Event.pod:type_name -> informer.PodInfo
	1, // 4: informer.Event.ip_info:type_name -> informer.IPInfo
	4, // 5: informer.EventStreamService.Subscribe:input_type -> informer.SubscribeMessage
	3, // 6: informer.EventStreamService.Subscribe:output_type -> informer.Event
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_informer_proto_init() }
func file_proto_informer_proto_init() {
	if File_proto_informer_proto != nil {
		return
	}
	file_proto_informer_proto_msgTypes[2].OneofWrappers = []any{
		(*Event_Pod)(nil),
		(*Event_IpInfo)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_informer_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_informer_proto_goTypes,
		DependencyIndexes: file_proto_informer_proto_depIdxs,
		EnumInfos:         file_proto_informer_proto_enumTypes,
		MessageInfos:      file_proto_informer_proto_msgTypes,
	}.Build()
	File_proto_informer_proto = out.File
	file_proto_informer_proto_rawDesc = nil
	file_proto_informer_proto_goTypes = nil
	file_proto_informer_proto_depIdxs = nil
}