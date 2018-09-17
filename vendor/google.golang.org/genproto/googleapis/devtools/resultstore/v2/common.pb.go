// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/devtools/resultstore/v2/common.proto

package resultstore

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// These correspond to the prefix of the rule name. Eg cc_test has language CC.
type Language int32

const (
	// Language unspecified or not listed here.
	Language_LANGUAGE_UNSPECIFIED Language = 0
	// Not related to any particular language
	Language_NONE Language = 1
	// Android
	Language_ANDROID Language = 2
	// ActionScript (Flash)
	Language_AS Language = 3
	// C++ or C
	Language_CC Language = 4
	// Cascading-Style-Sheets
	Language_CSS Language = 5
	// Dart
	Language_DART Language = 6
	// Go
	Language_GO Language = 7
	// Google-Web-Toolkit
	Language_GWT Language = 8
	// Haskell
	Language_HASKELL Language = 9
	// Java
	Language_JAVA Language = 10
	// Javascript
	Language_JS Language = 11
	// Lisp
	Language_LISP Language = 12
	// Objective-C
	Language_OBJC Language = 13
	// Python
	Language_PY Language = 14
	// Shell (Typically Bash)
	Language_SH Language = 15
	// Swift
	Language_SWIFT Language = 16
	// Typescript
	Language_TS Language = 18
	// Webtesting
	Language_WEB Language = 19
)

var Language_name = map[int32]string{
	0:  "LANGUAGE_UNSPECIFIED",
	1:  "NONE",
	2:  "ANDROID",
	3:  "AS",
	4:  "CC",
	5:  "CSS",
	6:  "DART",
	7:  "GO",
	8:  "GWT",
	9:  "HASKELL",
	10: "JAVA",
	11: "JS",
	12: "LISP",
	13: "OBJC",
	14: "PY",
	15: "SH",
	16: "SWIFT",
	18: "TS",
	19: "WEB",
}

var Language_value = map[string]int32{
	"LANGUAGE_UNSPECIFIED": 0,
	"NONE":                 1,
	"ANDROID":              2,
	"AS":                   3,
	"CC":                   4,
	"CSS":                  5,
	"DART":                 6,
	"GO":                   7,
	"GWT":                  8,
	"HASKELL":              9,
	"JAVA":                 10,
	"JS":                   11,
	"LISP":                 12,
	"OBJC":                 13,
	"PY":                   14,
	"SH":                   15,
	"SWIFT":                16,
	"TS":                   18,
	"WEB":                  19,
}

func (x Language) String() string {
	return proto.EnumName(Language_name, int32(x))
}

func (Language) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ff56b05a77242216, []int{0}
}

// Status of a resource.
type Status int32

const (
	// The implicit default enum value. Should never be set.
	Status_STATUS_UNSPECIFIED Status = 0
	// Displays as "Building". Means the target is compiling, linking, etc.
	Status_BUILDING Status = 1
	// Displays as "Built". Means the target was built successfully.
	// If testing was requested, it should never reach this status: it should go
	// straight from BUILDING to TESTING.
	Status_BUILT Status = 2
	// Displays as "Broken". Means build failure such as compile error.
	Status_FAILED_TO_BUILD Status = 3
	// Displays as "Testing". Means the test is running.
	Status_TESTING Status = 4
	// Displays as "Passed". Means the test was run and passed.
	Status_PASSED Status = 5
	// Displays as "Failed". Means the test was run and failed.
	Status_FAILED Status = 6
	// Displays as "Timed out". Means the test didn't finish in time.
	Status_TIMED_OUT Status = 7
	// Displays as "Cancelled". Means the build or test was cancelled.
	// E.g. User hit control-C.
	Status_CANCELLED Status = 8
	// Displays as "Tool Failed". Means the build or test had internal tool
	// failure.
	Status_TOOL_FAILED Status = 9
	// Displays as "Incomplete". Means the build or test did not complete.  This
	// might happen when a build breakage or test failure causes the tool to stop
	// trying to build anything more or run any more tests, with the default
	// bazel --nokeep_going option or the --notest_keep_going option.
	Status_INCOMPLETE Status = 10
	// Displays as "Flaky". Means the aggregate status contains some runs that
	// were successful, and some that were not.
	Status_FLAKY Status = 11
	// Displays as "Unknown". Means the tool uploading to the server died
	// mid-upload or does not know the state.
	Status_UNKNOWN Status = 12
	// Displays as "Skipped". Means building and testing were skipped.
	// (E.g. Restricted to a different configuration.)
	Status_SKIPPED Status = 13
)

