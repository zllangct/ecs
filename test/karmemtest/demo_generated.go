package karmemtest

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_demo = [18859]byte{}
var _NullReader_demo = karmem.NewReader(_Null_demo[:])

func init() {

}

type (
	XA uint8
)

const (
	XAAAAAAAAAAA XA = 0
	XAB          XA = 1
)

type (
	XB uint16
)

const (
	XBA                        XB = 0
	XBBBBBBBBBBBBBBBBBBBBBBBBB XB = 100
)

type (
	XC uint32
)

const (
	XCAAAAAAAA              XC = 0
	XCBBBBBBBBBBBBBBBBBBBBB XC = 10
)

type (
	XD uint64
)

const (
	XDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA XD = 0
	XDBBBBBBBBB                             XD = 10000000
)

type (
	IA int8
)

const (
	IAA IA = 0
	IAB IA = 1
)

type (
	IB int16
)

const (
	IBA                IB = 0
	IBBBBBBBBBBBBBBBBB IB = 300
)

type (
	IC int32
)

const (
	ICA IC = 0
	ICB IC = 1
)

type (
	ID int64
)

const (
	IDA ID = 0
	IDB ID = 1
)

const (
	PacketIdentifierSimpleNumbers       = 12663584113194877154
	PacketIdentifierSimpleNumbersPacked = 4896517088533716004
	PacketIdentifierComplexPacked       = 14296238997608355174
)

type SimpleNumbers struct {
	N8    uint8
	N16   uint16
	N32   uint32
	N64   uint64
	M8    int8
	M16   int16
	M32   int32
	M64   int64
	OF32  float32
	OF64  float64
	B1    bool
	NN8   []uint8
	NN16  []uint16
	NN32  []uint32
	NN64  []uint64
	NM8   []int8
	NM16  []int16
	NM32  []int32
	NM64  []int64
	NOF32 []float32
	NOF64 []float64
}

func NewSimpleNumbers() SimpleNumbers {
	return SimpleNumbers{}
}

func (x *SimpleNumbers) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierSimpleNumbers
}

func (x *SimpleNumbers) Reset() {
	x.Read((*SimpleNumbersViewer)(unsafe.Pointer(&_Null_demo[0])), _NullReader_demo)
}

func (x *SimpleNumbers) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *SimpleNumbers) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(123)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__N8Offset := offset + 0
	writer.Write1At(__N8Offset, *(*uint8)(unsafe.Pointer(&x.N8)))
	__N16Offset := offset + 1
	writer.Write2At(__N16Offset, *(*uint16)(unsafe.Pointer(&x.N16)))
	__N32Offset := offset + 3
	writer.Write4At(__N32Offset, *(*uint32)(unsafe.Pointer(&x.N32)))
	__N64Offset := offset + 7
	writer.Write8At(__N64Offset, *(*uint64)(unsafe.Pointer(&x.N64)))
	__M8Offset := offset + 15
	writer.Write1At(__M8Offset, *(*uint8)(unsafe.Pointer(&x.M8)))
	__M16Offset := offset + 16
	writer.Write2At(__M16Offset, *(*uint16)(unsafe.Pointer(&x.M16)))
	__M32Offset := offset + 18
	writer.Write4At(__M32Offset, *(*uint32)(unsafe.Pointer(&x.M32)))
	__M64Offset := offset + 22
	writer.Write8At(__M64Offset, *(*uint64)(unsafe.Pointer(&x.M64)))
	__OF32Offset := offset + 30
	writer.Write4At(__OF32Offset, *(*uint32)(unsafe.Pointer(&x.OF32)))
	__OF64Offset := offset + 34
	writer.Write8At(__OF64Offset, *(*uint64)(unsafe.Pointer(&x.OF64)))
	__B1Offset := offset + 42
	writer.Write1At(__B1Offset, *(*uint8)(unsafe.Pointer(&x.B1)))
	__NN8Size := uint(1 * len(x.NN8))
	__NN8Offset, err := writer.Alloc(__NN8Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+43, uint32(__NN8Offset))
	writer.Write4At(offset+43+4, uint32(__NN8Size))
	__NN8Slice := *(*[3]uint)(unsafe.Pointer(&x.NN8))
	__NN8Slice[1] = __NN8Size
	__NN8Slice[2] = __NN8Size
	writer.WriteAt(__NN8Offset, *(*[]byte)(unsafe.Pointer(&__NN8Slice)))
	__NN16Size := uint(2 * len(x.NN16))
	__NN16Offset, err := writer.Alloc(__NN16Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+51, uint32(__NN16Offset))
	writer.Write4At(offset+51+4, uint32(__NN16Size))
	__NN16Slice := *(*[3]uint)(unsafe.Pointer(&x.NN16))
	__NN16Slice[1] = __NN16Size
	__NN16Slice[2] = __NN16Size
	writer.WriteAt(__NN16Offset, *(*[]byte)(unsafe.Pointer(&__NN16Slice)))
	__NN32Size := uint(4 * len(x.NN32))
	__NN32Offset, err := writer.Alloc(__NN32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+59, uint32(__NN32Offset))
	writer.Write4At(offset+59+4, uint32(__NN32Size))
	__NN32Slice := *(*[3]uint)(unsafe.Pointer(&x.NN32))
	__NN32Slice[1] = __NN32Size
	__NN32Slice[2] = __NN32Size
	writer.WriteAt(__NN32Offset, *(*[]byte)(unsafe.Pointer(&__NN32Slice)))
	__NN64Size := uint(8 * len(x.NN64))
	__NN64Offset, err := writer.Alloc(__NN64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+67, uint32(__NN64Offset))
	writer.Write4At(offset+67+4, uint32(__NN64Size))
	__NN64Slice := *(*[3]uint)(unsafe.Pointer(&x.NN64))
	__NN64Slice[1] = __NN64Size
	__NN64Slice[2] = __NN64Size
	writer.WriteAt(__NN64Offset, *(*[]byte)(unsafe.Pointer(&__NN64Slice)))
	__NM8Size := uint(1 * len(x.NM8))
	__NM8Offset, err := writer.Alloc(__NM8Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+75, uint32(__NM8Offset))
	writer.Write4At(offset+75+4, uint32(__NM8Size))
	__NM8Slice := *(*[3]uint)(unsafe.Pointer(&x.NM8))
	__NM8Slice[1] = __NM8Size
	__NM8Slice[2] = __NM8Size
	writer.WriteAt(__NM8Offset, *(*[]byte)(unsafe.Pointer(&__NM8Slice)))
	__NM16Size := uint(2 * len(x.NM16))
	__NM16Offset, err := writer.Alloc(__NM16Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+83, uint32(__NM16Offset))
	writer.Write4At(offset+83+4, uint32(__NM16Size))
	__NM16Slice := *(*[3]uint)(unsafe.Pointer(&x.NM16))
	__NM16Slice[1] = __NM16Size
	__NM16Slice[2] = __NM16Size
	writer.WriteAt(__NM16Offset, *(*[]byte)(unsafe.Pointer(&__NM16Slice)))
	__NM32Size := uint(4 * len(x.NM32))
	__NM32Offset, err := writer.Alloc(__NM32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+91, uint32(__NM32Offset))
	writer.Write4At(offset+91+4, uint32(__NM32Size))
	__NM32Slice := *(*[3]uint)(unsafe.Pointer(&x.NM32))
	__NM32Slice[1] = __NM32Size
	__NM32Slice[2] = __NM32Size
	writer.WriteAt(__NM32Offset, *(*[]byte)(unsafe.Pointer(&__NM32Slice)))
	__NM64Size := uint(8 * len(x.NM64))
	__NM64Offset, err := writer.Alloc(__NM64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+99, uint32(__NM64Offset))
	writer.Write4At(offset+99+4, uint32(__NM64Size))
	__NM64Slice := *(*[3]uint)(unsafe.Pointer(&x.NM64))
	__NM64Slice[1] = __NM64Size
	__NM64Slice[2] = __NM64Size
	writer.WriteAt(__NM64Offset, *(*[]byte)(unsafe.Pointer(&__NM64Slice)))
	__NOF32Size := uint(4 * len(x.NOF32))
	__NOF32Offset, err := writer.Alloc(__NOF32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+107, uint32(__NOF32Offset))
	writer.Write4At(offset+107+4, uint32(__NOF32Size))
	__NOF32Slice := *(*[3]uint)(unsafe.Pointer(&x.NOF32))
	__NOF32Slice[1] = __NOF32Size
	__NOF32Slice[2] = __NOF32Size
	writer.WriteAt(__NOF32Offset, *(*[]byte)(unsafe.Pointer(&__NOF32Slice)))
	__NOF64Size := uint(8 * len(x.NOF64))
	__NOF64Offset, err := writer.Alloc(__NOF64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+115, uint32(__NOF64Offset))
	writer.Write4At(offset+115+4, uint32(__NOF64Size))
	__NOF64Slice := *(*[3]uint)(unsafe.Pointer(&x.NOF64))
	__NOF64Slice[1] = __NOF64Size
	__NOF64Slice[2] = __NOF64Size
	writer.WriteAt(__NOF64Offset, *(*[]byte)(unsafe.Pointer(&__NOF64Slice)))

	return offset, nil
}

func (*SimpleNumbers) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(123)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_demo[:size-4])
	__N8Offset := offset + 0
	_ = __N8Offset
	__N16Offset := offset + 1
	_ = __N16Offset
	__N32Offset := offset + 3
	_ = __N32Offset
	__N64Offset := offset + 7
	_ = __N64Offset
	__M8Offset := offset + 15
	_ = __M8Offset
	__M16Offset := offset + 16
	_ = __M16Offset
	__M32Offset := offset + 18
	_ = __M32Offset
	__M64Offset := offset + 22
	_ = __M64Offset
	__OF32Offset := offset + 30
	_ = __OF32Offset
	__OF64Offset := offset + 34
	_ = __OF64Offset
	__B1Offset := offset + 42
	_ = __B1Offset
	__NN8Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+43, uint32(__NN8Offset))
	writer.Write4At(offset+43+4, 0)
	__NN16Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+51, uint32(__NN16Offset))
	writer.Write4At(offset+51+4, 0)
	__NN32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+59, uint32(__NN32Offset))
	writer.Write4At(offset+59+4, 0)
	__NN64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+67, uint32(__NN64Offset))
	writer.Write4At(offset+67+4, 0)
	__NM8Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+75, uint32(__NM8Offset))
	writer.Write4At(offset+75+4, 0)
	__NM16Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+83, uint32(__NM16Offset))
	writer.Write4At(offset+83+4, 0)
	__NM32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+91, uint32(__NM32Offset))
	writer.Write4At(offset+91+4, 0)
	__NM64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+99, uint32(__NM64Offset))
	writer.Write4At(offset+99+4, 0)
	__NOF32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+107, uint32(__NOF32Offset))
	writer.Write4At(offset+107+4, 0)
	__NOF64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+115, uint32(__NOF64Offset))
	writer.Write4At(offset+115+4, 0)

	return offset, nil
}

func (x *SimpleNumbers) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewSimpleNumbersViewer(reader, 0), reader)
}

func (x *SimpleNumbers) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewSimpleNumbersViewer(reader, offset), reader)
}

