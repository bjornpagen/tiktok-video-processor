// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package tiktokdb

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type User struct {
	_tab flatbuffers.Table
}

func GetRootAsUser(buf []byte, offset flatbuffers.UOffsetT) *User {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &User{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsUser(buf []byte, offset flatbuffers.UOffsetT) *User {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &User{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *User) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *User) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *User) LatestUsername() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) LastUsedMincursor() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) Videos(obj *UserAwemes) *UserAwemes {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(UserAwemes)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func UserStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func UserAddLatestUsername(builder *flatbuffers.Builder, latestUsername flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(latestUsername), 0)
}
func UserAddLastUsedMincursor(builder *flatbuffers.Builder, lastUsedMincursor flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(lastUsedMincursor), 0)
}
func UserAddVideos(builder *flatbuffers.Builder, videos flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(videos), 0)
}
func UserEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
