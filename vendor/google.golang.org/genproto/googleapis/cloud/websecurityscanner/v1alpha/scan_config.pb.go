// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/websecurityscanner/v1alpha/scan_config.proto

package websecurityscanner

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Type of user agents used for scanning.
type ScanConfig_UserAgent int32

const (
	// The user agent is unknown. Service will default to CHROME_LINUX.
	ScanConfig_USER_AGENT_UNSPECIFIED ScanConfig_UserAgent = 0
	// Chrome on Linux. This is the service default if unspecified.
	ScanConfig_CHROME_LINUX ScanConfig_UserAgent = 1
	// Chrome on Android.
	ScanConfig_CHROME_ANDROID ScanConfig_UserAgent = 2
	// Safari on IPhone.
	ScanConfig_SAFARI_IPHONE ScanConfig_UserAgent = 3
)

var ScanConfig_UserAgent_name = map[int32]string{
	0: "USER_AGENT_UNSPECIFIED",
	1: "CHROME_LINUX",
	2: "CHROME_ANDROID",
	3: "SAFARI_IPHONE",
}

var ScanConfig_UserAgent_value = map[string]int32{
	"USER_AGENT_UNSPECIFIED": 0,
	"CHROME_LINUX":           1,
	"CHROME_ANDROID":         2,
	"SAFARI_IPHONE":          3,
}

func (x ScanConfig_UserAgent) String() string {
	return proto.EnumName(ScanConfig_UserAgent_name, int32(x))
}

func (ScanConfig_UserAgent) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0, 0}
}

// Cloud platforms supported by Cloud Web Security Scanner.
type ScanConfig_TargetPlatform int32

const (
	// The target platform is unknown. Requests with this enum value will be
	// rejected with INVALID_ARGUMENT error.
	ScanConfig_TARGET_PLATFORM_UNSPECIFIED ScanConfig_TargetPlatform = 0
	// Google App Engine service.
	ScanConfig_APP_ENGINE ScanConfig_TargetPlatform = 1
	// Google Compute Engine service.
	ScanConfig_COMPUTE ScanConfig_TargetPlatform = 2
)

var ScanConfig_TargetPlatform_name = map[int32]string{
	0: "TARGET_PLATFORM_UNSPECIFIED",
	1: "APP_ENGINE",
	2: "COMPUTE",
}

var ScanConfig_TargetPlatform_value = map[string]int32{
	"TARGET_PLATFORM_UNSPECIFIED": 0,
	"APP_ENGINE":                  1,
	"COMPUTE":                     2,
}

func (x ScanConfig_TargetPlatform) String() string {
	return proto.EnumName(ScanConfig_TargetPlatform_name, int32(x))
}

func (ScanConfig_TargetPlatform) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0, 1}
}