var Status_name = map[int32]string{
	0:  "STATUS_UNSPECIFIED",
	1:  "BUILDING",
	2:  "BUILT",
	3:  "FAILED_TO_BUILD",
	4:  "TESTING",
	5:  "PASSED",
	6:  "FAILED",
	7:  "TIMED_OUT",
	8:  "CANCELLED",
	9:  "TOOL_FAILED",
	10: "INCOMPLETE",
	11: "FLAKY",
	12: "UNKNOWN",
	13: "SKIPPED",
}

var Status_value = map[string]int32{
	"STATUS_UNSPECIFIED": 0,
	"BUILDING":           1,
	"BUILT":              2,
	"FAILED_TO_BUILD":    3,
	"TESTING":            4,
	"PASSED":             5,
	"FAILED":             6,
	"TIMED_OUT":          7,
	"CANCELLED":          8,
	"TOOL_FAILED":        9,
	"INCOMPLETE":         10,
	"FLAKY":              11,
	"UNKNOWN":            12,
	"SKIPPED":            13,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ff56b05a77242216, []int{1}
}

// Describes the status of a resource in both enum and string form.
// Only use description when conveying additional info not captured in the enum
// name.
type StatusAttributes struct {
	// Enum representation of the status.
	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=google.devtools.resultstore.v2.Status" json:"status,omitempty"`
	// A longer description about the status.
	Description          string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusAttributes) Reset()         { *m = StatusAttributes{} }
func (m *StatusAttributes) String() string { return proto.CompactTextString(m) }
func (*StatusAttributes) ProtoMessage()    {}
func (*StatusAttributes) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff56b05a77242216, []int{0}
}
func (m *StatusAttributes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusAttributes.Unmarshal(m, b)
}
func (m *StatusAttributes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusAttributes.Marshal(b, m, deterministic)
}
func (m *StatusAttributes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusAttributes.Merge(m, src)
}
func (m *StatusAttributes) XXX_Size() int {
	return xxx_messageInfo_StatusAttributes.Size(m)
}
func (m *StatusAttributes) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusAttributes.DiscardUnknown(m)
}

var xxx_messageInfo_StatusAttributes proto.InternalMessageInfo

func (m *StatusAttributes) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_STATUS_UNSPECIFIED
}

func (m *StatusAttributes) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

// A generic key-value property definition.
type Property struct {
	// The key.
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// The value.
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Property) Reset()         { *m = Property{} }
func (m *Property) String() string { return proto.CompactTextString(m) }
func (*Property) ProtoMessage()    {}
func (*Property) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff56b05a77242216, []int{1}
}
func (m *Property) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Property.Unmarshal(m, b)
}
func (m *Property) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Property.Marshal(b, m, deterministic)
}
func (m *Property) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Property.Merge(m, src)
}
func (m *Property) XXX_Size() int {
	return xxx_messageInfo_Property.Size(m)
}
func (m *Property) XXX_DiscardUnknown() {
	xxx_messageInfo_Property.DiscardUnknown(m)
}

var xxx_messageInfo_Property proto.InternalMessageInfo

