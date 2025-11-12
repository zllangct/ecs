package testdata

import (
	"errors"
	ecs "github.com/zllangct/ecs"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_components = [64]byte{}
var _NullReader_components = karmem.NewReader(_Null_components[:])

func init() {

}

const (
	PacketIdentifierPoint    = 12809409023221035873
	PacketIdentifierPosition = 12809409023221035873
	PacketIdentifierName     = 7328248172266787231
)

type Point struct {
	X float32
	Y float32
	Z float32
}

func NewPoint() Point {
	return Point{}
}

func (x *Point) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierPoint
}

func (x *Point) Reset() {
	x.Read((*PointViewer)(unsafe.Pointer(&_Null_components[0])), _NullReader_components)
}

func (x *Point) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Point) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(16)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__XOffset := offset + 0
	writer.Write4At(__XOffset, *(*uint32)(unsafe.Pointer(&x.X)))
	__YOffset := offset + 4
	writer.Write4At(__YOffset, *(*uint32)(unsafe.Pointer(&x.Y)))
	__ZOffset := offset + 8
	writer.Write4At(__ZOffset, *(*uint32)(unsafe.Pointer(&x.Z)))

	return offset, nil
}

func (*Point) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(16)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_components[:size-4])
	__XOffset := offset + 0
	_ = __XOffset
	__YOffset := offset + 4
	_ = __YOffset
	__ZOffset := offset + 8
	_ = __ZOffset

	return offset, nil
}

func (x *Point) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewPointViewer(reader, 0), reader)
}

func (x *Point) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewPointViewer(reader, offset), reader)
}

func (x *Point) Read(viewer *PointViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.X = viewer.X()
	x.Y = viewer.Y()
	x.Z = viewer.Z()
}

//component extension

func (x Point) ComponentObjectIdentifier() {}

func (x *Point) GetComponentSeq() int32 {
	return 1
}

func (x *Point) NewComponentSet() ecs.ComponentSet {
	return ecs.NewCSet[Point]()
}

func (x *Point) IsNomadic() bool {
	return false
}

func (x *Point) IsDisposable() bool {
	return false
}

type Position struct {
	X float32
	Y float32
	Z float32
}

func NewPosition() Position {
	return Position{}
}

func (x *Position) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierPosition
}

func (x *Position) Reset() {
	x.Read((*PositionViewer)(unsafe.Pointer(&_Null_components[0])), _NullReader_components)
}

func (x *Position) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Position) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(16)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__XOffset := offset + 0
	writer.Write4At(__XOffset, *(*uint32)(unsafe.Pointer(&x.X)))
	__YOffset := offset + 4
	writer.Write4At(__YOffset, *(*uint32)(unsafe.Pointer(&x.Y)))
	__ZOffset := offset + 8
	writer.Write4At(__ZOffset, *(*uint32)(unsafe.Pointer(&x.Z)))

	return offset, nil
}

func (*Position) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(16)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_components[:size-4])
	__XOffset := offset + 0
	_ = __XOffset
	__YOffset := offset + 4
	_ = __YOffset
	__ZOffset := offset + 8
	_ = __ZOffset

	return offset, nil
}

func (x *Position) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewPositionViewer(reader, 0), reader)
}

func (x *Position) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewPositionViewer(reader, offset), reader)
}

func (x *Position) Read(viewer *PositionViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.X = viewer.X()
	x.Y = viewer.Y()
	x.Z = viewer.Z()
}

//component extension

func (x Position) ComponentObjectIdentifier() {}

func (x *Position) GetComponentSeq() int32 {
	return 2
}

func (x *Position) NewComponentSet() ecs.ComponentSet {
	return ecs.NewCSet[Position]()
}

func (x *Position) IsNomadic() bool {
	return false
}

func (x *Position) IsDisposable() bool {
	return false
}

type Name struct {
	// fixed string of size 16
	Value  ecs.Fixed16
	Points [2]Point
	Arr    [2]int32
}

func NewName() Name {
	return Name{}
}

func (x *Name) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierName
}

func (x *Name) Reset() {
	x.Read((*NameViewer)(unsafe.Pointer(&_Null_components[0])), _NullReader_components)
}

func (x *Name) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Name) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(64)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__ValueOffset := offset + 0
	__ValueSize := uint(1 * x.Value.Len())
	writer.Write4At(__ValueOffset, uint32(__ValueSize))
	if __ValueSize > 0 {
		if __ValueSize > 16 {
			return 0, errors.New("out of array bounds")
		}
		writer.WriteAt(__ValueOffset+4, (*[16]byte)(unsafe.Pointer(&x.Value))[:])
	} else {
		writer.WriteAt(__ValueOffset+4, _Null_components[:16])
	}
	__PointsOffset := offset + 20
	for i := range x.Points {
		if _, err := x.Points[i].Write(writer, __PointsOffset); err != nil {
			return offset, err
		}
		__PointsOffset += 16
	}
	__ArrOffset := offset + 52
	writer.WriteAt(__ArrOffset, (*[8]byte)(unsafe.Pointer(&x.Arr))[:])

	return offset, nil
}

func (*Name) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(64)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_components[:size-4])
	__ValueOffset := offset + 0
	_ = __ValueOffset
	writer.Write4At(__ValueOffset, 0)
	__PointsOffset := offset + 20
	_ = __PointsOffset
	__ArrOffset := offset + 52
	_ = __ArrOffset

	return offset, nil
}

func (x *Name) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewNameViewer(reader, 0), reader)
}

func (x *Name) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewNameViewer(reader, offset), reader)
}