// A ScanConfig resource contains the configurations to launch a scan.
type ScanConfig struct {
	// The resource name of the ScanConfig. The name follows the format of
	// 'projects/{projectId}/scanConfigs/{scanConfigId}'. The ScanConfig IDs are
	// generated by the system.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Required.
	// The user provided display name of the ScanConfig.
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// The maximum QPS during scanning. A valid value ranges from 5 to 20
	// inclusively. If the field is unspecified or its value is set 0, server will
	// default to 15. Other values outside of [5, 20] range will be rejected with
	// INVALID_ARGUMENT error.
	MaxQps int32 `protobuf:"varint,3,opt,name=max_qps,json=maxQps,proto3" json:"max_qps,omitempty"`
	// Required.
	// The starting URLs from which the scanner finds site pages.
	StartingUrls []string `protobuf:"bytes,4,rep,name=starting_urls,json=startingUrls,proto3" json:"starting_urls,omitempty"`
	// The authentication configuration. If specified, service will use the
	// authentication configuration during scanning.
	Authentication *ScanConfig_Authentication `protobuf:"bytes,5,opt,name=authentication,proto3" json:"authentication,omitempty"`
	// The user agent used during scanning.
	UserAgent ScanConfig_UserAgent `protobuf:"varint,6,opt,name=user_agent,json=userAgent,proto3,enum=google.cloud.websecurityscanner.v1alpha.ScanConfig_UserAgent" json:"user_agent,omitempty"`
	// The blacklist URL patterns as described in
	// https://cloud.google.com/security-scanner/docs/excluded-urls
	BlacklistPatterns []string `protobuf:"bytes,7,rep,name=blacklist_patterns,json=blacklistPatterns,proto3" json:"blacklist_patterns,omitempty"`
	// The schedule of the ScanConfig.
	Schedule *ScanConfig_Schedule `protobuf:"bytes,8,opt,name=schedule,proto3" json:"schedule,omitempty"`
	// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be
	// used as a default.
	TargetPlatforms      []ScanConfig_TargetPlatform `protobuf:"varint,9,rep,packed,name=target_platforms,json=targetPlatforms,proto3,enum=google.cloud.websecurityscanner.v1alpha.ScanConfig_TargetPlatform" json:"target_platforms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ScanConfig) Reset()         { *m = ScanConfig{} }
func (m *ScanConfig) String() string { return proto.CompactTextString(m) }
func (*ScanConfig) ProtoMessage()    {}
func (*ScanConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0}
}
func (m *ScanConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanConfig.Unmarshal(m, b)
}
func (m *ScanConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanConfig.Marshal(b, m, deterministic)
}
func (m *ScanConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanConfig.Merge(m, src)
}
func (m *ScanConfig) XXX_Size() int {
	return xxx_messageInfo_ScanConfig.Size(m)
}
func (m *ScanConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ScanConfig proto.InternalMessageInfo

func (m *ScanConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ScanConfig) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *ScanConfig) GetMaxQps() int32 {
	if m != nil {
		return m.MaxQps
	}
	return 0
}

func (m *ScanConfig) GetStartingUrls() []string {
	if m != nil {
		return m.StartingUrls
	}
	return nil
}

func (m *ScanConfig) GetAuthentication() *ScanConfig_Authentication {
	if m != nil {
		return m.Authentication
	}
	return nil
}

func (m *ScanConfig) GetUserAgent() ScanConfig_UserAgent {
	if m != nil {
		return m.UserAgent
	}
	return ScanConfig_USER_AGENT_UNSPECIFIED
}

func (m *ScanConfig) GetBlacklistPatterns() []string {
	if m != nil {
		return m.BlacklistPatterns
	}
	return nil
}

func (m *ScanConfig) GetSchedule() *ScanConfig_Schedule {
	if m != nil {
		return m.Schedule
	}
	return nil
}

func (m *ScanConfig) GetTargetPlatforms() []ScanConfig_TargetPlatform {
	if m != nil {
		return m.TargetPlatforms
	}
	return nil
}

// Scan authentication configuration.
type ScanConfig_Authentication struct {
	// Required.
	// Authentication configuration
	//
	// Types that are valid to be assigned to Authentication:
	//	*ScanConfig_Authentication_GoogleAccount_
	//	*ScanConfig_Authentication_CustomAccount_
	Authentication       isScanConfig_Authentication_Authentication `protobuf_oneof:"authentication"`
	XXX_NoUnkeyedLiteral struct{}                                   `json:"-"`
	XXX_unrecognized     []byte                                     `json:"-"`
	XXX_sizecache        int32                                      `json:"-"`
}

func (m *ScanConfig_Authentication) Reset()         { *m = ScanConfig_Authentication{} }
func (m *ScanConfig_Authentication) String() string { return proto.CompactTextString(m) }
func (*ScanConfig_Authentication) ProtoMessage()    {}
func (*ScanConfig_Authentication) Descriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0, 0}
}
func (m *ScanConfig_Authentication) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanConfig_Authentication.Unmarshal(m, b)
}
func (m *ScanConfig_Authentication) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanConfig_Authentication.Marshal(b, m, deterministic)
}
func (m *ScanConfig_Authentication) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanConfig_Authentication.Merge(m, src)
}
func (m *ScanConfig_Authentication) XXX_Size() int {
	return xxx_messageInfo_ScanConfig_Authentication.Size(m)
}
func (m *ScanConfig_Authentication) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanConfig_Authentication.DiscardUnknown(m)
}

var xxx_messageInfo_ScanConfig_Authentication proto.InternalMessageInfo

type isScanConfig_Authentication_Authentication interface {
	isScanConfig_Authentication_Authentication()
}

type ScanConfig_Authentication_GoogleAccount_ struct {
	GoogleAccount *ScanConfig_Authentication_GoogleAccount `protobuf:"bytes,1,opt,name=google_account,json=googleAccount,proto3,oneof"`
}

type ScanConfig_Authentication_CustomAccount_ struct {
	CustomAccount *ScanConfig_Authentication_CustomAccount `protobuf:"bytes,2,opt,name=custom_account,json=customAccount,proto3,oneof"`
}

func (*ScanConfig_Authentication_GoogleAccount_) isScanConfig_Authentication_Authentication() {}

func (*ScanConfig_Authentication_CustomAccount_) isScanConfig_Authentication_Authentication() {}

func (m *ScanConfig_Authentication) GetAuthentication() isScanConfig_Authentication_Authentication {
	if m != nil {
		return m.Authentication
	}
	return nil
}

func (m *ScanConfig_Authentication) GetGoogleAccount() *ScanConfig_Authentication_GoogleAccount {
	if x, ok := m.GetAuthentication().(*ScanConfig_Authentication_GoogleAccount_); ok {
		return x.GoogleAccount
	}
	return nil
}

func (m *ScanConfig_Authentication) GetCustomAccount() *ScanConfig_Authentication_CustomAccount {
	if x, ok := m.GetAuthentication().(*ScanConfig_Authentication_CustomAccount_); ok {
		return x.CustomAccount
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ScanConfig_Authentication) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ScanConfig_Authentication_OneofMarshaler, _ScanConfig_Authentication_OneofUnmarshaler, _ScanConfig_Authentication_OneofSizer, []interface{}{
		(*ScanConfig_Authentication_GoogleAccount_)(nil),
		(*ScanConfig_Authentication_CustomAccount_)(nil),
	}
}

func _ScanConfig_Authentication_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ScanConfig_Authentication)
	// authentication
	switch x := m.Authentication.(type) {
	case *ScanConfig_Authentication_GoogleAccount_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GoogleAccount); err != nil {
			return err
		}
	case *ScanConfig_Authentication_CustomAccount_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CustomAccount); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ScanConfig_Authentication.Authentication has unexpected type %T", x)
	}
	return nil
}

func _ScanConfig_Authentication_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ScanConfig_Authentication)
	switch tag {
	case 1: // authentication.google_account
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ScanConfig_Authentication_GoogleAccount)
		err := b.DecodeMessage(msg)
		m.Authentication = &ScanConfig_Authentication_GoogleAccount_{msg}
		return true, err
	case 2: // authentication.custom_account
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ScanConfig_Authentication_CustomAccount)
		err := b.DecodeMessage(msg)
		m.Authentication = &ScanConfig_Authentication_CustomAccount_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ScanConfig_Authentication_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ScanConfig_Authentication)
	// authentication
	switch x := m.Authentication.(type) {
	case *ScanConfig_Authentication_GoogleAccount_:
		s := proto.Size(x.GoogleAccount)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ScanConfig_Authentication_CustomAccount_:
		s := proto.Size(x.CustomAccount)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Describes authentication configuration that uses a Google account.
type ScanConfig_Authentication_GoogleAccount struct {
	// Required.
	// The user name of the Google account.
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Input only.
	// Required.
	// The password of the Google account. The credential is stored encrypted
	// and not returned in any response.
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScanConfig_Authentication_GoogleAccount) Reset() {
	*m = ScanConfig_Authentication_GoogleAccount{}
}
func (m *ScanConfig_Authentication_GoogleAccount) String() string { return proto.CompactTextString(m) }
func (*ScanConfig_Authentication_GoogleAccount) ProtoMessage()    {}
func (*ScanConfig_Authentication_GoogleAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0, 0, 0}
}
func (m *ScanConfig_Authentication_GoogleAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanConfig_Authentication_GoogleAccount.Unmarshal(m, b)
}
func (m *ScanConfig_Authentication_GoogleAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanConfig_Authentication_GoogleAccount.Marshal(b, m, deterministic)
}
func (m *ScanConfig_Authentication_GoogleAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanConfig_Authentication_GoogleAccount.Merge(m, src)
}
func (m *ScanConfig_Authentication_GoogleAccount) XXX_Size() int {
	return xxx_messageInfo_ScanConfig_Authentication_GoogleAccount.Size(m)
}
func (m *ScanConfig_Authentication_GoogleAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanConfig_Authentication_GoogleAccount.DiscardUnknown(m)
}

var xxx_messageInfo_ScanConfig_Authentication_GoogleAccount proto.InternalMessageInfo

func (m *ScanConfig_Authentication_GoogleAccount) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ScanConfig_Authentication_GoogleAccount) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

// Describes authentication configuration that uses a custom account.
type ScanConfig_Authentication_CustomAccount struct {
	// Required.
	// The user name of the custom account.
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Input only.
	// Required.
	// The password of the custom account. The credential is stored encrypted
	// and not returned in any response.
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// Required.
	// The login form URL of the website.
	LoginUrl             string   `protobuf:"bytes,3,opt,name=login_url,json=loginUrl,proto3" json:"login_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScanConfig_Authentication_CustomAccount) Reset() {
	*m = ScanConfig_Authentication_CustomAccount{}
}
func (m *ScanConfig_Authentication_CustomAccount) String() string { return proto.CompactTextString(m) }
func (*ScanConfig_Authentication_CustomAccount) ProtoMessage()    {}
func (*ScanConfig_Authentication_CustomAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0, 0, 1}
}
func (m *ScanConfig_Authentication_CustomAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanConfig_Authentication_CustomAccount.Unmarshal(m, b)
}
func (m *ScanConfig_Authentication_CustomAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanConfig_Authentication_CustomAccount.Marshal(b, m, deterministic)
}
func (m *ScanConfig_Authentication_CustomAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanConfig_Authentication_CustomAccount.Merge(m, src)
}
func (m *ScanConfig_Authentication_CustomAccount) XXX_Size() int {
	return xxx_messageInfo_ScanConfig_Authentication_CustomAccount.Size(m)
}
func (m *ScanConfig_Authentication_CustomAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanConfig_Authentication_CustomAccount.DiscardUnknown(m)
}

