package karmemtest

import (
	"errors"
	ecs "github.com/zllangct/ecs"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_components = [36]byte{}
var _NullReader_components = karmem.NewReader(_Null_components[:])

func init() {

}

const (
	PacketIdentifierPoint     = 12809409023221035873
	PacketIdentifierContainer = 8006837439125283466
	PacketIdentifierPosition  = 14274740599321850262
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
	size := uint(12)
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
	size := uint(12)
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

type Container struct {
	Comp [3]Point
}

func NewContainer() Container {
	return Container{}
}

func (x *Container) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierContainer
}

func (x *Container) Reset() {
	x.Read((*ContainerViewer)(unsafe.Pointer(&_Null_components[0])), _NullReader_components)
}

func (x *Container) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Container) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(36)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__CompOffset := offset + 0
	for i := range x.Comp {
		if _, err := x.Comp[i].Write(writer, __CompOffset); err != nil {
			return offset, err
		}
		__CompOffset += 12
	}

	return offset, nil
}

func (*Container) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(36)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_components[:size-4])
	__CompOffset := offset + 0
	_ = __CompOffset

	return offset, nil
}

func (x *Container) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewContainerViewer(reader, 0), reader)
}

func (x *Container) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewContainerViewer(reader, offset), reader)
}

func (x *Container) Read(viewer *ContainerViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	__CompSlice := viewer.Comp()
	__CompLen := len(__CompSlice)
	if __CompLen > 3 {
		__CompLen = 3
	}
	for i := 0; i < __CompLen; i++ {
		x.Comp[i].Read(&__CompSlice[i], reader)
	}
	for i := __CompLen; i < len(x.Comp); i++ {
		x.Comp[i].Reset()
	}
}

//component extension

func (x Container) ComponentObjectIdentifier() {}

func (x *Container) GetComponentSeq() int32 {
	return 2
}

func (x *Container) NewComponentSet() ecs.ComponentSet {
	return ecs.NewCSet[Container]()
}

func (x *Container) IsNomadic() bool {
	return false
}

func (x *Container) IsDisposable() bool {
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
	size := uint(12)
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
	size := uint(12)
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
	return 3
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

type PointViewer [12]byte

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
	return 12
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

type ContainerViewer [36]byte

type ContainerSource struct {
	*ContainerViewer
}

func NewContainerViewer(reader *karmem.Reader, offset uint32) (v *ContainerViewer) {
	if !reader.IsValidOffset(offset, 36) {
		return (*ContainerViewer)(unsafe.Pointer(&_Null_components[0]))
	}
	v = (*ContainerViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewContainerSource(reader *karmem.Reader, offset uint32) (v ContainerSource) {
	return ContainerSource{NewContainerViewer(reader, offset)}
}

func newContainerSourceByViewer(viewer *ContainerViewer) (v ContainerSource) {
	return ContainerSource{viewer}
}

func (x *ContainerViewer) size() uint32 {
	return 36
}

func (x *ContainerViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *ContainerViewer) Comp() (v []PointViewer) {
	slice := [3]uintptr{
		uintptr(unsafe.Add(unsafe.Pointer(x), 0)), 3, 3,
	}
	return *(*[]PointViewer)(unsafe.Pointer(&slice))
}

func (x *ContainerSource) SetComp(v [3]Point) {
	writer := karmem.NewFixedWriter(x.ContainerViewer[:])
	__CompOffset := uint(0)
	for i := range v {
		_, _ = v[i].Write(writer, __CompOffset)
		__CompOffset += 12
	}
}

func (x *ContainerSource) CompSource(index int) (PointSource, error) {
	v := x.Comp()
	if index > len(v)-1 {
		return PointSource{}, errors.New("out of array bounds")
	}
	return newPointSourceByViewer(&v[index]), nil
}

type PositionViewer [12]byte

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
	return 12
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