func (x *SimpleNumbers) Read(viewer *SimpleNumbersViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.N8 = viewer.N8()
	x.N16 = viewer.N16()
	x.N32 = viewer.N32()
	x.N64 = viewer.N64()
	x.M8 = viewer.M8()
	x.M16 = viewer.M16()
	x.M32 = viewer.M32()
	x.M64 = viewer.M64()
	x.OF32 = viewer.OF32()
	x.OF64 = viewer.OF64()
	x.B1 = viewer.B1()
	__NN8Slice := viewer.NN8(reader)
	__NN8Len := len(__NN8Slice)
	if __NN8Len > cap(x.NN8) {
		x.NN8 = append(x.NN8, make([]uint8, __NN8Len-len(x.NN8))...)
	}
	if __NN8Len > len(x.NN8) {
		x.NN8 = x.NN8[:__NN8Len]
	}
	copy(x.NN8, __NN8Slice)
	x.NN8 = x.NN8[:__NN8Len]
	__NN16Slice := viewer.NN16(reader)
	__NN16Len := len(__NN16Slice)
	if __NN16Len > cap(x.NN16) {
		x.NN16 = append(x.NN16, make([]uint16, __NN16Len-len(x.NN16))...)
	}
	if __NN16Len > len(x.NN16) {
		x.NN16 = x.NN16[:__NN16Len]
	}
	copy(x.NN16, __NN16Slice)
	x.NN16 = x.NN16[:__NN16Len]
	__NN32Slice := viewer.NN32(reader)
	__NN32Len := len(__NN32Slice)
	if __NN32Len > cap(x.NN32) {
		x.NN32 = append(x.NN32, make([]uint32, __NN32Len-len(x.NN32))...)
	}
	if __NN32Len > len(x.NN32) {
		x.NN32 = x.NN32[:__NN32Len]
	}
	copy(x.NN32, __NN32Slice)
	x.NN32 = x.NN32[:__NN32Len]
	__NN64Slice := viewer.NN64(reader)
	__NN64Len := len(__NN64Slice)
	if __NN64Len > cap(x.NN64) {
		x.NN64 = append(x.NN64, make([]uint64, __NN64Len-len(x.NN64))...)
	}
	if __NN64Len > len(x.NN64) {
		x.NN64 = x.NN64[:__NN64Len]
	}
	copy(x.NN64, __NN64Slice)
	x.NN64 = x.NN64[:__NN64Len]
	__NM8Slice := viewer.NM8(reader)
	__NM8Len := len(__NM8Slice)
	if __NM8Len > cap(x.NM8) {
		x.NM8 = append(x.NM8, make([]int8, __NM8Len-len(x.NM8))...)
	}
	if __NM8Len > len(x.NM8) {
		x.NM8 = x.NM8[:__NM8Len]
	}
	copy(x.NM8, __NM8Slice)
	x.NM8 = x.NM8[:__NM8Len]
	__NM16Slice := viewer.NM16(reader)
	__NM16Len := len(__NM16Slice)
	if __NM16Len > cap(x.NM16) {
		x.NM16 = append(x.NM16, make([]int16, __NM16Len-len(x.NM16))...)
	}
	if __NM16Len > len(x.NM16) {
		x.NM16 = x.NM16[:__NM16Len]
	}
	copy(x.NM16, __NM16Slice)
	x.NM16 = x.NM16[:__NM16Len]
	__NM32Slice := viewer.NM32(reader)
	__NM32Len := len(__NM32Slice)
	if __NM32Len > cap(x.NM32) {
		x.NM32 = append(x.NM32, make([]int32, __NM32Len-len(x.NM32))...)
	}
	if __NM32Len > len(x.NM32) {
		x.NM32 = x.NM32[:__NM32Len]
	}
	copy(x.NM32, __NM32Slice)
	x.NM32 = x.NM32[:__NM32Len]
	__NM64Slice := viewer.NM64(reader)
	__NM64Len := len(__NM64Slice)
	if __NM64Len > cap(x.NM64) {
		x.NM64 = append(x.NM64, make([]int64, __NM64Len-len(x.NM64))...)
	}
	if __NM64Len > len(x.NM64) {
		x.NM64 = x.NM64[:__NM64Len]
	}
	copy(x.NM64, __NM64Slice)
	x.NM64 = x.NM64[:__NM64Len]
	__NOF32Slice := viewer.NOF32(reader)
	__NOF32Len := len(__NOF32Slice)
	if __NOF32Len > cap(x.NOF32) {
		x.NOF32 = append(x.NOF32, make([]float32, __NOF32Len-len(x.NOF32))...)
	}
	if __NOF32Len > len(x.NOF32) {
		x.NOF32 = x.NOF32[:__NOF32Len]
	}
	copy(x.NOF32, __NOF32Slice)
	x.NOF32 = x.NOF32[:__NOF32Len]
	__NOF64Slice := viewer.NOF64(reader)
	__NOF64Len := len(__NOF64Slice)
	if __NOF64Len > cap(x.NOF64) {
		x.NOF64 = append(x.NOF64, make([]float64, __NOF64Len-len(x.NOF64))...)
	}
	if __NOF64Len > len(x.NOF64) {
		x.NOF64 = x.NOF64[:__NOF64Len]
	}
	copy(x.NOF64, __NOF64Slice)
	x.NOF64 = x.NOF64[:__NOF64Len]
}

type SimpleNumbersPacked struct {
	N8    uint8
	N16   uint16
	N32   uint32
	N64   uint64
	M8    int8
	M16   int16
	M32   int32
	M64   int64
	OF32  float32
	OF64  float64
	B1    bool
	NN8   []uint8
	NN16  []uint16
	NN32  []uint32
	NN64  []uint64
	NM8   []int8
	NM16  []int16
	NM32  []int32
	NM64  []int64
	NOF32 []float32
	NOF64 []float64
}

func NewSimpleNumbersPacked() SimpleNumbersPacked {
	return SimpleNumbersPacked{}
}

func (x *SimpleNumbersPacked) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierSimpleNumbersPacked
}

func (x *SimpleNumbersPacked) Reset() {
	x.Read((*SimpleNumbersPackedViewer)(unsafe.Pointer(&_Null_demo[0])), _NullReader_demo)
}

func (x *SimpleNumbersPacked) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *SimpleNumbersPacked) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(123)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__N8Offset := offset + 0
	writer.Write1At(__N8Offset, *(*uint8)(unsafe.Pointer(&x.N8)))
	__N16Offset := offset + 1
	writer.Write2At(__N16Offset, *(*uint16)(unsafe.Pointer(&x.N16)))
	__N32Offset := offset + 3
	writer.Write4At(__N32Offset, *(*uint32)(unsafe.Pointer(&x.N32)))
	__N64Offset := offset + 7
	writer.Write8At(__N64Offset, *(*uint64)(unsafe.Pointer(&x.N64)))
	__M8Offset := offset + 15
	writer.Write1At(__M8Offset, *(*uint8)(unsafe.Pointer(&x.M8)))
	__M16Offset := offset + 16
	writer.Write2At(__M16Offset, *(*uint16)(unsafe.Pointer(&x.M16)))
	__M32Offset := offset + 18
	writer.Write4At(__M32Offset, *(*uint32)(unsafe.Pointer(&x.M32)))
	__M64Offset := offset + 22
	writer.Write8At(__M64Offset, *(*uint64)(unsafe.Pointer(&x.M64)))
	__OF32Offset := offset + 30
	writer.Write4At(__OF32Offset, *(*uint32)(unsafe.Pointer(&x.OF32)))
	__OF64Offset := offset + 34
	writer.Write8At(__OF64Offset, *(*uint64)(unsafe.Pointer(&x.OF64)))
	__B1Offset := offset + 42
	writer.Write1At(__B1Offset, *(*uint8)(unsafe.Pointer(&x.B1)))
	__NN8Size := uint(1 * len(x.NN8))
	__NN8Offset, err := writer.Alloc(__NN8Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+43, uint32(__NN8Offset))
	writer.Write4At(offset+43+4, uint32(__NN8Size))
	__NN8Slice := *(*[3]uint)(unsafe.Pointer(&x.NN8))
	__NN8Slice[1] = __NN8Size
	__NN8Slice[2] = __NN8Size
	writer.WriteAt(__NN8Offset, *(*[]byte)(unsafe.Pointer(&__NN8Slice)))
	__NN16Size := uint(2 * len(x.NN16))
	__NN16Offset, err := writer.Alloc(__NN16Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+51, uint32(__NN16Offset))
	writer.Write4At(offset+51+4, uint32(__NN16Size))
	__NN16Slice := *(*[3]uint)(unsafe.Pointer(&x.NN16))
	__NN16Slice[1] = __NN16Size
	__NN16Slice[2] = __NN16Size
	writer.WriteAt(__NN16Offset, *(*[]byte)(unsafe.Pointer(&__NN16Slice)))
	__NN32Size := uint(4 * len(x.NN32))
	__NN32Offset, err := writer.Alloc(__NN32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+59, uint32(__NN32Offset))
	writer.Write4At(offset+59+4, uint32(__NN32Size))
	__NN32Slice := *(*[3]uint)(unsafe.Pointer(&x.NN32))
	__NN32Slice[1] = __NN32Size
	__NN32Slice[2] = __NN32Size
	writer.WriteAt(__NN32Offset, *(*[]byte)(unsafe.Pointer(&__NN32Slice)))
	__NN64Size := uint(8 * len(x.NN64))
	__NN64Offset, err := writer.Alloc(__NN64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+67, uint32(__NN64Offset))
	writer.Write4At(offset+67+4, uint32(__NN64Size))
	__NN64Slice := *(*[3]uint)(unsafe.Pointer(&x.NN64))
	__NN64Slice[1] = __NN64Size
	__NN64Slice[2] = __NN64Size
	writer.WriteAt(__NN64Offset, *(*[]byte)(unsafe.Pointer(&__NN64Slice)))
	__NM8Size := uint(1 * len(x.NM8))
	__NM8Offset, err := writer.Alloc(__NM8Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+75, uint32(__NM8Offset))
	writer.Write4At(offset+75+4, uint32(__NM8Size))
	__NM8Slice := *(*[3]uint)(unsafe.Pointer(&x.NM8))
	__NM8Slice[1] = __NM8Size
	__NM8Slice[2] = __NM8Size
	writer.WriteAt(__NM8Offset, *(*[]byte)(unsafe.Pointer(&__NM8Slice)))
	__NM16Size := uint(2 * len(x.NM16))
	__NM16Offset, err := writer.Alloc(__NM16Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+83, uint32(__NM16Offset))
	writer.Write4At(offset+83+4, uint32(__NM16Size))
	__NM16Slice := *(*[3]uint)(unsafe.Pointer(&x.NM16))
	__NM16Slice[1] = __NM16Size
	__NM16Slice[2] = __NM16Size
	writer.WriteAt(__NM16Offset, *(*[]byte)(unsafe.Pointer(&__NM16Slice)))
	__NM32Size := uint(4 * len(x.NM32))
	__NM32Offset, err := writer.Alloc(__NM32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+91, uint32(__NM32Offset))
	writer.Write4At(offset+91+4, uint32(__NM32Size))
	__NM32Slice := *(*[3]uint)(unsafe.Pointer(&x.NM32))
	__NM32Slice[1] = __NM32Size
	__NM32Slice[2] = __NM32Size
	writer.WriteAt(__NM32Offset, *(*[]byte)(unsafe.Pointer(&__NM32Slice)))
	__NM64Size := uint(8 * len(x.NM64))
	__NM64Offset, err := writer.Alloc(__NM64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+99, uint32(__NM64Offset))
	writer.Write4At(offset+99+4, uint32(__NM64Size))
	__NM64Slice := *(*[3]uint)(unsafe.Pointer(&x.NM64))
	__NM64Slice[1] = __NM64Size
	__NM64Slice[2] = __NM64Size
	writer.WriteAt(__NM64Offset, *(*[]byte)(unsafe.Pointer(&__NM64Slice)))
	__NOF32Size := uint(4 * len(x.NOF32))
	__NOF32Offset, err := writer.Alloc(__NOF32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+107, uint32(__NOF32Offset))
	writer.Write4At(offset+107+4, uint32(__NOF32Size))
	__NOF32Slice := *(*[3]uint)(unsafe.Pointer(&x.NOF32))
	__NOF32Slice[1] = __NOF32Size
	__NOF32Slice[2] = __NOF32Size
	writer.WriteAt(__NOF32Offset, *(*[]byte)(unsafe.Pointer(&__NOF32Slice)))
	__NOF64Size := uint(8 * len(x.NOF64))
	__NOF64Offset, err := writer.Alloc(__NOF64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+115, uint32(__NOF64Offset))
	writer.Write4At(offset+115+4, uint32(__NOF64Size))
	__NOF64Slice := *(*[3]uint)(unsafe.Pointer(&x.NOF64))
	__NOF64Slice[1] = __NOF64Size
	__NOF64Slice[2] = __NOF64Size
	writer.WriteAt(__NOF64Offset, *(*[]byte)(unsafe.Pointer(&__NOF64Slice)))

	return offset, nil
}

