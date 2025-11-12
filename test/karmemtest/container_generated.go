package karmemtest

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_container = [256]byte{}
var _NullReader_container = karmem.NewReader(_Null_container[:])

func init() {

}

const (
	PacketIdentifierBase   = 9956155639035161441
	PacketIdentifierBaseA  = 9956155639035161441
	PacketIdentifierBaseAP = 9956155639035161441
	PacketIdentifierB1     = 16039755183289441794
	PacketIdentifierB2     = 1668427110309968365
	PacketIdentifierR1     = 9956155639035161441
	PacketIdentifierR2     = 9956155639035161441
	PacketIdentifierBase2  = 9956155639035161441
	PacketIdentifierBase3  = 9956155639035161441
	PacketIdentifierRef2   = 12539141064469261489
	PacketIdentifierRef3   = 12539141064469261489
)

type Base struct {
	Data [10]int64
	// fixed string of size 10
	ANC string
}

func NewBase() Base {
	return Base{}
}

func (x *Base) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBase
}

func (x *Base) Reset() {
	x.Read((*BaseViewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *Base) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Base) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(104)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(98))
	__DataOffset := offset + 4
	writer.WriteAt(__DataOffset, (*[80]byte)(unsafe.Pointer(&x.Data))[:])
	__ANCOffset := offset + 84
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 10 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[10]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_container[:10])
	}

	return offset, nil
}

func (*Base) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(104)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(98))
	writer.WriteAt(offset, _Null_container[:size-4])
	__DataOffset := offset + 4
	_ = __DataOffset
	__ANCOffset := offset + 84
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)

	return offset, nil
}

func (x *Base) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBaseViewer(reader, 0), reader)
}

func (x *Base) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBaseViewer(reader, offset), reader)
}

func (x *Base) Read(viewer *BaseViewer, reader *karmem.Reader) {
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
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type BaseA struct {
	Data []int64
	ANC  string
}

func NewBaseA() BaseA {
	return BaseA{}
}

func (x *BaseA) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBaseA
}

func (x *BaseA) Reset() {
	x.Read((*BaseAViewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *BaseA) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BaseA) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(32)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(32))
	__DataSize := uint(8 * len(x.Data))
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	writer.Write4At(offset+4+4, uint32(__DataSize))
	writer.Write4At(offset+4+4+4, 8)
	__DataSlice := *(*[3]uint)(unsafe.Pointer(&x.Data))
	__DataSlice[1] = __DataSize
	__DataSlice[2] = __DataSize
	writer.WriteAt(__DataOffset, *(*[]byte)(unsafe.Pointer(&__DataSlice)))
	__ANCSize := uint(1 * len(x.ANC))
	__ANCOffset, err := writer.Alloc(__ANCSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+16, uint32(__ANCOffset))
	writer.Write4At(offset+16+4, uint32(__ANCSize))
	writer.Write4At(offset+16+4+4, 1)
	__ANCSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANC)), __ANCSize, __ANCSize}
	writer.WriteAt(__ANCOffset, *(*[]byte)(unsafe.Pointer(&__ANCSlice)))

	return offset, nil
}

func (*BaseA) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(32)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(32))
	writer.WriteAt(offset, _Null_container[:size-4])
	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	writer.Write4At(offset+4+4, 0)
	writer.Write4At(offset+4+4+4, 8)
	__ANCOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+16, uint32(__ANCOffset))
	writer.Write4At(offset+16+4, 0)
	writer.Write4At(offset+16+4+4, 1)

	return offset, nil
}

func (x *BaseA) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBaseAViewer(reader, 0), reader)
}

func (x *BaseA) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBaseAViewer(reader, offset), reader)
}

func (x *BaseA) Read(viewer *BaseAViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__DataSlice := viewer.Data(reader)
	__DataLen := len(__DataSlice)
	if __DataLen > cap(x.Data) {
		x.Data = append(x.Data, make([]int64, __DataLen-len(x.Data))...)
	}
	if __DataLen > len(x.Data) {
		x.Data = x.Data[:__DataLen]
	}
	copy(x.Data, __DataSlice)
	x.Data = x.Data[:__DataLen]
	__ANCString := viewer.ANC(reader)
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type BaseAP struct {
	Data []int64
	ANC  string
}

func NewBaseAP() BaseAP {
	return BaseAP{}
}

func (x *BaseAP) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBaseAP
}

