// Code generated by protoc-gen-go. DO NOT EDIT.
// source: agent.proto

/*
Package msg is a generated protocol buffer package.

It is generated from these files:
	agent.proto

It has these top-level messages:
	Packet
	AllActionConfigs
	Hello
	KeepAlive
	NodeInfo
	ScheduleActionConfig
	UnscheduleActionConfig
	RunActionConfig
	ActionConfig
	Result
*/
package msg

import proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Packet_Sender int32

const (
	Packet_AGENT      Packet_Sender = 0
	Packet_CONTROLLER Packet_Sender = 1
)

var Packet_Sender_name = map[int32]string{
	0: "AGENT",
	1: "CONTROLLER",
}
var Packet_Sender_value = map[string]int32{
	"AGENT":      0,
	"CONTROLLER": 1,
}

func (x Packet_Sender) String() string {
	return proto.EnumName(Packet_Sender_name, int32(x))
}
func (Packet_Sender) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Packet struct {
	Sender Packet_Sender `protobuf:"varint,1,opt,name=sender,enum=msg.Packet_Sender" json:"sender,omitempty"`
	Code   int32         `protobuf:"varint,2,opt,name=code" json:"code,omitempty"`
	Msg    string        `protobuf:"bytes,3,opt,name=msg" json:"msg,omitempty"`
	// Types that are valid to be assigned to Packet:
	//	*Packet_Hello
	//	*Packet_Schedule
	//	*Packet_Result
	//	*Packet_Keepalive
	//	*Packet_Unschedule
	//	*Packet_Run
	//	*Packet_AllConfigs
	Packet isPacket_Packet `protobuf_oneof:"packet"`
}

func (m *Packet) Reset()                    { *m = Packet{} }
func (m *Packet) String() string            { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()               {}
func (*Packet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isPacket_Packet interface {
	isPacket_Packet()
}

type Packet_Hello struct {
	Hello *Hello `protobuf:"bytes,4,opt,name=hello,oneof"`
}
type Packet_Schedule struct {
	Schedule *ScheduleActionConfig `protobuf:"bytes,5,opt,name=schedule,oneof"`
}
type Packet_Result struct {
	Result *Result `protobuf:"bytes,6,opt,name=result,oneof"`
}
type Packet_Keepalive struct {
	Keepalive *KeepAlive `protobuf:"bytes,7,opt,name=keepalive,oneof"`
}
type Packet_Unschedule struct {
	Unschedule *UnscheduleActionConfig `protobuf:"bytes,8,opt,name=unschedule,oneof"`
}
type Packet_Run struct {
	Run *RunActionConfig `protobuf:"bytes,9,opt,name=run,oneof"`
}
type Packet_AllConfigs struct {
	AllConfigs *AllActionConfigs `protobuf:"bytes,10,opt,name=allConfigs,oneof"`
}

func (*Packet_Hello) isPacket_Packet()      {}
func (*Packet_Schedule) isPacket_Packet()   {}
func (*Packet_Result) isPacket_Packet()     {}
func (*Packet_Keepalive) isPacket_Packet()  {}
func (*Packet_Unschedule) isPacket_Packet() {}
func (*Packet_Run) isPacket_Packet()        {}
func (*Packet_AllConfigs) isPacket_Packet() {}

func (m *Packet) GetPacket() isPacket_Packet {
	if m != nil {
		return m.Packet
	}
	return nil
}

func (m *Packet) GetSender() Packet_Sender {
	if m != nil {
		return m.Sender
	}
	return Packet_AGENT
}

func (m *Packet) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Packet) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *Packet) GetHello() *Hello {
	if x, ok := m.GetPacket().(*Packet_Hello); ok {
		return x.Hello
	}
	return nil
}

func (m *Packet) GetSchedule() *ScheduleActionConfig {
	if x, ok := m.GetPacket().(*Packet_Schedule); ok {
		return x.Schedule
	}
	return nil
}

func (m *Packet) GetResult() *Result {
	if x, ok := m.GetPacket().(*Packet_Result); ok {
		return x.Result
	}
	return nil
}

func (m *Packet) GetKeepalive() *KeepAlive {
	if x, ok := m.GetPacket().(*Packet_Keepalive); ok {
		return x.Keepalive
	}
	return nil
}