func (*SimpleNumbersPacked) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(123)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_demo[:size-4])
	__N8Offset := offset + 0
	_ = __N8Offset
	__N16Offset := offset + 1
	_ = __N16Offset
	__N32Offset := offset + 3
	_ = __N32Offset
	__N64Offset := offset + 7
	_ = __N64Offset
	__M8Offset := offset + 15
	_ = __M8Offset
	__M16Offset := offset + 16
	_ = __M16Offset
	__M32Offset := offset + 18
	_ = __M32Offset
	__M64Offset := offset + 22
	_ = __M64Offset
	__OF32Offset := offset + 30
	_ = __OF32Offset
	__OF64Offset := offset + 34
	_ = __OF64Offset
	__B1Offset := offset + 42
	_ = __B1Offset
	__NN8Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+43, uint32(__NN8Offset))
	writer.Write4At(offset+43+4, 0)
	__NN16Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+51, uint32(__NN16Offset))
	writer.Write4At(offset+51+4, 0)
	__NN32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+59, uint32(__NN32Offset))
	writer.Write4At(offset+59+4, 0)
	__NN64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+67, uint32(__NN64Offset))
	writer.Write4At(offset+67+4, 0)
	__NM8Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+75, uint32(__NM8Offset))
	writer.Write4At(offset+75+4, 0)
	__NM16Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+83, uint32(__NM16Offset))
	writer.Write4At(offset+83+4, 0)
	__NM32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+91, uint32(__NM32Offset))
	writer.Write4At(offset+91+4, 0)
	__NM64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+99, uint32(__NM64Offset))
	writer.Write4At(offset+99+4, 0)
	__NOF32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+107, uint32(__NOF32Offset))
	writer.Write4At(offset+107+4, 0)
	__NOF64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+115, uint32(__NOF64Offset))
	writer.Write4At(offset+115+4, 0)

	return offset, nil
}

func (x *SimpleNumbersPacked) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewSimpleNumbersPackedViewer(reader, 0), reader)
}

func (x *SimpleNumbersPacked) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewSimpleNumbersPackedViewer(reader, offset), reader)
}

func (x *SimpleNumbersPacked) Read(viewer *SimpleNumbersPackedViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.N8 = viewer.N8()
	x.N16 = viewer.N16()
	x.N32 = viewer.N32()
	x.N64 = viewer.N64()
	x.M8 = viewer.M8()
	x.M16 = viewer.M16()
	x.M32 = viewer.M32()
	x.M64 = viewer.M64()
	x.OF32 = viewer.OF32()
	x.OF64 = viewer.OF64()
	x.B1 = viewer.B1()
	__NN8Slice := viewer.NN8(reader)
	__NN8Len := len(__NN8Slice)
	if __NN8Len > cap(x.NN8) {
		x.NN8 = append(x.NN8, make([]uint8, __NN8Len-len(x.NN8))...)
	}
	if __NN8Len > len(x.NN8) {
		x.NN8 = x.NN8[:__NN8Len]
	}
	copy(x.NN8, __NN8Slice)
	x.NN8 = x.NN8[:__NN8Len]
	__NN16Slice := viewer.NN16(reader)
	__NN16Len := len(__NN16Slice)
	if __NN16Len > cap(x.NN16) {
		x.NN16 = append(x.NN16, make([]uint16, __NN16Len-len(x.NN16))...)
	}
	if __NN16Len > len(x.NN16) {
		x.NN16 = x.NN16[:__NN16Len]
	}
	copy(x.NN16, __NN16Slice)
	x.NN16 = x.NN16[:__NN16Len]
	__NN32Slice := viewer.NN32(reader)
	__NN32Len := len(__NN32Slice)
	if __NN32Len > cap(x.NN32) {
		x.NN32 = append(x.NN32, make([]uint32, __NN32Len-len(x.NN32))...)
	}
	if __NN32Len > len(x.NN32) {
		x.NN32 = x.NN32[:__NN32Len]
	}
	copy(x.NN32, __NN32Slice)
	x.NN32 = x.NN32[:__NN32Len]
	__NN64Slice := viewer.NN64(reader)
	__NN64Len := len(__NN64Slice)
	if __NN64Len > cap(x.NN64) {
		x.NN64 = append(x.NN64, make([]uint64, __NN64Len-len(x.NN64))...)
	}
	if __NN64Len > len(x.NN64) {
		x.NN64 = x.NN64[:__NN64Len]
	}
	copy(x.NN64, __NN64Slice)
	x.NN64 = x.NN64[:__NN64Len]
	__NM8Slice := viewer.NM8(reader)
	__NM8Len := len(__NM8Slice)
	if __NM8Len > cap(x.NM8) {
		x.NM8 = append(x.NM8, make([]int8, __NM8Len-len(x.NM8))...)
	}
	if __NM8Len > len(x.NM8) {
		x.NM8 = x.NM8[:__NM8Len]
	}
	copy(x.NM8, __NM8Slice)
	x.NM8 = x.NM8[:__NM8Len]
	__NM16Slice := viewer.NM16(reader)
	__NM16Len := len(__NM16Slice)
	if __NM16Len > cap(x.NM16) {
		x.NM16 = append(x.NM16, make([]int16, __NM16Len-len(x.NM16))...)
	}
	if __NM16Len > len(x.NM16) {
		x.NM16 = x.NM16[:__NM16Len]
	}
	copy(x.NM16, __NM16Slice)
	x.NM16 = x.NM16[:__NM16Len]
	__NM32Slice := viewer.NM32(reader)
	__NM32Len := len(__NM32Slice)
	if __NM32Len > cap(x.NM32) {
		x.NM32 = append(x.NM32, make([]int32, __NM32Len-len(x.NM32))...)
	}
	if __NM32Len > len(x.NM32) {
		x.NM32 = x.NM32[:__NM32Len]
	}
	copy(x.NM32, __NM32Slice)
	x.NM32 = x.NM32[:__NM32Len]
	__NM64Slice := viewer.NM64(reader)
	__NM64Len := len(__NM64Slice)
	if __NM64Len > cap(x.NM64) {
		x.NM64 = append(x.NM64, make([]int64, __NM64Len-len(x.NM64))...)
	}
	if __NM64Len > len(x.NM64) {
		x.NM64 = x.NM64[:__NM64Len]
	}
	copy(x.NM64, __NM64Slice)
	x.NM64 = x.NM64[:__NM64Len]
	__NOF32Slice := viewer.NOF32(reader)
	__NOF32Len := len(__NOF32Slice)
	if __NOF32Len > cap(x.NOF32) {
		x.NOF32 = append(x.NOF32, make([]float32, __NOF32Len-len(x.NOF32))...)
	}
	if __NOF32Len > len(x.NOF32) {
		x.NOF32 = x.NOF32[:__NOF32Len]
	}
	copy(x.NOF32, __NOF32Slice)
	x.NOF32 = x.NOF32[:__NOF32Len]
	__NOF64Slice := viewer.NOF64(reader)
	__NOF64Len := len(__NOF64Slice)
	if __NOF64Len > cap(x.NOF64) {
		x.NOF64 = append(x.NOF64, make([]float64, __NOF64Len-len(x.NOF64))...)
	}
	if __NOF64Len > len(x.NOF64) {
		x.NOF64 = x.NOF64[:__NOF64Len]
	}
	copy(x.NOF64, __NOF64Slice)
	x.NOF64 = x.NOF64[:__NOF64Len]
}

type ComplexPacked struct {
	N8     uint8
	N16    uint16
	N32    uint32
	N64    uint64
	M8     int8
	M16    int16
	M32    int32
	M64    int64
	OF32   float32
	OF64   float64
	B1     bool
	NN8    []uint8
	NN16   []uint16
	NN32   []uint32
	NN64   []uint64
	NM8    []int8
	NM16   []int16
	NM32   []int32
	NM64   []int64
	NOF32  []float32
	NOF64  []float64
	NB     []bool
	NC     string
	NIL    []SimpleNumbers
	NIP    []SimpleNumbersPacked
	ANN8   [64]uint8
	ANN16  [64]uint16
	ANN32  [64]uint32
	ANN64  [64]uint64
	ANM8   [64]int8
	ANM16  [64]int16
	ANM32  [64]int32
	ANM64  [64]int64
	ANOF32 [64]float32
	ANOF64 [64]float64
	ANB    [64]bool
	// fixed string of size 64
	ANC string
	// fixed string of size 128
	ANC1 string
	ANIL [64]SimpleNumbers
	ANIP [64]SimpleNumbersPacked
}

func NewComplexPacked() ComplexPacked {
	return ComplexPacked{}
}

func (x *ComplexPacked) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierComplexPacked
}

func (x *ComplexPacked) Reset() {
	x.Read((*ComplexPackedViewer)(unsafe.Pointer(&_Null_demo[0])), _NullReader_demo)
}