func (x *BaseAP) Reset() {
	x.Read((*BaseAPViewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *BaseAP) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BaseAP) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(24)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(24))
	__DataSize := uint(8 * len(x.Data))
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	writer.Write4At(offset+4+4, uint32(__DataSize))
	__DataSlice := *(*[3]uint)(unsafe.Pointer(&x.Data))
	__DataSlice[1] = __DataSize
	__DataSlice[2] = __DataSize
	writer.WriteAt(__DataOffset, *(*[]byte)(unsafe.Pointer(&__DataSlice)))
	__ANCSize := uint(1 * len(x.ANC))
	__ANCOffset, err := writer.Alloc(__ANCSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+12, uint32(__ANCOffset))
	writer.Write4At(offset+12+4, uint32(__ANCSize))
	__ANCSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANC)), __ANCSize, __ANCSize}
	writer.WriteAt(__ANCOffset, *(*[]byte)(unsafe.Pointer(&__ANCSlice)))

	return offset, nil
}

func (*BaseAP) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(24)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(24))
	writer.WriteAt(offset, _Null_container[:size-4])
	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	writer.Write4At(offset+4+4, 0)
	__ANCOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+12, uint32(__ANCOffset))
	writer.Write4At(offset+12+4, 0)

	return offset, nil
}

func (x *BaseAP) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBaseAPViewer(reader, 0), reader)
}

func (x *BaseAP) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBaseAPViewer(reader, offset), reader)
}

func (x *BaseAP) Read(viewer *BaseAPViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__DataSlice := viewer.Data(reader)
	__DataLen := len(__DataSlice)
	if __DataLen > cap(x.Data) {
		x.Data = append(x.Data, make([]int64, __DataLen-len(x.Data))...)
	}
	if __DataLen > len(x.Data) {
		x.Data = x.Data[:__DataLen]
	}
	copy(x.Data, __DataSlice)
	x.Data = x.Data[:__DataLen]
	__ANCString := viewer.ANC(reader)
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type B1 struct {
	Data *Base
}

func NewB1() B1 {
	return B1{}
}

func (x *B1) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierB1
}

func (x *B1) Reset() {
	x.Read((*B1Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *B1) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *B1) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(16)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(9))
	__DataSize := uint(104)
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	if x.Data == nil {
		writer.Write1At(offset+4+4, 0)
	} else {
		writer.Write1At(offset+4+4, 1)
	}
	if x.Data == nil {
		if _, err := x.Data.WriteDefault(writer, __DataOffset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Data.Write(writer, __DataOffset); err != nil {
			return offset, err
		}
	}

	return offset, nil
}

func (*B1) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(16)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(9))
	writer.WriteAt(offset, _Null_container[:size-4])

	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	writer.Write1At(__DataOffset, 0)
	if _, err := (*Base)(nil).WriteDefault(writer, __DataOffset+1); err != nil {
		return offset, err
	}

	return offset, nil
}

func (x *B1) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewB1Viewer(reader, 0), reader)
}

func (x *B1) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewB1Viewer(reader, offset), reader)
}

func (x *B1) Read(viewer *B1Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	if !viewer.isDataNil() {
		DataViewer := viewer.Data(reader)
		if x.Data == nil {
			x.Data = new(Base)
		}
		x.Data.Read(DataViewer, reader)
	}
}

type B2 struct {
	Data *Base
}

func NewB2() B2 {
	return B2{}
}

func (x *B2) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierB2
}

func (x *B2) Reset() {
	x.Read((*B2Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *B2) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *B2) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(8)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__DataSize := uint(104)
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__DataOffset))
	if x.Data == nil {
		writer.Write1At(offset+0+4, 0)
	} else {
		writer.Write1At(offset+0+4, 1)
	}
	if x.Data == nil {
		if _, err := x.Data.WriteDefault(writer, __DataOffset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Data.Write(writer, __DataOffset); err != nil {
			return offset, err
		}
	}

	return offset, nil
}

func (*B2) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(8)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_container[:size-4])

	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__DataOffset))
	writer.Write1At(__DataOffset, 0)
	if _, err := (*Base)(nil).WriteDefault(writer, __DataOffset+1); err != nil {
		return offset, err
	}

	return offset, nil
}

func (x *B2) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewB2Viewer(reader, 0), reader)
}

func (x *B2) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewB2Viewer(reader, offset), reader)
}

func (x *B2) Read(viewer *B2Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	if !viewer.isDataNil() {
		DataViewer := viewer.Data(reader)
		if x.Data == nil {
			x.Data = new(Base)
		}
		x.Data.Read(DataViewer, reader)
	}
}

type R1 struct {
	Data *Base
	ANC  string
}

func NewR1() R1 {
	return R1{}
}

func (x *R1) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierR1
}

