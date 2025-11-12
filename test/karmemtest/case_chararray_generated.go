package karmemtest

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_case_chararray = [48]byte{}
var _NullReader_case_chararray = karmem.NewReader(_Null_case_chararray[:])

func init() {

}

const (
	PacketIdentifierB = 9956155639035161441
)

type B struct {
	// fixed string of size 20
	ANC  string
	ANCS string
}

func NewB() B {
	return B{}
}

func (x *B) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierB
}

func (x *B) Reset() {
	x.Read((*BViewer)(unsafe.Pointer(&_Null_case_chararray[0])), _NullReader_case_chararray)
}

func (x *B) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *B) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(48)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(44))
	__ANCOffset := offset + 4
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 20 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[20]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_case_chararray[:20])
	}
	__ANCSSize := uint(1 * len(x.ANCS))
	__ANCSOffset, err := writer.Alloc(__ANCSSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+28, uint32(__ANCSOffset))
	writer.Write4At(offset+28+4, uint32(__ANCSSize))
	writer.Write4At(offset+28+4+4, 1)
	__ANCSSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANCS)), __ANCSSize, __ANCSSize}
	writer.WriteAt(__ANCSOffset, *(*[]byte)(unsafe.Pointer(&__ANCSSlice)))

	return offset, nil
}

func (*B) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(48)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(44))
	writer.WriteAt(offset, _Null_case_chararray[:size-4])
	__ANCOffset := offset + 4
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)
	__ANCSOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+28, uint32(__ANCSOffset))
	writer.Write4At(offset+28+4, 0)
	writer.Write4At(offset+28+4+4, 1)

	return offset, nil
}

func (x *B) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBViewer(reader, 0), reader)
}

func (x *B) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBViewer(reader, offset), reader)
}

func (x *B) Read(viewer *BViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
	__ANCSString := viewer.ANCS(reader)
	if x.ANCS != __ANCSString {
		__ANCSStringCopy := make([]byte, len(__ANCSString))
		copy(__ANCSStringCopy, __ANCSString)
		x.ANCS = *(*string)(unsafe.Pointer(&__ANCSStringCopy))
	}
}

type BViewer [48]byte

type BSource struct {
	*BViewer
}

func NewBViewer(reader *karmem.Reader, offset uint32) (v *BViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*BViewer)(unsafe.Pointer(&_Null_case_chararray[0]))
	}
	v = (*BViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BViewer)(unsafe.Pointer(&_Null_case_chararray[0]))
	}
	return v
}

func NewBSource(reader *karmem.Reader, offset uint32) (v BSource) {
	return BSource{NewBViewer(reader, offset)}
}

func newBSourceByViewer(viewer *BViewer) (v BSource) {
	return BSource{viewer}
}

func (x *BViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BViewer) ANC() (v string) {
	if 4+20 > x.size() {
		return v
	}
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+4)), uintptr(size), 20,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BSource) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.BViewer)) = uint32(len(v))
	copy((*[20]byte)(unsafe.Add(unsafe.Pointer(x.BViewer), 4+4))[:], v)
}

func (x *BViewer) ANCS(reader *karmem.Reader) (v string) {
	if 28+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 28))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 28+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BSource) SetANCS(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BViewer), 28))
	__ANCSSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BViewer), 28+4))
	if uint32(__ANCSSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSSize, __ANCSSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BViewer), offset)), __ANCSSize, __ANCSSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BViewer), 28+4)) = uint32(__ANCSSize)
	return nil
}
