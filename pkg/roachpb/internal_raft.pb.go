// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: roachpb/internal_raft.proto

package roachpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import hlc "github.com/weisslj/cockroach/pkg/util/hlc"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// RaftTruncatedState contains metadata about the truncated portion of the raft log.
// Raft requires access to the term of the last truncated log entry even after the
// rest of the entry has been discarded.
type RaftTruncatedState struct {
	// The highest index that has been removed from the log.
	Index uint64 `protobuf:"varint,1,opt,name=index" json:"index"`
	// The term corresponding to 'index'.
	Term                 uint64   `protobuf:"varint,2,opt,name=term" json:"term"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RaftTruncatedState) Reset()         { *m = RaftTruncatedState{} }
func (m *RaftTruncatedState) String() string { return proto.CompactTextString(m) }
func (*RaftTruncatedState) ProtoMessage()    {}
func (*RaftTruncatedState) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_raft_8611ae4325a43066, []int{0}
}
func (m *RaftTruncatedState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RaftTruncatedState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *RaftTruncatedState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RaftTruncatedState.Merge(dst, src)
}
func (m *RaftTruncatedState) XXX_Size() int {
	return m.Size()
}
func (m *RaftTruncatedState) XXX_DiscardUnknown() {
	xxx_messageInfo_RaftTruncatedState.DiscardUnknown(m)
}

var xxx_messageInfo_RaftTruncatedState proto.InternalMessageInfo

// RaftTombstone contains information about a replica that has been deleted.
type RaftTombstone struct {
	NextReplicaID        ReplicaID `protobuf:"varint,1,opt,name=next_replica_id,json=nextReplicaId,casttype=ReplicaID" json:"next_replica_id"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *RaftTombstone) Reset()         { *m = RaftTombstone{} }
func (m *RaftTombstone) String() string { return proto.CompactTextString(m) }
func (*RaftTombstone) ProtoMessage()    {}
func (*RaftTombstone) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_raft_8611ae4325a43066, []int{1}
}
func (m *RaftTombstone) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RaftTombstone) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *RaftTombstone) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RaftTombstone.Merge(dst, src)
}
func (m *RaftTombstone) XXX_Size() int {
	return m.Size()
}
func (m *RaftTombstone) XXX_DiscardUnknown() {
	xxx_messageInfo_RaftTombstone.DiscardUnknown(m)
}

var xxx_messageInfo_RaftTombstone proto.InternalMessageInfo

// RaftSnapshotData is the payload of a raftpb.Snapshot. It contains a raw copy of
// all of the range's data and metadata, including the raft log, abort span, etc.
type RaftSnapshotData struct {
	// The latest RangeDescriptor
	RangeDescriptor RangeDescriptor             `protobuf:"bytes,1,opt,name=range_descriptor,json=rangeDescriptor" json:"range_descriptor"`
	KV              []RaftSnapshotData_KeyValue `protobuf:"bytes,2,rep,name=KV" json:"KV"`
	// These are really raftpb.Entry, but we model them as raw bytes to avoid
	// roundtripping through memory.
	LogEntries           [][]byte `protobuf:"bytes,3,rep,name=log_entries,json=logEntries" json:"log_entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RaftSnapshotData) Reset()         { *m = RaftSnapshotData{} }
func (m *RaftSnapshotData) String() string { return proto.CompactTextString(m) }
func (*RaftSnapshotData) ProtoMessage()    {}
func (*RaftSnapshotData) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_raft_8611ae4325a43066, []int{2}
}
func (m *RaftSnapshotData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RaftSnapshotData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *RaftSnapshotData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RaftSnapshotData.Merge(dst, src)
}
func (m *RaftSnapshotData) XXX_Size() int {
	return m.Size()
}
func (m *RaftSnapshotData) XXX_DiscardUnknown() {
	xxx_messageInfo_RaftSnapshotData.DiscardUnknown(m)
}

var xxx_messageInfo_RaftSnapshotData proto.InternalMessageInfo

type RaftSnapshotData_KeyValue struct {
	Key                  []byte        `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value                []byte        `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Timestamp            hlc.Timestamp `protobuf:"bytes,3,opt,name=timestamp" json:"timestamp"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RaftSnapshotData_KeyValue) Reset()         { *m = RaftSnapshotData_KeyValue{} }
