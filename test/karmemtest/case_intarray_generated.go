package karmemtest

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_case_intarray = [176]byte{}
var _NullReader_case_intarray = karmem.NewReader(_Null_case_intarray[:])

func init() {

}

const (
	PacketIdentifierBI = 12307522458191621139
)

type BI struct {
	ANC  [20]int64
	ANCS string
}

func NewBI() BI {
	return BI{}
}

func (x *BI) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBI
}

func (x *BI) Reset() {
	x.Read((*BIViewer)(unsafe.Pointer(&_Null_case_intarray[0])), _NullReader_case_intarray)
}

func (x *BI) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BI) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(176)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(176))
	__ANCOffset := offset + 4
	writer.WriteAt(__ANCOffset, (*[160]byte)(unsafe.Pointer(&x.ANC))[:])
	__ANCSSize := uint(1 * len(x.ANCS))
	__ANCSOffset, err := writer.Alloc(__ANCSSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+164, uint32(__ANCSOffset))
	writer.Write4At(offset+164+4, uint32(__ANCSSize))
	__ANCSSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANCS)), __ANCSSize, __ANCSSize}
	writer.WriteAt(__ANCSOffset, *(*[]byte)(unsafe.Pointer(&__ANCSSlice)))

	return offset, nil
}

func (*BI) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(176)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(176))
	writer.WriteAt(offset, _Null_case_intarray[:size-4])
	__ANCOffset := offset + 4
	_ = __ANCOffset
	__ANCSOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+164, uint32(__ANCSOffset))
	writer.Write4At(offset+164+4, 0)

	return offset, nil
}

func (x *BI) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBIViewer(reader, 0), reader)
}

func (x *BI) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBIViewer(reader, offset), reader)
}

func (x *BI) Read(viewer *BIViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__ANCSlice := viewer.ANC()
	__ANCLen := len(__ANCSlice)
	if __ANCLen > 20 {
		__ANCLen = 20
	}
	copy(x.ANC[:], __ANCSlice)
	for i := __ANCLen; i < len(x.ANC); i++ {
		x.ANC[i] = 0
	}
	__ANCSString := viewer.ANCS(reader)
	if x.ANCS != __ANCSString {
		__ANCSStringCopy := make([]byte, len(__ANCSString))
		copy(__ANCSStringCopy, __ANCSString)
		x.ANCS = *(*string)(unsafe.Pointer(&__ANCSStringCopy))
	}
}

type BIViewer [176]byte

type BISource struct {
	*BIViewer
}

func NewBIViewer(reader *karmem.Reader, offset uint32) (v *BIViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*BIViewer)(unsafe.Pointer(&_Null_case_intarray[0]))
	}
	v = (*BIViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BIViewer)(unsafe.Pointer(&_Null_case_intarray[0]))
	}
	return v
}

func NewBISource(reader *karmem.Reader, offset uint32) (v BISource) {
	return BISource{NewBIViewer(reader, offset)}
}

func newBISourceByViewer(viewer *BIViewer) (v BISource) {
	return BISource{viewer}
}

func (x *BIViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BIViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BIViewer) ANC() (v []int64) {
	if 4+160 > x.size() {
		return []int64{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4)), 20, 20,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *BISource) SetANC(v [20]int64) {
	copy((*(*[20]int64)(unsafe.Add(unsafe.Pointer(x.BIViewer), 4)))[:], v[:])
}

func (x *BIViewer) ANCS(reader *karmem.Reader) (v string) {
	if 164+8 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 164))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 164+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BISource) SetANCS(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BIViewer), 164))
	__ANCSSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BIViewer), 164+4))
	if uint32(__ANCSSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSSize, __ANCSSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BIViewer), offset)), __ANCSSize, __ANCSSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BIViewer), 164+4)) = uint32(__ANCSSize)
	return nil
}
