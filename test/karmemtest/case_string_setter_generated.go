package karmemtest

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_case_string_setter = [128]byte{}
var _NullReader_case_string_setter = karmem.NewReader(_Null_case_string_setter[:])

func init() {

}

const (
	PacketIdentifierBaseS = 9956155639035161441
)

type BaseS struct {
	Data     [10]int64
	ANC      string
	ANE      []uint8
	Username string
	AND      string
}

func NewBaseS() BaseS {
	return BaseS{}
}

func (x *BaseS) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBaseS
}

func (x *BaseS) Reset() {
	x.Read((*BaseSViewer)(unsafe.Pointer(&_Null_case_string_setter[0])), _NullReader_case_string_setter)
}

func (x *BaseS) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BaseS) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(128)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(128))
	__DataOffset := offset + 4
	writer.WriteAt(__DataOffset, (*[80]byte)(unsafe.Pointer(&x.Data))[:])
	__ANCSize := uint(1 * len(x.ANC))
	__ANCOffset, err := writer.Alloc(__ANCSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+84, uint32(__ANCOffset))
	writer.Write4At(offset+84+4, uint32(__ANCSize))
	__ANCSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANC)), __ANCSize, __ANCSize}
	writer.WriteAt(__ANCOffset, *(*[]byte)(unsafe.Pointer(&__ANCSlice)))
	__ANESize := uint(1 * len(x.ANE))
	__ANEOffset, err := writer.Alloc(__ANESize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+96, uint32(__ANEOffset))
	writer.Write4At(offset+96+4, uint32(__ANESize))
	__ANESlice := *(*[3]uint)(unsafe.Pointer(&x.ANE))
	__ANESlice[1] = __ANESize
	__ANESlice[2] = __ANESize
	writer.WriteAt(__ANEOffset, *(*[]byte)(unsafe.Pointer(&__ANESlice)))
	__UsernameSize := uint(1 * len(x.Username))
	__UsernameOffset, err := writer.Alloc(__UsernameSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+104, uint32(__UsernameOffset))
	writer.Write4At(offset+104+4, uint32(__UsernameSize))
	__UsernameSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Username)), __UsernameSize, __UsernameSize}
	writer.WriteAt(__UsernameOffset, *(*[]byte)(unsafe.Pointer(&__UsernameSlice)))
	__ANDSize := uint(1 * len(x.AND))
	__ANDOffset, err := writer.Alloc(__ANDSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+116, uint32(__ANDOffset))
	writer.Write4At(offset+116+4, uint32(__ANDSize))
	__ANDSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.AND)), __ANDSize, __ANDSize}
	writer.WriteAt(__ANDOffset, *(*[]byte)(unsafe.Pointer(&__ANDSlice)))

	return offset, nil
}

func (*BaseS) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(128)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(128))
	writer.WriteAt(offset, _Null_case_string_setter[:size-4])
	__DataOffset := offset + 4
	_ = __DataOffset
	__ANCOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+84, uint32(__ANCOffset))
	writer.Write4At(offset+84+4, 0)
	__ANEOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+96, uint32(__ANEOffset))
	writer.Write4At(offset+96+4, 0)
	__UsernameOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+104, uint32(__UsernameOffset))
	writer.Write4At(offset+104+4, 0)
	__ANDOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+116, uint32(__ANDOffset))
	writer.Write4At(offset+116+4, 0)

	return offset, nil
}

func (x *BaseS) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBaseSViewer(reader, 0), reader)
}

func (x *BaseS) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBaseSViewer(reader, offset), reader)
}