func (m *Property) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Property) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// The timing of a particular Invocation, Action, etc. The start_time is
// specified, stop time can be calculated by adding duration to start_time.
type Timing struct {
	// The time the resource started running. This is in UTC Epoch time.
	StartTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// The duration for which the resource ran.
	Duration             *duration.Duration `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Timing) Reset()         { *m = Timing{} }
func (m *Timing) String() string { return proto.CompactTextString(m) }
func (*Timing) ProtoMessage()    {}
func (*Timing) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff56b05a77242216, []int{2}
}
func (m *Timing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Timing.Unmarshal(m, b)
}
func (m *Timing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Timing.Marshal(b, m, deterministic)
}
func (m *Timing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Timing.Merge(m, src)
}
func (m *Timing) XXX_Size() int {
	return xxx_messageInfo_Timing.Size(m)
}
func (m *Timing) XXX_DiscardUnknown() {
	xxx_messageInfo_Timing.DiscardUnknown(m)
}

var xxx_messageInfo_Timing proto.InternalMessageInfo

func (m *Timing) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *Timing) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

// Represents a dependency of a resource on another resource. This can be used
// to define a graph or a workflow paradigm through resources.
type Dependency struct {
	// The resource depended upon. It may be a Target, ConfiguredTarget, or
	// Action.
	//
	// Types that are valid to be assigned to Resource:
	//	*Dependency_Target
	//	*Dependency_ConfiguredTarget
	//	*Dependency_Action
	Resource isDependency_Resource `protobuf_oneof:"resource"`
	// A label describing this dependency.
	// The label "Root Cause" is handled specially. It is used to point to the
	// exact resource that caused a resource to fail.
	Label                string   `protobuf:"bytes,4,opt,name=label,proto3" json:"label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dependency) Reset()         { *m = Dependency{} }
func (m *Dependency) String() string { return proto.CompactTextString(m) }
func (*Dependency) ProtoMessage()    {}
func (*Dependency) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff56b05a77242216, []int{3}
}
func (m *Dependency) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dependency.Unmarshal(m, b)
}
func (m *Dependency) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dependency.Marshal(b, m, deterministic)
}
func (m *Dependency) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dependency.Merge(m, src)
}
func (m *Dependency) XXX_Size() int {
	return xxx_messageInfo_Dependency.Size(m)
}
func (m *Dependency) XXX_DiscardUnknown() {
	xxx_messageInfo_Dependency.DiscardUnknown(m)
}

var xxx_messageInfo_Dependency proto.InternalMessageInfo

type isDependency_Resource interface {
	isDependency_Resource()
}

type Dependency_Target struct {
	Target string `protobuf:"bytes,1,opt,name=target,proto3,oneof"`
}

type Dependency_ConfiguredTarget struct {
	ConfiguredTarget string `protobuf:"bytes,2,opt,name=configured_target,json=configuredTarget,proto3,oneof"`
}

type Dependency_Action struct {
	Action string `protobuf:"bytes,3,opt,name=action,proto3,oneof"`
}

func (*Dependency_Target) isDependency_Resource() {}

func (*Dependency_ConfiguredTarget) isDependency_Resource() {}

func (*Dependency_Action) isDependency_Resource() {}

func (m *Dependency) GetResource() isDependency_Resource {
	if m != nil {
		return m.Resource
	}
	return nil
}

func (m *Dependency) GetTarget() string {
	if x, ok := m.GetResource().(*Dependency_Target); ok {
		return x.Target
	}
	return ""
}

func (m *Dependency) GetConfiguredTarget() string {
	if x, ok := m.GetResource().(*Dependency_ConfiguredTarget); ok {
		return x.ConfiguredTarget
	}
	return ""
}

func (m *Dependency) GetAction() string {
	if x, ok := m.GetResource().(*Dependency_Action); ok {
		return x.Action
	}
	return ""
}