var xxx_messageInfo_ScanConfig_Authentication_CustomAccount proto.InternalMessageInfo

func (m *ScanConfig_Authentication_CustomAccount) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ScanConfig_Authentication_CustomAccount) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *ScanConfig_Authentication_CustomAccount) GetLoginUrl() string {
	if m != nil {
		return m.LoginUrl
	}
	return ""
}

// Scan schedule configuration.
type ScanConfig_Schedule struct {
	// A timestamp indicates when the next run will be scheduled. The value is
	// refreshed by the server after each run. If unspecified, it will default
	// to current server time, which means the scan will be scheduled to start
	// immediately.
	ScheduleTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=schedule_time,json=scheduleTime,proto3" json:"schedule_time,omitempty"`
	// Required.
	// The duration of time between executions in days.
	IntervalDurationDays int32    `protobuf:"varint,2,opt,name=interval_duration_days,json=intervalDurationDays,proto3" json:"interval_duration_days,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScanConfig_Schedule) Reset()         { *m = ScanConfig_Schedule{} }
func (m *ScanConfig_Schedule) String() string { return proto.CompactTextString(m) }
func (*ScanConfig_Schedule) ProtoMessage()    {}
func (*ScanConfig_Schedule) Descriptor() ([]byte, []int) {
	return fileDescriptor_50b1b6d7cca97898, []int{0, 1}
}
func (m *ScanConfig_Schedule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanConfig_Schedule.Unmarshal(m, b)
}
func (m *ScanConfig_Schedule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanConfig_Schedule.Marshal(b, m, deterministic)
}
func (m *ScanConfig_Schedule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanConfig_Schedule.Merge(m, src)
}
func (m *ScanConfig_Schedule) XXX_Size() int {
	return xxx_messageInfo_ScanConfig_Schedule.Size(m)
}
func (m *ScanConfig_Schedule) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanConfig_Schedule.DiscardUnknown(m)
}

var xxx_messageInfo_ScanConfig_Schedule proto.InternalMessageInfo

func (m *ScanConfig_Schedule) GetScheduleTime() *timestamp.Timestamp {
	if m != nil {
		return m.ScheduleTime
	}
	return nil
}

func (m *ScanConfig_Schedule) GetIntervalDurationDays() int32 {
	if m != nil {
		return m.IntervalDurationDays
	}
	return 0
}

func init() {
	proto.RegisterType((*ScanConfig)(nil), "google.cloud.websecurityscanner.v1alpha.ScanConfig")
	proto.RegisterType((*ScanConfig_Authentication)(nil), "google.cloud.websecurityscanner.v1alpha.ScanConfig.Authentication")
	proto.RegisterType((*ScanConfig_Authentication_GoogleAccount)(nil), "google.cloud.websecurityscanner.v1alpha.ScanConfig.Authentication.GoogleAccount")
	proto.RegisterType((*ScanConfig_Authentication_CustomAccount)(nil), "google.cloud.websecurityscanner.v1alpha.ScanConfig.Authentication.CustomAccount")
	proto.RegisterType((*ScanConfig_Schedule)(nil), "google.cloud.websecurityscanner.v1alpha.ScanConfig.Schedule")
	proto.RegisterEnum("google.cloud.websecurityscanner.v1alpha.ScanConfig_UserAgent", ScanConfig_UserAgent_name, ScanConfig_UserAgent_value)
	proto.RegisterEnum("google.cloud.websecurityscanner.v1alpha.ScanConfig_TargetPlatform", ScanConfig_TargetPlatform_name, ScanConfig_TargetPlatform_value)
}

func init() {
	proto.RegisterFile("google/cloud/websecurityscanner/v1alpha/scan_config.proto", fileDescriptor_50b1b6d7cca97898)
}

var fileDescriptor_50b1b6d7cca97898 = []byte{
	// 748 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x51, 0x6f, 0xe3, 0x44,
	0x10, 0x3e, 0x37, 0xd7, 0x36, 0x99, 0x36, 0x39, 0xdf, 0x0a, 0x1d, 0x91, 0x0f, 0xe9, 0x42, 0x79,
	0x20, 0x12, 0xc2, 0x16, 0x85, 0x17, 0x04, 0x08, 0xb9, 0x89, 0x9b, 0x46, 0xba, 0x3a, 0x66, 0x93,
	0x48, 0x07, 0x42, 0x5a, 0xb6, 0xce, 0xd6, 0x35, 0xd8, 0xbb, 0x66, 0x77, 0x7d, 0x77, 0x79, 0xe4,
	0x77, 0xf0, 0x73, 0x78, 0xe0, 0x6f, 0x21, 0x6f, 0xec, 0x5c, 0xd3, 0x7b, 0xa0, 0x2a, 0xbc, 0x79,
	0xe6, 0xf3, 0x7c, 0xdf, 0xec, 0xf8, 0x9b, 0x35, 0x7c, 0x9d, 0x08, 0x91, 0x64, 0xcc, 0x8b, 0x33,
	0x51, 0xae, 0xbc, 0x37, 0xec, 0x4a, 0xb1, 0xb8, 0x94, 0xa9, 0x5e, 0xab, 0x98, 0x72, 0xce, 0xa4,
	0xf7, 0xfa, 0x0b, 0x9a, 0x15, 0x37, 0xd4, 0xab, 0x62, 0x12, 0x0b, 0x7e, 0x9d, 0x26, 0x6e, 0x21,
	0x85, 0x16, 0xe8, 0xd3, 0x4d, 0xa9, 0x6b, 0x4a, 0xdd, 0xf7, 0x4b, 0xdd, 0xba, 0xd4, 0xf9, 0xa8,
	0xd6, 0xa0, 0x45, 0xea, 0x51, 0xce, 0x85, 0xa6, 0x3a, 0x15, 0x5c, 0x6d, 0x68, 0x9c, 0x17, 0x35,
	0x6a, 0xa2, 0xab, 0xf2, 0xda, 0xd3, 0x69, 0xce, 0x94, 0xa6, 0x79, 0xb1, 0x79, 0xe1, 0xe4, 0x2f,
	0x00, 0x98, 0xc7, 0x94, 0x8f, 0x8c, 0x38, 0x42, 0xf0, 0x98, 0xd3, 0x9c, 0xf5, 0xad, 0x81, 0x35,
	0xec, 0x60, 0xf3, 0x8c, 0x3e, 0x86, 0xe3, 0x55, 0xaa, 0x8a, 0x8c, 0xae, 0x89, 0xc1, 0xf6, 0x0c,
	0x76, 0x54, 0xe7, 0xc2, 0xea, 0x95, 0x0f, 0xe1, 0x30, 0xa7, 0x6f, 0xc9, 0xef, 0x85, 0xea, 0xb7,
	0x06, 0xd6, 0x70, 0x1f, 0x1f, 0xe4, 0xf4, 0xed, 0x0f, 0x85, 0x42, 0x9f, 0x40, 0x57, 0x69, 0x2a,
	0x75, 0xca, 0x13, 0x52, 0xca, 0x4c, 0xf5, 0x1f, 0x0f, 0x5a, 0xc3, 0x0e, 0x3e, 0x6e, 0x92, 0x4b,
	0x99, 0x29, 0xf4, 0x2b, 0xf4, 0x68, 0xa9, 0x6f, 0x18, 0xd7, 0x69, 0x6c, 0xba, 0xef, 0xef, 0x0f,
	0xac, 0xe1, 0xd1, 0xe9, 0x99, 0x7b, 0xcf, 0x21, 0xb8, 0xef, 0x4e, 0xe0, 0xfa, 0x3b, 0x4c, 0xf8,
	0x0e, 0x33, 0xfa, 0x19, 0xa0, 0x54, 0x4c, 0x12, 0x9a, 0x30, 0xae, 0xfb, 0x07, 0x03, 0x6b, 0xd8,
	0x3b, 0xfd, 0xee, 0x21, 0x3a, 0x4b, 0xc5, 0xa4, 0x5f, 0x91, 0xe0, 0x4e, 0xd9, 0x3c, 0xa2, 0xcf,
	0x01, 0x5d, 0x65, 0x34, 0xfe, 0x2d, 0x4b, 0x95, 0x26, 0x05, 0xd5, 0x9a, 0x49, 0xae, 0xfa, 0x87,
	0xe6, 0xcc, 0x4f, 0xb7, 0x48, 0x54, 0x03, 0xe8, 0x15, 0xb4, 0x55, 0x7c, 0xc3, 0x56, 0x65, 0xc6,
	0xfa, 0x6d, 0x73, 0xe4, 0x6f, 0x1f, 0xd2, 0xca, 0xbc, 0xe6, 0xc0, 0x5b, 0x36, 0x94, 0x83, 0xad,
	0xa9, 0x4c, 0x98, 0x26, 0x45, 0x46, 0xf5, 0xb5, 0x90, 0xb9, 0xea, 0x77, 0x06, 0xad, 0x61, 0xef,
	0x61, 0x43, 0x5d, 0x18, 0xae, 0xa8, 0xa6, 0xc2, 0x4f, 0xf4, 0x4e, 0xac, 0x9c, 0xbf, 0x5b, 0xd0,
	0xdb, 0x1d, 0x3c, 0x5a, 0x43, 0x6f, 0x23, 0x44, 0x68, 0x1c, 0x8b, 0x92, 0x6b, 0xe3, 0xa9, 0xa3,
	0xd3, 0xe8, 0xbf, 0x7f, 0x54, 0x77, 0x62, 0x18, 0xfc, 0x0d, 0xef, 0xc5, 0x23, 0xdc, 0x4d, 0x6e,
	0x27, 0x2a, 0xe9, 0xb8, 0x54, 0x5a, 0xe4, 0x5b, 0xe9, 0xbd, 0xff, 0x4d, 0x7a, 0x64, 0x88, 0x6f,
	0x49, 0xc7, 0xb7, 0x13, 0xce, 0x04, 0xba, 0x3b, 0xcd, 0x21, 0x07, 0xda, 0x95, 0x3d, 0x6e, 0x2d,
	0xd5, 0x36, 0xae, 0xb0, 0x82, 0x2a, 0xf5, 0x46, 0xc8, 0x55, 0xbd, 0x54, 0xdb, 0xd8, 0x59, 0x41,
	0x77, 0x47, 0xea, 0xa1, 0x44, 0xe8, 0x39, 0x74, 0x32, 0x91, 0xa4, 0xbc, 0x5a, 0x3f, 0xb3, 0x9c,
	0x1d, 0xdc, 0x36, 0x89, 0xa5, 0xcc, 0xce, 0xec, 0xbb, 0x9b, 0xe7, 0xfc, 0x61, 0x41, 0xbb, 0xf1,
	0x13, 0xfa, 0x1e, 0xba, 0x8d, 0xa3, 0x48, 0x75, 0x71, 0xd4, 0x9f, 0xd0, 0x69, 0xe6, 0xd8, 0xdc,
	0x2a, 0xee, 0xa2, 0xb9, 0x55, 0xf0, 0x71, 0x53, 0x50, 0xa5, 0xd0, 0x57, 0xf0, 0x2c, 0xe5, 0x9a,
	0xc9, 0xd7, 0x34, 0x23, 0xab, 0x52, 0x1a, 0x09, 0xb2, 0xa2, 0x6b, 0x65, 0xda, 0xdc, 0xc7, 0x1f,
	0x34, 0xe8, 0xb8, 0x06, 0xc7, 0x74, 0xad, 0x4e, 0x7e, 0x81, 0xce, 0x76, 0xbb, 0x90, 0x03, 0xcf,
	0x96, 0xf3, 0x00, 0x13, 0x7f, 0x12, 0x84, 0x0b, 0xb2, 0x0c, 0xe7, 0x51, 0x30, 0x9a, 0x9e, 0x4f,
	0x83, 0xb1, 0xfd, 0x08, 0xd9, 0x70, 0x3c, 0xba, 0xc0, 0xb3, 0xcb, 0x80, 0xbc, 0x9c, 0x86, 0xcb,
	0x57, 0xb6, 0x85, 0x10, 0xf4, 0xea, 0x8c, 0x1f, 0x8e, 0xf1, 0x6c, 0x3a, 0xb6, 0xf7, 0xd0, 0x53,
	0xe8, 0xce, 0xfd, 0x73, 0x1f, 0x4f, 0xc9, 0x34, 0xba, 0x98, 0x85, 0x81, 0xdd, 0x3a, 0x09, 0xa1,
	0xb7, 0x6b, 0x69, 0xf4, 0x02, 0x9e, 0x2f, 0x7c, 0x3c, 0x09, 0x16, 0x24, 0x7a, 0xe9, 0x2f, 0xce,
	0x67, 0xf8, 0xf2, 0x8e, 0x56, 0x0f, 0xc0, 0x8f, 0x22, 0x12, 0x84, 0x93, 0x69, 0x18, 0xd8, 0x16,
	0x3a, 0x82, 0xc3, 0xd1, 0xec, 0x32, 0x5a, 0x2e, 0x02, 0x7b, 0xef, 0xec, 0x4f, 0x0b, 0x3e, 0x8b,
	0x45, 0x7e, 0x5f, 0x7f, 0x9d, 0x3d, 0x79, 0x67, 0xb0, 0xa8, 0x9a, 0x61, 0x64, 0xfd, 0xf4, 0x63,
	0x5d, 0x9b, 0x88, 0x8c, 0xf2, 0xc4, 0x15, 0x32, 0xf1, 0x12, 0xc6, 0xcd, 0x84, 0xbd, 0x0d, 0x44,
	0x8b, 0x54, 0xfd, 0xeb, 0xaf, 0xe4, 0x9b, 0xf7, 0xa1, 0xab, 0x03, 0xc3, 0xf2, 0xe5, 0x3f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x7d, 0xa3, 0x6b, 0x23, 0x8f, 0x06, 0x00, 0x00,
}
