package usepackbus

import (
	"encoding/binary"
	"math"
)

type DataTypeDefInterFace interface{
	Out_Byte() byte
	Out_Uint2() uint16
	Out_Uint4() uint
	Out_Int1() int8
	Out_Int2() int16
	Out_Int4() int
	Out_FP2() uint16
	Out_FP3() string
	Out_FP4() string
	Out_IEEE4B() float32
	Out_IEEE8B() float64
	Out_Bool8() byte
	Out_Bool() byte
	Out_Bool2() uint16
	Out_Bool4() uint
	Out_Sec() int
	Out_USec() string
	Out_NSec() int64
	Out_ASCII() string
	Out_ASCIIZ() string
	Out_Short() int16
	Out_Long() int
	Out_UShort() uint16
	Out_ULong() uint
	Out_IEEE4L() float32
	Out_IEEE8L() float64
	Out_SecNano() int64
}

func (t *TableDef) Out_Byte(result []byte) (byte, int){
	return result[0], 1
}

func (t *TableDef) Out_Uint2(result []byte) (uint16, int){
	return binary.BigEndian.Uint16(result[:2]), 2
}

func (t *TableDef) Out_Uint4(result []byte) (uint, int){
	return uint(binary.BigEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_Int1(result []byte) (int8, int){
	return int8(result[0]), 1
}

func (t *TableDef) Out_Int2(result []byte) (int16, int){
	return int16(binary.BigEndian.Uint16(result[:2])), 2
}

func (t *TableDef) Out_Int4(result []byte) (int, int){
	return int(binary.BigEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_FP2(result []byte) (float32, int){
	fp2 := binary.BigEndian.Uint16(result[:2])
	mant := fp2 & 0x1FFF
	exp := fp2 >> 13 & 0x3
	sign := fp2 >> 15
	value := (float32(math.Pow(-1, float64(sign))) * float32(mant)) / float32(math.Pow(10,float64(exp)))
	return value, 2
}

func (t *TableDef) Out_FP3(result []byte) (string, int){
	return string(result[:3]), 3
}

func (t *TableDef) Out_FP4(result []byte) (string, int){
	return string(result[:4]), 4
}

func (t *TableDef) Out_IEEE4B(result []byte) (float32, int){
	return math.Float32frombits(binary.BigEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_IEEE8B(result []byte) (float64, int){
	return math.Float64frombits(binary.BigEndian.Uint64(result[:8])) , 8
}

func (t *TableDef) Out_Bool8(result []byte) (byte, int){
	return result[0], 1
}

func (t *TableDef) Out_Bool(result []byte) (byte, int){
	return result[0], 1
}

func (t *TableDef) Out_Bool2(result []byte) (uint16, int){
	return binary.BigEndian.Uint16(result[:2]), 2
}

func (t *TableDef) Out_Bool4(result []byte) (uint, int){
	return uint(binary.BigEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_Sec(result []byte) (int, int){
	return int(binary.BigEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_USec(result []byte) (string, int){
	return string(result[:6]), 6
}

func (t *TableDef) Out_NSec(result []byte) (int64, int){
	return int64(binary.BigEndian.Uint64(result[:8])), 8
}

func (t *TableDef) Out_ASCII(result []byte, length int) (string, int){
	return string(result[:length]), length
}

func (t *TableDef) Out_ASCIIZ(result []byte) (string, int){
	var nul int
	for i := 0; i < len(result); i++ {
		if '\000' == result[i] {
			nul = i
			break
		}
	}
	value := (result[:nul])
	size := len(value) + 1
	return string(value), size
}

func (t *TableDef) Out_Short(result []byte) (int16, int){
	return int16(binary.LittleEndian.Uint16(result[:2])), 2
}

func (t *TableDef) Out_Long(result []byte) (int, int){
	return int(binary.LittleEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_UShort(result []byte) (uint16, int){
	return uint16(binary.LittleEndian.Uint16(result[:2])), 2
}

func (t *TableDef) Out_ULong(result []byte) (uint, int){
	return uint(binary.LittleEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_IEEE4L(result []byte) (float32, int){
	return math.Float32frombits(binary.LittleEndian.Uint32(result[:4])), 4
}

func (t *TableDef) Out_IEEE8L(result []byte) (float64, int){
	return math.Float64frombits(binary.LittleEndian.Uint64(result[:8])) , 8
}

func (t *TableDef) Out_SecNano(result []byte) (int64, int){
	return int64(binary.LittleEndian.Uint64(result[:8])), 8
}

func Decode_bin(Type string, buff []byte, length int) (interface{}, int) {
	t := TableDef{}
	switch Type {
	case "Byte":
		return t.Out_Byte(buff)
	case "UInt2":
		return t.Out_Uint2(buff)
	case "UInt4":
		return t.Out_Uint4(buff)
	case "Int1":
		return t.Out_Int1(buff)
	case "Int2":
		return t.Out_Int2(buff)
	case "Int4":
		return t.Out_Int4(buff)
	case "FP2":
		return t.Out_FP2(buff)
	case "FP3":
		return t.Out_FP3(buff)
	case "FP4":
		return t.Out_FP4(buff)
	case "IEEE4B":
		return t.Out_IEEE4B(buff)
	case "IEEE8B":
		return t.Out_IEEE8B(buff)
	case "Bool8":
		return t.Out_Bool8(buff)
	case "Bool":
		return t.Out_Bool(buff)
	case "Bool2":
		return t.Out_Bool2(buff)
	case "Bool4":
		return t.Out_Bool4(buff)
	case "Sec":
		return t.Out_Sec(buff)
	case "USec":
		return t.Out_USec(buff)
	case "NSec":
		return t.Out_NSec(buff)
	case "ASCII":
		return t.Out_ASCII(buff, length)
	case "ASCIIZ":
		return t.Out_ASCIIZ(buff)
	case "Short":
		return t.Out_Short(buff)
	case "Long":
		return t.Out_Long(buff)
	case "UShort":
		return t.Out_UShort(buff)
	case "ULong":
		return t.Out_ULong(buff)
	case "IEEE4L":
		return t.Out_IEEE4L(buff)
	case "IEEE8L":
		return t.Out_IEEE8L(buff)
	case "SecNano":
		return t.Out_SecNano(buff)
	default:
		return nil, 0
	}

}