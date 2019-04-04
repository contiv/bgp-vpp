// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pod.proto

package pod

/*
Package pod defines data model for Kubernetes Pod.
*/

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Protocol defines network protocols supported for container ports.
type Pod_Container_Port_Protocol int32

const (
	Pod_Container_Port_TCP Pod_Container_Port_Protocol = 0
	Pod_Container_Port_UDP Pod_Container_Port_Protocol = 1
)

var Pod_Container_Port_Protocol_name = map[int32]string{
	0: "TCP",
	1: "UDP",
}
var Pod_Container_Port_Protocol_value = map[string]int32{
	"TCP": 0,
	"UDP": 1,
}

func (x Pod_Container_Port_Protocol) String() string {
	return proto.EnumName(Pod_Container_Port_Protocol_name, int32(x))
}
func (Pod_Container_Port_Protocol) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_pod_cb1d4279c891eed3, []int{0, 1, 0, 0}
}

// Pod is a collection of containers that can run on a host.
// This resource is created by clients and scheduled onto hosts.
type Pod struct {
	// Name of the pod unique within the namespace.
	// Cannot be updated.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Namespace the pod is inserted into.
	// An empty namespace is equivalent to the "default" namespace, but "default"
	// is the canonical representation.
	// Cannot be updated.
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// A list of labels attached to this pod.
	// +optional
	Label []*Pod_Label `protobuf:"bytes,3,rep,name=label" json:"label,omitempty"`
	// IP address allocated to the pod. Routable at least within the cluster.
	// Empty if not yet allocated.
	// +optional
	IpAddress string `protobuf:"bytes,4,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	// IP address of the host to which the pod is assigned.
	// Empty if not yet scheduled.
	// +optional
	HostIpAddress string `protobuf:"bytes,5,opt,name=host_ip_address,json=hostIpAddress,proto3" json:"host_ip_address,omitempty"`
	// List of containers belonging to the pod.
	// Containers cannot currently be added or removed.
	// There must be at least one container in a Pod.
	// Cannot be updated.
	Container            []*Pod_Container `protobuf:"bytes,6,rep,name=container" json:"container,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Pod) Reset()         { *m = Pod{} }
func (m *Pod) String() string { return proto.CompactTextString(m) }
func (*Pod) ProtoMessage()    {}
func (*Pod) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_cb1d4279c891eed3, []int{0}
}
func (m *Pod) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pod.Unmarshal(m, b)
}
func (m *Pod) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pod.Marshal(b, m, deterministic)
}
func (dst *Pod) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pod.Merge(dst, src)
}
func (m *Pod) XXX_Size() int {
	return xxx_messageInfo_Pod.Size(m)
}
func (m *Pod) XXX_DiscardUnknown() {
	xxx_messageInfo_Pod.DiscardUnknown(m)
}

var xxx_messageInfo_Pod proto.InternalMessageInfo

func (m *Pod) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Pod) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Pod) GetLabel() []*Pod_Label {
	if m != nil {
		return m.Label
	}
	return nil
}

func (m *Pod) GetIpAddress() string {
	if m != nil {
		return m.IpAddress
	}
	return ""
}

func (m *Pod) GetHostIpAddress() string {
	if m != nil {
		return m.HostIpAddress
	}
	return ""
}

func (m *Pod) GetContainer() []*Pod_Container {
	if m != nil {
		return m.Container
	}
	return nil
}

// Label is a key/value pair attached to an object (pod in this case).
// Labels are used to organize and to select subsets of objects.
type Pod_Label struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pod_Label) Reset()         { *m = Pod_Label{} }
func (m *Pod_Label) String() string { return proto.CompactTextString(m) }
func (*Pod_Label) ProtoMessage()    {}
func (*Pod_Label) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_cb1d4279c891eed3, []int{0, 0}
}
func (m *Pod_Label) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pod_Label.Unmarshal(m, b)
}
func (m *Pod_Label) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pod_Label.Marshal(b, m, deterministic)
}
func (dst *Pod_Label) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pod_Label.Merge(dst, src)
}
func (m *Pod_Label) XXX_Size() int {
	return xxx_messageInfo_Pod_Label.Size(m)
}
func (m *Pod_Label) XXX_DiscardUnknown() {
	xxx_messageInfo_Pod_Label.DiscardUnknown(m)
}

var xxx_messageInfo_Pod_Label proto.InternalMessageInfo