func (x *ComplexPacked) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *ComplexPacked) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(18859)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(18859))
	__N8Offset := offset + 4
	writer.Write1At(__N8Offset, *(*uint8)(unsafe.Pointer(&x.N8)))
	__N16Offset := offset + 5
	writer.Write2At(__N16Offset, *(*uint16)(unsafe.Pointer(&x.N16)))
	__N32Offset := offset + 7
	writer.Write4At(__N32Offset, *(*uint32)(unsafe.Pointer(&x.N32)))
	__N64Offset := offset + 11
	writer.Write8At(__N64Offset, *(*uint64)(unsafe.Pointer(&x.N64)))
	__M8Offset := offset + 19
	writer.Write1At(__M8Offset, *(*uint8)(unsafe.Pointer(&x.M8)))
	__M16Offset := offset + 20
	writer.Write2At(__M16Offset, *(*uint16)(unsafe.Pointer(&x.M16)))
	__M32Offset := offset + 22
	writer.Write4At(__M32Offset, *(*uint32)(unsafe.Pointer(&x.M32)))
	__M64Offset := offset + 26
	writer.Write8At(__M64Offset, *(*uint64)(unsafe.Pointer(&x.M64)))
	__OF32Offset := offset + 34
	writer.Write4At(__OF32Offset, *(*uint32)(unsafe.Pointer(&x.OF32)))
	__OF64Offset := offset + 38
	writer.Write8At(__OF64Offset, *(*uint64)(unsafe.Pointer(&x.OF64)))
	__B1Offset := offset + 46
	writer.Write1At(__B1Offset, *(*uint8)(unsafe.Pointer(&x.B1)))
	__NN8Size := uint(1 * len(x.NN8))
	__NN8Offset, err := writer.Alloc(__NN8Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+47, uint32(__NN8Offset))
	writer.Write4At(offset+47+4, uint32(__NN8Size))
	__NN8Slice := *(*[3]uint)(unsafe.Pointer(&x.NN8))
	__NN8Slice[1] = __NN8Size
	__NN8Slice[2] = __NN8Size
	writer.WriteAt(__NN8Offset, *(*[]byte)(unsafe.Pointer(&__NN8Slice)))
	__NN16Size := uint(2 * len(x.NN16))
	__NN16Offset, err := writer.Alloc(__NN16Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+55, uint32(__NN16Offset))
	writer.Write4At(offset+55+4, uint32(__NN16Size))
	__NN16Slice := *(*[3]uint)(unsafe.Pointer(&x.NN16))
	__NN16Slice[1] = __NN16Size
	__NN16Slice[2] = __NN16Size
	writer.WriteAt(__NN16Offset, *(*[]byte)(unsafe.Pointer(&__NN16Slice)))
	__NN32Size := uint(4 * len(x.NN32))
	__NN32Offset, err := writer.Alloc(__NN32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+63, uint32(__NN32Offset))
	writer.Write4At(offset+63+4, uint32(__NN32Size))
	__NN32Slice := *(*[3]uint)(unsafe.Pointer(&x.NN32))
	__NN32Slice[1] = __NN32Size
	__NN32Slice[2] = __NN32Size
	writer.WriteAt(__NN32Offset, *(*[]byte)(unsafe.Pointer(&__NN32Slice)))
	__NN64Size := uint(8 * len(x.NN64))
	__NN64Offset, err := writer.Alloc(__NN64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+71, uint32(__NN64Offset))
	writer.Write4At(offset+71+4, uint32(__NN64Size))
	__NN64Slice := *(*[3]uint)(unsafe.Pointer(&x.NN64))
	__NN64Slice[1] = __NN64Size
	__NN64Slice[2] = __NN64Size
	writer.WriteAt(__NN64Offset, *(*[]byte)(unsafe.Pointer(&__NN64Slice)))
	__NM8Size := uint(1 * len(x.NM8))
	__NM8Offset, err := writer.Alloc(__NM8Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+79, uint32(__NM8Offset))
	writer.Write4At(offset+79+4, uint32(__NM8Size))
	__NM8Slice := *(*[3]uint)(unsafe.Pointer(&x.NM8))
	__NM8Slice[1] = __NM8Size
	__NM8Slice[2] = __NM8Size
	writer.WriteAt(__NM8Offset, *(*[]byte)(unsafe.Pointer(&__NM8Slice)))
	__NM16Size := uint(2 * len(x.NM16))
	__NM16Offset, err := writer.Alloc(__NM16Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+87, uint32(__NM16Offset))
	writer.Write4At(offset+87+4, uint32(__NM16Size))
	__NM16Slice := *(*[3]uint)(unsafe.Pointer(&x.NM16))
	__NM16Slice[1] = __NM16Size
	__NM16Slice[2] = __NM16Size
	writer.WriteAt(__NM16Offset, *(*[]byte)(unsafe.Pointer(&__NM16Slice)))
	__NM32Size := uint(4 * len(x.NM32))
	__NM32Offset, err := writer.Alloc(__NM32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+95, uint32(__NM32Offset))
	writer.Write4At(offset+95+4, uint32(__NM32Size))
	__NM32Slice := *(*[3]uint)(unsafe.Pointer(&x.NM32))
	__NM32Slice[1] = __NM32Size
	__NM32Slice[2] = __NM32Size
	writer.WriteAt(__NM32Offset, *(*[]byte)(unsafe.Pointer(&__NM32Slice)))
	__NM64Size := uint(8 * len(x.NM64))
	__NM64Offset, err := writer.Alloc(__NM64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+103, uint32(__NM64Offset))
	writer.Write4At(offset+103+4, uint32(__NM64Size))
	__NM64Slice := *(*[3]uint)(unsafe.Pointer(&x.NM64))
	__NM64Slice[1] = __NM64Size
	__NM64Slice[2] = __NM64Size
	writer.WriteAt(__NM64Offset, *(*[]byte)(unsafe.Pointer(&__NM64Slice)))
	__NOF32Size := uint(4 * len(x.NOF32))
	__NOF32Offset, err := writer.Alloc(__NOF32Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+111, uint32(__NOF32Offset))
	writer.Write4At(offset+111+4, uint32(__NOF32Size))
	__NOF32Slice := *(*[3]uint)(unsafe.Pointer(&x.NOF32))
	__NOF32Slice[1] = __NOF32Size
	__NOF32Slice[2] = __NOF32Size
	writer.WriteAt(__NOF32Offset, *(*[]byte)(unsafe.Pointer(&__NOF32Slice)))
	__NOF64Size := uint(8 * len(x.NOF64))
	__NOF64Offset, err := writer.Alloc(__NOF64Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+119, uint32(__NOF64Offset))
	writer.Write4At(offset+119+4, uint32(__NOF64Size))
	__NOF64Slice := *(*[3]uint)(unsafe.Pointer(&x.NOF64))
	__NOF64Slice[1] = __NOF64Size
	__NOF64Slice[2] = __NOF64Size
	writer.WriteAt(__NOF64Offset, *(*[]byte)(unsafe.Pointer(&__NOF64Slice)))
	__NBSize := uint(1 * len(x.NB))
	__NBOffset, err := writer.Alloc(__NBSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+127, uint32(__NBOffset))
	writer.Write4At(offset+127+4, uint32(__NBSize))
	__NBSlice := *(*[3]uint)(unsafe.Pointer(&x.NB))
	__NBSlice[1] = __NBSize
	__NBSlice[2] = __NBSize
	writer.WriteAt(__NBOffset, *(*[]byte)(unsafe.Pointer(&__NBSlice)))
	__NCSize := uint(1 * len(x.NC))
	__NCOffset, err := writer.Alloc(__NCSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+135, uint32(__NCOffset))
	writer.Write4At(offset+135+4, uint32(__NCSize))
	__NCSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.NC)), __NCSize, __NCSize}
	writer.WriteAt(__NCOffset, *(*[]byte)(unsafe.Pointer(&__NCSlice)))
	__NILSize := uint(123 * len(x.NIL))
	__NILOffset, err := writer.Alloc(__NILSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+147, uint32(__NILOffset))
	writer.Write4At(offset+147+4, uint32(__NILSize))
	for i := range x.NIL {
		if _, err := x.NIL[i].Write(writer, __NILOffset); err != nil {
			return offset, err
		}
		__NILOffset += 123
	}
	__NIPSize := uint(123 * len(x.NIP))
	__NIPOffset, err := writer.Alloc(__NIPSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+155, uint32(__NIPOffset))
	writer.Write4At(offset+155+4, uint32(__NIPSize))
	for i := range x.NIP {
		if _, err := x.NIP[i].Write(writer, __NIPOffset); err != nil {
			return offset, err
		}
		__NIPOffset += 123
	}
	__ANN8Offset := offset + 163
	writer.WriteAt(__ANN8Offset, (*[64]byte)(unsafe.Pointer(&x.ANN8))[:])
	__ANN16Offset := offset + 227
	writer.WriteAt(__ANN16Offset, (*[128]byte)(unsafe.Pointer(&x.ANN16))[:])
	__ANN32Offset := offset + 355
	writer.WriteAt(__ANN32Offset, (*[256]byte)(unsafe.Pointer(&x.ANN32))[:])
	__ANN64Offset := offset + 611
	writer.WriteAt(__ANN64Offset, (*[512]byte)(unsafe.Pointer(&x.ANN64))[:])
	__ANM8Offset := offset + 1123
	writer.WriteAt(__ANM8Offset, (*[64]byte)(unsafe.Pointer(&x.ANM8))[:])
	__ANM16Offset := offset + 1187
	writer.WriteAt(__ANM16Offset, (*[128]byte)(unsafe.Pointer(&x.ANM16))[:])
	__ANM32Offset := offset + 1315
	writer.WriteAt(__ANM32Offset, (*[256]byte)(unsafe.Pointer(&x.ANM32))[:])
	__ANM64Offset := offset + 1571
	writer.WriteAt(__ANM64Offset, (*[512]byte)(unsafe.Pointer(&x.ANM64))[:])
	__ANOF32Offset := offset + 2083
	writer.WriteAt(__ANOF32Offset, (*[256]byte)(unsafe.Pointer(&x.ANOF32))[:])
	__ANOF64Offset := offset + 2339
	writer.WriteAt(__ANOF64Offset, (*[512]byte)(unsafe.Pointer(&x.ANOF64))[:])
	__ANBOffset := offset + 2851
	writer.WriteAt(__ANBOffset, (*[64]byte)(unsafe.Pointer(&x.ANB))[:])
	__ANCOffset := offset + 2915
	__ANCSize := uint(1 * len(x.ANC))
	writer.Write4At(__ANCOffset, uint32(__ANCSize))
	if __ANCSize > 0 {
		if __ANCSize > 64 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANCOffset+4, (*[64]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC))))[:])
	} else {
		writer.WriteAt(__ANCOffset+4, _Null_demo[:64])
	}
	__ANC1Offset := offset + 2983
	__ANC1Size := uint(1 * len(x.ANC1))
	writer.Write4At(__ANC1Offset, uint32(__ANC1Size))
	if __ANC1Size > 0 {
		if __ANC1Size > 128 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ANC1Offset+4, (*[128]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&x.ANC1))))[:])
	} else {
		writer.WriteAt(__ANC1Offset+4, _Null_demo[:128])
	}
	__ANILOffset := offset + 3115
	for i := range x.ANIL {
		if _, err := x.ANIL[i].Write(writer, __ANILOffset); err != nil {
			return offset, err
		}
		__ANILOffset += 123
	}
	__ANIPOffset := offset + 10987
	for i := range x.ANIP {
		if _, err := x.ANIP[i].Write(writer, __ANIPOffset); err != nil {
			return offset, err
		}
		__ANIPOffset += 123
	}

	return offset, nil
}

