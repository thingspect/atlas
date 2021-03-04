// Code generated by protoc-gen-go. DO NOT EDIT.
// source: as/integration/integration.proto

package integration

import (
	fmt "fmt"
	common "github.com/brocaar/chirpstack-api/go/v3/common"
	gw "github.com/brocaar/chirpstack-api/go/v3/gw"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ErrorType int32

const (
	// Unknown type.
	ErrorType_UNKNOWN ErrorType = 0
	// Error related to the downlink payload size.
	// Usually seen when the payload exceeded the maximum allowed payload size.
	ErrorType_DOWNLINK_PAYLOAD_SIZE ErrorType = 1
	// Error related to the downlink frame-counter.
	// Usually seen when the frame-counter has already been used.
	ErrorType_DOWNLINK_FCNT ErrorType = 2
	// Uplink codec error.
	ErrorType_UPLINK_CODEC ErrorType = 3
	// Downlink codec error.
	ErrorType_DOWNLINK_CODEC ErrorType = 4
	// OTAA error.
	ErrorType_OTAA ErrorType = 5
	// Uplink frame-counter was reset.
	ErrorType_UPLINK_FCNT_RESET ErrorType = 6
	// Uplink MIC error.
	ErrorType_UPLINK_MIC ErrorType = 7
	// Uplink frame-counter retransmission.
	ErrorType_UPLINK_FCNT_RETRANSMISSION ErrorType = 8
	// Downlink gateway error.
	ErrorType_DOWNLINK_GATEWAY ErrorType = 9
)

var ErrorType_name = map[int32]string{
	0: "UNKNOWN",
	1: "DOWNLINK_PAYLOAD_SIZE",
	2: "DOWNLINK_FCNT",
	3: "UPLINK_CODEC",
	4: "DOWNLINK_CODEC",
	5: "OTAA",
	6: "UPLINK_FCNT_RESET",
	7: "UPLINK_MIC",
	8: "UPLINK_FCNT_RETRANSMISSION",
	9: "DOWNLINK_GATEWAY",
}

var ErrorType_value = map[string]int32{
	"UNKNOWN":                    0,
	"DOWNLINK_PAYLOAD_SIZE":      1,
	"DOWNLINK_FCNT":              2,
	"UPLINK_CODEC":               3,
	"DOWNLINK_CODEC":             4,
	"OTAA":                       5,
	"UPLINK_FCNT_RESET":          6,
	"UPLINK_MIC":                 7,
	"UPLINK_FCNT_RETRANSMISSION": 8,
	"DOWNLINK_GATEWAY":           9,
}

func (x ErrorType) String() string {
	return proto.EnumName(ErrorType_name, int32(x))
}

func (ErrorType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{0}
}

// UplinkEvent is the message sent when an uplink payload has been received.
type UplinkEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// Receiving gateway RX info.
	RxInfo []*gw.UplinkRXInfo `protobuf:"bytes,5,rep,name=rx_info,json=rxInfo,proto3" json:"rx_info,omitempty"`
	// TX info.
	TxInfo *gw.UplinkTXInfo `protobuf:"bytes,6,opt,name=tx_info,json=txInfo,proto3" json:"tx_info,omitempty"`
	// Device has ADR bit set.
	Adr bool `protobuf:"varint,7,opt,name=adr,proto3" json:"adr,omitempty"`
	// Data-rate.
	Dr uint32 `protobuf:"varint,8,opt,name=dr,proto3" json:"dr,omitempty"`
	// Frame counter.
	FCnt uint32 `protobuf:"varint,9,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	// Frame port.
	FPort uint32 `protobuf:"varint,10,opt,name=f_port,json=fPort,proto3" json:"f_port,omitempty"`
	// FRMPayload data.
	Data []byte `protobuf:"bytes,11,opt,name=data,proto3" json:"data,omitempty"`
	// JSON string containing the decoded object.
	// Note that this is only set when a codec is configured in the Device Profile.
	ObjectJson string `protobuf:"bytes,12,opt,name=object_json,json=objectJSON,proto3" json:"object_json,omitempty"`
	// User-defined device tags.
	Tags map[string]string `protobuf:"bytes,13,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Uplink was of type confirmed.
	ConfirmedUplink bool `protobuf:"varint,14,opt,name=confirmed_uplink,json=confirmedUplink,proto3" json:"confirmed_uplink,omitempty"`
	// Device address.
	DevAddr              []byte   `protobuf:"bytes,15,opt,name=dev_addr,json=devAddr,proto3" json:"dev_addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UplinkEvent) Reset()         { *m = UplinkEvent{} }
func (m *UplinkEvent) String() string { return proto.CompactTextString(m) }
func (*UplinkEvent) ProtoMessage()    {}
func (*UplinkEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{0}
}

func (m *UplinkEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UplinkEvent.Unmarshal(m, b)
}
func (m *UplinkEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UplinkEvent.Marshal(b, m, deterministic)
}
func (m *UplinkEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UplinkEvent.Merge(m, src)
}
func (m *UplinkEvent) XXX_Size() int {
	return xxx_messageInfo_UplinkEvent.Size(m)
}
func (m *UplinkEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_UplinkEvent.DiscardUnknown(m)
}

var xxx_messageInfo_UplinkEvent proto.InternalMessageInfo

func (m *UplinkEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *UplinkEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *UplinkEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *UplinkEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *UplinkEvent) GetRxInfo() []*gw.UplinkRXInfo {
	if m != nil {
		return m.RxInfo
	}
	return nil
}

func (m *UplinkEvent) GetTxInfo() *gw.UplinkTXInfo {
	if m != nil {
		return m.TxInfo
	}
	return nil
}

func (m *UplinkEvent) GetAdr() bool {
	if m != nil {
		return m.Adr
	}
	return false
}

func (m *UplinkEvent) GetDr() uint32 {
	if m != nil {
		return m.Dr
	}
	return 0
}

func (m *UplinkEvent) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

func (m *UplinkEvent) GetFPort() uint32 {
	if m != nil {
		return m.FPort
	}
	return 0
}

func (m *UplinkEvent) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *UplinkEvent) GetObjectJson() string {
	if m != nil {
		return m.ObjectJson
	}
	return ""
}

func (m *UplinkEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *UplinkEvent) GetConfirmedUplink() bool {
	if m != nil {
		return m.ConfirmedUplink
	}
	return false
}

func (m *UplinkEvent) GetDevAddr() []byte {
	if m != nil {
		return m.DevAddr
	}
	return nil
}

// JoinEvent is the message sent when a device joined the network.
// Note that this is only sent after the first received uplink after the
// device (re)activation.
type JoinEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// Device address.
	DevAddr []byte `protobuf:"bytes,5,opt,name=dev_addr,json=devAddr,proto3" json:"dev_addr,omitempty"`
	// Receiving gateway RX info.
	RxInfo []*gw.UplinkRXInfo `protobuf:"bytes,6,rep,name=rx_info,json=rxInfo,proto3" json:"rx_info,omitempty"`
	// TX info.
	TxInfo *gw.UplinkTXInfo `protobuf:"bytes,7,opt,name=tx_info,json=txInfo,proto3" json:"tx_info,omitempty"`
	// Data-rate.
	Dr uint32 `protobuf:"varint,8,opt,name=dr,proto3" json:"dr,omitempty"`
	// User-defined device tags.
	Tags                 map[string]string `protobuf:"bytes,9,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *JoinEvent) Reset()         { *m = JoinEvent{} }
func (m *JoinEvent) String() string { return proto.CompactTextString(m) }
func (*JoinEvent) ProtoMessage()    {}
func (*JoinEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{1}
}

func (m *JoinEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinEvent.Unmarshal(m, b)
}
func (m *JoinEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinEvent.Marshal(b, m, deterministic)
}
func (m *JoinEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinEvent.Merge(m, src)
}
func (m *JoinEvent) XXX_Size() int {
	return xxx_messageInfo_JoinEvent.Size(m)
}
func (m *JoinEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinEvent.DiscardUnknown(m)
}

var xxx_messageInfo_JoinEvent proto.InternalMessageInfo

func (m *JoinEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *JoinEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *JoinEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *JoinEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *JoinEvent) GetDevAddr() []byte {
	if m != nil {
		return m.DevAddr
	}
	return nil
}

func (m *JoinEvent) GetRxInfo() []*gw.UplinkRXInfo {
	if m != nil {
		return m.RxInfo
	}
	return nil
}

func (m *JoinEvent) GetTxInfo() *gw.UplinkTXInfo {
	if m != nil {
		return m.TxInfo
	}
	return nil
}

func (m *JoinEvent) GetDr() uint32 {
	if m != nil {
		return m.Dr
	}
	return 0
}

func (m *JoinEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

// AckEvent is the message sent when a confirmation on a confirmed downlink
// has been received -or- when the downlink timed out.
type AckEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// Frame was acknowledged.
	Acknowledged bool `protobuf:"varint,5,opt,name=acknowledged,proto3" json:"acknowledged,omitempty"`
	// Downlink frame counter to which the acknowledgement relates.
	FCnt uint32 `protobuf:"varint,6,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	// User-defined device tags.
	Tags                 map[string]string `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AckEvent) Reset()         { *m = AckEvent{} }
func (m *AckEvent) String() string { return proto.CompactTextString(m) }
func (*AckEvent) ProtoMessage()    {}
func (*AckEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{2}
}

func (m *AckEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckEvent.Unmarshal(m, b)
}
func (m *AckEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckEvent.Marshal(b, m, deterministic)
}
func (m *AckEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckEvent.Merge(m, src)
}
func (m *AckEvent) XXX_Size() int {
	return xxx_messageInfo_AckEvent.Size(m)
}
func (m *AckEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_AckEvent.DiscardUnknown(m)
}

var xxx_messageInfo_AckEvent proto.InternalMessageInfo

func (m *AckEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *AckEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *AckEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *AckEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *AckEvent) GetAcknowledged() bool {
	if m != nil {
		return m.Acknowledged
	}
	return false
}

func (m *AckEvent) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

func (m *AckEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

// TxAckEvent is the message sent when a downlink was acknowledged by the gateway
// for transmission. As a downlink can be scheduled in the future, this event
// does not confirm that the message has already been transmitted.
type TxAckEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// Downlink frame-counter.
	FCnt uint32 `protobuf:"varint,5,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	// User-defined device tags.
	Tags map[string]string `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Gateway ID.
	GatewayId []byte `protobuf:"bytes,7,opt,name=gateway_id,json=gatewayID,proto3" json:"gateway_id,omitempty"`
	// TX info.
	TxInfo               *gw.DownlinkTXInfo `protobuf:"bytes,8,opt,name=tx_info,json=txInfo,proto3" json:"tx_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TxAckEvent) Reset()         { *m = TxAckEvent{} }
func (m *TxAckEvent) String() string { return proto.CompactTextString(m) }
func (*TxAckEvent) ProtoMessage()    {}
func (*TxAckEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{3}
}

func (m *TxAckEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxAckEvent.Unmarshal(m, b)
}
func (m *TxAckEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxAckEvent.Marshal(b, m, deterministic)
}
func (m *TxAckEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxAckEvent.Merge(m, src)
}
func (m *TxAckEvent) XXX_Size() int {
	return xxx_messageInfo_TxAckEvent.Size(m)
}
func (m *TxAckEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_TxAckEvent.DiscardUnknown(m)
}

var xxx_messageInfo_TxAckEvent proto.InternalMessageInfo

func (m *TxAckEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *TxAckEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *TxAckEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *TxAckEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *TxAckEvent) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

func (m *TxAckEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *TxAckEvent) GetGatewayId() []byte {
	if m != nil {
		return m.GatewayId
	}
	return nil
}

func (m *TxAckEvent) GetTxInfo() *gw.DownlinkTXInfo {
	if m != nil {
		return m.TxInfo
	}
	return nil
}

// ErrorEvent is the message sent when an error occurred.
type ErrorEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// Error type.
	Type ErrorType `protobuf:"varint,5,opt,name=type,proto3,enum=integration.ErrorType" json:"type,omitempty"`
	// Error message.
	Error string `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
	// Downlink frame-counter (in case the downlink is related to a scheduled downlink).
	FCnt uint32 `protobuf:"varint,7,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	// User-defined device tags.
	Tags                 map[string]string `protobuf:"bytes,8,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ErrorEvent) Reset()         { *m = ErrorEvent{} }
func (m *ErrorEvent) String() string { return proto.CompactTextString(m) }
func (*ErrorEvent) ProtoMessage()    {}
func (*ErrorEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{4}
}

func (m *ErrorEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorEvent.Unmarshal(m, b)
}
func (m *ErrorEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorEvent.Marshal(b, m, deterministic)
}
func (m *ErrorEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorEvent.Merge(m, src)
}
func (m *ErrorEvent) XXX_Size() int {
	return xxx_messageInfo_ErrorEvent.Size(m)
}
func (m *ErrorEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorEvent proto.InternalMessageInfo

func (m *ErrorEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *ErrorEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *ErrorEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *ErrorEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *ErrorEvent) GetType() ErrorType {
	if m != nil {
		return m.Type
	}
	return ErrorType_UNKNOWN
}

func (m *ErrorEvent) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *ErrorEvent) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

func (m *ErrorEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

// StatusEvent is the message sent when a device-status mac-command was sent
// by the device.
type StatusEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// The demodulation signal-to-noise ratio in dB for the last successfully
	// received device-status request by the Network Server.
	Margin int32 `protobuf:"varint,5,opt,name=margin,proto3" json:"margin,omitempty"`
	// Device is connected to an external power source.
	ExternalPowerSource bool `protobuf:"varint,6,opt,name=external_power_source,json=externalPowerSource,proto3" json:"external_power_source,omitempty"`
	// Battery level is not available.
	BatteryLevelUnavailable bool `protobuf:"varint,7,opt,name=battery_level_unavailable,json=batteryLevelUnavailable,proto3" json:"battery_level_unavailable,omitempty"`
	// Battery level.
	BatteryLevel float32 `protobuf:"fixed32,8,opt,name=battery_level,json=batteryLevel,proto3" json:"battery_level,omitempty"`
	// User-defined device tags.
	Tags                 map[string]string `protobuf:"bytes,9,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *StatusEvent) Reset()         { *m = StatusEvent{} }
func (m *StatusEvent) String() string { return proto.CompactTextString(m) }
func (*StatusEvent) ProtoMessage()    {}
func (*StatusEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{5}
}

func (m *StatusEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusEvent.Unmarshal(m, b)
}
func (m *StatusEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusEvent.Marshal(b, m, deterministic)
}
func (m *StatusEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusEvent.Merge(m, src)
}
func (m *StatusEvent) XXX_Size() int {
	return xxx_messageInfo_StatusEvent.Size(m)
}
func (m *StatusEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusEvent.DiscardUnknown(m)
}

var xxx_messageInfo_StatusEvent proto.InternalMessageInfo

func (m *StatusEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *StatusEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *StatusEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *StatusEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *StatusEvent) GetMargin() int32 {
	if m != nil {
		return m.Margin
	}
	return 0
}

func (m *StatusEvent) GetExternalPowerSource() bool {
	if m != nil {
		return m.ExternalPowerSource
	}
	return false
}

func (m *StatusEvent) GetBatteryLevelUnavailable() bool {
	if m != nil {
		return m.BatteryLevelUnavailable
	}
	return false
}

func (m *StatusEvent) GetBatteryLevel() float32 {
	if m != nil {
		return m.BatteryLevel
	}
	return 0
}

func (m *StatusEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

// LocationEvent is the message sent when a geolocation resolve was returned.
type LocationEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// Location.
	Location *common.Location `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	// User-defined device tags.
	Tags map[string]string `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Uplink IDs used for geolocation.
	// This is set in case the geolocation is based on the uplink meta-data.
	UplinkIds [][]byte `protobuf:"bytes,7,rep,name=uplink_ids,json=uplinkIDs,proto3" json:"uplink_ids,omitempty"`
	// Frame counter (in case the geolocation is based on the payload).
	// This is set in case the geolocation is based on the uplink payload content.
	FCnt                 uint32   `protobuf:"varint,8,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LocationEvent) Reset()         { *m = LocationEvent{} }
func (m *LocationEvent) String() string { return proto.CompactTextString(m) }
func (*LocationEvent) ProtoMessage()    {}
func (*LocationEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{6}
}

func (m *LocationEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocationEvent.Unmarshal(m, b)
}
func (m *LocationEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocationEvent.Marshal(b, m, deterministic)
}
func (m *LocationEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocationEvent.Merge(m, src)
}
func (m *LocationEvent) XXX_Size() int {
	return xxx_messageInfo_LocationEvent.Size(m)
}
func (m *LocationEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_LocationEvent.DiscardUnknown(m)
}

var xxx_messageInfo_LocationEvent proto.InternalMessageInfo

func (m *LocationEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *LocationEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *LocationEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *LocationEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *LocationEvent) GetLocation() *common.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *LocationEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *LocationEvent) GetUplinkIds() [][]byte {
	if m != nil {
		return m.UplinkIds
	}
	return nil
}

func (m *LocationEvent) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

// IntegrationEvent is the message that can be sent by an integration.
// It allows for sending events which are provided by an external integration
// which are "not native" to ChirpStack.
type IntegrationEvent struct {
	// Application ID.
	ApplicationId uint64 `protobuf:"varint,1,opt,name=application_id,json=applicationID,proto3" json:"application_id,omitempty"`
	// Application name.
	ApplicationName string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// Device name.
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	// Device EUI.
	DevEui []byte `protobuf:"bytes,4,opt,name=dev_eui,json=devEUI,proto3" json:"dev_eui,omitempty"`
	// User-defined device tags.
	Tags map[string]string `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Integration name.
	IntegrationName string `protobuf:"bytes,6,opt,name=integration_name,json=integrationName,proto3" json:"integration_name,omitempty"`
	// Event type.
	EventType string `protobuf:"bytes,7,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	// JSON string containing the event object.
	ObjectJson           string   `protobuf:"bytes,8,opt,name=object_json,json=objectJSON,proto3" json:"object_json,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntegrationEvent) Reset()         { *m = IntegrationEvent{} }
func (m *IntegrationEvent) String() string { return proto.CompactTextString(m) }
func (*IntegrationEvent) ProtoMessage()    {}
func (*IntegrationEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba82681be587591, []int{7}
}

func (m *IntegrationEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntegrationEvent.Unmarshal(m, b)
}
func (m *IntegrationEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntegrationEvent.Marshal(b, m, deterministic)
}
func (m *IntegrationEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntegrationEvent.Merge(m, src)
}
func (m *IntegrationEvent) XXX_Size() int {
	return xxx_messageInfo_IntegrationEvent.Size(m)
}
func (m *IntegrationEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_IntegrationEvent.DiscardUnknown(m)
}

var xxx_messageInfo_IntegrationEvent proto.InternalMessageInfo

func (m *IntegrationEvent) GetApplicationId() uint64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *IntegrationEvent) GetApplicationName() string {
	if m != nil {
		return m.ApplicationName
	}
	return ""
}

func (m *IntegrationEvent) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *IntegrationEvent) GetDevEui() []byte {
	if m != nil {
		return m.DevEui
	}
	return nil
}

func (m *IntegrationEvent) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *IntegrationEvent) GetIntegrationName() string {
	if m != nil {
		return m.IntegrationName
	}
	return ""
}

func (m *IntegrationEvent) GetEventType() string {
	if m != nil {
		return m.EventType
	}
	return ""
}

func (m *IntegrationEvent) GetObjectJson() string {
	if m != nil {
		return m.ObjectJson
	}
	return ""
}

func init() {
	proto.RegisterEnum("integration.ErrorType", ErrorType_name, ErrorType_value)
	proto.RegisterType((*UplinkEvent)(nil), "integration.UplinkEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.UplinkEvent.TagsEntry")
	proto.RegisterType((*JoinEvent)(nil), "integration.JoinEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.JoinEvent.TagsEntry")
	proto.RegisterType((*AckEvent)(nil), "integration.AckEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.AckEvent.TagsEntry")
	proto.RegisterType((*TxAckEvent)(nil), "integration.TxAckEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.TxAckEvent.TagsEntry")
	proto.RegisterType((*ErrorEvent)(nil), "integration.ErrorEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.ErrorEvent.TagsEntry")
	proto.RegisterType((*StatusEvent)(nil), "integration.StatusEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.StatusEvent.TagsEntry")
	proto.RegisterType((*LocationEvent)(nil), "integration.LocationEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.LocationEvent.TagsEntry")
	proto.RegisterType((*IntegrationEvent)(nil), "integration.IntegrationEvent")
	proto.RegisterMapType((map[string]string)(nil), "integration.IntegrationEvent.TagsEntry")
}

func init() {
	proto.RegisterFile("as/integration/integration.proto", fileDescriptor_5ba82681be587591)
}

var fileDescriptor_5ba82681be587591 = []byte{
	// 1081 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x57, 0xcd, 0x6e, 0xdb, 0x46,
	0x17, 0xfd, 0x44, 0xfd, 0x91, 0x57, 0x3f, 0x66, 0xc6, 0x71, 0x42, 0x1b, 0xc8, 0x17, 0xd5, 0x6d,
	0x51, 0x25, 0x6d, 0x25, 0xc0, 0x6e, 0xd3, 0x20, 0x5d, 0x29, 0x96, 0x5a, 0x30, 0x71, 0x24, 0x81,
	0x92, 0xe1, 0x26, 0x1b, 0x62, 0x44, 0x8e, 0x18, 0xc6, 0x12, 0x87, 0x18, 0x8d, 0x24, 0xeb, 0x09,
	0xfa, 0x2c, 0x7d, 0x82, 0xae, 0xfa, 0x14, 0x7d, 0x80, 0xae, 0xbb, 0xeb, 0x1b, 0xb4, 0xe0, 0x0c,
	0x2d, 0x91, 0x8e, 0x81, 0x16, 0x86, 0x17, 0x5a, 0x69, 0xe6, 0xcc, 0xd1, 0xcc, 0xbd, 0xe7, 0xdc,
	0xf9, 0x21, 0xd4, 0xf0, 0xac, 0xe9, 0x07, 0x9c, 0x78, 0x0c, 0x73, 0x9f, 0x06, 0xc9, 0x76, 0x23,
	0x64, 0x94, 0x53, 0x54, 0x4a, 0x40, 0x07, 0xbb, 0x0e, 0x9d, 0x4e, 0x69, 0xd0, 0x94, 0x3f, 0x92,
	0x71, 0x50, 0xf2, 0x96, 0x4d, 0x6f, 0x29, 0x3b, 0x87, 0xbf, 0xe4, 0xa0, 0x74, 0x16, 0x4e, 0xfc,
	0xe0, 0xa2, 0xb3, 0x20, 0x01, 0x47, 0x9f, 0x43, 0x15, 0x87, 0xe1, 0xc4, 0x77, 0xc4, 0x04, 0xb6,
	0xef, 0x1a, 0x99, 0x5a, 0xa6, 0x9e, 0xb3, 0x2a, 0x09, 0xd4, 0x6c, 0xa3, 0x27, 0xa0, 0x27, 0x69,
	0x01, 0x9e, 0x12, 0x43, 0xa9, 0x65, 0xea, 0x9a, 0xb5, 0x93, 0xc0, 0xbb, 0x78, 0x4a, 0xd0, 0x63,
	0x28, 0xb9, 0x64, 0xe1, 0x3b, 0x44, 0xb2, 0xb2, 0x82, 0x05, 0x12, 0x12, 0x84, 0x87, 0x50, 0x74,
	0xc9, 0xc2, 0x26, 0x73, 0xdf, 0xc8, 0xd5, 0x32, 0xf5, 0xb2, 0x55, 0x70, 0xc9, 0xa2, 0x73, 0x66,
	0xa2, 0x27, 0x50, 0x64, 0x97, 0xb6, 0x1f, 0x8c, 0xa9, 0x91, 0xaf, 0x65, 0xeb, 0xa5, 0x23, 0xbd,
	0xe1, 0x2d, 0x1b, 0x32, 0x5a, 0xeb, 0x27, 0x33, 0x18, 0x53, 0xab, 0xc0, 0x2e, 0xa3, 0xdf, 0x88,
	0xca, 0x63, 0x6a, 0xa1, 0x96, 0x49, 0x53, 0x87, 0x31, 0x95, 0x4b, 0xaa, 0x0e, 0x59, 0xec, 0x32,
	0xa3, 0x58, 0xcb, 0xd4, 0x55, 0x2b, 0x6a, 0xa2, 0x2a, 0x28, 0x2e, 0x33, 0xd4, 0x5a, 0xa6, 0x5e,
	0xb1, 0x14, 0x97, 0xa1, 0x5d, 0xc8, 0x8f, 0x6d, 0x27, 0xe0, 0x86, 0x26, 0xa0, 0xdc, 0xf8, 0x24,
	0xe0, 0x68, 0x0f, 0x0a, 0x63, 0x3b, 0xa4, 0x8c, 0x1b, 0x20, 0xd0, 0xfc, 0xb8, 0x4f, 0x19, 0x47,
	0x08, 0x72, 0x2e, 0xe6, 0xd8, 0x28, 0x89, 0xc8, 0x45, 0x3b, 0xca, 0x98, 0x8e, 0x3e, 0x10, 0x87,
	0xdb, 0x1f, 0x66, 0x34, 0x30, 0xca, 0x32, 0x63, 0x09, 0xbd, 0x1a, 0xf4, 0xba, 0xe8, 0x19, 0xe4,
	0x38, 0xf6, 0x66, 0x46, 0x45, 0x64, 0x75, 0xd8, 0x48, 0xba, 0x98, 0x30, 0xa3, 0x31, 0xc4, 0xde,
	0xac, 0x13, 0x70, 0xb6, 0xb2, 0x04, 0x3f, 0x52, 0xdd, 0xa1, 0xc1, 0xd8, 0x67, 0x53, 0xe2, 0xda,
	0x73, 0x41, 0x34, 0xaa, 0x22, 0x8f, 0x9d, 0x35, 0x2e, 0xff, 0x8f, 0xf6, 0x41, 0x8d, 0x44, 0xc5,
	0xae, 0xcb, 0x8c, 0x1d, 0x11, 0x5b, 0x24, 0x72, 0xcb, 0x75, 0xd9, 0xc1, 0x77, 0xa0, 0xad, 0x27,
	0x8e, 0xd4, 0xb8, 0x20, 0x2b, 0x61, 0xb2, 0x66, 0x45, 0x4d, 0x74, 0x1f, 0xf2, 0x0b, 0x3c, 0x99,
	0x5f, 0xf9, 0x29, 0x3b, 0x2f, 0x94, 0xe7, 0x99, 0xc3, 0x9f, 0xb3, 0xa0, 0xbd, 0xa2, 0x7e, 0xb0,
	0x7d, 0x95, 0x92, 0xcc, 0x36, 0x9f, 0xca, 0x36, 0x59, 0x44, 0x85, 0xff, 0x5e, 0x44, 0xc5, 0x7f,
	0x29, 0xa2, 0xeb, 0x25, 0xf3, 0x4d, 0xec, 0xa8, 0x26, 0x96, 0xa8, 0xa5, 0x1c, 0x5d, 0x4b, 0x76,
	0xdd, 0xcf, 0xdb, 0x3b, 0xf1, 0x9b, 0x02, 0x6a, 0xcb, 0xd9, 0xc2, 0x2d, 0x7b, 0x08, 0x65, 0xec,
	0x5c, 0x04, 0x74, 0x39, 0x21, 0xae, 0x47, 0x5c, 0x61, 0x86, 0x6a, 0xa5, 0xb0, 0xcd, 0xf6, 0x2a,
	0x24, 0xb6, 0xd7, 0x71, 0x2c, 0x60, 0x51, 0x08, 0xf8, 0x38, 0x25, 0xe0, 0x55, 0xa6, 0x77, 0xa7,
	0xdf, 0x9f, 0x0a, 0xc0, 0xf0, 0x72, 0x2b, 0x15, 0x5c, 0xab, 0x93, 0x4f, 0xa8, 0xf3, 0x6d, 0xac,
	0x8e, 0xac, 0xe0, 0x4f, 0x52, 0xea, 0x6c, 0xf2, 0xf8, 0xe8, 0xbc, 0x78, 0x04, 0xe0, 0x61, 0x4e,
	0x96, 0x78, 0x15, 0xe5, 0x54, 0x14, 0xeb, 0x68, 0x31, 0x62, 0xb6, 0xd1, 0x97, 0x9b, 0x7a, 0x57,
	0x45, 0xbd, 0xa3, 0xa8, 0xde, 0xdb, 0x74, 0x19, 0x7c, 0x5c, 0xf1, 0xb7, 0xd7, 0xfa, 0x0f, 0x05,
	0xa0, 0xc3, 0x18, 0x65, 0xdb, 0xa7, 0xf5, 0x53, 0xc8, 0xf1, 0x55, 0x48, 0x84, 0xd4, 0xd5, 0xa3,
	0x07, 0x29, 0x59, 0x45, 0xc8, 0xc3, 0x55, 0x48, 0x2c, 0xc1, 0x89, 0x12, 0x24, 0x11, 0x24, 0xaa,
	0x56, 0xb3, 0x64, 0x67, 0xe3, 0x56, 0xf1, 0x06, 0xb7, 0xd4, 0x1b, 0xdc, 0xda, 0x28, 0x71, 0x77,
	0xd5, 0xfc, 0x6b, 0x16, 0x4a, 0x03, 0x8e, 0xf9, 0x7c, 0xb6, 0x7d, 0x12, 0x3f, 0x80, 0xc2, 0x14,
	0x33, 0xcf, 0x0f, 0x84, 0xc8, 0x79, 0x2b, 0xee, 0xa1, 0x23, 0xd8, 0x23, 0x97, 0x9c, 0xb0, 0x00,
	0x4f, 0xec, 0x90, 0x2e, 0x09, 0xb3, 0x67, 0x74, 0xce, 0x1c, 0x22, 0xe4, 0x55, 0xad, 0xdd, 0xab,
	0xc1, 0x7e, 0x34, 0x36, 0x10, 0x43, 0xe8, 0x05, 0xec, 0x8f, 0x30, 0xe7, 0x84, 0xad, 0xec, 0x09,
	0x59, 0x90, 0x89, 0x3d, 0x0f, 0xf0, 0x02, 0xfb, 0x13, 0x3c, 0x9a, 0x90, 0xf8, 0x3e, 0x7f, 0x18,
	0x13, 0x4e, 0xa3, 0xf1, 0xb3, 0xcd, 0x30, 0xfa, 0x14, 0x2a, 0xa9, 0xff, 0x8a, 0x8a, 0x57, 0xac,
	0x72, 0x92, 0xbf, 0xbe, 0x97, 0xb5, 0x1b, 0xee, 0xe5, 0x84, 0xc0, 0x77, 0xe7, 0xdc, 0x5f, 0x0a,
	0x54, 0x4e, 0xa9, 0x14, 0x7a, 0xfb, 0xbc, 0xfb, 0x0a, 0xd4, 0x49, 0x1c, 0x9c, 0x70, 0x2f, 0xba,
	0x10, 0xe3, 0x97, 0xe4, 0x55, 0xd0, 0xd6, 0x9a, 0x81, 0x9e, 0xa7, 0xce, 0xa8, 0xcf, 0x52, 0xe2,
	0xa5, 0x72, 0xbc, 0xe9, 0x98, 0x92, 0x8f, 0x19, 0xdb, 0x77, 0xe5, 0x0d, 0x50, 0xb6, 0x34, 0x89,
	0x98, 0xed, 0xd9, 0x66, 0x8f, 0xa9, 0x9b, 0x3d, 0x76, 0x7b, 0xc9, 0xff, 0x56, 0x40, 0x37, 0x37,
	0xa1, 0x6d, 0x9f, 0xea, 0xdf, 0xc7, 0x3a, 0xca, 0x27, 0xef, 0x17, 0x29, 0x1d, 0xaf, 0x07, 0x7e,
	0xd3, 0x0b, 0x31, 0xc1, 0x97, 0x6b, 0xcb, 0x03, 0x6b, 0x27, 0x81, 0x8b, 0x00, 0x1e, 0x01, 0x90,
	0x68, 0x0e, 0x5b, 0x1c, 0x81, 0x45, 0x41, 0xd2, 0x04, 0x12, 0x9d, 0x7a, 0xd7, 0x1f, 0xb1, 0xea,
	0xf5, 0x47, 0xec, 0xad, 0x1d, 0x78, 0xfa, 0x7b, 0x06, 0xb4, 0xf5, 0xe9, 0x8a, 0x4a, 0x50, 0x3c,
	0xeb, 0xbe, 0xee, 0xf6, 0xce, 0xbb, 0xfa, 0xff, 0xd0, 0x3e, 0xec, 0xb5, 0x7b, 0xe7, 0xdd, 0x53,
	0xb3, 0xfb, 0xda, 0xee, 0xb7, 0xde, 0x9e, 0xf6, 0x5a, 0x6d, 0x7b, 0x60, 0xbe, 0xeb, 0xe8, 0x19,
	0x74, 0x0f, 0x2a, 0xeb, 0xa1, 0x1f, 0x4e, 0xba, 0x43, 0x5d, 0x41, 0x3a, 0x94, 0xcf, 0xfa, 0x02,
	0x38, 0xe9, 0xb5, 0x3b, 0x27, 0x7a, 0x16, 0x21, 0xa8, 0xae, 0x49, 0x12, 0xcb, 0x21, 0x15, 0x72,
	0xbd, 0x61, 0xab, 0xa5, 0xe7, 0xd1, 0x1e, 0xdc, 0x8b, 0xf9, 0xd1, 0x04, 0xb6, 0xd5, 0x19, 0x74,
	0x86, 0x7a, 0x01, 0x55, 0x01, 0x62, 0xf8, 0x8d, 0x79, 0xa2, 0x17, 0xd1, 0xff, 0xe1, 0x20, 0x4d,
	0x1b, 0x5a, 0xad, 0xee, 0xe0, 0x8d, 0x39, 0x18, 0x98, 0xbd, 0xae, 0xae, 0xa2, 0xfb, 0xa0, 0xaf,
	0x17, 0xf9, 0xb1, 0x35, 0xec, 0x9c, 0xb7, 0xde, 0xea, 0xda, 0xcb, 0x00, 0x6a, 0x3e, 0x6d, 0x38,
	0xef, 0x7d, 0x16, 0xce, 0x38, 0x76, 0x2e, 0x1a, 0x38, 0xf4, 0x1b, 0x78, 0x96, 0xb4, 0xef, 0x65,
	0xb2, 0xf0, 0xfa, 0xd1, 0xe7, 0x57, 0x3f, 0xf3, 0xee, 0x99, 0xe7, 0xf3, 0xf7, 0xf3, 0x51, 0xb4,
	0xb1, 0x9a, 0x23, 0x46, 0x1d, 0x8c, 0x59, 0x73, 0x33, 0xcb, 0xd7, 0x38, 0xf4, 0x9b, 0x1e, 0x6d,
	0x2e, 0x8e, 0x9b, 0xe9, 0x2f, 0xbf, 0x51, 0x41, 0x7c, 0xbf, 0x1d, 0xff, 0x13, 0x00, 0x00, 0xff,
	0xff, 0x5a, 0x47, 0xe2, 0x65, 0x12, 0x0e, 0x00, 0x00,
}