func (m *RaftSnapshotData_KeyValue) String() string { return proto.CompactTextString(m) }
func (*RaftSnapshotData_KeyValue) ProtoMessage()    {}
func (*RaftSnapshotData_KeyValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_raft_8611ae4325a43066, []int{2, 0}
}
func (m *RaftSnapshotData_KeyValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RaftSnapshotData_KeyValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *RaftSnapshotData_KeyValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RaftSnapshotData_KeyValue.Merge(dst, src)
}
func (m *RaftSnapshotData_KeyValue) XXX_Size() int {
	return m.Size()
}
func (m *RaftSnapshotData_KeyValue) XXX_DiscardUnknown() {
	xxx_messageInfo_RaftSnapshotData_KeyValue.DiscardUnknown(m)
}

var xxx_messageInfo_RaftSnapshotData_KeyValue proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RaftTruncatedState)(nil), "cockroach.roachpb.RaftTruncatedState")
	proto.RegisterType((*RaftTombstone)(nil), "cockroach.roachpb.RaftTombstone")
	proto.RegisterType((*RaftSnapshotData)(nil), "cockroach.roachpb.RaftSnapshotData")
	proto.RegisterType((*RaftSnapshotData_KeyValue)(nil), "cockroach.roachpb.RaftSnapshotData.KeyValue")
}
func (this *RaftTruncatedState) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RaftTruncatedState)
	if !ok {
		that2, ok := that.(RaftTruncatedState)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Index != that1.Index {
		return false
	}
	if this.Term != that1.Term {
		return false
	}
	return true
}
func (m *RaftTruncatedState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftTruncatedState) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintInternalRaft(dAtA, i, uint64(m.Index))
	dAtA[i] = 0x10
	i++
	i = encodeVarintInternalRaft(dAtA, i, uint64(m.Term))
	return i, nil
}

func (m *RaftTombstone) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftTombstone) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintInternalRaft(dAtA, i, uint64(m.NextReplicaID))
	return i, nil
}

func (m *RaftSnapshotData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftSnapshotData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintInternalRaft(dAtA, i, uint64(m.RangeDescriptor.Size()))
	n1, err := m.RangeDescriptor.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.KV) > 0 {
		for _, msg := range m.KV {
			dAtA[i] = 0x12
			i++
			i = encodeVarintInternalRaft(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.LogEntries) > 0 {
		for _, b := range m.LogEntries {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintInternalRaft(dAtA, i, uint64(len(b)))
			i += copy(dAtA[i:], b)
		}
	}
	return i, nil
}

func (m *RaftSnapshotData_KeyValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftSnapshotData_KeyValue) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Key != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintInternalRaft(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	if m.Value != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintInternalRaft(dAtA, i, uint64(len(m.Value)))
		i += copy(dAtA[i:], m.Value)
	}
	dAtA[i] = 0x1a
	i++
	i = encodeVarintInternalRaft(dAtA, i, uint64(m.Timestamp.Size()))
	n2, err := m.Timestamp.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	return i, nil
}