func (m *Packet) GetUnschedule() *UnscheduleActionConfig {
	if x, ok := m.GetPacket().(*Packet_Unschedule); ok {
		return x.Unschedule
	}
	return nil
}

func (m *Packet) GetRun() *RunActionConfig {
	if x, ok := m.GetPacket().(*Packet_Run); ok {
		return x.Run
	}
	return nil
}

func (m *Packet) GetAllConfigs() *AllActionConfigs {
	if x, ok := m.GetPacket().(*Packet_AllConfigs); ok {
		return x.AllConfigs
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Packet) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Packet_OneofMarshaler, _Packet_OneofUnmarshaler, _Packet_OneofSizer, []interface{}{
		(*Packet_Hello)(nil),
		(*Packet_Schedule)(nil),
		(*Packet_Result)(nil),
		(*Packet_Keepalive)(nil),
		(*Packet_Unschedule)(nil),
		(*Packet_Run)(nil),
		(*Packet_AllConfigs)(nil),
	}
}

func _Packet_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Packet)
	// packet
	switch x := m.Packet.(type) {
	case *Packet_Hello:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Hello); err != nil {
			return err
		}
	case *Packet_Schedule:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Schedule); err != nil {
			return err
		}
	case *Packet_Result:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Result); err != nil {
			return err
		}
	case *Packet_Keepalive:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Keepalive); err != nil {
			return err
		}
	case *Packet_Unschedule:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Unschedule); err != nil {
			return err
		}
	case *Packet_Run:
		b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Run); err != nil {
			return err
		}
	case *Packet_AllConfigs:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AllConfigs); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Packet.Packet has unexpected type %T", x)
	}
	return nil
}

