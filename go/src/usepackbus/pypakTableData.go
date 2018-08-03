package usepackbus

import (
	"fmt"
	"encoding/binary"
	"unsafe"
)

type TableDef struct {
	FlsVersion  byte
	TabelHeader []TableHeader
	Table_Count int
}
type TableHeader struct {
	FlsVersion  byte
	TableName   string
	TableSize   uint
	TimeType    byte
	TblTimeInto []uint32
	TblInterval []uint32
	Table_Sig   uint16

	TableField        []TableFld
	Table_Field_Count int
}
type TableFld struct {
	FieldType   string
	ReadOnly    byte
	FieldName   string
	AliasName   []string
	Processing  string
	Units       string
	Description string
	BegIdx      uint
	Dimension   uint
	Subdim      []uint
}

type Frag struct {
	TableNbr   uint16
	BegRecNbr  uint32
	TableName  string
	IsOffset   byte
	ByteOffset uint32
	RecFrag    []byte

	NbrOfRecs   uint16
	TimeOfRec   uint64
	Recode      []RecodeDef
	RecodeCount int
}

type RecodeDef struct {
	RecNbr     uint32
	TimeOfRec  []uint32
	FieldCount uint32
	Field      []FieldDef
}

type FieldDef struct {
	FieldName string
	FieldType string
	Dimension uint
	Raw interface{}
}

func (t *TableDef) Print_TableDef() {
	for i := 0; i < t.Table_Count; i++ {
		fmt.Println("+++Table Def+++")
		fmt.Println("Table Name : ", t.TabelHeader[i].TableName)
		fmt.Printf("Table Sig : 0x%x\n", t.TabelHeader[i].Table_Sig)
		fmt.Println("Field Count : ", t.TabelHeader[i].Table_Field_Count)
		fmt.Println("Field Data Type", t.TabelHeader[i].TableField[0].FieldType)
	}
}

