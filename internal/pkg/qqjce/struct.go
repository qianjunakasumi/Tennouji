package qqjce

const (
	Byte    = 0
	Bool    = 0
	Int16   = 1
	Int32   = 2
	Int64   = 3
	Float32 = 4
	Float64 = 5
	String  = 6
	String2 = 7
	Map     = 8
	Slice   = 9
	Begin   = 10
	End     = 11
	Zero    = 12
	Bytes   = 13
)

type (
	Packet struct { // Packet 包
		Version    int64 `jce:"1"`
		PacketType int64
		MsgType    int64
		ReqID      int64
		Controller string
		Method     string
		Data       []byte
		Timeout    int64
		Context    map[string][]byte
		Status     map[string][]byte
	}

	DataV3 map[string][]byte            // DataV3 请求数据 三代
	DataV2 map[string]map[string][]byte // DataV2 请求数据 二代
)