func _Packet_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Packet)
	switch tag {
	case 4: // packet.hello
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Hello)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_Hello{msg}
		return true, err
	case 5: // packet.schedule
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ScheduleActionConfig)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_Schedule{msg}
		return true, err
	case 6: // packet.result
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Result)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_Result{msg}
		return true, err
	case 7: // packet.keepalive
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(KeepAlive)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_Keepalive{msg}
		return true, err
	case 8: // packet.unschedule
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UnscheduleActionConfig)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_Unschedule{msg}
		return true, err
	case 9: // packet.run
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RunActionConfig)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_Run{msg}
		return true, err
	case 10: // packet.allConfigs
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(AllActionConfigs)
		err := b.DecodeMessage(msg)
		m.Packet = &Packet_AllConfigs{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Packet_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Packet)
	// packet
	switch x := m.Packet.(type) {
	case *Packet_Hello:
		s := proto.Size(x.Hello)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Packet_Schedule:
		s := proto.Size(x.Schedule)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Packet_Result:
		s := proto.Size(x.Result)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Packet_Keepalive:
		s := proto.Size(x.Keepalive)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Packet_Unschedule:
		s := proto.Size(x.Unschedule)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Packet_Run:
		s := proto.Size(x.Run)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Packet_AllConfigs:
		s := proto.Size(x.AllConfigs)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type AllActionConfigs struct {
	Configs []*ScheduleActionConfig `protobuf:"bytes,1,rep,name=configs" json:"configs,omitempty"`
}

func (m *AllActionConfigs) Reset()                    { *m = AllActionConfigs{} }
func (m *AllActionConfigs) String() string            { return proto.CompactTextString(m) }
func (*AllActionConfigs) ProtoMessage()               {}
func (*AllActionConfigs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AllActionConfigs) GetConfigs() []*ScheduleActionConfig {
	if m != nil {
		return m.Configs
	}
	return nil
}

type Hello struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *Hello) Reset()                    { *m = Hello{} }
func (m *Hello) String() string            { return proto.CompactTextString(m) }
func (*Hello) ProtoMessage()               {}
func (*Hello) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Hello) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Hello) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type KeepAlive struct {
	Info      *NodeInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	Timestamp int64     `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *KeepAlive) Reset()                    { *m = KeepAlive{} }
func (m *KeepAlive) String() string            { return proto.CompactTextString(m) }
func (*KeepAlive) ProtoMessage()               {}
func (*KeepAlive) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *KeepAlive) GetInfo() *NodeInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *KeepAlive) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type NodeInfo struct {
	Id           string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Os           string `protobuf:"bytes,2,opt,name=os" json:"os,omitempty"`
	Kernel       string `protobuf:"bytes,3,opt,name=kernel" json:"kernel,omitempty"`
	Platform     string `protobuf:"bytes,4,opt,name=platform" json:"platform,omitempty"`
	Ip           string `protobuf:"bytes,5,opt,name=ip" json:"ip,omitempty"`
	CoreCount    int32  `protobuf:"varint,6,opt,name=coreCount" json:"coreCount,omitempty"`
	Hostname     string `protobuf:"bytes,7,opt,name=hostname" json:"hostname,omitempty"`
	MajorVersion int32  `protobuf:"varint,8,opt,name=majorVersion" json:"majorVersion,omitempty"`
	MinorVersion int32  `protobuf:"varint,9,opt,name=minorVersion" json:"minorVersion,omitempty"`
	PatchVersion int32  `protobuf:"varint,10,opt,name=patchVersion" json:"patchVersion,omitempty"`
	StartTime    int64  `protobuf:"varint,11,opt,name=startTime" json:"startTime,omitempty"`
}

func (m *NodeInfo) Reset()                    { *m = NodeInfo{} }
func (m *NodeInfo) String() string            { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()               {}
func (*NodeInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *NodeInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NodeInfo) GetOs() string {
	if m != nil {
		return m.Os
	}
	return ""
}

func (m *NodeInfo) GetKernel() string {
	if m != nil {
		return m.Kernel
	}
	return ""
}

func (m *NodeInfo) GetPlatform() string {
	if m != nil {
		return m.Platform
	}
	return ""
}

func (m *NodeInfo) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *NodeInfo) GetCoreCount() int32 {
	if m != nil {
		return m.CoreCount
	}
	return 0
}

func (m *NodeInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *NodeInfo) GetMajorVersion() int32 {
	if m != nil {
		return m.MajorVersion
	}
	return 0
}

func (m *NodeInfo) GetMinorVersion() int32 {
	if m != nil {
		return m.MinorVersion
	}
	return 0
}

func (m *NodeInfo) GetPatchVersion() int32 {
	if m != nil {
		return m.PatchVersion
	}
	return 0
}

func (m *NodeInfo) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

type ScheduleActionConfig struct {
	Config   *ActionConfig `protobuf:"bytes,1,opt,name=config" json:"config,omitempty"`
	Schedule string        `protobuf:"bytes,2,opt,name=schedule" json:"schedule,omitempty"`
}

func (m *ScheduleActionConfig) Reset()                    { *m = ScheduleActionConfig{} }
func (m *ScheduleActionConfig) String() string            { return proto.CompactTextString(m) }
func (*ScheduleActionConfig) ProtoMessage()               {}
func (*ScheduleActionConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ScheduleActionConfig) GetConfig() *ActionConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *ScheduleActionConfig) GetSchedule() string {
	if m != nil {
		return m.Schedule
	}
	return ""
}

type UnscheduleActionConfig struct {
	Config *ActionConfig `protobuf:"bytes,1,opt,name=config" json:"config,omitempty"`
}

func (m *UnscheduleActionConfig) Reset()                    { *m = UnscheduleActionConfig{} }
func (m *UnscheduleActionConfig) String() string            { return proto.CompactTextString(m) }
func (*UnscheduleActionConfig) ProtoMessage()               {}
func (*UnscheduleActionConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *UnscheduleActionConfig) GetConfig() *ActionConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

type RunActionConfig struct {
	Config  *ActionConfig `protobuf:"bytes,1,opt,name=config" json:"config,omitempty"`
	Runtime int64         `protobuf:"varint,2,opt,name=runtime" json:"runtime,omitempty"`
}

func (m *RunActionConfig) Reset()                    { *m = RunActionConfig{} }
func (m *RunActionConfig) String() string            { return proto.CompactTextString(m) }
func (*RunActionConfig) ProtoMessage()               {}
func (*RunActionConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *RunActionConfig) GetConfig() *ActionConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *RunActionConfig) GetRuntime() int64 {
	if m != nil {
		return m.Runtime
	}
	return 0
}

type ActionConfig struct {
	Id         string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Action     int32             `protobuf:"varint,2,opt,name=action" json:"action,omitempty"`
	Parameters map[string]string `protobuf:"bytes,3,rep,name=parameters" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ActionConfig) Reset()                    { *m = ActionConfig{} }