func (x *R1) Reset() {
	x.Read((*R1Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *R1) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *R1) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(24)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__DataSize := uint(104)
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__DataOffset))
	if x.Data == nil {
		writer.Write1At(offset+0+4, 0)
	} else {
		writer.Write1At(offset+0+4, 1)
	}
	if x.Data == nil {
		if _, err := x.Data.WriteDefault(writer, __DataOffset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Data.Write(writer, __DataOffset); err != nil {
			return offset, err
		}
	}
	__ANCSize := uint(1 * len(x.ANC))
	__ANCOffset, err := writer.Alloc(__ANCSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+5, uint32(__ANCOffset))
	writer.Write4At(offset+5+4, uint32(__ANCSize))
	writer.Write4At(offset+5+4+4, 1)
	__ANCSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANC)), __ANCSize, __ANCSize}
	writer.WriteAt(__ANCOffset, *(*[]byte)(unsafe.Pointer(&__ANCSlice)))

	return offset, nil
}

func (*R1) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(24)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_container[:size-4])

	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__DataOffset))
	writer.Write1At(__DataOffset, 0)
	if _, err := (*Base)(nil).WriteDefault(writer, __DataOffset+1); err != nil {
		return offset, err
	}
	__ANCOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+5, uint32(__ANCOffset))
	writer.Write4At(offset+5+4, 0)
	writer.Write4At(offset+5+4+4, 1)

	return offset, nil
}

func (x *R1) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewR1Viewer(reader, 0), reader)
}

func (x *R1) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewR1Viewer(reader, offset), reader)
}

func (x *R1) Read(viewer *R1Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	if !viewer.isDataNil() {
		DataViewer := viewer.Data(reader)
		if x.Data == nil {
			x.Data = new(Base)
		}
		x.Data.Read(DataViewer, reader)
	}
	__ANCString := viewer.ANC(reader)
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type R2 struct {
	Data Base
	ANC  string
}

func NewR2() R2 {
	return R2{}
}

func (x *R2) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierR2
}

func (x *R2) Reset() {
	x.Read((*R2Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *R2) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *R2) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(24)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__DataSize := uint(104)
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__DataOffset))
	if _, err := x.Data.Write(writer, __DataOffset); err != nil {
		return offset, err
	}
	__ANCSize := uint(1 * len(x.ANC))
	__ANCOffset, err := writer.Alloc(__ANCSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+5, uint32(__ANCOffset))
	writer.Write4At(offset+5+4, uint32(__ANCSize))
	writer.Write4At(offset+5+4+4, 1)
	__ANCSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ANC)), __ANCSize, __ANCSize}
	writer.WriteAt(__ANCOffset, *(*[]byte)(unsafe.Pointer(&__ANCSlice)))

	return offset, nil
}

func (*R2) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(24)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_container[:size-4])

	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__DataOffset))
	if _, err := (*Base)(nil).WriteDefault(writer, __DataOffset); err != nil {
		return offset, err
	}
	__ANCOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+5, uint32(__ANCOffset))
	writer.Write4At(offset+5+4, 0)
	writer.Write4At(offset+5+4+4, 1)

	return offset, nil
}

func (x *R2) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewR2Viewer(reader, 0), reader)
}

func (x *R2) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewR2Viewer(reader, offset), reader)
}

func (x *R2) Read(viewer *R2Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.Data.Read(viewer.Data(reader), reader)
	__ANCString := viewer.ANC(reader)
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type Base2 struct {
	Data [10]int64
	// fixed string of size 64
	ANC string
}

func NewBase2() Base2 {
	return Base2{}
}

func (x *Base2) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBase2
}

func (x *Base2) Reset() {
	x.Read((*Base2Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *Base2) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Base2) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(152)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__DataOffset := offset + 0
	writer.WriteAt(__DataOffset, (*[80]byte)(unsafe.Pointer(&x.Data))[:])
	__ANCOffset := offset + 80
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 64 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[64]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_container[:64])
	}

	return offset, nil
}

func (*Base2) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(152)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_container[:size-4])
	__DataOffset := offset + 0
	_ = __DataOffset
	__ANCOffset := offset + 80
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)

	return offset, nil
}

func (x *Base2) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBase2Viewer(reader, 0), reader)
}

func (x *Base2) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBase2Viewer(reader, offset), reader)
}

func (x *Base2) Read(viewer *Base2Viewer, reader *karmem.Reader) {
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
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type Base3 struct {
	Value int64
	Data  *Base
	// fixed string of size 64
	ANC string
}

func NewBase3() Base3 {
	return Base3{}
}

func (x *Base3) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBase3
}

func (x *Base3) Reset() {
	x.Read((*Base3Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *Base3) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Base3) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(88)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__ValueOffset := offset + 0
	writer.Write8At(__ValueOffset, *(*uint64)(unsafe.Pointer(&x.Value)))
	__DataSize := uint(104)
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+8, uint32(__DataOffset))
	if x.Data == nil {
		writer.Write1At(offset+8+4, 0)
	} else {
		writer.Write1At(offset+8+4, 1)
	}
	if x.Data == nil {
		if _, err := x.Data.WriteDefault(writer, __DataOffset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Data.Write(writer, __DataOffset); err != nil {
			return offset, err
		}
	}
	__ANCOffset := offset + 13
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 64 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[64]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_container[:64])
	}

	return offset, nil
}

