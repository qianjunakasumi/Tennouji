package qqtlv

import (
	"bytes"
	"encoding/binary"
)

// writer 写入器
type writer struct {
	b    *bytes.Buffer // 缓冲区
	Type uint16        // 类型
}

// NewWriter 返回一个写入器
func NewWriter(Type uint16) *writer { return &writer{bytes.NewBuffer(nil), Type} }

// WriteByte 写入 byte
func (w *writer) WriteByte(b byte) *writer {
	w.b.WriteByte(b)
	return w
}

// WriteBool 写入 bool
func (w *writer) WriteBool(b bool) *writer {
	if b {
		w.b.WriteByte(1)
		return w
	}
	w.b.WriteByte(0)
	return w
}

// WriteBytes 写入 []byte
func (w *writer) WriteBytes(b []byte) *writer {
	w.b.Write(b)
	return w
}

// WriteUint16 写入 uint16
func (w *writer) WriteUint16(u uint16) *writer {
	_ = binary.Write(w.b, binary.BigEndian, u)
	return w
}

// WriteUint32 写入 uint32
func (w *writer) WriteUint32(u uint32) *writer {
	_ = binary.Write(w.b, binary.BigEndian, u)
	return w
}

// WriteUint64 写入 uint64
func (w *writer) WriteUint64(u uint64) *writer {
	_ = binary.Write(w.b, binary.BigEndian, u)
	return w
}

// WriteString 写入 string
func (w *writer) WriteString(s string) *writer {
	w.WriteUint16(uint16(len(s)))
	w.b.WriteString(s)
	return w
}

// WriteLongString 写入长 string
func (w *writer) WriteLongString(s string) *writer {
	w.WriteUint32(uint32(len(s) + 4))
	w.b.WriteString(s)
	return w
}

// BytesWithTLV 返回 TLV 结构 []byte
func (w *writer) BytesWithTLV() (b []byte) {
	r := w.b.Bytes()
	_ = binary.Write(w.b, binary.BigEndian, w.Type)         // T
	_ = binary.Write(w.b, binary.BigEndian, uint32(len(r))) // L
	_ = binary.Write(w.b, binary.BigEndian, r)              // V
	return
}

// BytesWithLV 返回 LV 结构 []byte
func (w *writer) BytesWithLV() []byte {
	b := bytes.NewBuffer(nil)
	r := w.b.Bytes()
	_ = binary.Write(b, binary.BigEndian, uint32(len(r))) // L
	_ = binary.Write(b, binary.BigEndian, r)              // V
	return b.Bytes()
}