func (m *ActionConfig) String() string            { return proto.CompactTextString(m) }
func (*ActionConfig) ProtoMessage()               {}
func (*ActionConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ActionConfig) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ActionConfig) GetAction() int32 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *ActionConfig) GetParameters() map[string]string {
	if m != nil {
		return m.Parameters
	}
	return nil
}

type Result struct {
	Id        string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ConfigId  string            `protobuf:"bytes,2,opt,name=configId" json:"configId,omitempty"`
	DataCode  int32             `protobuf:"varint,3,opt,name=dataCode" json:"dataCode,omitempty"`
	Data      []byte            `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Tags      map[string]string `protobuf:"bytes,5,rep,name=tags" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Timestamp int64             `protobuf:"varint,6,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Result) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Result) GetConfigId() string {
	if m != nil {
		return m.ConfigId
	}
	return ""
}

func (m *Result) GetDataCode() int32 {
	if m != nil {
		return m.DataCode
	}
	return 0
}

func (m *Result) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Result) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Result) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*Packet)(nil), "msg.Packet")
	proto.RegisterType((*AllActionConfigs)(nil), "msg.AllActionConfigs")
	proto.RegisterType((*Hello)(nil), "msg.Hello")
	proto.RegisterType((*KeepAlive)(nil), "msg.KeepAlive")
	proto.RegisterType((*NodeInfo)(nil), "msg.NodeInfo")
	proto.RegisterType((*ScheduleActionConfig)(nil), "msg.ScheduleActionConfig")
	proto.RegisterType((*UnscheduleActionConfig)(nil), "msg.UnscheduleActionConfig")
	proto.RegisterType((*RunActionConfig)(nil), "msg.RunActionConfig")
	proto.RegisterType((*ActionConfig)(nil), "msg.ActionConfig")
	proto.RegisterType((*Result)(nil), "msg.Result")
	proto.RegisterEnum("msg.Packet_Sender", Packet_Sender_name, Packet_Sender_value)
}

func init() { proto.RegisterFile("agent.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 766 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xc1, 0x8e, 0xdb, 0x36,
	0x10, 0xb5, 0x2c, 0x4b, 0xb6, 0xc6, 0x5b, 0xc7, 0x25, 0x36, 0x0b, 0xd5, 0xed, 0xc1, 0x51, 0x51,
	0xc0, 0xe9, 0xc1, 0x87, 0xcd, 0x61, 0x8b, 0x02, 0x41, 0xe1, 0x1a, 0x8b, 0x38, 0xe8, 0x62, 0x13,
	0x30, 0x6e, 0x6e, 0x3d, 0xb0, 0x16, 0x6d, 0xab, 0x96, 0x48, 0x81, 0xa4, 0x52, 0xe4, 0xb7, 0x7a,
	0xea, 0xdf, 0xf4, 0xd6, 0xef, 0x28, 0x38, 0xa2, 0x64, 0xaf, 0x13, 0x14, 0xd8, 0x9b, 0x66, 0xe6,
	0xbd, 0xe1, 0x90, 0xf3, 0x66, 0x04, 0x43, 0xb6, 0xe3, 0xc2, 0xcc, 0x4b, 0x25, 0x8d, 0x24, 0x7e,
	0xa1, 0x77, 0xc9, 0x3f, 0x3e, 0x84, 0x6f, 0xd9, 0xe6, 0xc0, 0x0d, 0xf9, 0x1e, 0x42, 0xcd, 0x45,
	0xca, 0x55, 0xec, 0x4d, 0xbd, 0xd9, 0xe8, 0x9a, 0xcc, 0x0b, 0xbd, 0x9b, 0xd7, 0xc1, 0xf9, 0x3b,
	0x8c, 0x50, 0x87, 0x20, 0x04, 0x7a, 0x1b, 0x99, 0xf2, 0xb8, 0x3b, 0xf5, 0x66, 0x01, 0xc5, 0x6f,
	0x32, 0x06, 0x9b, 0x31, 0xf6, 0xa7, 0xde, 0x2c, 0xa2, 0xf6, 0x93, 0x24, 0x10, 0xec, 0x79, 0x9e,
	0xcb, 0xb8, 0x37, 0xf5, 0x66, 0xc3, 0x6b, 0xc0, 0x84, 0x2b, 0xeb, 0x59, 0x75, 0x68, 0x1d, 0x22,
	0x37, 0x30, 0xd0, 0x9b, 0x3d, 0x4f, 0xab, 0x9c, 0xc7, 0x01, 0xc2, 0xbe, 0x42, 0xd8, 0x3b, 0xe7,
	0x5c, 0x6c, 0x4c, 0x26, 0xc5, 0x52, 0x8a, 0x6d, 0xb6, 0x5b, 0x75, 0x68, 0x0b, 0x26, 0xdf, 0x41,
	0xa8, 0xb8, 0xae, 0x72, 0x13, 0x87, 0x48, 0x1b, 0x22, 0x8d, 0xa2, 0x6b, 0xd5, 0xa1, 0x2e, 0x48,
	0xe6, 0x10, 0x1d, 0x38, 0x2f, 0x59, 0x9e, 0x7d, 0xe0, 0x71, 0x1f, 0x91, 0x23, 0x44, 0xfe, 0xc2,
	0x79, 0xb9, 0xb0, 0xde, 0x55, 0x87, 0x1e, 0x21, 0xe4, 0x25, 0x40, 0x25, 0xda, 0x8a, 0x06, 0x48,
	0xf8, 0x1a, 0x09, 0xbf, 0xb6, 0xee, 0xb3, 0x9a, 0x4e, 0x08, 0x64, 0x06, 0xbe, 0xaa, 0x44, 0x1c,
	0x21, 0xef, 0xb2, 0x2e, 0xa9, 0x12, 0x67, 0x04, 0x0b, 0x21, 0x37, 0x00, 0x2c, 0xcf, 0x6b, 0x9f,
	0x8e, 0x01, 0x09, 0x4f, 0x91, 0xb0, 0xc8, 0xf3, 0x53, 0x82, 0xb6, 0x47, 0x1c, 0xa1, 0xc9, 0xb7,
	0x10, 0xd6, 0xdd, 0x20, 0x11, 0x04, 0x8b, 0x57, 0xb7, 0xf7, 0xeb, 0x71, 0x87, 0x8c, 0x00, 0x96,
	0x6f, 0xee, 0xd7, 0xf4, 0xcd, 0xdd, 0xdd, 0x2d, 0x1d, 0x7b, 0x3f, 0x0f, 0x20, 0x2c, 0xb1, 0x73,
	0xc9, 0x2b, 0x18, 0x9f, 0x27, 0x24, 0x2f, 0xa0, 0xbf, 0x71, 0x07, 0x7b, 0x53, 0xff, 0x7f, 0xdf,
	0x9c, 0x36, 0xc8, 0xe4, 0x27, 0x08, 0xb0, 0x77, 0x64, 0x02, 0x83, 0x4a, 0x73, 0x25, 0x58, 0xc1,
	0x51, 0x2a, 0x11, 0x6d, 0x6d, 0x1b, 0x2b, 0x99, 0xd6, 0x7f, 0x4a, 0x95, 0xa2, 0x38, 0x22, 0xda,
	0xda, 0xc9, 0x1d, 0x44, 0xed, 0xa3, 0x93, 0x67, 0xd0, 0xcb, 0xc4, 0x56, 0x62, 0x82, 0xe1, 0xf5,
	0x17, 0x78, 0xfe, 0xbd, 0x4c, 0xf9, 0x6b, 0xb1, 0x95, 0x14, 0x43, 0xe4, 0x1b, 0x88, 0x4c, 0x56,
	0x70, 0x6d, 0x58, 0x51, 0x62, 0x32, 0x9f, 0x1e, 0x1d, 0xc9, 0x5f, 0x5d, 0x18, 0x34, 0x04, 0x32,
	0x82, 0x6e, 0x96, 0xba, 0x62, 0xba, 0x59, 0x6a, 0x6d, 0xa9, 0x5d, 0x01, 0x5d, 0xa9, 0xc9, 0x15,
	0x84, 0x07, 0xae, 0x04, 0xcf, 0x9d, 0x3c, 0x9d, 0x85, 0xe5, 0xe6, 0xcc, 0x6c, 0xa5, 0x2a, 0x50,
	0xa4, 0xb6, 0x5c, 0x67, 0x63, 0xce, 0x12, 0x35, 0x69, 0x73, 0x96, 0xb6, 0x9c, 0x8d, 0x54, 0x7c,
	0x29, 0x2b, 0x51, 0x6b, 0x2e, 0xa0, 0x47, 0x87, 0xcd, 0xb4, 0x97, 0xda, 0xe0, 0xa3, 0xf4, 0xeb,
	0x4c, 0x8d, 0x4d, 0x12, 0xb8, 0x28, 0xd8, 0x1f, 0x52, 0xbd, 0xe7, 0x4a, 0x67, 0x52, 0xa0, 0xaa,
	0x02, 0xfa, 0xc0, 0x87, 0x98, 0x4c, 0x1c, 0x31, 0x91, 0xc3, 0x9c, 0xf8, 0x2c, 0xa6, 0x64, 0x66,
	0xb3, 0x6f, 0x30, 0x50, 0x63, 0x4e, 0x7d, 0xb6, 0x4a, 0x6d, 0x98, 0x32, 0xeb, 0xac, 0xe0, 0xf1,
	0xb0, 0x7e, 0xb4, 0xd6, 0x91, 0xfc, 0x06, 0x97, 0x9f, 0x6b, 0x32, 0x79, 0x0e, 0x61, 0xdd, 0x66,
	0xd7, 0x8f, 0x2f, 0x6b, 0x21, 0x9e, 0xea, 0xc0, 0x01, 0xec, 0x45, 0xdb, 0xf1, 0x70, 0x1d, 0x6e,
	0xec, 0x64, 0x09, 0x57, 0x9f, 0x9f, 0x92, 0x47, 0x1c, 0x90, 0xbc, 0x87, 0x27, 0x67, 0x23, 0xf3,
	0x98, 0xf2, 0x62, 0xe8, 0xab, 0x4a, 0x58, 0x99, 0x38, 0xc9, 0x34, 0x66, 0xf2, 0xb7, 0x07, 0x17,
	0x0f, 0xb2, 0x9e, 0x8b, 0xe6, 0x0a, 0x42, 0x86, 0x71, 0xb7, 0xd6, 0x9c, 0x45, 0x16, 0x00, 0x25,
	0x53, 0xac, 0xe0, 0x86, 0x2b, 0x1d, 0xfb, 0x38, 0x30, 0xcf, 0x3e, 0xa9, 0x60, 0xfe, 0xb6, 0xc5,
	0xdc, 0x0a, 0xa3, 0x3e, 0xd2, 0x13, 0xd2, 0xe4, 0x25, 0x3c, 0x39, 0x0b, 0xdb, 0x75, 0x79, 0xe0,
	0x1f, 0xdd, 0xf1, 0xf6, 0x93, 0x5c, 0x42, 0xf0, 0x81, 0xe5, 0x55, 0xf3, 0xac, 0xb5, 0xf1, 0x63,
	0xf7, 0x07, 0x2f, 0xf9, 0xd7, 0x83, 0xb0, 0xde, 0x6c, 0x9f, 0x14, 0x3d, 0x81, 0x41, 0x7d, 0xf3,
	0xd7, 0xed, 0xc0, 0x35, 0xb6, 0x8d, 0xa5, 0xcc, 0xb0, 0xa5, 0xdd, 0xd4, 0x3e, 0x5e, 0xa9, 0xb5,
	0xed, 0x06, 0xb7, 0xdf, 0xa8, 0xfa, 0x0b, 0x8a, 0xdf, 0xe4, 0x39, 0xf4, 0x0c, 0xdb, 0xe9, 0x38,
	0xc0, 0x2b, 0x3e, 0x3d, 0x59, 0xa8, 0xf3, 0x35, 0xdb, 0xb9, 0x6b, 0x21, 0xe4, 0xe1, 0x6c, 0x86,
	0x67, 0xb3, 0x39, 0xb9, 0x81, 0xa8, 0x25, 0x3c, 0xe6, 0xa2, 0xbf, 0x87, 0xf8, 0x6b, 0x7a, 0xf1,
	0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0f, 0xd4, 0x5d, 0x1b, 0xa9, 0x06, 0x00, 0x00,
}