func (t *TableDef) Parse_TableDef(raw []byte) {

	var offset = 0
	t.FlsVersion = raw[offset]
	offset += 1
	t.Table_Count = 0

	t.TabelHeader = make([]TableHeader, 0)

	fmt.Println("Raw Size : ", len(raw))

	for offset < len(raw) {
		var start int
		start = offset

		TableHeader_Buffer := TableHeader{}
		decode_buffer, size := Decode_bin("ASCIIZ", raw[offset:], 1)
		TableHeader_Buffer.TableName = decode_buffer.(string)
		offset += size
		decode_buffer, size = Decode_bin("UInt4", raw[offset:], 1)
		TableHeader_Buffer.TableSize = decode_buffer.(uint)
		offset += size
		//fmt.Println("offset", offset)
		TableHeader_Buffer.TimeType = raw[offset]
		offset += 1
		decode_buffer, size = Decode_bin("NSec", raw[offset:], 1)
		timeinto := make([]uint32, 2)
		TableHeader_Buffer.TblTimeInto = make([]uint32, 2)
		timeinto[0] = binary.BigEndian.Uint32(raw[offset:offset+4])
		timeinto[1] = binary.BigEndian.Uint32(raw[offset+4:offset+8])
		TableHeader_Buffer.TblTimeInto = timeinto
		offset += size

		decode_buffer, size = Decode_bin("NSec", raw[offset:], 1)
		interval := make([]uint32, 2)
		TableHeader_Buffer.TblTimeInto = make([]uint32, 2)
		interval[0] = binary.BigEndian.Uint32(raw[offset:offset+4])
		interval[1] = binary.BigEndian.Uint32(raw[offset+4:offset+8])
		TableHeader_Buffer.TblInterval = interval
		offset += size

		TableHeader_Buffer.TableField = make([]TableFld, 0)
		TableHeader_Buffer.Table_Field_Count = 0
		for {
			//fmt.Println("---Field Number : ", t.table_Fld_count)
			TableFld_buff := TableFld{}
			fieldtype := raw[offset]
			offset += 1
			if fieldtype == 0 {
				break
			}

			TableFld_buff.ReadOnly = fieldtype >> 7

			for i := 0; i < len(GlobalDataType); i++ {
				if (GlobalDataType[i].Code == fieldtype&0x7F) {
					TableFld_buff.FieldType = GlobalDataType[i].Name
				}
			}

			//데이터 타입 정의해야 함
			decode_buffer, size = Decode_bin("ASCIIZ", raw[offset:], 1)
			TableFld_buff.FieldName = decode_buffer.(string)
			offset += size
			//fmt.Println("FieldName : ", TableFld_buff.FieldName)

			TableFld_buff.AliasName = make([]string, 1)
			for i := 0; ; i++ {
				decode_buffer, size = Decode_bin("ASCIIZ", raw[offset:], 1)
				aliasname := decode_buffer.(string)
				offset += size
				if aliasname == "" {
					break
				}
				TableFld_buff.AliasName = append(TableFld_buff.AliasName, aliasname)
			}
			//fmt.Println("Aliasname : ", TableFld_buff.AliasName)

			decode_buffer, size = Decode_bin("ASCIIZ", raw[offset:], 1)
			TableFld_buff.Processing = decode_buffer.(string)
			offset += size
			//fmt.Println("Processing : ", TableFld_buff.Processing)

			decode_buffer, size = Decode_bin("ASCIIZ", raw[offset:], 1)
			TableFld_buff.Units = decode_buffer.(string)
			offset += size
			//fmt.Println("Units : ", TableFld_buff.Units)

			decode_buffer, size = Decode_bin("ASCIIZ", raw[offset:], 1)
			TableFld_buff.Description = decode_buffer.(string)
			offset += size
			//fmt.Println("Description : ", TableFld_buff.Description)

			decode_buffer, size = Decode_bin("UInt4", raw[offset:], 1)
			TableFld_buff.BegIdx = decode_buffer.(uint)
			offset += size
			//fmt.Println("BegIdx : ", TableFld_buff.BegIdx)

			decode_buffer, size = Decode_bin("UInt4", raw[offset:], 1)
			TableFld_buff.Dimension = decode_buffer.(uint)
			offset += size
			//fmt.Println("Dimension : ", TableFld_buff.Dimension)

			TableFld_buff.Subdim = make([]uint, 1)
			for i := 0; ; i++ {
				decode_buffer, size = Decode_bin("UInt4", raw[offset:], 1)
				offset += size
				var subdim uint
				subdim = decode_buffer.(uint)
				if subdim == 0 {
					break
				}
				TableFld_buff.Subdim = append(TableFld_buff.Subdim, subdim)

			}
			//fmt.Println("Subdim : ", TableFld_buff.Subdim)

			TableHeader_Buffer.TableField = append(TableHeader_Buffer.TableField, TableFld_buff)
			TableHeader_Buffer.Table_Field_Count++
		}
		TableHeader_Buffer.Table_Sig = CalcSigFor(raw[start:offset], 0xAAAA)
		t.TabelHeader = append(t.TabelHeader, TableHeader_Buffer)
		t.Table_Count++
	}
	fmt.Println("End offset ", offset)
}

func (t *TableDef) GetTableNbr(str string) (uint16) {
	for i := 0; i < len(t.TabelHeader); i++ {
		if t.TabelHeader[i].TableName == str {
			return uint16(i)
			break
		}
	}
	return 0
}

