package ecs

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_ecs = [40]byte{}
var _NullReader_ecs = karmem.NewReader(_Null_ecs[:])

func init() {

}

const (
	PacketIdentifierSerializableUSetData = 14494556342195307757
)

type SerializableUSetData struct {
	EleSize  uint64
	Len      int64
	InitSize int64
	Data     []byte
}

func NewSerializableUSetData() SerializableUSetData {
	return SerializableUSetData{}
}

func (x *SerializableUSetData) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierSerializableUSetData
}

func (x *SerializableUSetData) Reset() {
	x.Read((*SerializableUSetDataViewer)(unsafe.Pointer(&_Null_ecs[0])), _NullReader_ecs)
}

func (x *SerializableUSetData) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *SerializableUSetData) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(40)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__EleSizeOffset := offset + 0
	writer.Write8At(__EleSizeOffset, *(*uint64)(unsafe.Pointer(&x.EleSize)))
	__LenOffset := offset + 8
	writer.Write8At(__LenOffset, *(*uint64)(unsafe.Pointer(&x.Len)))
	__InitSizeOffset := offset + 16
	writer.Write8At(__InitSizeOffset, *(*uint64)(unsafe.Pointer(&x.InitSize)))
	__DataSize := uint(1 * len(x.Data))
	__DataOffset, err := writer.Alloc(__DataSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+24, uint32(__DataOffset))
	writer.Write4At(offset+24+4, uint32(__DataSize))
	writer.Write4At(offset+24+4+4, 1)
	__DataSlice := *(*[3]uint)(unsafe.Pointer(&x.Data))
	__DataSlice[1] = __DataSize
	__DataSlice[2] = __DataSize
	writer.WriteAt(__DataOffset, *(*[]byte)(unsafe.Pointer(&__DataSlice)))

	return offset, nil
}

func (*SerializableUSetData) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(40)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_ecs[:size-4])
	__EleSizeOffset := offset + 0
	_ = __EleSizeOffset
	__LenOffset := offset + 8
	_ = __LenOffset
	__InitSizeOffset := offset + 16
	_ = __InitSizeOffset
	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+24, uint32(__DataOffset))
	writer.Write4At(offset+24+4, 0)
	writer.Write4At(offset+24+4+4, 1)

	return offset, nil
}

func (x *SerializableUSetData) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewSerializableUSetDataViewer(reader, 0), reader)
}

func (x *SerializableUSetData) Read(viewer *SerializableUSetDataViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.EleSize = viewer.EleSize()
	x.Len = viewer.Len()
	x.InitSize = viewer.InitSize()
	__DataSlice := viewer.Data(reader)
	__DataLen := len(__DataSlice)
	if __DataLen > cap(x.Data) {
		x.Data = append(x.Data, make([]byte, __DataLen-len(x.Data))...)
	}
	if __DataLen > len(x.Data) {
		x.Data = x.Data[:__DataLen]
	}
	copy(x.Data, __DataSlice)
	x.Data = x.Data[:__DataLen]
}

type SerializableUSetDataViewer [40]byte

func NewSerializableUSetDataViewer(reader *karmem.Reader, offset uint32) (v *SerializableUSetDataViewer) {
	if !reader.IsValidOffset(offset, 36) {
		return (*SerializableUSetDataViewer)(unsafe.Pointer(&_Null_ecs[0]))
	}
	v = (*SerializableUSetDataViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func (x *SerializableUSetDataViewer) size() uint32 {
	return 40
}

func (x *SerializableUSetDataViewer) EleSize() (v uint64) {
	return *(*uint64)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *SerializableUSetDataViewer) Len() (v int64) {
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 8))
}

func (x *SerializableUSetDataViewer) InitSize() (v int64) {
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 16))
}

func (x *SerializableUSetDataViewer) Data(reader *karmem.Reader) (v []byte) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 24))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 24+4))
	if !reader.IsValidOffset(offset, size) {
		return []byte{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]byte)(unsafe.Pointer(&slice))
}