func (m *Pod_Label) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Pod_Label) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// A single application container run within a pod.
type Pod_Container struct {
	// Name of the container.
	// Each container in a pod has a unique name.
	// Cannot be updated.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// List of ports to expose from the container. Exposing a port here gives
	// the system additional information about the network connections a
	// container uses, but it is primarily informational. Not specifying a port
	// here DOES NOT prevent that port from being exposed. Any port which is
	// listening on the default "0.0.0.0" address inside a container will be
	// accessible from the network.
	// Cannot be updated.
	// +optional
	Port                 []*Pod_Container_Port `protobuf:"bytes,2,rep,name=port" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Pod_Container) Reset()         { *m = Pod_Container{} }
func (m *Pod_Container) String() string { return proto.CompactTextString(m) }
func (*Pod_Container) ProtoMessage()    {}
func (*Pod_Container) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_cb1d4279c891eed3, []int{0, 1}
}
func (m *Pod_Container) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pod_Container.Unmarshal(m, b)
}
func (m *Pod_Container) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pod_Container.Marshal(b, m, deterministic)
}
func (dst *Pod_Container) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pod_Container.Merge(dst, src)
}
func (m *Pod_Container) XXX_Size() int {
	return xxx_messageInfo_Pod_Container.Size(m)
}
func (m *Pod_Container) XXX_DiscardUnknown() {
	xxx_messageInfo_Pod_Container.DiscardUnknown(m)
}

var xxx_messageInfo_Pod_Container proto.InternalMessageInfo

func (m *Pod_Container) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Pod_Container) GetPort() []*Pod_Container_Port {
	if m != nil {
		return m.Port
	}
	return nil
}

// Port represents a network port in a single container.
type Pod_Container_Port struct {
	// An IANA_SVC_NAME formatted port name, unique within the pod.
	// The name can be referred to by services, policies, ...
	// +optional
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Port number to expose on the host.
	// The port number is in the range: 0 < x < 65536.
	// If pod is in the host network namespace, this must match container_port.
	// Most containers do not need this.
	// +optional
	HostPort int32 `protobuf:"varint,2,opt,name=host_port,json=hostPort,proto3" json:"host_port,omitempty"`
	// Port number to expose on the pod's IP address.
	// The port number is in the range: 0 < x < 65536.
	ContainerPort int32 `protobuf:"varint,3,opt,name=container_port,json=containerPort,proto3" json:"container_port,omitempty"`
	// Protocol for port. Must be UDP or TCP.
	// Defaults to "TCP".
	// +optional
	Protocol Pod_Container_Port_Protocol `protobuf:"varint,4,opt,name=protocol,proto3,enum=pod.Pod_Container_Port_Protocol" json:"protocol,omitempty"`
	// What host IP to bind the external port to.
	// +optional
	HostIpAddress        string   `protobuf:"bytes,5,opt,name=host_ip_address,json=hostIpAddress,proto3" json:"host_ip_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pod_Container_Port) Reset()         { *m = Pod_Container_Port{} }
func (m *Pod_Container_Port) String() string { return proto.CompactTextString(m) }
func (*Pod_Container_Port) ProtoMessage()    {}
func (*Pod_Container_Port) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_cb1d4279c891eed3, []int{0, 1, 0}
}
func (m *Pod_Container_Port) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pod_Container_Port.Unmarshal(m, b)
}
func (m *Pod_Container_Port) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pod_Container_Port.Marshal(b, m, deterministic)
}
func (dst *Pod_Container_Port) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pod_Container_Port.Merge(dst, src)
}
func (m *Pod_Container_Port) XXX_Size() int {
	return xxx_messageInfo_Pod_Container_Port.Size(m)
}
func (m *Pod_Container_Port) XXX_DiscardUnknown() {
	xxx_messageInfo_Pod_Container_Port.DiscardUnknown(m)
}

var xxx_messageInfo_Pod_Container_Port proto.InternalMessageInfo

func (m *Pod_Container_Port) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Pod_Container_Port) GetHostPort() int32 {
	if m != nil {
		return m.HostPort
	}
	return 0
}

func (m *Pod_Container_Port) GetContainerPort() int32 {
	if m != nil {
		return m.ContainerPort
	}
	return 0
}

func (m *Pod_Container_Port) GetProtocol() Pod_Container_Port_Protocol {
	if m != nil {
		return m.Protocol
	}
	return Pod_Container_Port_TCP
}

func (m *Pod_Container_Port) GetHostIpAddress() string {
	if m != nil {
		return m.HostIpAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*Pod)(nil), "pod.Pod")
	proto.RegisterType((*Pod_Label)(nil), "pod.Pod.Label")
	proto.RegisterType((*Pod_Container)(nil), "pod.Pod.Container")
	proto.RegisterType((*Pod_Container_Port)(nil), "pod.Pod.Container.Port")
	proto.RegisterEnum("pod.Pod_Container_Port_Protocol", Pod_Container_Port_Protocol_name, Pod_Container_Port_Protocol_value)
}