func (*ComplexPacked) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(18859)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(18859))
	writer.WriteAt(offset, _Null_demo[:size-4])
	__N8Offset := offset + 4
	_ = __N8Offset
	__N16Offset := offset + 5
	_ = __N16Offset
	__N32Offset := offset + 7
	_ = __N32Offset
	__N64Offset := offset + 11
	_ = __N64Offset
	__M8Offset := offset + 19
	_ = __M8Offset
	__M16Offset := offset + 20
	_ = __M16Offset
	__M32Offset := offset + 22
	_ = __M32Offset
	__M64Offset := offset + 26
	_ = __M64Offset
	__OF32Offset := offset + 34
	_ = __OF32Offset
	__OF64Offset := offset + 38
	_ = __OF64Offset
	__B1Offset := offset + 46
	_ = __B1Offset
	__NN8Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+47, uint32(__NN8Offset))
	writer.Write4At(offset+47+4, 0)
	__NN16Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+55, uint32(__NN16Offset))
	writer.Write4At(offset+55+4, 0)
	__NN32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+63, uint32(__NN32Offset))
	writer.Write4At(offset+63+4, 0)
	__NN64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+71, uint32(__NN64Offset))
	writer.Write4At(offset+71+4, 0)
	__NM8Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+79, uint32(__NM8Offset))
	writer.Write4At(offset+79+4, 0)
	__NM16Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+87, uint32(__NM16Offset))
	writer.Write4At(offset+87+4, 0)
	__NM32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+95, uint32(__NM32Offset))
	writer.Write4At(offset+95+4, 0)
	__NM64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+103, uint32(__NM64Offset))
	writer.Write4At(offset+103+4, 0)
	__NOF32Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+111, uint32(__NOF32Offset))
	writer.Write4At(offset+111+4, 0)
	__NOF64Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+119, uint32(__NOF64Offset))
	writer.Write4At(offset+119+4, 0)
	__NBOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+127, uint32(__NBOffset))
	writer.Write4At(offset+127+4, 0)
	__NCOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+135, uint32(__NCOffset))
	writer.Write4At(offset+135+4, 0)
	__NILOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+147, uint32(__NILOffset))
	writer.Write4At(offset+147+4, 0)
	__NIPOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+155, uint32(__NIPOffset))
	writer.Write4At(offset+155+4, 0)
	__ANN8Offset := offset + 163
	_ = __ANN8Offset
	__ANN16Offset := offset + 227
	_ = __ANN16Offset
	__ANN32Offset := offset + 355
	_ = __ANN32Offset
	__ANN64Offset := offset + 611
	_ = __ANN64Offset
	__ANM8Offset := offset + 1123
	_ = __ANM8Offset
	__ANM16Offset := offset + 1187
	_ = __ANM16Offset
	__ANM32Offset := offset + 1315
	_ = __ANM32Offset
	__ANM64Offset := offset + 1571
	_ = __ANM64Offset
	__ANOF32Offset := offset + 2083
	_ = __ANOF32Offset
	__ANOF64Offset := offset + 2339
	_ = __ANOF64Offset
	__ANBOffset := offset + 2851
	_ = __ANBOffset
	__ANCOffset := offset + 2915
	_ = __ANCOffset
	writer.Write4At(__ANCOffset, 0)
	__ANC1Offset := offset + 2983
	_ = __ANC1Offset
	writer.Write4At(__ANC1Offset, 0)
	__ANILOffset := offset + 3115
	_ = __ANILOffset
	__ANIPOffset := offset + 10987
	_ = __ANIPOffset

	return offset, nil
}

func (x *ComplexPacked) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewComplexPackedViewer(reader, 0), reader)
}

func (x *ComplexPacked) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewComplexPackedViewer(reader, offset), reader)
}

func (x *ComplexPacked) Read(viewer *ComplexPackedViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.N8 = viewer.N8()
	x.N16 = viewer.N16()
	x.N32 = viewer.N32()
	x.N64 = viewer.N64()
	x.M8 = viewer.M8()
	x.M16 = viewer.M16()
	x.M32 = viewer.M32()
	x.M64 = viewer.M64()
	x.OF32 = viewer.OF32()
	x.OF64 = viewer.OF64()
	x.B1 = viewer.B1()
	__NN8Slice := viewer.NN8(reader)
	__NN8Len := len(__NN8Slice)
	if __NN8Len > cap(x.NN8) {
		x.NN8 = append(x.NN8, make([]uint8, __NN8Len-len(x.NN8))...)
	}
	if __NN8Len > len(x.NN8) {
		x.NN8 = x.NN8[:__NN8Len]
	}
	copy(x.NN8, __NN8Slice)
	x.NN8 = x.NN8[:__NN8Len]
	__NN16Slice := viewer.NN16(reader)
	__NN16Len := len(__NN16Slice)
	if __NN16Len > cap(x.NN16) {
		x.NN16 = append(x.NN16, make([]uint16, __NN16Len-len(x.NN16))...)
	}
	if __NN16Len > len(x.NN16) {
		x.NN16 = x.NN16[:__NN16Len]
	}
	copy(x.NN16, __NN16Slice)
	x.NN16 = x.NN16[:__NN16Len]
	__NN32Slice := viewer.NN32(reader)
	__NN32Len := len(__NN32Slice)
	if __NN32Len > cap(x.NN32) {
		x.NN32 = append(x.NN32, make([]uint32, __NN32Len-len(x.NN32))...)
	}
	if __NN32Len > len(x.NN32) {
		x.NN32 = x.NN32[:__NN32Len]
	}
	copy(x.NN32, __NN32Slice)
	x.NN32 = x.NN32[:__NN32Len]
	__NN64Slice := viewer.NN64(reader)
	__NN64Len := len(__NN64Slice)
	if __NN64Len > cap(x.NN64) {
		x.NN64 = append(x.NN64, make([]uint64, __NN64Len-len(x.NN64))...)
	}
	if __NN64Len > len(x.NN64) {
		x.NN64 = x.NN64[:__NN64Len]
	}
	copy(x.NN64, __NN64Slice)
	x.NN64 = x.NN64[:__NN64Len]
	__NM8Slice := viewer.NM8(reader)
	__NM8Len := len(__NM8Slice)
	if __NM8Len > cap(x.NM8) {
		x.NM8 = append(x.NM8, make([]int8, __NM8Len-len(x.NM8))...)
	}
	if __NM8Len > len(x.NM8) {
		x.NM8 = x.NM8[:__NM8Len]
	}
	copy(x.NM8, __NM8Slice)
	x.NM8 = x.NM8[:__NM8Len]
	__NM16Slice := viewer.NM16(reader)
	__NM16Len := len(__NM16Slice)
	if __NM16Len > cap(x.NM16) {
		x.NM16 = append(x.NM16, make([]int16, __NM16Len-len(x.NM16))...)
	}
	if __NM16Len > len(x.NM16) {
		x.NM16 = x.NM16[:__NM16Len]
	}
	copy(x.NM16, __NM16Slice)
	x.NM16 = x.NM16[:__NM16Len]
	__NM32Slice := viewer.NM32(reader)
	__NM32Len := len(__NM32Slice)
	if __NM32Len > cap(x.NM32) {
		x.NM32 = append(x.NM32, make([]int32, __NM32Len-len(x.NM32))...)
	}
	if __NM32Len > len(x.NM32) {
		x.NM32 = x.NM32[:__NM32Len]
	}
	copy(x.NM32, __NM32Slice)
	x.NM32 = x.NM32[:__NM32Len]
	__NM64Slice := viewer.NM64(reader)
	__NM64Len := len(__NM64Slice)
	if __NM64Len > cap(x.NM64) {
		x.NM64 = append(x.NM64, make([]int64, __NM64Len-len(x.NM64))...)
	}
	if __NM64Len > len(x.NM64) {
		x.NM64 = x.NM64[:__NM64Len]
	}
	copy(x.NM64, __NM64Slice)
	x.NM64 = x.NM64[:__NM64Len]
	__NOF32Slice := viewer.NOF32(reader)
	__NOF32Len := len(__NOF32Slice)
	if __NOF32Len > cap(x.NOF32) {
		x.NOF32 = append(x.NOF32, make([]float32, __NOF32Len-len(x.NOF32))...)
	}
	if __NOF32Len > len(x.NOF32) {
		x.NOF32 = x.NOF32[:__NOF32Len]
	}
	copy(x.NOF32, __NOF32Slice)
	x.NOF32 = x.NOF32[:__NOF32Len]
	__NOF64Slice := viewer.NOF64(reader)
	__NOF64Len := len(__NOF64Slice)
	if __NOF64Len > cap(x.NOF64) {
		x.NOF64 = append(x.NOF64, make([]float64, __NOF64Len-len(x.NOF64))...)
	}
	if __NOF64Len > len(x.NOF64) {
		x.NOF64 = x.NOF64[:__NOF64Len]
	}
	copy(x.NOF64, __NOF64Slice)
	x.NOF64 = x.NOF64[:__NOF64Len]
	__NBSlice := viewer.NB(reader)
	__NBLen := len(__NBSlice)
	if __NBLen > cap(x.NB) {
		x.NB = append(x.NB, make([]bool, __NBLen-len(x.NB))...)
	}
	if __NBLen > len(x.NB) {
		x.NB = x.NB[:__NBLen]
	}
	copy(x.NB, __NBSlice)
	x.NB = x.NB[:__NBLen]
	__NCString := viewer.NC(reader)
	if x.NC != __NCString {
		__NCStringCopy := make([]byte, len(__NCString))
		copy(__NCStringCopy, __NCString)
		x.NC = *(*string)(unsafe.Pointer(&__NCStringCopy))
	}
	__NILSlice := viewer.NIL(reader)
	__NILLen := len(__NILSlice)
	if __NILLen > cap(x.NIL) {
		x.NIL = append(x.NIL, make([]SimpleNumbers, __NILLen-len(x.NIL))...)
	}
	if __NILLen > len(x.NIL) {
		x.NIL = x.NIL[:__NILLen]
	}
	for i := 0; i < __NILLen; i++ {
		x.NIL[i].Read(&__NILSlice[i], reader)
	}
	x.NIL = x.NIL[:__NILLen]
	__NIPSlice := viewer.NIP(reader)
	__NIPLen := len(__NIPSlice)
	if __NIPLen > cap(x.NIP) {
		x.NIP = append(x.NIP, make([]SimpleNumbersPacked, __NIPLen-len(x.NIP))...)
	}
	if __NIPLen > len(x.NIP) {
		x.NIP = x.NIP[:__NIPLen]
	}
	for i := 0; i < __NIPLen; i++ {
		x.NIP[i].Read(&__NIPSlice[i], reader)
	}
	x.NIP = x.NIP[:__NIPLen]
	__ANN8Slice := viewer.ANN8()
	__ANN8Len := len(__ANN8Slice)
	if __ANN8Len > 64 {
		__ANN8Len = 64
	}
	copy(x.ANN8[:], __ANN8Slice)
	for i := __ANN8Len; i < len(x.ANN8); i++ {
		x.ANN8[i] = 0
	}
	__ANN16Slice := viewer.ANN16()
	__ANN16Len := len(__ANN16Slice)
	if __ANN16Len > 64 {
		__ANN16Len = 64
	}
	copy(x.ANN16[:], __ANN16Slice)
	for i := __ANN16Len; i < len(x.ANN16); i++ {
		x.ANN16[i] = 0
	}
	__ANN32Slice := viewer.ANN32()
	__ANN32Len := len(__ANN32Slice)
	if __ANN32Len > 64 {
		__ANN32Len = 64
	}
	copy(x.ANN32[:], __ANN32Slice)
	for i := __ANN32Len; i < len(x.ANN32); i++ {
		x.ANN32[i] = 0
	}
	__ANN64Slice := viewer.ANN64()
	__ANN64Len := len(__ANN64Slice)
	if __ANN64Len > 64 {
		__ANN64Len = 64
	}
	copy(x.ANN64[:], __ANN64Slice)
	for i := __ANN64Len; i < len(x.ANN64); i++ {
		x.ANN64[i] = 0
	}
	__ANM8Slice := viewer.ANM8()
	__ANM8Len := len(__ANM8Slice)
	if __ANM8Len > 64 {
		__ANM8Len = 64
	}
	copy(x.ANM8[:], __ANM8Slice)
	for i := __ANM8Len; i < len(x.ANM8); i++ {
		x.ANM8[i] = 0
	}
	__ANM16Slice := viewer.ANM16()
	__ANM16Len := len(__ANM16Slice)
	if __ANM16Len > 64 {
		__ANM16Len = 64
	}
	copy(x.ANM16[:], __ANM16Slice)
	for i := __ANM16Len; i < len(x.ANM16); i++ {
		x.ANM16[i] = 0
	}
	__ANM32Slice := viewer.ANM32()
	__ANM32Len := len(__ANM32Slice)
	if __ANM32Len > 64 {
		__ANM32Len = 64
	}
	copy(x.ANM32[:], __ANM32Slice)
	for i := __ANM32Len; i < len(x.ANM32); i++ {
		x.ANM32[i] = 0
	}
	__ANM64Slice := viewer.ANM64()
	__ANM64Len := len(__ANM64Slice)
	if __ANM64Len > 64 {
		__ANM64Len = 64
	}
	copy(x.ANM64[:], __ANM64Slice)
	for i := __ANM64Len; i < len(x.ANM64); i++ {
		x.ANM64[i] = 0
	}
	__ANOF32Slice := viewer.ANOF32()
	__ANOF32Len := len(__ANOF32Slice)
	if __ANOF32Len > 64 {
		__ANOF32Len = 64
	}
	copy(x.ANOF32[:], __ANOF32Slice)
	for i := __ANOF32Len; i < len(x.ANOF32); i++ {
		x.ANOF32[i] = 0
	}
	__ANOF64Slice := viewer.ANOF64()
	__ANOF64Len := len(__ANOF64Slice)
	if __ANOF64Len > 64 {
		__ANOF64Len = 64
	}
	copy(x.ANOF64[:], __ANOF64Slice)
	for i := __ANOF64Len; i < len(x.ANOF64); i++ {
		x.ANOF64[i] = 0
	}
	__ANBSlice := viewer.ANB()
	__ANBLen := len(__ANBSlice)
	if __ANBLen > 64 {
		__ANBLen = 64
	}
	copy(x.ANB[:], __ANBSlice)
	for i := __ANBLen; i < len(x.ANB); i++ {
		x.ANB[i] = false
	}
	__ANCString := viewer.ANC()
	if x.ANC != __ANCString {
		__ANCStringCopy := make([]byte, len(__ANCString))
		copy(__ANCStringCopy, __ANCString)
		x.ANC = *(*string)(unsafe.Pointer(&__ANCStringCopy))
	}
	__ANC1String := viewer.ANC1()
	if x.ANC1 != __ANC1String {
		__ANC1StringCopy := make([]byte, len(__ANC1String))
		copy(__ANC1StringCopy, __ANC1String)
		x.ANC1 = *(*string)(unsafe.Pointer(&__ANC1StringCopy))
	}
	__ANILSlice := viewer.ANIL()
	__ANILLen := len(__ANILSlice)
	if __ANILLen > 64 {
		__ANILLen = 64
	}
	for i := 0; i < __ANILLen; i++ {
		x.ANIL[i].Read(&__ANILSlice[i], reader)
	}
	for i := __ANILLen; i < len(x.ANIL); i++ {
		x.ANIL[i].Reset()
	}
	__ANIPSlice := viewer.ANIP()
	__ANIPLen := len(__ANIPSlice)
	if __ANIPLen > 64 {
		__ANIPLen = 64
	}
	for i := 0; i < __ANIPLen; i++ {
		x.ANIP[i].Read(&__ANIPSlice[i], reader)
	}
	for i := __ANIPLen; i < len(x.ANIP); i++ {
		x.ANIP[i].Reset()
	}
}