func (m *Dependency) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Dependency) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Dependency_OneofMarshaler, _Dependency_OneofUnmarshaler, _Dependency_OneofSizer, []interface{}{
		(*Dependency_Target)(nil),
		(*Dependency_ConfiguredTarget)(nil),
		(*Dependency_Action)(nil),
	}
}

func _Dependency_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Dependency)
	// resource
	switch x := m.Resource.(type) {
	case *Dependency_Target:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Target)
	case *Dependency_ConfiguredTarget:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.ConfiguredTarget)
	case *Dependency_Action:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Action)
	case nil:
	default:
		return fmt.Errorf("Dependency.Resource has unexpected type %T", x)
	}
	return nil
}

func _Dependency_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Dependency)
	switch tag {
	case 1: // resource.target
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Resource = &Dependency_Target{x}
		return true, err
	case 2: // resource.configured_target
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Resource = &Dependency_ConfiguredTarget{x}
		return true, err
	case 3: // resource.action
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Resource = &Dependency_Action{x}
		return true, err
	default:
		return false, nil
	}
}

func _Dependency_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Dependency)
	// resource
	switch x := m.Resource.(type) {
	case *Dependency_Target:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Target)))
		n += len(x.Target)
	case *Dependency_ConfiguredTarget:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.ConfiguredTarget)))
		n += len(x.ConfiguredTarget)
	case *Dependency_Action:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Action)))
		n += len(x.Action)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*StatusAttributes)(nil), "google.devtools.resultstore.v2.StatusAttributes")
	proto.RegisterType((*Property)(nil), "google.devtools.resultstore.v2.Property")
	proto.RegisterType((*Timing)(nil), "google.devtools.resultstore.v2.Timing")
	proto.RegisterType((*Dependency)(nil), "google.devtools.resultstore.v2.Dependency")
	proto.RegisterEnum("google.devtools.resultstore.v2.Language", Language_name, Language_value)
	proto.RegisterEnum("google.devtools.resultstore.v2.Status", Status_name, Status_value)
}

func init() {
	proto.RegisterFile("google/devtools/resultstore/v2/common.proto", fileDescriptor_ff56b05a77242216)
}

