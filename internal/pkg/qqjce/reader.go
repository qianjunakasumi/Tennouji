package qqjce

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
)

// reader 读取器
type reader struct {
	r *bytes.Reader // 读取器
}

// NewReader 返回一个读取器
func NewReader(b []byte) *reader { return &reader{bytes.NewReader(b)} }

// Read 读取 结构体
func (r *reader) Read(inter interface{}) {

	var (
		Type  = reflect.TypeOf(inter).Elem()
		Value = reflect.ValueOf(inter).Elem()
	)

	for i := 0; i < Value.NumField(); i++ {

		Data, Tag, _ := r.ReadAny()

		if j, ok := Type.Field(i).Tag.Lookup("jce"); ok {
			t, _ := strconv.ParseUint(j, 10, 8)
			if Tag < uint8(t) {
				i--
				continue
			}
		}

		Field := Value.Field(i)
		switch v := reflect.ValueOf(Data); v.Kind() {
		case reflect.Uint8:
			Field.SetInt(int64(v.Uint()))
		case reflect.Int16, reflect.Int32, reflect.Int64:
			Field.SetInt(v.Int())
		default:
			Field.Set(v)
		}
	}
}

// readKey 读取键
func (r *reader) readKey() (Type byte, tag uint8, o []byte) {
	b, _ := r.r.ReadByte()
	o = append(o, b)
	Type = b & 0xF
	tag = (b & 0xF0) >> 4
	if tag == 0xF {
		b, _ = r.r.ReadByte()
		o = append(o, b)
		tag = b & 0xFF
		return
	}
	return
}

// ReadByte 读取 Byte
func (r *reader) ReadByte() (b byte, o []byte) {
	b, _ = r.r.ReadByte()
	return b, []byte{b}
}

// ReadInt16 读取 Int16
func (r *reader) ReadInt16() (i int16, o []byte) {
	b := make([]byte, 2)
	_, _ = r.r.Read(b)
	_ = binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	return i, b
}

// ReadInt32 读取 Int32
func (r *reader) ReadInt32() (i int32, o []byte) {
	b := make([]byte, 4)
	_, _ = r.r.Read(b)
	_ = binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	return i, b
}

// ReadInt64 读取 Int64
func (r *reader) ReadInt64() (i int64, o []byte) {
	b := make([]byte, 8)
	_, _ = r.r.Read(b)
	_ = binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	return i, b
}

// ReadFloat32 读取 Float32
func (r *reader) ReadFloat32() (f float32, o []byte) {
	b := make([]byte, 4)
	_, _ = r.r.Read(b)
	_ = binary.Read(bytes.NewReader(b), binary.BigEndian, &f)
	return f, b
}

// ReadFloat64 读取 Float64
func (r *reader) ReadFloat64() (f float64, o []byte) {
	b := make([]byte, 8)
	_, _ = r.r.Read(b)
	_ = binary.Read(bytes.NewReader(b), binary.BigEndian, &f)
	return f, b
}

// ReadString 读取 String
func (r *reader) ReadString() (s string, o []byte) { // TODO 提高复用
	l, or := r.ReadByte()
	o = append(o, or...)
	b := make([]byte, l)
	_, _ = r.r.Read(b)
	o = append(o, b...)
	return string(b), o
}

// ReadString2 读取 String2
func (r *reader) ReadString2() (s string, o []byte) {
	l, or := r.ReadInt32()
	o = append(o, or...)
	b := make([]byte, l)
	_, _ = r.r.Read(b)
	o = append(o, b...)
	return string(b), o
}

// readMap 读取 Map
func (r *reader) readMap() (m map[string][]byte, o []byte) {
	m = map[string][]byte{}
	l, or := r.readLength()
	o = append(o, or...)
	for i := 0; i < int(l); i++ {
		k, _, v := r.ReadAny()
		o = append(o, v...)
		_, _, v = r.ReadAny()
		o = append(o, v...)
		m[reflect.ValueOf(k).String()] = v
	}
	return
}

// readSlice 读取 Slice
func (r *reader) readSlice() (s [][]byte, o []byte) {
	l, or := r.readLength()
	o = append(o, or...)
	for i := 0; i < int(l); i++ {
		_, _, v := r.ReadAny()
		s = append(s, v)
		o = append(o, v...)
	}
	return
}

// readStruct 读取 Struct
func (r *reader) readStruct() (b []byte, o []byte) {
	var d []byte
	for {
		_, t, or := r.ReadAny()
		o = append(o, or...)
		if t == 0 {
			b = d
			return
		}
		d = append(d, or...)
	}
}

// readBytes 读取 Bytes
func (r *reader) readBytes() (b []byte, o []byte) {
	_, _, or := r.readKey()
	l, orr := r.readLength()
	o = append(or, orr...)
	b = make([]byte, l)
	_, _ = r.r.Read(b)
	o = append(o, b...)
	return
}

// ReadAny 读取任意类型
func (r *reader) ReadAny() (d interface{}, tag uint8, o []byte) {

	var (
		Type byte
		or   []byte
	)

	Type, tag, o = r.readKey()
	switch Type {
	case Byte: // 和 Bool
		d, or = r.ReadByte()
	case Int16:
		d, or = r.ReadInt16()
	case Int32:
		d, or = r.ReadInt32()
	case Int64:
		d, or = r.ReadInt64()
	case Float32:
		d, or = r.ReadFloat32()
	case Float64:
		d, or = r.ReadFloat64()
	case String:
		d, or = r.ReadString()
	case String2:
		d, or = r.ReadString2()
	case Map:
		d, or = r.readMap()
	case Slice:
		d, or = r.readSlice()
	case Begin:
		d, or = r.readStruct()
	case End: // 空白
	case Bytes:
		d, or = r.readBytes()
	case Zero:
		d = byte(0)
	}
	o = append(o, or...)
	return
}

// readLength 读取长度
func (r *reader) readLength() (u uint32, o []byte) {

	t, _, or := r.readKey()
	o = append(o, or...)

	switch t {
	case Byte:
		b, _ := r.ReadByte()
		o = append(o, b)
		u = uint32(b)
	case Int16:
		i, orr := r.ReadInt16()
		o = append(o, orr...)
		u = uint32(i)
	case Int32:
		i, orr := r.ReadInt32()
		o = append(o, orr...)
		u = uint32(i)
	default:
		u = 0
	}
	return
}