func (t *TableDef) Parse_CollectData(raw []byte, FieldNbr int) (Frag) {
	offset := 0

	var frag Frag

	frag.RecodeCount = 0

	fmt.Println("Start Parse Collect Data!!", len(raw))
	for offset < len(raw)-1 {
		decode_buffer, size := Decode_bin("UInt2", raw[offset:], 1)
		frag.TableNbr = decode_buffer.(uint16)
		offset += size
		fmt.Println("Frag.TableNbr : ", frag.TableNbr)

		decode_buffer, size = Decode_bin("UInt4", raw[offset:], 1)
		frag.BegRecNbr = uint32(decode_buffer.(uint))
		offset += size
		fmt.Println("Ffrag.BegRecNbr : ", frag.BegRecNbr)

		fmt.Println("Check Header Size : ", len(t.TabelHeader))
		frag.TableName = t.TabelHeader[frag.TableNbr-1].TableName

		temp := raw[offset]
		frag.IsOffset = temp >> 7

		fmt.Println("IsOffset : ", frag.IsOffset)
		//동적 할당
		frag.Recode = make([]RecodeDef, 0)

		if frag.IsOffset != 0 {
			decode_buffer, size = Decode_bin("UInt4", raw[offset:], 1)
			byteoffset := decode_buffer.(uint)
			frag.ByteOffset = uint32(byteoffset & 0x7FFFFFFF)
			frag.NbrOfRecs = '\000'
			offset += size
			fmt.Println("frag.ByteOffset : ", frag.ByteOffset)
			frag.RecFrag = raw[offset : len(raw)-1]
			offset += len(frag.RecFrag)
			fmt.Println("frag.RecFrag : ", frag.RecFrag)
		} else {
			decode_buffer, size := Decode_bin("UInt2", raw[offset:], 1)
			nbrofrecs := decode_buffer.(uint16)
			frag.NbrOfRecs = nbrofrecs & 0x7FFF
			frag.ByteOffset = '\000'
			offset += size
			fmt.Println("frag.NbrOfRecs : ", frag.NbrOfRecs)
			interval := make([]uint32, 2)
			interval[0] = t.TabelHeader[frag.TableNbr-1].TblInterval[0]
			interval[1] = t.TabelHeader[frag.TableNbr-1].TblInterval[1]
			fmt.Println("Interval : ", interval)

			timeofrec := make([]uint32, 2)
			if interval[0] == 0 && interval[1] == 0 {
				timeofrec[0] = '\000'
				timeofrec[1] = '\000'
			} else {
				decode_buffer, size = Decode_bin("NSec", raw[offset:], 1)
				timeofrec[0] = binary.BigEndian.Uint32(raw[offset:offset+4])
				timeofrec[1] = binary.BigEndian.Uint32(raw[offset+4:offset+8])
				offset += size
			}
			fmt.Println("NsecCheck 1 : ", decode_buffer)

			recode := RecodeDef{}
			for i := 0; i < int(frag.NbrOfRecs); i++ {
				recode.RecNbr = frag.BegRecNbr + uint32(i)
				if timeofrec[0] != 0 || timeofrec[1] != 0 {
					recode.TimeOfRec = make([]uint32, 2)
					recode.TimeOfRec[0] = timeofrec[0] + uint32(i)*interval[0]
					recode.TimeOfRec[1] = timeofrec[1] + uint32(i)*interval[1]
				} else {
					decode_buffer, size = Decode_bin("NSec", raw[offset:], 1)
					recode.TimeOfRec = make([]uint32, 2)
					recode.TimeOfRec[0] = binary.BigEndian.Uint32(raw[offset:offset+4])
					recode.TimeOfRec[1] = binary.BigEndian.Uint32(raw[offset+4:offset+8])
					offset += size
				}
			}
			fmt.Println("NsecCheck 2 : ", recode.TimeOfRec)

			recode.Field = make([]FieldDef, 0)
			field := FieldDef{}
			for i := 0; i < t.TabelHeader[frag.TableNbr-1].Table_Field_Count; i++ {
				field.FieldName = t.TabelHeader[frag.TableNbr-1].TableField[i].FieldName
				field.FieldType = t.TabelHeader[frag.TableNbr-1].TableField[i].FieldType
				field.Dimension = t.TabelHeader[frag.TableNbr-1].TableField[i].Dimension
				fmt.Println("Field Data : ", field)

				if field.FieldType == "ASCII" {
					decode_buffer, size = Decode_bin("ASCII", raw[offset:], int(field.Dimension))
					field.Raw = decode_buffer
					fmt.Println("DecodeBuffer size :", unsafe.Sizeof(decode_buffer))
					offset += size
				} else {
					decode_buffer, size = Decode_bin(field.FieldType, raw[offset:], int(field.Dimension))
					field.Raw = decode_buffer
					fmt.Println("DecodeBuffer size :", unsafe.Sizeof(decode_buffer))
					offset += size
				}
				recode.Field = append(recode.Field, field)
				recode.FieldCount++
				fmt.Println("Check Field Data : ", field)
			}
			frag.Recode = append(frag.Recode, recode)
			frag.RecodeCount++
			fmt.Println("Check Field Data : ", recode.Field)
		}

	}

	return frag

}



func Print_Frag(frag Frag) {
	for i := 0; i < frag.RecodeCount; i++ {
		fmt.Println("+++++Print Collect Data+++++")
		fmt.Println(frag)
	}
}