type SimpleNumbersViewer [123]byte

type SimpleNumbersSource struct {
	*SimpleNumbersViewer
}

func NewSimpleNumbersViewer(reader *karmem.Reader, offset uint32) (v *SimpleNumbersViewer) {
	if !reader.IsValidOffset(offset, 123) {
		return (*SimpleNumbersViewer)(unsafe.Pointer(&_Null_demo[0]))
	}
	v = (*SimpleNumbersViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewSimpleNumbersSource(reader *karmem.Reader, offset uint32) (v SimpleNumbersSource) {
	return SimpleNumbersSource{NewSimpleNumbersViewer(reader, offset)}
}

func newSimpleNumbersSourceByViewer(viewer *SimpleNumbersViewer) (v SimpleNumbersSource) {
	return SimpleNumbersSource{viewer}
}

func (x *SimpleNumbersViewer) size() uint32 {
	return 123
}

func (x *SimpleNumbersViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *SimpleNumbersViewer) N8() (v uint8) {
	return *(*uint8)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *SimpleNumbersSource) SetN8(v uint8) {
	*(*uint8)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 0)) = v
}

func (x *SimpleNumbersViewer) N16() (v uint16) {
	return *(*uint16)(unsafe.Add(unsafe.Pointer(x), 1))
}

func (x *SimpleNumbersSource) SetN16(v uint16) {
	*(*uint16)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 1)) = v
}

func (x *SimpleNumbersViewer) N32() (v uint32) {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(x), 3))
}

func (x *SimpleNumbersSource) SetN32(v uint32) {
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 3)) = v
}

func (x *SimpleNumbersViewer) N64() (v uint64) {
	return *(*uint64)(unsafe.Add(unsafe.Pointer(x), 7))
}

func (x *SimpleNumbersSource) SetN64(v uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 7)) = v
}

func (x *SimpleNumbersViewer) M8() (v int8) {
	return *(*int8)(unsafe.Add(unsafe.Pointer(x), 15))
}

func (x *SimpleNumbersSource) SetM8(v int8) {
	*(*int8)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 15)) = v
}

func (x *SimpleNumbersViewer) M16() (v int16) {
	return *(*int16)(unsafe.Add(unsafe.Pointer(x), 16))
}

func (x *SimpleNumbersSource) SetM16(v int16) {
	*(*int16)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 16)) = v
}

func (x *SimpleNumbersViewer) M32() (v int32) {
	return *(*int32)(unsafe.Add(unsafe.Pointer(x), 18))
}

func (x *SimpleNumbersSource) SetM32(v int32) {
	*(*int32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 18)) = v
}

func (x *SimpleNumbersViewer) M64() (v int64) {
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 22))
}

func (x *SimpleNumbersSource) SetM64(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 22)) = v
}

func (x *SimpleNumbersViewer) OF32() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 30))
}

func (x *SimpleNumbersSource) SetOF32(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 30)) = v
}

func (x *SimpleNumbersViewer) OF64() (v float64) {
	return *(*float64)(unsafe.Add(unsafe.Pointer(x), 34))
}

func (x *SimpleNumbersSource) SetOF64(v float64) {
	*(*float64)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 34)) = v
}

func (x *SimpleNumbersViewer) B1() (v bool) {
	return *(*bool)(unsafe.Add(unsafe.Pointer(x), 42))
}

func (x *SimpleNumbersSource) SetB1(v bool) {
	*(*bool)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 42)) = v
}

func (x *SimpleNumbersViewer) NN8(reader *karmem.Reader) (v []uint8) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 43))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 43+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint8)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNN8(v []uint8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 43))
	__NN8Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 43+4))
	if uint32(__NN8Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN8Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN8Size, __NN8Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NN8Size, __NN8Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN8Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 43+4)) = uint32(__NN8Size)
	return nil
}

func (x *SimpleNumbersViewer) NN16(reader *karmem.Reader) (v []uint16) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 51))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 51+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint16{}
	}
	length := uintptr(size / 2)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint16)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNN16(v []uint16) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 51))
	__NN16Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 51+4))
	if uint32(__NN16Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN16Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN16Size, __NN16Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NN16Size, __NN16Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN16Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 51+4)) = uint32(__NN16Size)
	return nil
}

func (x *SimpleNumbersViewer) NN32(reader *karmem.Reader) (v []uint32) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 59))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 59+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint32)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNN32(v []uint32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 59))
	__NN32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 59+4))
	if uint32(__NN32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN32Size, __NN32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NN32Size, __NN32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 59+4)) = uint32(__NN32Size)
	return nil
}

func (x *SimpleNumbersViewer) NN64(reader *karmem.Reader) (v []uint64) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 67))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 67+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint64)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNN64(v []uint64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 67))
	__NN64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 67+4))
	if uint32(__NN64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN64Size, __NN64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NN64Size, __NN64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 67+4)) = uint32(__NN64Size)
	return nil
}

func (x *SimpleNumbersViewer) NM8(reader *karmem.Reader) (v []int8) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 75))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 75+4))
	if !reader.IsValidOffset(offset, size) {
		return []int8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int8)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNM8(v []int8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 75))
	__NM8Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 75+4))
	if uint32(__NM8Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM8Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM8Size, __NM8Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NM8Size, __NM8Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM8Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 75+4)) = uint32(__NM8Size)
	return nil
}

func (x *SimpleNumbersViewer) NM16(reader *karmem.Reader) (v []int16) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 83))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 83+4))
	if !reader.IsValidOffset(offset, size) {
		return []int16{}
	}
	length := uintptr(size / 2)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int16)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNM16(v []int16) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 83))
	__NM16Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 83+4))
	if uint32(__NM16Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM16Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM16Size, __NM16Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NM16Size, __NM16Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM16Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 83+4)) = uint32(__NM16Size)
	return nil
}

func (x *SimpleNumbersViewer) NM32(reader *karmem.Reader) (v []int32) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 91))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 91+4))
	if !reader.IsValidOffset(offset, size) {
		return []int32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int32)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNM32(v []int32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 91))
	__NM32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 91+4))
	if uint32(__NM32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM32Size, __NM32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NM32Size, __NM32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 91+4)) = uint32(__NM32Size)
	return nil
}

func (x *SimpleNumbersViewer) NM64(reader *karmem.Reader) (v []int64) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 99))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 99+4))
	if !reader.IsValidOffset(offset, size) {
		return []int64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNM64(v []int64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 99))
	__NM64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 99+4))
	if uint32(__NM64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM64Size, __NM64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NM64Size, __NM64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 99+4)) = uint32(__NM64Size)
	return nil
}

func (x *SimpleNumbersViewer) NOF32(reader *karmem.Reader) (v []float32) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 107))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 107+4))
	if !reader.IsValidOffset(offset, size) {
		return []float32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]float32)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNOF32(v []float32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 107))
	__NOF32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 107+4))
	if uint32(__NOF32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NOF32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NOF32Size, __NOF32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NOF32Size, __NOF32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NOF32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 107+4)) = uint32(__NOF32Size)
	return nil
}

func (x *SimpleNumbersViewer) NOF64(reader *karmem.Reader) (v []float64) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 115))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 115+4))
	if !reader.IsValidOffset(offset, size) {
		return []float64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]float64)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersSource) SetNOF64(v []float64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 115))
	__NOF64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 115+4))
	if uint32(__NOF64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NOF64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NOF64Size, __NOF64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), offset)), __NOF64Size, __NOF64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NOF64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersViewer), 115+4)) = uint32(__NOF64Size)
	return nil
}

type SimpleNumbersPackedViewer [123]byte

type SimpleNumbersPackedSource struct {
	*SimpleNumbersPackedViewer
}

func NewSimpleNumbersPackedViewer(reader *karmem.Reader, offset uint32) (v *SimpleNumbersPackedViewer) {
	if !reader.IsValidOffset(offset, 123) {
		return (*SimpleNumbersPackedViewer)(unsafe.Pointer(&_Null_demo[0]))
	}
	v = (*SimpleNumbersPackedViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewSimpleNumbersPackedSource(reader *karmem.Reader, offset uint32) (v SimpleNumbersPackedSource) {
	return SimpleNumbersPackedSource{NewSimpleNumbersPackedViewer(reader, offset)}
}

func newSimpleNumbersPackedSourceByViewer(viewer *SimpleNumbersPackedViewer) (v SimpleNumbersPackedSource) {
	return SimpleNumbersPackedSource{viewer}
}

func (x *SimpleNumbersPackedViewer) size() uint32 {
	return 123
}

func (x *SimpleNumbersPackedViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *SimpleNumbersPackedViewer) N8() (v uint8) {
	return *(*uint8)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *SimpleNumbersPackedSource) SetN8(v uint8) {
	*(*uint8)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 0)) = v
}