var fileDescriptor_ff56b05a77242216 = []byte{
	// 690 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0xcd, 0x8e, 0xe3, 0x44,
	0x10, 0xc7, 0xd7, 0xf9, 0x70, 0x9c, 0xca, 0x7c, 0x14, 0xbd, 0x2b, 0x14, 0xe6, 0x00, 0xa3, 0x1c,
	0xd0, 0x6a, 0x11, 0xb6, 0x14, 0xc4, 0x01, 0x21, 0x21, 0x39, 0xb6, 0x93, 0xf1, 0xc4, 0x6b, 0x5b,
	0xee, 0x0e, 0xd1, 0x72, 0x89, 0x9c, 0xa4, 0xd7, 0xb2, 0x48, 0xec, 0x60, 0xb7, 0x23, 0x0d, 0x6f,
	0xc1, 0x8b, 0x71, 0xe5, 0x75, 0x50, 0xdb, 0x0e, 0x8c, 0x06, 0x01, 0xa7, 0xee, 0x7f, 0xfd, 0x7f,
	0x55, 0x5d, 0xd5, 0x6a, 0x1b, 0xbe, 0x4a, 0xf2, 0x3c, 0x39, 0x70, 0x63, 0xcf, 0xcf, 0x22, 0xcf,
	0x0f, 0xa5, 0x51, 0xf0, 0xb2, 0x3a, 0x88, 0x52, 0xe4, 0x05, 0x37, 0xce, 0x53, 0x63, 0x97, 0x1f,
	0x8f, 0x79, 0xa6, 0x9f, 0x8a, 0x5c, 0xe4, 0xe4, 0xf3, 0x06, 0xd6, 0x2f, 0xb0, 0xfe, 0x0c, 0xd6,
	0xcf, 0xd3, 0xbb, 0xd6, 0x37, 0x6a, 0x7a, 0x5b, 0x7d, 0x34, 0xf6, 0x55, 0x11, 0x8b, 0xf4, 0x92,
	0x7f, 0xf7, 0xc5, 0x4b, 0x5f, 0xa4, 0x47, 0x5e, 0x8a, 0xf8, 0x78, 0x6a, 0x80, 0x89, 0x00, 0xa4,
	0x22, 0x16, 0x55, 0x69, 0x0a, 0x51, 0xa4, 0xdb, 0x4a, 0xf0, 0x92, 0xfc, 0x00, 0x6a, 0x59, 0xc7,
	0xc6, 0xca, 0xbd, 0xf2, 0xf6, 0x66, 0xfa, 0xa5, 0xfe, 0xdf, 0x5d, 0xe8, 0x4d, 0x85, 0xa8, 0xcd,
	0x22, 0xf7, 0x30, 0xda, 0xf3, 0x72, 0x57, 0xa4, 0x27, 0xd9, 0xc9, 0xb8, 0x73, 0xaf, 0xbc, 0x1d,
	0x46, 0xcf, 0x43, 0x93, 0x29, 0x68, 0x61, 0x91, 0x9f, 0x78, 0x21, 0x9e, 0x08, 0x42, 0xf7, 0x67,
	0xfe, 0x54, 0x1f, 0x35, 0x8c, 0xe4, 0x96, 0xbc, 0x81, 0xfe, 0x39, 0x3e, 0x54, 0xbc, 0xcd, 0x6c,
	0xc4, 0xe4, 0x57, 0x50, 0x59, 0x7a, 0x4c, 0xb3, 0x84, 0x7c, 0x07, 0x50, 0x8a, 0xb8, 0x10, 0x1b,
	0x39, 0x4c, 0x9d, 0x38, 0x9a, 0xde, 0x5d, 0x7a, 0xbc, 0x4c, 0xaa, 0xb3, 0xcb, 0xa4, 0xd1, 0xb0,
	0xa6, 0xa5, 0x26, 0xdf, 0x82, 0x76, 0xb9, 0xa1, 0xba, 0xfa, 0x68, 0xfa, 0xd9, 0x3f, 0x12, 0xed,
	0x16, 0x88, 0xfe, 0x42, 0x27, 0xbf, 0x29, 0x00, 0x36, 0x3f, 0xf1, 0x6c, 0xcf, 0xb3, 0xdd, 0x13,
	0x19, 0x83, 0x2a, 0xe2, 0x22, 0xe1, 0xa2, 0xe9, 0xfa, 0xe1, 0x55, 0xd4, 0x6a, 0xf2, 0x35, 0x7c,
	0xb2, 0xcb, 0xb3, 0x8f, 0x69, 0x52, 0x15, 0x7c, 0xbf, 0x69, 0xa1, 0x4e, 0x0b, 0xe1, 0xdf, 0x16,
	0x6b, 0xf0, 0x31, 0xa8, 0xf1, 0xae, 0x6e, 0xa6, 0x7b, 0x29, 0xd4, 0x68, 0x79, 0x07, 0x87, 0x78,
	0xcb, 0x0f, 0xe3, 0x5e, 0x73, 0x07, 0xb5, 0x98, 0x01, 0x68, 0x05, 0x2f, 0xf3, 0xaa, 0xd8, 0xf1,
	0x77, 0xbf, 0x2b, 0xa0, 0x79, 0x71, 0x96, 0x54, 0x71, 0xc2, 0xc9, 0x18, 0xde, 0x78, 0xa6, 0xbf,
	0x58, 0x99, 0x0b, 0x67, 0xb3, 0xf2, 0x69, 0xe8, 0x58, 0xee, 0xdc, 0x75, 0x6c, 0x7c, 0x45, 0x34,
	0xe8, 0xf9, 0x81, 0xef, 0xa0, 0x42, 0x46, 0x30, 0x30, 0x7d, 0x3b, 0x0a, 0x5c, 0x1b, 0x3b, 0x44,
	0x85, 0x8e, 0x49, 0xb1, 0x2b, 0x57, 0xcb, 0xc2, 0x1e, 0x19, 0x40, 0xd7, 0xa2, 0x14, 0xfb, 0x92,
	0xb7, 0xcd, 0x88, 0xa1, 0x2a, 0xad, 0x45, 0x80, 0x03, 0x69, 0x2d, 0xd6, 0x0c, 0x35, 0x59, 0xe0,
	0xc1, 0xa4, 0x4b, 0xc7, 0xf3, 0x70, 0x28, 0xb9, 0x47, 0xf3, 0x47, 0x13, 0x41, 0x72, 0x8f, 0x14,
	0x47, 0x32, 0xe2, 0xb9, 0x34, 0xc4, 0x2b, 0xb9, 0x0b, 0x66, 0x8f, 0x16, 0x5e, 0x4b, 0x2f, 0xfc,
	0x80, 0x37, 0x72, 0xa5, 0x0f, 0x78, 0x4b, 0x86, 0xd0, 0xa7, 0x6b, 0x77, 0xce, 0x10, 0x65, 0x88,
	0x51, 0x24, 0xb2, 0xfc, 0xda, 0x99, 0xe1, 0xeb, 0x77, 0x7f, 0x28, 0xa0, 0x36, 0x2f, 0x89, 0x7c,
	0x0a, 0x84, 0x32, 0x93, 0xad, 0xe8, 0x8b, 0x61, 0xae, 0x40, 0x9b, 0xad, 0x5c, 0xcf, 0x76, 0xfd,
	0x05, 0x2a, 0xb2, 0x98, 0x54, 0x0c, 0x3b, 0xe4, 0x35, 0xdc, 0xce, 0x4d, 0xd7, 0x73, 0xec, 0x0d,
	0x0b, 0x36, 0x35, 0x82, 0x5d, 0xd9, 0x2f, 0x73, 0x28, 0x93, 0x70, 0x8f, 0x00, 0xa8, 0xa1, 0x49,
	0xa9, 0x63, 0x63, 0x5f, 0xee, 0x1b, 0x1a, 0x55, 0x72, 0x0d, 0x43, 0xe6, 0xbe, 0x77, 0xec, 0x4d,
	0xb0, 0x62, 0x38, 0x90, 0xd2, 0x32, 0x7d, 0xcb, 0xf1, 0xa4, 0xab, 0x91, 0x5b, 0x18, 0xb1, 0x20,
	0xf0, 0x36, 0x2d, 0x3e, 0x24, 0x37, 0x00, 0xae, 0x6f, 0x05, 0xef, 0x43, 0xcf, 0x61, 0x0e, 0x82,
	0xec, 0x61, 0xee, 0x99, 0xcb, 0x0f, 0x38, 0x92, 0xc7, 0xad, 0xfc, 0xa5, 0x1f, 0xac, 0x7d, 0xbc,
	0x92, 0x82, 0x2e, 0xdd, 0x30, 0x74, 0x6c, 0xbc, 0x9e, 0xfd, 0x02, 0x93, 0x5d, 0x7e, 0xfc, 0x9f,
	0xaf, 0x28, 0x54, 0x7e, 0x72, 0x5b, 0x22, 0xc9, 0x0f, 0x71, 0x96, 0xe8, 0x79, 0x91, 0x18, 0x09,
	0xcf, 0xea, 0x87, 0x69, 0x34, 0x56, 0x7c, 0x4a, 0xcb, 0x7f, 0xfb, 0x73, 0x7c, 0xff, 0x4c, 0x6e,
	0xd5, 0x3a, 0xeb, 0x9b, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xdb, 0xae, 0x2d, 0xfb, 0x6e, 0x04,
	0x00, 0x00,
}