func (*Base3) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(88)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_container[:size-4])
	__ValueOffset := offset + 0
	_ = __ValueOffset

	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+8, uint32(__DataOffset))
	writer.Write1At(__DataOffset, 0)
	if _, err := (*Base)(nil).WriteDefault(writer, __DataOffset+1); err != nil {
		return offset, err
	}
	__ANCOffset := offset + 13
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)

	return offset, nil
}

func (x *Base3) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBase3Viewer(reader, 0), reader)
}

func (x *Base3) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBase3Viewer(reader, offset), reader)
}

func (x *Base3) Read(viewer *Base3Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.Value = viewer.Value()
	if !viewer.isDataNil() {
		DataViewer := viewer.Data(reader)
		if x.Data == nil {
			x.Data = new(Base)
		}
		x.Data.Read(DataViewer, reader)
	}
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type Ref2 struct {
	ID   int64
	P    int64
	Arr  *Base2
	Arr2 []Base2
	// fixed string of size 64
	ANC string
}

func NewRef2() Ref2 {
	return Ref2{}
}

func (x *Ref2) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierRef2
}

func (x *Ref2) Reset() {
	x.Read((*Ref2Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *Ref2) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Ref2) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(256)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(253))
	__IDOffset := offset + 4
	writer.Write8At(__IDOffset, *(*uint64)(unsafe.Pointer(&x.ID)))
	__POffset := offset + 12
	writer.Write8At(__POffset, *(*uint64)(unsafe.Pointer(&x.P)))
	__ArrOffset := offset + 20
	if x.Arr == nil {
		writer.Write1At(offset+20, 0)
	} else {
		writer.Write1At(offset+20, 1)
	}
	__ArrOffset += 1
	if x.Arr == nil {
		if _, err := x.Arr.WriteDefault(writer, __ArrOffset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Arr.Write(writer, __ArrOffset); err != nil {
			return offset, err
		}
	}
	__Arr2Size := uint(152 * len(x.Arr2))
	__Arr2Offset, err := writer.Alloc(__Arr2Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+173, uint32(__Arr2Offset))
	writer.Write4At(offset+173+4, uint32(__Arr2Size))
	writer.Write4At(offset+173+4+4, 152)
	for i := range x.Arr2 {
		if _, err := x.Arr2[i].Write(writer, __Arr2Offset); err != nil {
			return offset, err
		}
		__Arr2Offset += 152
	}
	__ANCOffset := offset + 185
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 64 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[64]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_container[:64])
	}

	return offset, nil
}

func (*Ref2) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(256)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(253))
	writer.WriteAt(offset, _Null_container[:size-4])
	__IDOffset := offset + 4
	_ = __IDOffset
	__POffset := offset + 12
	_ = __POffset
	__ArrOffset := offset + 20
	_ = __ArrOffset
	writer.Write1At(__ArrOffset, 0)
	if _, err := (*Base2)(nil).WriteDefault(writer, __ArrOffset+1); err != nil {
		return offset, err
	}
	__Arr2Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+173, uint32(__Arr2Offset))
	writer.Write4At(offset+173+4, 0)
	writer.Write4At(offset+173+4+4, 152)
	__ANCOffset := offset + 185
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)

	return offset, nil
}

func (x *Ref2) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewRef2Viewer(reader, 0), reader)
}

func (x *Ref2) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewRef2Viewer(reader, offset), reader)
}

func (x *Ref2) Read(viewer *Ref2Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.ID = viewer.ID()
	x.P = viewer.P()
	if !viewer.isArrNil() {
		ArrViewer := viewer.Arr()
		if x.Arr == nil {
			x.Arr = new(Base2)
		}
		x.Arr.Read(ArrViewer, reader)
	}
	__Arr2Slice := viewer.Arr2(reader)
	__Arr2Len := len(__Arr2Slice)
	if __Arr2Len > cap(x.Arr2) {
		x.Arr2 = append(x.Arr2, make([]Base2, __Arr2Len-len(x.Arr2))...)
	}
	if __Arr2Len > len(x.Arr2) {
		x.Arr2 = x.Arr2[:__Arr2Len]
	}
	for i := 0; i < __Arr2Len; i++ {
		x.Arr2[i].Read(&__Arr2Slice[i], reader)
	}
	x.Arr2 = x.Arr2[:__Arr2Len]
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type Ref3 struct {
	ID   int64
	Arr  *Base3
	Arr2 []Base3
	// fixed string of size 64
	ANC string
}

func NewRef3() Ref3 {
	return Ref3{}
}

func (x *Ref3) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierRef3
}