func (x *SimpleNumbersPackedViewer) N16() (v uint16) {
	return *(*uint16)(unsafe.Add(unsafe.Pointer(x), 1))
}

func (x *SimpleNumbersPackedSource) SetN16(v uint16) {
	*(*uint16)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 1)) = v
}

func (x *SimpleNumbersPackedViewer) N32() (v uint32) {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(x), 3))
}

func (x *SimpleNumbersPackedSource) SetN32(v uint32) {
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 3)) = v
}

func (x *SimpleNumbersPackedViewer) N64() (v uint64) {
	return *(*uint64)(unsafe.Add(unsafe.Pointer(x), 7))
}

func (x *SimpleNumbersPackedSource) SetN64(v uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 7)) = v
}

func (x *SimpleNumbersPackedViewer) M8() (v int8) {
	return *(*int8)(unsafe.Add(unsafe.Pointer(x), 15))
}

func (x *SimpleNumbersPackedSource) SetM8(v int8) {
	*(*int8)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 15)) = v
}

func (x *SimpleNumbersPackedViewer) M16() (v int16) {
	return *(*int16)(unsafe.Add(unsafe.Pointer(x), 16))
}

func (x *SimpleNumbersPackedSource) SetM16(v int16) {
	*(*int16)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 16)) = v
}

func (x *SimpleNumbersPackedViewer) M32() (v int32) {
	return *(*int32)(unsafe.Add(unsafe.Pointer(x), 18))
}

func (x *SimpleNumbersPackedSource) SetM32(v int32) {
	*(*int32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 18)) = v
}

func (x *SimpleNumbersPackedViewer) M64() (v int64) {
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 22))
}

func (x *SimpleNumbersPackedSource) SetM64(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 22)) = v
}

func (x *SimpleNumbersPackedViewer) OF32() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 30))
}

func (x *SimpleNumbersPackedSource) SetOF32(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 30)) = v
}

func (x *SimpleNumbersPackedViewer) OF64() (v float64) {
	return *(*float64)(unsafe.Add(unsafe.Pointer(x), 34))
}

func (x *SimpleNumbersPackedSource) SetOF64(v float64) {
	*(*float64)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 34)) = v
}

func (x *SimpleNumbersPackedViewer) B1() (v bool) {
	return *(*bool)(unsafe.Add(unsafe.Pointer(x), 42))
}

func (x *SimpleNumbersPackedSource) SetB1(v bool) {
	*(*bool)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 42)) = v
}

func (x *SimpleNumbersPackedViewer) NN8(reader *karmem.Reader) (v []uint8) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 43))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 43+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint8)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNN8(v []uint8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 43))
	__NN8Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 43+4))
	if uint32(__NN8Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN8Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN8Size, __NN8Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NN8Size, __NN8Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN8Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 43+4)) = uint32(__NN8Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NN16(reader *karmem.Reader) (v []uint16) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 51))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 51+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint16{}
	}
	length := uintptr(size / 2)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint16)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNN16(v []uint16) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 51))
	__NN16Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 51+4))
	if uint32(__NN16Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN16Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN16Size, __NN16Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NN16Size, __NN16Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN16Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 51+4)) = uint32(__NN16Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NN32(reader *karmem.Reader) (v []uint32) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 59))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 59+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint32)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNN32(v []uint32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 59))
	__NN32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 59+4))
	if uint32(__NN32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN32Size, __NN32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NN32Size, __NN32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 59+4)) = uint32(__NN32Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NN64(reader *karmem.Reader) (v []uint64) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 67))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 67+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint64)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNN64(v []uint64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 67))
	__NN64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 67+4))
	if uint32(__NN64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN64Size, __NN64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NN64Size, __NN64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 67+4)) = uint32(__NN64Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NM8(reader *karmem.Reader) (v []int8) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 75))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 75+4))
	if !reader.IsValidOffset(offset, size) {
		return []int8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int8)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNM8(v []int8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 75))
	__NM8Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 75+4))
	if uint32(__NM8Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM8Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM8Size, __NM8Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NM8Size, __NM8Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM8Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 75+4)) = uint32(__NM8Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NM16(reader *karmem.Reader) (v []int16) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 83))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 83+4))
	if !reader.IsValidOffset(offset, size) {
		return []int16{}
	}
	length := uintptr(size / 2)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int16)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNM16(v []int16) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 83))
	__NM16Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 83+4))
	if uint32(__NM16Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM16Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM16Size, __NM16Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NM16Size, __NM16Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM16Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 83+4)) = uint32(__NM16Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NM32(reader *karmem.Reader) (v []int32) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 91))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 91+4))
	if !reader.IsValidOffset(offset, size) {
		return []int32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int32)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNM32(v []int32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 91))
	__NM32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 91+4))
	if uint32(__NM32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM32Size, __NM32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NM32Size, __NM32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 91+4)) = uint32(__NM32Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NM64(reader *karmem.Reader) (v []int64) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 99))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 99+4))
	if !reader.IsValidOffset(offset, size) {
		return []int64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNM64(v []int64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 99))
	__NM64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 99+4))
	if uint32(__NM64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM64Size, __NM64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NM64Size, __NM64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 99+4)) = uint32(__NM64Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NOF32(reader *karmem.Reader) (v []float32) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 107))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 107+4))
	if !reader.IsValidOffset(offset, size) {
		return []float32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]float32)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNOF32(v []float32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 107))
	__NOF32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 107+4))
	if uint32(__NOF32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NOF32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NOF32Size, __NOF32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NOF32Size, __NOF32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NOF32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 107+4)) = uint32(__NOF32Size)
	return nil
}

func (x *SimpleNumbersPackedViewer) NOF64(reader *karmem.Reader) (v []float64) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 115))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 115+4))
	if !reader.IsValidOffset(offset, size) {
		return []float64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]float64)(unsafe.Pointer(&slice))
}

func (x *SimpleNumbersPackedSource) SetNOF64(v []float64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 115))
	__NOF64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 115+4))
	if uint32(__NOF64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NOF64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NOF64Size, __NOF64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), offset)), __NOF64Size, __NOF64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NOF64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.SimpleNumbersPackedViewer), 115+4)) = uint32(__NOF64Size)
	return nil
}

type ComplexPackedViewer [18859]byte

type ComplexPackedSource struct {
	*ComplexPackedViewer
}

func NewComplexPackedViewer(reader *karmem.Reader, offset uint32) (v *ComplexPackedViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*ComplexPackedViewer)(unsafe.Pointer(&_Null_demo[0]))
	}
	v = (*ComplexPackedViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*ComplexPackedViewer)(unsafe.Pointer(&_Null_demo[0]))
	}
	return v
}

func NewComplexPackedSource(reader *karmem.Reader, offset uint32) (v ComplexPackedSource) {
	return ComplexPackedSource{NewComplexPackedViewer(reader, offset)}
}

func newComplexPackedSourceByViewer(viewer *ComplexPackedViewer) (v ComplexPackedSource) {
	return ComplexPackedSource{viewer}
}

func (x *ComplexPackedViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *ComplexPackedViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *ComplexPackedViewer) N8() (v uint8) {
	if 4+1 > x.size() {
		return v
	}
	return *(*uint8)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *ComplexPackedSource) SetN8(v uint8) {
	*(*uint8)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 4)) = v
}

func (x *ComplexPackedViewer) N16() (v uint16) {
	if 5+2 > x.size() {
		return v
	}
	return *(*uint16)(unsafe.Add(unsafe.Pointer(x), 5))
}

func (x *ComplexPackedSource) SetN16(v uint16) {
	*(*uint16)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 5)) = v
}

func (x *ComplexPackedViewer) N32() (v uint32) {
	if 7+4 > x.size() {
		return v
	}
	return *(*uint32)(unsafe.Add(unsafe.Pointer(x), 7))
}

func (x *ComplexPackedSource) SetN32(v uint32) {
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 7)) = v
}

func (x *ComplexPackedViewer) N64() (v uint64) {
	if 11+8 > x.size() {
		return v
	}
	return *(*uint64)(unsafe.Add(unsafe.Pointer(x), 11))
}

func (x *ComplexPackedSource) SetN64(v uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 11)) = v
}

func (x *ComplexPackedViewer) M8() (v int8) {
	if 19+1 > x.size() {
		return v
	}
	return *(*int8)(unsafe.Add(unsafe.Pointer(x), 19))
}

func (x *ComplexPackedSource) SetM8(v int8) {
	*(*int8)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 19)) = v
}

func (x *ComplexPackedViewer) M16() (v int16) {
	if 20+2 > x.size() {
		return v
	}
	return *(*int16)(unsafe.Add(unsafe.Pointer(x), 20))
}

func (x *ComplexPackedSource) SetM16(v int16) {
	*(*int16)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 20)) = v
}

func (x *ComplexPackedViewer) M32() (v int32) {
	if 22+4 > x.size() {
		return v
	}
	return *(*int32)(unsafe.Add(unsafe.Pointer(x), 22))
}

func (x *ComplexPackedSource) SetM32(v int32) {
	*(*int32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 22)) = v
}

func (x *ComplexPackedViewer) M64() (v int64) {
	if 26+8 > x.size() {
		return v
	}
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 26))
}

func (x *ComplexPackedSource) SetM64(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 26)) = v
}

func (x *ComplexPackedViewer) OF32() (v float32) {
	if 34+4 > x.size() {
		return v
	}
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 34))
}

func (x *ComplexPackedSource) SetOF32(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 34)) = v
}

func (x *ComplexPackedViewer) OF64() (v float64) {
	if 38+8 > x.size() {
		return v
	}
	return *(*float64)(unsafe.Add(unsafe.Pointer(x), 38))
}

func (x *ComplexPackedSource) SetOF64(v float64) {
	*(*float64)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 38)) = v
}

func (x *ComplexPackedViewer) B1() (v bool) {
	if 46+1 > x.size() {
		return v
	}
	return *(*bool)(unsafe.Add(unsafe.Pointer(x), 46))
}

func (x *ComplexPackedSource) SetB1(v bool) {
	*(*bool)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 46)) = v
}

func (x *ComplexPackedViewer) NN8(reader *karmem.Reader) (v []uint8) {
	if 47+8 > x.size() {
		return []uint8{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 47))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 47+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint8)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNN8(v []uint8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 47))
	__NN8Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 47+4))
	if uint32(__NN8Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN8Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN8Size, __NN8Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NN8Size, __NN8Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN8Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 47+4)) = uint32(__NN8Size)
	return nil
}

func (x *ComplexPackedViewer) NN16(reader *karmem.Reader) (v []uint16) {
	if 55+8 > x.size() {
		return []uint16{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 55))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 55+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint16{}
	}
	length := uintptr(size / 2)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint16)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNN16(v []uint16) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 55))
	__NN16Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 55+4))
	if uint32(__NN16Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN16Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN16Size, __NN16Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NN16Size, __NN16Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN16Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 55+4)) = uint32(__NN16Size)
	return nil
}

func (x *ComplexPackedViewer) NN32(reader *karmem.Reader) (v []uint32) {
	if 63+8 > x.size() {
		return []uint32{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 63))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 63+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint32)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNN32(v []uint32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 63))
	__NN32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 63+4))
	if uint32(__NN32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN32Size, __NN32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NN32Size, __NN32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 63+4)) = uint32(__NN32Size)
	return nil
}

func (x *ComplexPackedViewer) NN64(reader *karmem.Reader) (v []uint64) {
	if 71+8 > x.size() {
		return []uint64{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 71))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 71+4))
	if !reader.IsValidOffset(offset, size) {
		return []uint64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]uint64)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNN64(v []uint64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 71))
	__NN64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 71+4))
	if uint32(__NN64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NN64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NN64Size, __NN64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NN64Size, __NN64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NN64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 71+4)) = uint32(__NN64Size)
	return nil
}

