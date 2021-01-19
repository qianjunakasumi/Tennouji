package jce

import (
	"bytes"
	"encoding/binary"
)

// writer 写入器
type writer struct {
	b   *bytes.Buffer // 缓冲区
	tag uint8         // 标签
}

// NewWriter 返回一个写入器
func NewWriter(tag uint8) *writer { return &writer{bytes.NewBuffer(nil), tag} }

// writeKey 写入键
func (w *writer) writeKey(Type byte) {
	if w.tag > 14 {
		w.b.WriteByte(Type | 0xF0)
		w.b.WriteByte(w.tag)
	} else {
		w.b.WriteByte(w.tag<<4 | Type)
	}
	w.tag++
}

// WriteByte 写入 Byte
func (w *writer) WriteByte(b byte) *writer {
	if b == 0 {
		w.writeKey(Zero)
		return w
	}
	w.writeKey(Byte)
	w.b.WriteByte(b)
	return w
}

// WriteBool 写入 Bool
func (w *writer) WriteBool(b bool) *writer {
	if b {
		w.writeKey(Bool)
		w.b.WriteByte(1)
		return w
	}
	w.writeKey(Zero)
	return w
}

// WriteInt16 写入 Int16
func (w *writer) WriteInt16(i int16) *writer {
	if i < 128 && i > -129 {
		return w.WriteByte(byte(i))
	}
	w.writeKey(Int16)
	_ = binary.Write(w.b, binary.BigEndian, i)
	return w
}

// WriteInt32 写入 Int32
func (w *writer) WriteInt32(i int32) *writer {
	if i < 32768 && i > -32769 {
		return w.WriteInt16(int16(i))
	}
	w.writeKey(Int32)
	_ = binary.Write(w.b, binary.BigEndian, i)
	return w
}

// WriteInt64 写入 Int64
func (w *writer) WriteInt64(i int64) *writer {
	if i < 2147483648 && i > -2147483649 {
		return w.WriteInt32(int32(i))
	}
	w.writeKey(Int64)
	_ = binary.Write(w.b, binary.BigEndian, i)
	return w
}

// WriteFloat32 写入 Float32
func (w *writer) WriteFloat32(i float32) *writer {
	w.writeKey(Float32)
	_ = binary.Write(w.b, binary.BigEndian, i)
	return w
}

// WriteFloat64 写入 Float64
func (w *writer) WriteFloat64(i float64) *writer {
	w.writeKey(Float64)
	_ = binary.Write(w.b, binary.BigEndian, i)
	return w
}

// WriteString 写入 String
func (w *writer) WriteString(s string) *writer {
	if len(s) > 255 {
		w.writeKey(String2)
		_ = binary.Write(w.b, binary.BigEndian, len(s))
		w.b.WriteString(s)
		return w
	}
	w.writeKey(String)
	w.b.WriteByte(byte(len(s)))
	w.b.WriteString(s)
	return w
}