func (x *Ref3) Reset() {
	x.Read((*Ref3Viewer)(unsafe.Pointer(&_Null_container[0])), _NullReader_container)
}

func (x *Ref3) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Ref3) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(184)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(181))
	__IDOffset := offset + 4
	writer.Write8At(__IDOffset, *(*uint64)(unsafe.Pointer(&x.ID)))
	__ArrOffset := offset + 12
	if x.Arr == nil {
		writer.Write1At(offset+12, 0)
	} else {
		writer.Write1At(offset+12, 1)
	}
	__ArrOffset += 1
	if x.Arr == nil {
		if _, err := x.Arr.WriteDefault(writer, __ArrOffset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Arr.Write(writer, __ArrOffset); err != nil {
			return offset, err
		}
	}
	__Arr2Size := uint(88 * len(x.Arr2))
	__Arr2Offset, err := writer.Alloc(__Arr2Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+101, uint32(__Arr2Offset))
	writer.Write4At(offset+101+4, uint32(__Arr2Size))
	writer.Write4At(offset+101+4+4, 88)
	for i := range x.Arr2 {
		if _, err := x.Arr2[i].Write(writer, __Arr2Offset); err != nil {
			return offset, err
		}
		__Arr2Offset += 88
	}
	__ANCOffset := offset + 113
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 64 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[64]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_container[:64])
	}

	return offset, nil
}

func (*Ref3) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(184)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(181))
	writer.WriteAt(offset, _Null_container[:size-4])
	__IDOffset := offset + 4
	_ = __IDOffset
	__ArrOffset := offset + 12
	_ = __ArrOffset
	writer.Write1At(__ArrOffset, 0)
	if _, err := (*Base3)(nil).WriteDefault(writer, __ArrOffset+1); err != nil {
		return offset, err
	}
	__Arr2Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+101, uint32(__Arr2Offset))
	writer.Write4At(offset+101+4, 0)
	writer.Write4At(offset+101+4+4, 88)
	__ANCOffset := offset + 113
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)

	return offset, nil
}

func (x *Ref3) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewRef3Viewer(reader, 0), reader)
}

func (x *Ref3) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewRef3Viewer(reader, offset), reader)
}

func (x *Ref3) Read(viewer *Ref3Viewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.ID = viewer.ID()
	if !viewer.isArrNil() {
		ArrViewer := viewer.Arr()
		if x.Arr == nil {
			x.Arr = new(Base3)
		}
		x.Arr.Read(ArrViewer, reader)
	}
	__Arr2Slice := viewer.Arr2(reader)
	__Arr2Len := len(__Arr2Slice)
	if __Arr2Len > cap(x.Arr2) {
		x.Arr2 = append(x.Arr2, make([]Base3, __Arr2Len-len(x.Arr2))...)
	}
	if __Arr2Len > len(x.Arr2) {
		x.Arr2 = x.Arr2[:__Arr2Len]
	}
	for i := 0; i < __Arr2Len; i++ {
		x.Arr2[i].Read(&__Arr2Slice[i], reader)
	}
	x.Arr2 = x.Arr2[:__Arr2Len]
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
}

type BaseViewer [104]byte

type BaseSource struct {
	*BaseViewer
}

func NewBaseViewer(reader *karmem.Reader, offset uint32) (v *BaseViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*BaseViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*BaseViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BaseViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return v
}

func NewBaseSource(reader *karmem.Reader, offset uint32) (v BaseSource) {
	return BaseSource{NewBaseViewer(reader, offset)}
}

func newBaseSourceByViewer(viewer *BaseViewer) (v BaseSource) {
	return BaseSource{viewer}
}

func (x *BaseViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BaseViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BaseViewer) Data() (v []int64) {
	if 4+80 > x.size() {
		return []int64{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4)), 10, 10,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *BaseSource) SetData(v [10]int64) {
	copy((*(*[10]int64)(unsafe.Add(unsafe.Pointer(x.BaseViewer), 4)))[:], v[:])
}

func (x *BaseViewer) ANC() (v string) {
	if 84+10 > x.size() {
		return v
	}
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 84))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+84)), uintptr(size), 10,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BaseSource) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.BaseViewer)) = uint32(len(v))
	copy((*[10]byte)(unsafe.Add(unsafe.Pointer(x.BaseViewer), 4+84))[:], v)
}

type BaseAViewer [32]byte

type BaseASource struct {
	*BaseAViewer
}

func NewBaseAViewer(reader *karmem.Reader, offset uint32) (v *BaseAViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*BaseAViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*BaseAViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BaseAViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return v
}