func (x *ComplexPackedViewer) NM8(reader *karmem.Reader) (v []int8) {
	if 79+8 > x.size() {
		return []int8{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 79))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 79+4))
	if !reader.IsValidOffset(offset, size) {
		return []int8{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int8)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNM8(v []int8) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 79))
	__NM8Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 79+4))
	if uint32(__NM8Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM8Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM8Size, __NM8Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NM8Size, __NM8Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM8Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 79+4)) = uint32(__NM8Size)
	return nil
}

func (x *ComplexPackedViewer) NM16(reader *karmem.Reader) (v []int16) {
	if 87+8 > x.size() {
		return []int16{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 87))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 87+4))
	if !reader.IsValidOffset(offset, size) {
		return []int16{}
	}
	length := uintptr(size / 2)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int16)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNM16(v []int16) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 87))
	__NM16Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 87+4))
	if uint32(__NM16Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM16Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM16Size, __NM16Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NM16Size, __NM16Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM16Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 87+4)) = uint32(__NM16Size)
	return nil
}

func (x *ComplexPackedViewer) NM32(reader *karmem.Reader) (v []int32) {
	if 95+8 > x.size() {
		return []int32{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 95))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 95+4))
	if !reader.IsValidOffset(offset, size) {
		return []int32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int32)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNM32(v []int32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 95))
	__NM32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 95+4))
	if uint32(__NM32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM32Size, __NM32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NM32Size, __NM32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 95+4)) = uint32(__NM32Size)
	return nil
}

func (x *ComplexPackedViewer) NM64(reader *karmem.Reader) (v []int64) {
	if 103+8 > x.size() {
		return []int64{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 103))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 103+4))
	if !reader.IsValidOffset(offset, size) {
		return []int64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNM64(v []int64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 103))
	__NM64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 103+4))
	if uint32(__NM64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NM64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NM64Size, __NM64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NM64Size, __NM64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NM64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 103+4)) = uint32(__NM64Size)
	return nil
}

func (x *ComplexPackedViewer) NOF32(reader *karmem.Reader) (v []float32) {
	if 111+8 > x.size() {
		return []float32{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 111))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 111+4))
	if !reader.IsValidOffset(offset, size) {
		return []float32{}
	}
	length := uintptr(size / 4)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]float32)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNOF32(v []float32) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 111))
	__NOF32Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 111+4))
	if uint32(__NOF32Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NOF32Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NOF32Size, __NOF32Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NOF32Size, __NOF32Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NOF32Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 111+4)) = uint32(__NOF32Size)
	return nil
}

func (x *ComplexPackedViewer) NOF64(reader *karmem.Reader) (v []float64) {
	if 119+8 > x.size() {
		return []float64{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 119))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 119+4))
	if !reader.IsValidOffset(offset, size) {
		return []float64{}
	}
	length := uintptr(size / 8)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]float64)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNOF64(v []float64) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 119))
	__NOF64Size := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 119+4))
	if uint32(__NOF64Size) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NOF64Slice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NOF64Size, __NOF64Size}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NOF64Size, __NOF64Size}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NOF64Slice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 119+4)) = uint32(__NOF64Size)
	return nil
}

func (x *ComplexPackedViewer) NB(reader *karmem.Reader) (v []bool) {
	if 127+8 > x.size() {
		return []bool{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 127))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 127+4))
	if !reader.IsValidOffset(offset, size) {
		return []bool{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]bool)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNB(v []bool) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 127))
	__NBSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 127+4))
	if uint32(__NBSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NBSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NBSize, __NBSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NBSize, __NBSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NBSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 127+4)) = uint32(__NBSize)
	return nil
}

func (x *ComplexPackedViewer) NC(reader *karmem.Reader) (v string) {
	if 135+8 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 135))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 135+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNC(v string) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 135))
	__NCSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 135+4))
	if uint32(__NCSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NCSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NCSize, __NCSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NCSize, __NCSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NCSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 135+4)) = uint32(__NCSize)
	return nil
}

func (x *ComplexPackedViewer) NIL(reader *karmem.Reader) (v []SimpleNumbersViewer) {
	if 147+8 > x.size() {
		return []SimpleNumbersViewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 147))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 147+4))
	if !reader.IsValidOffset(offset, size) {
		return []SimpleNumbersViewer{}
	}
	length := uintptr(size / 123)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]SimpleNumbersViewer)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNIL(v []SimpleNumbers) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 147))
	__NILSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 147+4))
	if uint32(__NILSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NILSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NILSize, __NILSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NILSize, __NILSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NILSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 147+4)) = uint32(__NILSize)
	return nil
}

func (x *ComplexPackedViewer) NIP(reader *karmem.Reader) (v []SimpleNumbersPackedViewer) {
	if 155+8 > x.size() {
		return []SimpleNumbersPackedViewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 155))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 155+4))
	if !reader.IsValidOffset(offset, size) {
		return []SimpleNumbersPackedViewer{}
	}
	length := uintptr(size / 123)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]SimpleNumbersPackedViewer)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetNIP(v []SimpleNumbersPacked) error {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 155))
	__NIPSize := uintptr(len(v))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 155+4))
	if uint32(__NIPSize) > size {
		return errors.New("invalid size, new size must less than old")
	}
	__NIPSlice := [3]uintptr{*(*uintptr)(unsafe.Pointer(&v)), __NIPSize, __NIPSize}
	__d := [3]uintptr{uintptr(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), offset)), __NIPSize, __NIPSize}
	copy(*(*[]byte)(unsafe.Pointer(&__d)), *(*[]byte)(unsafe.Pointer(&__NIPSlice)))
	*(*uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 155+4)) = uint32(__NIPSize)
	return nil
}

func (x *ComplexPackedViewer) ANN8() (v []uint8) {
	if 163+64 > x.size() {
		return []uint8{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 163)), 64, 64,
	}
	return *(*[]uint8)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANN8(v [64]uint8) {
	copy((*(*[64]uint8)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 163)))[:], v[:])
}

func (x *ComplexPackedViewer) ANN16() (v []uint16) {
	if 227+128 > x.size() {
		return []uint16{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 227)), 64, 64,
	}
	return *(*[]uint16)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANN16(v [64]uint16) {
	copy((*(*[64]uint16)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 227)))[:], v[:])
}

func (x *ComplexPackedViewer) ANN32() (v []uint32) {
	if 355+256 > x.size() {
		return []uint32{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 355)), 64, 64,
	}
	return *(*[]uint32)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANN32(v [64]uint32) {
	copy((*(*[64]uint32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 355)))[:], v[:])
}

func (x *ComplexPackedViewer) ANN64() (v []uint64) {
	if 611+512 > x.size() {
		return []uint64{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 611)), 64, 64,
	}
	return *(*[]uint64)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANN64(v [64]uint64) {
	copy((*(*[64]uint64)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 611)))[:], v[:])
}

func (x *ComplexPackedViewer) ANM8() (v []int8) {
	if 1123+64 > x.size() {
		return []int8{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 1123)), 64, 64,
	}
	return *(*[]int8)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANM8(v [64]int8) {
	copy((*(*[64]int8)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 1123)))[:], v[:])
}

func (x *ComplexPackedViewer) ANM16() (v []int16) {
	if 1187+128 > x.size() {
		return []int16{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 1187)), 64, 64,
	}
	return *(*[]int16)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANM16(v [64]int16) {
	copy((*(*[64]int16)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 1187)))[:], v[:])
}

func (x *ComplexPackedViewer) ANM32() (v []int32) {
	if 1315+256 > x.size() {
		return []int32{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 1315)), 64, 64,
	}
	return *(*[]int32)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANM32(v [64]int32) {
	copy((*(*[64]int32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 1315)))[:], v[:])
}

func (x *ComplexPackedViewer) ANM64() (v []int64) {
	if 1571+512 > x.size() {
		return []int64{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 1571)), 64, 64,
	}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANM64(v [64]int64) {
	copy((*(*[64]int64)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 1571)))[:], v[:])
}

func (x *ComplexPackedViewer) ANOF32() (v []float32) {
	if 2083+256 > x.size() {
		return []float32{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 2083)), 64, 64,
	}
	return *(*[]float32)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANOF32(v [64]float32) {
	copy((*(*[64]float32)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 2083)))[:], v[:])
}

func (x *ComplexPackedViewer) ANOF64() (v []float64) {
	if 2339+512 > x.size() {
		return []float64{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 2339)), 64, 64,
	}
	return *(*[]float64)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANOF64(v [64]float64) {
	copy((*(*[64]float64)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 2339)))[:], v[:])
}

func (x *ComplexPackedViewer) ANB() (v []bool) {
	if 2851+64 > x.size() {
		return []bool{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 2851)), 64, 64,
	}
	return *(*[]bool)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANB(v [64]bool) {
	copy((*(*[64]bool)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 2851)))[:], v[:])
}

func (x *ComplexPackedViewer) ANC() (v string) {
	if 2915+64 > x.size() {
		return v
	}
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 2915))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+2915)), uintptr(size), 64,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANC(v string) {
	*(*uint32)(unsafe.Pointer(x.ComplexPackedViewer)) = uint32(len(v))
	copy((*[64]byte)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 4+2915))[:], v)
}

func (x *ComplexPackedViewer) ANC1() (v string) {
	if 2983+128 > x.size() {
		return v
	}
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 2983))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+2983)), uintptr(size), 128,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANC1(v string) {
	*(*uint32)(unsafe.Pointer(x.ComplexPackedViewer)) = uint32(len(v))
	copy((*[128]byte)(unsafe.Add(unsafe.Pointer(x.ComplexPackedViewer), 4+2983))[:], v)
}

func (x *ComplexPackedViewer) ANIL() (v []SimpleNumbersViewer) {
	if 3115+7872 > x.size() {
		return []SimpleNumbersViewer{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 3115)), 64, 64,
	}
	return *(*[]SimpleNumbersViewer)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANIL(v [64]SimpleNumbers) {
	writer := karmem.NewFixedWriter(x.ComplexPackedViewer[:])
	__ANILOffset := uint(3115)
	for i := range v {
		_, _ = v[i].Write(writer, __ANILOffset)
		__ANILOffset += 123
	}
}

func (x *ComplexPackedSource) ANILSource(index int) (SimpleNumbersSource, error) {
	v := x.ANIL()
	if index > len(v)-1 {
		return SimpleNumbersSource{}, errors.New("out of array bounds")
	}
	return newSimpleNumbersSourceByViewer(&v[index]), nil
}

func (x *ComplexPackedViewer) ANIP() (v []SimpleNumbersPackedViewer) {
	if 10987+7872 > x.size() {
		return []SimpleNumbersPackedViewer{}
	}
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 10987)), 64, 64,
	}
	return *(*[]SimpleNumbersPackedViewer)(unsafe.Pointer(&slice))
}

func (x *ComplexPackedSource) SetANIP(v [64]SimpleNumbersPacked) {
	writer := karmem.NewFixedWriter(x.ComplexPackedViewer[:])
	__ANIPOffset := uint(10987)
	for i := range v {
		_, _ = v[i].Write(writer, __ANIPOffset)
		__ANIPOffset += 123
	}
}

func (x *ComplexPackedSource) ANIPSource(index int) (SimpleNumbersPackedSource, error) {
	v := x.ANIP()
	if index > len(v)-1 {
		return SimpleNumbersPackedSource{}, errors.New("out of array bounds")
	}
	return newSimpleNumbersPackedSourceByViewer(&v[index]), nil
}