func encodeVarintInternalRaft(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedRaftTruncatedState(r randyInternalRaft, easy bool) *RaftTruncatedState {
	this := &RaftTruncatedState{}
	this.Index = uint64(uint64(r.Uint32()))
	this.Term = uint64(uint64(r.Uint32()))
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyInternalRaft interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneInternalRaft(r randyInternalRaft) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringInternalRaft(r randyInternalRaft) string {
	v1 := r.Intn(100)
	tmps := make([]rune, v1)
	for i := 0; i < v1; i++ {
		tmps[i] = randUTF8RuneInternalRaft(r)
	}
	return string(tmps)
}
func randUnrecognizedInternalRaft(r randyInternalRaft, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldInternalRaft(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldInternalRaft(dAtA []byte, r randyInternalRaft, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateInternalRaft(dAtA, uint64(key))
		v2 := r.Int63()
		if r.Intn(2) == 0 {
			v2 *= -1
		}
		dAtA = encodeVarintPopulateInternalRaft(dAtA, uint64(v2))
	case 1:
		dAtA = encodeVarintPopulateInternalRaft(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateInternalRaft(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateInternalRaft(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateInternalRaft(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateInternalRaft(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *RaftTruncatedState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovInternalRaft(uint64(m.Index))
	n += 1 + sovInternalRaft(uint64(m.Term))
	return n
}

func (m *RaftTombstone) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovInternalRaft(uint64(m.NextReplicaID))
	return n
}

func (m *RaftSnapshotData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.RangeDescriptor.Size()
	n += 1 + l + sovInternalRaft(uint64(l))
	if len(m.KV) > 0 {
		for _, e := range m.KV {
			l = e.Size()
			n += 1 + l + sovInternalRaft(uint64(l))
		}
	}
	if len(m.LogEntries) > 0 {
		for _, b := range m.LogEntries {
			l = len(b)
			n += 1 + l + sovInternalRaft(uint64(l))
		}
	}
	return n
}

func (m *RaftSnapshotData_KeyValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Key != nil {
		l = len(m.Key)
		n += 1 + l + sovInternalRaft(uint64(l))
	}
	if m.Value != nil {
		l = len(m.Value)
		n += 1 + l + sovInternalRaft(uint64(l))
	}
	l = m.Timestamp.Size()
	n += 1 + l + sovInternalRaft(uint64(l))
	return n
}

func sovInternalRaft(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozInternalRaft(x uint64) (n int) {
	return sovInternalRaft(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RaftTruncatedState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternalRaft
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RaftTruncatedState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RaftTruncatedState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Term", wireType)
			}
			m.Term = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Term |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipInternalRaft(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternalRaft
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RaftTombstone) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternalRaft
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RaftTombstone: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RaftTombstone: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextReplicaID", wireType)
			}
			m.NextReplicaID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextReplicaID |= (ReplicaID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipInternalRaft(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternalRaft
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RaftSnapshotData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternalRaft
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RaftSnapshotData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RaftSnapshotData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RangeDescriptor", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInternalRaft
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RangeDescriptor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KV", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInternalRaft
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KV = append(m.KV, RaftSnapshotData_KeyValue{})
			if err := m.KV[len(m.KV)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LogEntries", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthInternalRaft
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LogEntries = append(m.LogEntries, make([]byte, postIndex-iNdEx))
			copy(m.LogEntries[len(m.LogEntries)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternalRaft(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternalRaft
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RaftSnapshotData_KeyValue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternalRaft
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: KeyValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeyValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthInternalRaft
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthInternalRaft
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = append(m.Value[:0], dAtA[iNdEx:postIndex]...)
			if m.Value == nil {
				m.Value = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInternalRaft
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternalRaft(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternalRaft
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipInternalRaft(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInternalRaft
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInternalRaft
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthInternalRaft
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowInternalRaft
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipInternalRaft(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthInternalRaft = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInternalRaft   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("roachpb/internal_raft.proto", fileDescriptor_internal_raft_8611ae4325a43066)
}

var fileDescriptor_internal_raft_8611ae4325a43066 = []byte{
	// 447 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xb1, 0x8e, 0xd3, 0x4c,
	0x10, 0x8e, 0xed, 0x44, 0xff, 0xdd, 0x26, 0xd1, 0xe5, 0x5f, 0x9d, 0x90, 0x15, 0x84, 0x1d, 0x5c,
	0xa5, 0x40, 0x8e, 0x74, 0x25, 0x1d, 0x51, 0x90, 0x80, 0x48, 0x14, 0x9b, 0x28, 0x05, 0x14, 0xd6,
	0x9e, 0x3d, 0xe7, 0x58, 0x67, 0xef, 0x5a, 0xeb, 0x09, 0xca, 0xbd, 0x05, 0x8f, 0x70, 0x8f, 0xc1,
	0x23, 0x84, 0x8e, 0x92, 0x2a, 0x02, 0xd3, 0xf0, 0x0c, 0x54, 0xc8, 0x6b, 0x27, 0xe4, 0x04, 0xdd,
	0xcc, 0xf7, 0x7d, 0xf3, 0x79, 0xbe, 0xf1, 0x92, 0xc7, 0x4a, 0xf2, 0x70, 0x9d, 0x5f, 0x4f, 0x12,
	0x81, 0xa0, 0x04, 0x4f, 0x03, 0xc5, 0x6f, 0xd0, 0xcf, 0x95, 0x44, 0x49, 0xff, 0x0f, 0x65, 0x78,
	0xab, 0x05, 0x7e, 0x23, 0x1b, 0x3e, 0x3a, 0xe8, 0x33, 0x40, 0x1e, 0x71, 0xe4, 0xb5, 0x74, 0x68,
	0x6f, 0x30, 0x49, 0x27, 0xeb, 0x34, 0x9c, 0x60, 0x92, 0x41, 0x81, 0x3c, 0xcb, 0x1b, 0xe6, 0x32,
	0x96, 0xb1, 0xd4, 0xe5, 0xa4, 0xaa, 0x6a, 0xd4, 0x5b, 0x12, 0xca, 0xf8, 0x0d, 0x2e, 0xd5, 0x46,
	0x84, 0x1c, 0x21, 0x5a, 0x20, 0x47, 0xa0, 0x43, 0xd2, 0x49, 0x44, 0x04, 0x5b, 0xdb, 0x18, 0x19,
	0xe3, 0xf6, 0xb4, 0xbd, 0xdb, 0xbb, 0x2d, 0x56, 0x43, 0xd4, 0x26, 0x6d, 0x04, 0x95, 0xd9, 0xe6,
	0x09, 0xa5, 0x91, 0xe7, 0x67, 0x9f, 0xee, 0x5d, 0xe3, 0xe7, 0xbd, 0x6b, 0x78, 0xef, 0x49, 0x5f,
	0xbb, 0xca, 0xec, 0xba, 0x40, 0x29, 0x80, 0xbe, 0x21, 0x17, 0x02, 0xb6, 0x18, 0x28, 0xc8, 0xd3,
	0x24, 0xe4, 0x41, 0x12, 0x69, 0xeb, 0xce, 0xd4, 0xab, 0xe6, 0xcb, 0xbd, 0xdb, 0x7f, 0x0b, 0x5b,
	0x64, 0x35, 0xfb, 0x7a, 0xf6, 0x6b, 0xef, 0x9e, 0x1f, 0x1b, 0xd6, 0x17, 0x27, 0x5c, 0xe4, 0x7d,
	0x36, 0xc9, 0xa0, 0x72, 0x5f, 0x08, 0x9e, 0x17, 0x6b, 0x89, 0x33, 0x8e, 0x9c, 0x2e, 0xc8, 0x40,
	0x71, 0x11, 0x43, 0x10, 0x41, 0x11, 0xaa, 0x24, 0x47, 0xa9, 0xf4, 0x17, 0xba, 0x57, 0x9e, 0xff,
	0xd7, 0xf5, 0x7c, 0x56, 0x49, 0x67, 0x47, 0x65, 0x93, 0xe2, 0x42, 0x3d, 0x84, 0xe9, 0x2b, 0x62,
	0xce, 0x57, 0xb6, 0x39, 0xb2, 0xc6, 0xdd, 0xab, 0x67, 0xff, 0xb4, 0x79, 0xb8, 0x85, 0x3f, 0x87,
	0xbb, 0x15, 0x4f, 0x37, 0x30, 0x25, 0x4d, 0x2c, 0x73, 0xbe, 0x62, 0xe6, 0x7c, 0x45, 0x5d, 0xd2,
	0x4d, 0x65, 0x1c, 0x80, 0x40, 0x95, 0x40, 0x61, 0x5b, 0x23, 0x6b, 0xdc, 0x63, 0x24, 0x95, 0xf1,
	0xcb, 0x1a, 0x19, 0x6e, 0xc8, 0xd9, 0x61, 0x98, 0x0e, 0x88, 0x75, 0x0b, 0x77, 0x7a, 0xfd, 0x1e,
	0xab, 0x4a, 0x7a, 0x49, 0x3a, 0x1f, 0x2a, 0x4a, 0x1f, 0xbd, 0xc7, 0xea, 0x86, 0xbe, 0x20, 0xe7,
	0xc7, 0x9f, 0x6c, 0x5b, 0x3a, 0xec, 0x93, 0x93, 0x2d, 0xab, 0x97, 0xe0, 0xaf, 0xd3, 0xd0, 0x5f,
	0x1e, 0x44, 0x4d, 0xce, 0x3f, 0x53, 0xd3, 0xa7, 0xbb, 0xef, 0x4e, 0x6b, 0x57, 0x3a, 0xc6, 0x97,
	0xd2, 0x31, 0xbe, 0x96, 0x8e, 0xf1, 0xad, 0x74, 0x8c, 0x8f, 0x3f, 0x9c, 0xd6, 0xbb, 0xff, 0x9a,
	0x90, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x40, 0xad, 0x88, 0x9a, 0x02, 0x00, 0x00,
}