func NewBaseASource(reader *karmem.Reader, offset uint32) (v BaseASource) {
	return BaseASource{NewBaseAViewer(reader, offset)}
}

func newBaseASourceByViewer(viewer *BaseAViewer) (v BaseASource) {
	return BaseASource{viewer}
}

func (x *BaseAViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BaseAViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BaseAViewer) Data(reader *karmem.Reader) (v []int64) {
	if 4+12 > x.size() {
		return []int64{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4+4))
	if !reader.IsValidOffset(offset, size) {
		return []int64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *BaseASource) SetData(v []int64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAViewer), 4))
	__DataSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAViewer), 4+4))
	if uint32(__DataSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__DataSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __DataSize, __DataSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseAViewer), offset)), __DataSize, __DataSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__DataSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAViewer), 4+4)) = uint32(__DataSize)
	return nil
}

func (x *BaseAViewer) ANC(reader *karmem.Reader) (v string) {
	if 16+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 16))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 16+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BaseASource) SetANC(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAViewer), 16))
	__ANCSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAViewer), 16+4))
	if uint32(__ANCSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSize, __ANCSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseAViewer), offset)), __ANCSize, __ANCSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAViewer), 16+4)) = uint32(__ANCSize)
	return nil
}

type BaseAPViewer [24]byte

type BaseAPSource struct {
	*BaseAPViewer
}

func NewBaseAPViewer(reader *karmem.Reader, offset uint32) (v *BaseAPViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*BaseAPViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*BaseAPViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BaseAPViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return v
}

func NewBaseAPSource(reader *karmem.Reader, offset uint32) (v BaseAPSource) {
	return BaseAPSource{NewBaseAPViewer(reader, offset)}
}

func newBaseAPSourceByViewer(viewer *BaseAPViewer) (v BaseAPSource) {
	return BaseAPSource{viewer}
}

func (x *BaseAPViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BaseAPViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BaseAPViewer) Data(reader *karmem.Reader) (v []int64) {
	if 4+8 > x.size() {
		return []int64{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4+4))
	if !reader.IsValidOffset(offset, size) {
		return []int64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *BaseAPSource) SetData(v []int64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), 4))
	__DataSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), 4+4))
	if uint32(__DataSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__DataSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __DataSize, __DataSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), offset)), __DataSize, __DataSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__DataSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), 4+4)) = uint32(__DataSize)
	return nil
}

func (x *BaseAPViewer) ANC(reader *karmem.Reader) (v string) {
	if 12+8 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 12))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 12+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *BaseAPSource) SetANC(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), 12))
	__ANCSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), 12+4))
	if uint32(__ANCSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSize, __ANCSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), offset)), __ANCSize, __ANCSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.BaseAPViewer), 12+4)) = uint32(__ANCSize)
	return nil
}

type B1Viewer [16]byte

type B1Source struct {
	*B1Viewer
}

func NewB1Viewer(reader *karmem.Reader, offset uint32) (v *B1Viewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*B1Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*B1Viewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*B1Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return v
}

func NewB1Source(reader *karmem.Reader, offset uint32) (v B1Source) {
	return B1Source{NewB1Viewer(reader, offset)}
}

func newB1SourceByViewer(viewer *B1Viewer) (v B1Source) {
	return B1Source{viewer}
}

func (x *B1Viewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *B1Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *B1Viewer) isDataNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 4+4))
}

func (x *B1Viewer) Data(reader *karmem.Reader) (v *BaseViewer) {
	if 4+4 > x.size() {
		return (*BaseViewer)(unsafe.Pointer(&_Null_container[0]))
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4))
	return NewBaseViewer(reader, offset)
}

type B2Viewer [8]byte

type B2Source struct {
	*B2Viewer
}

func NewB2Viewer(reader *karmem.Reader, offset uint32) (v *B2Viewer) {
	if !reader.IsValidOffset(offset, 5) {
		return (*B2Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*B2Viewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewB2Source(reader *karmem.Reader, offset uint32) (v B2Source) {
	return B2Source{NewB2Viewer(reader, offset)}
}

func newB2SourceByViewer(viewer *B2Viewer) (v B2Source) {
	return B2Source{viewer}
}

func (x *B2Viewer) size() uint32 {
	return 8
}

func (x *B2Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *B2Viewer) isDataNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 0+4))
}

func (x *B2Viewer) Data(reader *karmem.Reader) (v *BaseViewer) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 0))
	return NewBaseViewer(reader, offset)
}

type R1Viewer [24]byte

type R1Source struct {
	*R1Viewer
}

