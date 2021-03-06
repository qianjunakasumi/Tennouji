package qqjce

// ReadMap 读取 Map
func (r *reader) ReadMap() (m map[string][]byte) {
	r.readKey()
	m, _ = r.readMap()
	return
}

// Deprecated: ReadSlice 读取 Slice
func (r *reader) ReadSlice() (s [][]byte) {
	r.readKey()
	s, _ = r.readSlice()
	return
}

// ReadStruct 读取 Struct
func (r *reader) ReadStruct() (b []byte) {
	r.readKey()
	b, _ = r.readStruct()
	return
}

// ReadBytes 读取 Bytes
func (r *reader) ReadBytes() (s []byte) {
	r.readKey()
	s, _ = r.readBytes()
	return
}

// UnPack 解包
func UnPack(r *reader) []byte {
	b := r.ReadBytes()
	return b[1 : len(b)-1]
}

// ReadWithDataV2 读取 DataV2
func (r *reader) ReadWithDataV2(k1, k2 string) []byte {
	return UnPack(NewReader(NewReader(r.ReadMap()[k1]).ReadMap()[k2]))
}

// ReadWithDataV3 读取 DataV3
func (r *reader) ReadWithDataV3(k string) []byte {
	return UnPack(NewReader(r.ReadMap()[k]))
}
