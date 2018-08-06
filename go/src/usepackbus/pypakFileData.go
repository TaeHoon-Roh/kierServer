package usepackbus

import (
	"fmt"
	"time"
	"strings"
)

type FileDef struct {
	DirVersion      byte
	FileDef         []FileHeader
	FileDef_Counter int
}

type FileHeader struct {
	FileName   string
	FileSize   uint
	LastUpdate time.Time
	Attribute  []byte
}

func (f *FileDef) Print_FileDir () {

	for i := 0; i < f.FileDef_Counter; i++ {
		fmt.Println("+++File Dir Info+++")
		fmt.Print("Table Name : ", f.FileDef[i].FileName)
		fmt.Print("  Last Update : ", f.FileDef[i].LastUpdate)
		fmt.Println()
	}

}

func (f *FileDef) Parse_FileDir(raw []byte) {

	offset := 0
	f.DirVersion = raw[offset]
	offset += 1

	f.FileDef = make([]FileHeader, 0)
	f.FileDef_Counter = 0
	for {
		buffer := FileHeader{}
		decode_buffer, size := Decode_bin("ASCIIZ", raw[offset:], len(raw[offset:]))
		buffer.FileName = decode_buffer.(string)
		offset += size
		fmt.Println("FileName : ", buffer.FileName)

		if buffer.FileName == "" {
			break
		}
		decode_buffer, size = Decode_bin("UInt4", raw[offset:], len(raw[offset:]))
		buffer.FileSize = decode_buffer.(uint)
		offset += size
		fmt.Println("FileSize : ", buffer.FileSize)
		decode_buffer, size = Decode_bin("ASCIIZ", raw[offset:], len(raw[offset:]))
		aaa := decode_buffer.(string)
		aaa = strings.Replace(aaa, " ", "T", 1)
		aaa = aaa + "Z"
		t,_ := time.Parse(time.RFC3339, aaa)
		buffer.LastUpdate = t
		offset += size
		fmt.Println("LastUpdate : ", buffer.LastUpdate)

		buffer.Attribute = make([]byte, 0)
		for i := 0; i < 12; i++ {
			temp := raw[offset]
			offset += 1
			if temp == '\000' {
				break
			} else {
				buffer.Attribute = append(buffer.Attribute, temp)
			}

		}
		fmt.Println("Attribute : ", buffer.Attribute)
		f.FileDef = append(f.FileDef, buffer)
		f.FileDef_Counter++
	}
}
