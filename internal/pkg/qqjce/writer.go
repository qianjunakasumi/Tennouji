package qqjce

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
)

// writer 写入器
type writer struct {
	b   *bytes.Buffer // 缓冲区
	tag uint8         // 标签
}

// NewWriter 返回一个写入器
func NewWriter(tag ...uint8) *writer {
	var t uint8 = 0
	if len(tag) != 0 {
		t = tag[0]
	}
	return &writer{bytes.NewBuffer(nil), t}
}

// writeKey 写入键
func (w *writer) writeKey(Type byte) *writer {
	if w.tag > 14 {
		w.b.WriteByte(Type | 0xF0)
		w.b.WriteByte(w.tag)
	} else {
		w.b.WriteByte(w.tag<<4 | Type)
	}
	w.tag++
	return w
}

// SetTag 设置标签
func (w *writer) SetTag(u uint8) *writer {
	w.tag = u
	return w
}

// Write 写入 结构体
func (w *writer) Write(inter interface{}) *writer {
	Type := reflect.TypeOf(inter)
	for i := 0; i < Type.NumField(); i++ {
		if jce := Type.Field(i).Tag.Get("jce"); jce != "" {
			id, err := strconv.ParseUint(jce, 10, 8)
			if err != nil {
				continue
			}
			w.SetTag(uint8(id))
		}
		w.WriteAny(reflect.ValueOf(inter).Field(i).Interface())
	}
	return w
}

// WriteByte 写入 Byte
func (w *writer) WriteByte(b byte) *writer {
	if b == 0 {
		w.writeKey(Zero)
		return w
	}
	w.writeKey(Byte).b.WriteByte(b)
	return w
}

// WriteBool 写入 Bool
func (w *writer) WriteBool(b bool) *writer {
	if b {
		w.writeKey(Bool).b.WriteByte(1)
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
		_ = binary.Write(w.b, binary.BigEndian, int32(len(s)))
		w.b.WriteString(s)
		return w
	}
	w.writeKey(String).b.WriteByte(byte(len(s)))
	w.b.WriteString(s)
	return w
}

// WriteMap 写入 Map
func (w *writer) WriteMap(i interface{}) *writer {
	m := reflect.ValueOf(i)
	keys := m.MapKeys()
	w.writeKey(Map).b.Write(NewWriter(0).WriteInt32(int32(len(keys))).Bytes())
	for _, key := range keys {
		w.b.Write(NewWriter(0).WriteAny(key.Interface()).WriteAny(m.MapIndex(key).Interface()).Bytes())
	}
	return w
}

// WriteMap 写入 Slice
func (w *writer) WriteSlice(i interface{}) *writer {
	s := reflect.ValueOf(i)
	length := s.Len()
	w.writeKey(Slice).b.Write(NewWriter(0).WriteInt32(int32(length)).Bytes())
	for i := 0; i < length; i++ {
		w.b.Write(NewWriter(0).WriteAny(s.Index(i).Interface()).Bytes())
	}
	return w
}

// WriteStruct 写入 Struct
func (w *writer) WriteStruct(i interface{}) *writer {
	w.writeKey(Begin).Write(i).writeKey(End)
	return w
}

// WriteBytes 写入 Bytes
func (w *writer) WriteBytes(b []byte) *writer {
	w.writeKey(Bytes).b.Write(NewWriter(0).writeKey(Byte).SetTag(0).WriteInt32(int32(len(b))).Bytes())
	w.b.Write(b)
	return w
}

// Bytes 返回 []byte
func (w *writer) Bytes() []byte { return w.b.Bytes() }

// BytesWithPack 返回 包 []byte
func (w *writer) BytesWithPack() []byte {
	b := bytes.NewBuffer([]byte{0x0A})
	b.Write(w.b.Bytes())
	b.WriteByte(0x0B)
	return b.Bytes()
}

// WriteAny 写入任意类型
func (w *writer) WriteAny(i interface{}) *writer {
	switch o := i.(type) {
	case byte:
		w.WriteByte(o)
	case bool:
		w.WriteBool(o)
	case int16:
		w.WriteInt16(o)
	case int32:
		w.WriteInt32(o)
	case int64:
		w.WriteInt64(o)
	case float32:
		w.WriteFloat32(o)
	case float64:
		w.WriteFloat64(o)
	case string:
		w.WriteString(o)
	case []byte:
		w.WriteBytes(o)
	default:
		switch reflect.TypeOf(i).Kind() {
		case reflect.Map:
			w.WriteMap(o)
		case reflect.Slice:
			w.WriteSlice(o)
		case reflect.Struct:
			w.WriteStruct(o)
		}
	}
	return w
}