func NewR1Viewer(reader *karmem.Reader, offset uint32) (v *R1Viewer) {
	if !reader.IsValidOffset(offset, 21) {
		return (*R1Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*R1Viewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewR1Source(reader *karmem.Reader, offset uint32) (v R1Source) {
	return R1Source{NewR1Viewer(reader, offset)}
}

func newR1SourceByViewer(viewer *R1Viewer) (v R1Source) {
	return R1Source{viewer}
}

func (x *R1Viewer) size() uint32 {
	return 24
}

func (x *R1Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *R1Viewer) isDataNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 0+4))
}

func (x *R1Viewer) Data(reader *karmem.Reader) (v *BaseViewer) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 0))
	return NewBaseViewer(reader, offset)
}

func (x *R1Viewer) ANC(reader *karmem.Reader) (v string) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 5))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 5+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *R1Source) SetANC(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.R1Viewer), 5))
	__ANCSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.R1Viewer), 5+4))
	if uint32(__ANCSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSize, __ANCSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.R1Viewer), offset)), __ANCSize, __ANCSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.R1Viewer), 5+4)) = uint32(__ANCSize)
	return nil
}

type R2Viewer [24]byte

type R2Source struct {
	*R2Viewer
}

func NewR2Viewer(reader *karmem.Reader, offset uint32) (v *R2Viewer) {
	if !reader.IsValidOffset(offset, 21) {
		return (*R2Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*R2Viewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewR2Source(reader *karmem.Reader, offset uint32) (v R2Source) {
	return R2Source{NewR2Viewer(reader, offset)}
}

func newR2SourceByViewer(viewer *R2Viewer) (v R2Source) {
	return R2Source{viewer}
}

func (x *R2Viewer) size() uint32 {
	return 24
}

func (x *R2Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *R2Viewer) Data(reader *karmem.Reader) (v *BaseViewer) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 0))
	return NewBaseViewer(reader, offset)
}

func (x *R2Viewer) ANC(reader *karmem.Reader) (v string) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 5))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 5+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *R2Source) SetANC(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.R2Viewer), 5))
	__ANCSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.R2Viewer), 5+4))
	if uint32(__ANCSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__ANCSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __ANCSize, __ANCSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.R2Viewer), offset)), __ANCSize, __ANCSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__ANCSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.R2Viewer), 5+4)) = uint32(__ANCSize)
	return nil
}

type Base2Viewer [152]byte

type Base2Source struct {
	*Base2Viewer
}

func NewBase2Viewer(reader *karmem.Reader, offset uint32) (v *Base2Viewer) {
	if !reader.IsValidOffset(offset, 148) {
		return (*Base2Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*Base2Viewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewBase2Source(reader *karmem.Reader, offset uint32) (v Base2Source) {
	return Base2Source{NewBase2Viewer(reader, offset)}
}

func newBase2SourceByViewer(viewer *Base2Viewer) (v Base2Source) {
	return Base2Source{viewer}
}

func (x *Base2Viewer) size() uint32 {
	return 152
}

func (x *Base2Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *Base2Viewer) Data() (v []int64) {
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 0)), 10, 10,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *Base2Source) SetData(v [10]int64) {
	copy((*(*[10]int64)(unsafe.Add(unsafe.Pointer(x.Base2Viewer), 0)))[:], v[:])
}

func (x *Base2Viewer) ANC() (v string) {
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 80))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+80)), uintptr(size), 64,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *Base2Source) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.Base2Viewer)) = uint32(len(v))
	copy((*[64]byte)(unsafe.Add(unsafe.Pointer(x.Base2Viewer), 4+80))[:], v)
}

type Base3Viewer [88]byte

type Base3Source struct {
	*Base3Viewer
}

func NewBase3Viewer(reader *karmem.Reader, offset uint32) (v *Base3Viewer) {
	if !reader.IsValidOffset(offset, 81) {
		return (*Base3Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*Base3Viewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewBase3Source(reader *karmem.Reader, offset uint32) (v Base3Source) {
	return Base3Source{NewBase3Viewer(reader, offset)}
}

func newBase3SourceByViewer(viewer *Base3Viewer) (v Base3Source) {
	return Base3Source{viewer}
}

func (x *Base3Viewer) size() uint32 {
	return 88
}

func (x *Base3Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *Base3Viewer) Value() (v int64) {
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *Base3Source) SetValue(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.Base3Viewer), 0)) = v
}

func (x *Base3Viewer) isDataNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 8+4))
}

func (x *Base3Viewer) Data(reader *karmem.Reader) (v *BaseViewer) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 8))
	return NewBaseViewer(reader, offset)
}

func (x *Base3Viewer) ANC() (v string) {
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 13))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+13)), uintptr(size), 64,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *Base3Source) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.Base3Viewer)) = uint32(len(v))
	copy((*[64]byte)(unsafe.Add(unsafe.Pointer(x.Base3Viewer), 4+13))[:], v)
}

