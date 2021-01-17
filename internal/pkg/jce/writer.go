package jce

import "bytes"

// writer 写入器
type writer struct{ b *bytes.Buffer }

// NewWriter 返回一个写入器
func NewWriter() *writer { return &writer{bytes.NewBuffer(nil)} }

// writeKey 写入键
func (w *writer) writeKey(Type byte, tag uint8) {
	if tag > 14 {
		w.b.WriteByte(Type | 0xF0)
		w.b.WriteByte(tag)
		return
	}
	w.b.WriteByte(tag<<4 | Type)
}