func init() { proto.RegisterFile("pod.proto", fileDescriptor_pod_cb1d4279c891eed3) }

var fileDescriptor_pod_cb1d4279c891eed3 = []byte{
	// 322 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xdf, 0x4a, 0xc3, 0x30,
	0x14, 0xc6, 0xed, 0xd2, 0xce, 0xe5, 0xc8, 0xe6, 0x08, 0x82, 0x61, 0x4e, 0x18, 0x43, 0x65, 0x20,
	0x54, 0x99, 0xb7, 0xde, 0xc8, 0xbc, 0x11, 0xbc, 0x08, 0x41, 0xaf, 0x47, 0xb6, 0x04, 0x1c, 0xd6,
	0x25, 0xb4, 0x55, 0xf0, 0xb1, 0xbc, 0xf6, 0x6d, 0x7c, 0x12, 0xc9, 0x69, 0x9b, 0x09, 0x4e, 0xf0,
	0xaa, 0x27, 0xdf, 0xf9, 0x9d, 0xef, 0xfc, 0x29, 0x50, 0x67, 0x75, 0xea, 0x72, 0x5b, 0x5a, 0x46,
	0x9c, 0xd5, 0xe3, 0xcf, 0x18, 0x88, 0xb0, 0x9a, 0x31, 0x88, 0xd7, 0xea, 0xc5, 0xf0, 0x68, 0x14,
	0x4d, 0xa8, 0xc4, 0x98, 0x0d, 0x81, 0xfa, 0x6f, 0xe1, 0xd4, 0xd2, 0xf0, 0x16, 0x26, 0x36, 0x02,
	0x3b, 0x81, 0x24, 0x53, 0x0b, 0x93, 0x71, 0x32, 0x22, 0x93, 0xbd, 0x69, 0x2f, 0xf5, 0xce, 0xc2,
	0xea, 0xf4, 0xde, 0xab, 0xb2, 0x4a, 0xb2, 0x63, 0x80, 0x95, 0x9b, 0x2b, 0xad, 0x73, 0x53, 0x14,
	0x3c, 0xae, 0x4c, 0x56, 0xee, 0xa6, 0x12, 0xd8, 0x19, 0xec, 0x3f, 0xd9, 0xa2, 0x9c, 0xff, 0x60,
	0x12, 0x64, 0xba, 0x5e, 0xbe, 0x0b, 0xdc, 0x25, 0xd0, 0xa5, 0x5d, 0x97, 0x6a, 0xb5, 0x36, 0x39,
	0x6f, 0x63, 0x43, 0x16, 0x1a, 0xce, 0x9a, 0x8c, 0xdc, 0x40, 0x83, 0x0b, 0x48, 0x70, 0x10, 0xd6,
	0x07, 0xf2, 0x6c, 0xde, 0xeb, 0xc5, 0x7c, 0xc8, 0x0e, 0x20, 0x79, 0x53, 0xd9, 0x6b, 0xb3, 0x53,
	0xf5, 0x18, 0x7c, 0xb4, 0x80, 0x06, 0xa7, 0xad, 0xf7, 0x38, 0x87, 0xd8, 0xd9, 0xbc, 0xe4, 0x2d,
	0xec, 0x7f, 0xf8, 0xbb, 0x7f, 0x2a, 0x6c, 0x5e, 0x4a, 0x84, 0x06, 0x5f, 0x11, 0xc4, 0xfe, 0xb9,
	0xd5, 0xe9, 0x08, 0x28, 0xae, 0x5d, 0xdb, 0x45, 0x93, 0x44, 0x76, 0xbc, 0x80, 0x05, 0xa7, 0xd0,
	0x0b, 0x6b, 0x54, 0x04, 0x41, 0xa2, 0x1b, 0x54, 0xc4, 0xae, 0xa1, 0x83, 0xff, 0x71, 0x69, 0x33,
	0xbc, 0x6b, 0x6f, 0x3a, 0xfa, 0x63, 0xa2, 0x54, 0xd4, 0x9c, 0x0c, 0x15, 0xff, 0x3d, 0xfc, 0x78,
	0x08, 0x9d, 0xa6, 0x9a, 0xed, 0x02, 0x79, 0x98, 0x89, 0xfe, 0x8e, 0x0f, 0x1e, 0x6f, 0x45, 0x3f,
	0x5a, 0xb4, 0xd1, 0xef, 0xea, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x6d, 0x78, 0xcf, 0x56, 0x02,
	0x00, 0x00,
}