type Ref2Viewer [256]byte

type Ref2Source struct {
	*Ref2Viewer
}

func NewRef2Viewer(reader *karmem.Reader, offset uint32) (v *Ref2Viewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*Ref2Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*Ref2Viewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*Ref2Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return v
}

func NewRef2Source(reader *karmem.Reader, offset uint32) (v Ref2Source) {
	return Ref2Source{NewRef2Viewer(reader, offset)}
}

func newRef2SourceByViewer(viewer *Ref2Viewer) (v Ref2Source) {
	return Ref2Source{viewer}
}

func (x *Ref2Viewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *Ref2Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *Ref2Viewer) ID() (v int64) {
	if 4+8 > x.size() {
		return v
	}
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *Ref2Source) SetID(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), 4)) = v
}

func (x *Ref2Viewer) P() (v int64) {
	if 12+8 > x.size() {
		return v
	}
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 12))
}

func (x *Ref2Source) SetP(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), 12)) = v
}

func (x *Ref2Viewer) isArrNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 20))
}

func (x *Ref2Viewer) Arr() (v *Base2Viewer) {
	if 20+152 > x.size() {
		return (*Base2Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return (*Base2Viewer)(unsafe.Add(unsafe.Pointer(x), 20+1))
}

func (x *Ref2Viewer) Arr2(reader *karmem.Reader) (v []Base2Viewer) {
	if 173+12 > x.size() {
		return []Base2Viewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 173))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 173+4))
	if !reader.IsValidOffset(offset, size) {
		return []Base2Viewer{}
	}
	length := uintptr(size / 152)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]Base2Viewer)(unsafe.Pointer(&slice))
}

func (x *Ref2Source) SetArr2(v []Base2) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), 173))
	__Arr2Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), 173+4))
	if uint32(__Arr2Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__Arr2Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __Arr2Size, __Arr2Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), offset)), __Arr2Size, __Arr2Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__Arr2Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), 173+4)) = uint32(__Arr2Size)
	return nil
}

func (x *Ref2Viewer) ANC() (v string) {
	if 185+64 > x.size() {
		return v
	}
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 185))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+185)), uintptr(size), 64,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *Ref2Source) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.Ref2Viewer)) = uint32(len(v))
	copy((*[64]byte)(unsafe.Add(unsafe.Pointer(x.Ref2Viewer), 4+185))[:], v)
}

type Ref3Viewer [184]byte

type Ref3Source struct {
	*Ref3Viewer
}

func NewRef3Viewer(reader *karmem.Reader, offset uint32) (v *Ref3Viewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*Ref3Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	v = (*Ref3Viewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*Ref3Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return v
}

func NewRef3Source(reader *karmem.Reader, offset uint32) (v Ref3Source) {
	return Ref3Source{NewRef3Viewer(reader, offset)}
}

func newRef3SourceByViewer(viewer *Ref3Viewer) (v Ref3Source) {
	return Ref3Source{viewer}
}

func (x *Ref3Viewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *Ref3Viewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *Ref3Viewer) ID() (v int64) {
	if 4+8 > x.size() {
		return v
	}
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *Ref3Source) SetID(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.Ref3Viewer), 4)) = v
}

func (x *Ref3Viewer) isArrNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 12))
}

func (x *Ref3Viewer) Arr() (v *Base3Viewer) {
	if 12+88 > x.size() {
		return (*Base3Viewer)(unsafe.Pointer(&_Null_container[0]))
	}
	return (*Base3Viewer)(unsafe.Add(unsafe.Pointer(x), 12+1))
}

func (x *Ref3Viewer) Arr2(reader *karmem.Reader) (v []Base3Viewer) {
	if 101+12 > x.size() {
		return []Base3Viewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 101))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 101+4))
	if !reader.IsValidOffset(offset, size) {
		return []Base3Viewer{}
	}
	length := uintptr(size / 88)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]Base3Viewer)(unsafe.Pointer(&slice))
}

func (x *Ref3Source) SetArr2(v []Base3) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.Ref3Viewer), 101))
	__Arr2Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.Ref3Viewer), 101+4))
	if uint32(__Arr2Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__Arr2Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __Arr2Size, __Arr2Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.Ref3Viewer), offset)), __Arr2Size, __Arr2Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__Arr2Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.Ref3Viewer), 101+4)) = uint32(__Arr2Size)
	return nil
}

func (x *Ref3Viewer) ANC() (v string) {
	if 113+64 > x.size() {
		return v
	}
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 113))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+113)), uintptr(size), 64,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *Ref3Source) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.Ref3Viewer)) = uint32(len(v))
	copy((*[64]byte)(unsafe.Add(unsafe.Pointer(x.Ref3Viewer), 4+113))[:], v)
}