func (x *Name) Read(viewer *NameViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__ValueString := viewer.Value()
	if x.Value.String() != __ValueString {
		__ValueStringCopy := make([]byte, len(__ValueString))
		copy(__ValueStringCopy, __ValueString)
		x.Value.Set(*(*string)(unsafe.Pointer(&__ValueStringCopy)))
	}
	__PointsSlice := viewer.Points()
	__PointsLen := len(__PointsSlice)
	if __PointsLen > 2 {
		__PointsLen = 2
	}
	for i := 0; i < __PointsLen; i++ {
		x.Points[i].Read(&__PointsSlice[i], reader)
	}
	for i := __PointsLen; i < len(x.Points); i++ {
		x.Points[i].Reset()
	}
	__ArrSlice := viewer.Arr()
	__ArrLen := len(__ArrSlice)
	if __ArrLen > 2 {
		__ArrLen = 2
	}
	copy(x.Arr[:], __ArrSlice)
	for i := __ArrLen; i < len(x.Arr); i++ {
		x.Arr[i] = 0
	}
}

//component extension

func (x Name) ComponentObjectIdentifier() {}

func (x *Name) GetComponentSeq() int32 {
	return 3
}

func (x *Name) NewComponentSet() ecs.ComponentSet {
	return ecs.NewCSet[Name]()
}

func (x *Name) IsNomadic() bool {
	return false
}

func (x *Name) IsDisposable() bool {
	return false
}

type PointViewer [16]byte

type PointSource struct {
	*PointViewer
}

func NewPointViewer(reader *karmem.Reader, offset uint32) (v *PointViewer) {
	if !reader.IsValidOffset(offset, 12) {
		return (*PointViewer)(unsafe.Pointer(&_Null_components[0]))
	}
	v = (*PointViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewPointSource(reader *karmem.Reader, offset uint32) (v PointSource) {
	return PointSource{NewPointViewer(reader, offset)}
}

func newPointSourceByViewer(viewer *PointViewer) (v PointSource) {
	return PointSource{viewer}
}

func (x *PointViewer) size() uint32 {
	return 16
}

func (x *PointViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *PointViewer) X() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *PointSource) SetX(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.PointViewer), 0)) = v
}

func (x *PointViewer) Y() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *PointSource) SetY(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.PointViewer), 4)) = v
}

func (x *PointViewer) Z() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 8))
}

func (x *PointSource) SetZ(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.PointViewer), 8)) = v
}

type PositionViewer [16]byte

type PositionSource struct {
	*PositionViewer
}

func NewPositionViewer(reader *karmem.Reader, offset uint32) (v *PositionViewer) {
	if !reader.IsValidOffset(offset, 12) {
		return (*PositionViewer)(unsafe.Pointer(&_Null_components[0]))
	}
	v = (*PositionViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewPositionSource(reader *karmem.Reader, offset uint32) (v PositionSource) {
	return PositionSource{NewPositionViewer(reader, offset)}
}

func newPositionSourceByViewer(viewer *PositionViewer) (v PositionSource) {
	return PositionSource{viewer}
}

func (x *PositionViewer) size() uint32 {
	return 16
}

func (x *PositionViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *PositionViewer) X() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *PositionSource) SetX(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.PositionViewer), 0)) = v
}

func (x *PositionViewer) Y() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *PositionSource) SetY(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.PositionViewer), 4)) = v
}

func (x *PositionViewer) Z() (v float32) {
	return *(*float32)(unsafe.Add(unsafe.Pointer(x), 8))
}

func (x *PositionSource) SetZ(v float32) {
	*(*float32)(unsafe.Add(unsafe.Pointer(x.PositionViewer), 8)) = v
}

type NameViewer [64]byte

type NameSource struct {
	*NameViewer
}

func NewNameViewer(reader *karmem.Reader, offset uint32) (v *NameViewer) {
	if !reader.IsValidOffset(offset, 60) {
		return (*NameViewer)(unsafe.Pointer(&_Null_components[0]))
	}
	v = (*NameViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewNameSource(reader *karmem.Reader, offset uint32) (v NameSource) {
	return NameSource{NewNameViewer(reader, offset)}
}

func newNameSourceByViewer(viewer *NameViewer) (v NameSource) {
	return NameSource{viewer}
}

func (x *NameViewer) size() uint32 {
	return 64
}

func (x *NameViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *NameViewer) Value() (v string) {
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 0))
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 4+0)), uintptr(size), 16,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func (x *NameSource) SetValue(v string) {
	*(*uint32)(unsafe.Pointer(x.NameViewer)) = uint32(len(v))
	copy((*[16]byte)(unsafe.Add(unsafe.Pointer(x.NameViewer), 4+0))[:], v)
}

func (x *NameViewer) Points() (v []PointViewer) {
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 20)), 2, 2,
	}
	return *(*[]PointViewer)(unsafe.Pointer(&slice))
}

func (x *NameSource) SetPoints(v [2]Point) {
	writer := karmem.NewFixedWriter(x.NameViewer[:])
	__PointsOffset := uint(20)
	for i := range v {
		_, _ = v[i].Write(writer, __PointsOffset)
		__PointsOffset += 16
	}
}

func (x *NameSource) PointsSource(index int) (PointSource, error) {
	v := x.Points()
	if index > len(v)-1 {
		return PointSource{}, errors.New("out of array bounds")
	}
	return newPointSourceByViewer(&v[index]), nil
}

func (x *NameViewer) Arr() (v []int32) {
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 52)), 2, 2,
	}
	return *(*[]int32)(unsafe.Pointer(&slice))
}

func (x *NameSource) SetArr(v [2]int32) {
	copy((*(*[2]int32)(unsafe.Add(unsafe.Pointer(x.NameViewer), 52)))[:], v[:])
}
