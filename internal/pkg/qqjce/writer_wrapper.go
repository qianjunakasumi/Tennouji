package qqjce

import (
	"bytes"
)

// BytesWithPack 返回 包 []byte
func (w *writer) BytesWithPack() []byte {
	b := bytes.NewBuffer([]byte{0x0A})
	b.Write(w.b.Bytes())
	b.WriteByte(0x0B)
	return b.Bytes()
}

// WriteWithDataV2 写入 DataV2
func (w *writer) WriteWithDataV2(wr *writer, k1, k2 string) *writer {
	return NewWriter().WriteMap(DataV2{k1: {k2: wr.BytesWithPack()}})
}

// WriteWithDataV3 写入 DataV3
func (w *writer) WriteWithDataV3(wr *writer, k string) *writer {
	return NewWriter().WriteMap(DataV3{k: wr.BytesWithPack()})
}
