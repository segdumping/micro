package util

import (
	"reflect"
	"strconv"
	"unsafe"
)

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case int8:
		return strconv.Itoa(int(v))
	case uint8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case uint16:
		return strconv.Itoa(int(v))
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 2, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case bool:
		return strconv.FormatBool(v)
	}

	return ""
}

func Byte2String(b []byte) string {
	s := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h := reflect.StringHeader{Data: s.Data, Len: s.Len}
	return *(*string)(unsafe.Pointer(&h))
}

//return value must be readonly
func String2Byte(s string) []byte {
	b := (*reflect.StringHeader)(unsafe.Pointer(&s))
	h := reflect.SliceHeader{Data: b.Data, Len: b.Len, Cap: b.Len}

	return *(*[]byte)(unsafe.Pointer(&h))
}