func (x *BaseS) Read(viewer *BaseSViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__DataSlice := viewer.Data()
	__DataLen := len(__DataSlice)
	if __DataLen > 10 {
		__DataLen = 10
	}
	copy(x.Data[:], __DataSlice)
	for i := __DataLen; i < len(x.Data); i++ {
		x.Data[i] = 0
	}
	__ANCString := viewer.ANC(reader)
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
	__ANESlice := viewer.ANE(reader)
	__ANELen := len(__ANESlice)
	if __ANELen > cap(x.ANE) {
		x.ANE = append(x.ANE, make([]uint8, __ANELen-len(x.ANE))...)
	}
	if __ANELen > len(x.ANE) {
		x.ANE = x.ANE[:__ANELen]
	}
	copy(x.ANE, __ANESlice)
	x.ANE = x.ANE[:__ANELen]
	__UsernameString := viewer.Username(reader)
	if x.Username != __UsernameString {
		__UsernameStringCopy := make([]byte, len(__UsernameString))
		copy(__UsernameStringCopy, __UsernameString)
		x.Username = *(*string)(unsafe.Pointer(&__UsernameStringCopy))
	}
	__ANDString := viewer.AND(reader)
	if x.AND != __ANDString {
		__ANDStringCopy := make([]byte, len(__ANDString))
		copy(__ANDStringCopy, __ANDString)
		x.AND = *(*string)(unsafe.Pointer(&__ANDStringCopy))
	}
}

type BaseSViewer [128]byte

type BaseSSource struct {
	*BaseSViewer
}

func NewBaseSViewer(reader *karmem.Reader, offset uint32) (v *BaseSViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*BaseSViewer)(unsafe.Pointer(&_Null_case_string_setter[0]))
	}
	v = (*BaseSViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BaseSViewer)(unsafe.Pointer(&_Null_case_string_setter[0]))
	}
	return v
}

func NewBaseSSource(reader *karmem.Reader, offset uint32) (v BaseSSource) {
	return BaseSSource{NewBaseSViewer(reader, offset)}
}

func newBaseSSourceByViewer(viewer *BaseSViewer) (v BaseSSource) {
	return BaseSSource{viewer}
}

func (x *BaseSViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BaseSViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BaseSViewer) Data() (v []int64) {
	if 4+80 > x.size() {
		return []int64{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4)), 10, 10,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *BaseSSource) SetData(v [10]int64) {
	copy((*(*[10]int64)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 4)))[:], v[:])
}

func (x *BaseSViewer) ANC(reader *karmem.Reader) (v string) {
	if 84+8 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 84))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 84+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BaseSSource) SetANC(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 84))
	__ANCSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 84+4))
	if uint32(__ANCSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSize, __ANCSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseSViewer), offset)), __ANCSize, __ANCSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 84+4)) = uint32(__ANCSize)
	return nil
}

func (x *BaseSViewer) ANE(reader *karmem.Reader) (v []uint8) {
	if 96+8 > x.size() {
		return []uint8{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 96))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 96+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint8)(unsafe.Pointer(&slice))
}

func (x *BaseSSource) SetANE(v []uint8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 96))
	__ANESize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 96+4))
	if uint32(__ANESize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANESlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANESize, __ANESize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseSViewer), offset)), __ANESize, __ANESize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANESlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 96+4)) = uint32(__ANESize)
	return nil
}

func (x *BaseSViewer) Username(reader *karmem.Reader) (v string) {
	if 104+8 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 104))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 104+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	if length > 120 {
		length = 120
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BaseSSource) SetUsername(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 104))
	__UsernameSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 104+4))
	if uint32(__UsernameSize) > 120 {
		return errors.New("invalid size, new size must less than max length")
	}
	if uint32(__UsernameSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__UsernameSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __UsernameSize, __UsernameSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseSViewer), offset)), __UsernameSize, __UsernameSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__UsernameSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 104+4)) = uint32(__UsernameSize)
	return nil
}

func (x *BaseSViewer) AND(reader *karmem.Reader) (v string) {
	if 116+8 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 116))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 116+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BaseSSource) SetAND(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 116))
	__ANDSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 116+4))
	if uint32(__ANDSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANDSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANDSize, __ANDSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseSViewer), offset)), __ANDSize, __ANDSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANDSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseSViewer), 116+4)) = uint32(__ANDSize)
	return nil
}
