package karmemtest

import (
	"errors"
	karmem "github.com/zllangct/ecs/karmem"
	"unsafe"
)

var _ unsafe.Pointer
var _ = errors.New("")

var _Null_case_refone = [14]byte{}
var _NullReader_case_refone = karmem.NewReader(_Null_case_refone[:])

func init() {

}

const (
	PacketIdentifierBIR  = 626714986465469827
	PacketIdentifierBIRR = 9649917901262869840
	PacketIdentifierBII  = 4427636902920729320
	PacketIdentifierBIIR = 11194102834194897860
)

type BIR struct {
	Value int64
}

func NewBIR() BIR {
	return BIR{}
}

func (x *BIR) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBIR
}

func (x *BIR) Reset() {
	x.Read((*BIRViewer)(unsafe.Pointer(&_Null_case_refone[0])), _NullReader_case_refone)
}

func (x *BIR) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BIR) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(12)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(12))
	__ValueOffset := offset + 4
	writer.Write8At(__ValueOffset, *(*uint64)(unsafe.Pointer(&x.Value)))

	return offset, nil
}

func (*BIR) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(12)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(12))
	writer.WriteAt(offset, _Null_case_refone[:size-4])
	__ValueOffset := offset + 4
	_ = __ValueOffset

	return offset, nil
}

func (x *BIR) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBIRViewer(reader, 0), reader)
}

func (x *BIR) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBIRViewer(reader, offset), reader)
}

func (x *BIR) Read(viewer *BIRViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.Value = viewer.Value()
}

type BIRR struct {
	Data  *BIR
	Data2 *BIR
}

func NewBIRR() BIRR {
	return BIRR{}
}

func (x *BIRR) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBIRR
}

func (x *BIRR) Reset() {
	x.Read((*BIRRViewer)(unsafe.Pointer(&_Null_case_refone[0])), _NullReader_case_refone)
}

func (x *BIRR) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BIRR) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(14)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(14))
	__DataSize := uint(12)
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
	__Data2Size := uint(12)
	__Data2Offset, err := writer.Alloc(__Data2Size)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+9, uint32(__Data2Offset))
	if x.Data2 == nil {
		writer.Write1At(offset+9+4, 0)
	} else {
		writer.Write1At(offset+9+4, 1)
	}
	if x.Data2 == nil {
		if _, err := x.Data2.WriteDefault(writer, __Data2Offset); err != nil {
			return offset, err
		}
	} else {
		if _, err := x.Data2.Write(writer, __Data2Offset); err != nil {
			return offset, err
		}
	}

	return offset, nil
}

func (*BIRR) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(14)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(14))
	writer.WriteAt(offset, _Null_case_refone[:size-4])

	__DataOffset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__DataOffset))
	writer.Write1At(__DataOffset, 0)
	if _, err := (*BIR)(nil).WriteDefault(writer, __DataOffset+1); err != nil {
		return offset, err
	}

	__Data2Offset, err := writer.Alloc(0)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+9, uint32(__Data2Offset))
	writer.Write1At(__Data2Offset, 0)
	if _, err := (*BIR)(nil).WriteDefault(writer, __Data2Offset+1); err != nil {
		return offset, err
	}

	return offset, nil
}

func (x *BIRR) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBIRRViewer(reader, 0), reader)
}

func (x *BIRR) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBIRRViewer(reader, offset), reader)
}

func (x *BIRR) Read(viewer *BIRRViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	if !viewer.isDataNil() {
		DataViewer := viewer.Data(reader)
		if x.Data == nil {
			x.Data = new(BIR)
		}
		x.Data.Read(DataViewer, reader)
	}
	if !viewer.isData2Nil() {
		Data2Viewer := viewer.Data2(reader)
		if x.Data2 == nil {
			x.Data2 = new(BIR)
		}
		x.Data2.Read(Data2Viewer, reader)
	}
}

type BII struct {
	Value int64
}

func NewBII() BII {
	return BII{}
}

func (x *BII) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBII
}

func (x *BII) Reset() {
	x.Read((*BIIViewer)(unsafe.Pointer(&_Null_case_refone[0])), _NullReader_case_refone)
}

func (x *BII) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BII) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(8)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__ValueOffset := offset + 0
	writer.Write8At(__ValueOffset, *(*uint64)(unsafe.Pointer(&x.Value)))

	return offset, nil
}

func (*BII) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(8)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.WriteAt(offset, _Null_case_refone[:size-4])
	__ValueOffset := offset + 0
	_ = __ValueOffset

	return offset, nil
}

func (x *BII) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBIIViewer(reader, 0), reader)
}

func (x *BII) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBIIViewer(reader, offset), reader)
}

func (x *BII) Read(viewer *BIIViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	x.Value = viewer.Value()
}

type BIIR struct {
	Data *BII
}

func NewBIIR() BIIR {
	return BIIR{}
}

func (x *BIIR) PacketIdentifier() karmem.PacketIdentifier {
	return PacketIdentifierBIIR
}

func (x *BIIR) Reset() {
	x.Read((*BIIRViewer)(unsafe.Pointer(&_Null_case_refone[0])), _NullReader_case_refone)
}

func (x *BIIR) WriteAsRoot(writer karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *BIIR) Write(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(13)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	// struct size for table
	writer.Write4At(offset, uint32(13))
	__DataOffset := offset + 4
	if x.Data == nil {
		writer.Write1At(offset+4, 0)
	} else {
		writer.Write1At(offset+4, 1)
	}
	__DataOffset += 1
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

func (*BIIR) WriteDefault(writer karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(13)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(13))
	writer.WriteAt(offset, _Null_case_refone[:size-4])
	__DataOffset := offset + 4
	_ = __DataOffset
	writer.Write1At(__DataOffset, 0)
	if _, err := (*BII)(nil).WriteDefault(writer, __DataOffset+1); err != nil {
		return offset, err
	}

	return offset, nil
}

func (x *BIIR) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewBIIRViewer(reader, 0), reader)
}

func (x *BIIR) ReadWithOffset(reader *karmem.Reader, offset uint32) {
	x.Read(NewBIIRViewer(reader, offset), reader)
}

func (x *BIIR) Read(viewer *BIIRViewer, reader *karmem.Reader) {
	if viewer == nil {
		return
	}
	if !viewer.isDataNil() {
		DataViewer := viewer.Data()
		if x.Data == nil {
			x.Data = new(BII)
		}
		x.Data.Read(DataViewer, reader)
	}
}

type BIRViewer [12]byte

type BIRSource struct {
	*BIRViewer
}

func NewBIRViewer(reader *karmem.Reader, offset uint32) (v *BIRViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*BIRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	v = (*BIRViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BIRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	return v
}

func NewBIRSource(reader *karmem.Reader, offset uint32) (v BIRSource) {
	return BIRSource{NewBIRViewer(reader, offset)}
}

func newBIRSourceByViewer(viewer *BIRViewer) (v BIRSource) {
	return BIRSource{viewer}
}

func (x *BIRViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BIRViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BIRViewer) Value() (v int64) {
	if 4+8 > x.size() {
		return v
	}
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *BIRSource) SetValue(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.BIRViewer), 4)) = v
}

type BIRRViewer [14]byte

type BIRRSource struct {
	*BIRRViewer
}

func NewBIRRViewer(reader *karmem.Reader, offset uint32) (v *BIRRViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*BIRRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	v = (*BIRRViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BIRRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	return v
}

func NewBIRRSource(reader *karmem.Reader, offset uint32) (v BIRRSource) {
	return BIRRSource{NewBIRRViewer(reader, offset)}
}

func newBIRRSourceByViewer(viewer *BIRRViewer) (v BIRRSource) {
	return BIRRSource{viewer}
}

func (x *BIRRViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BIRRViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BIRRViewer) isDataNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 4+4))
}

func (x *BIRRViewer) Data(reader *karmem.Reader) (v *BIRViewer) {
	if 4+4 > x.size() {
		return (*BIRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 4))
	return NewBIRViewer(reader, offset)
}

func (x *BIRRViewer) isData2Nil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 9+4))
}

func (x *BIRRViewer) Data2(reader *karmem.Reader) (v *BIRViewer) {
	if 9+4 > x.size() {
		return (*BIRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(x), 9))
	return NewBIRViewer(reader, offset)
}

type BIIViewer [8]byte

type BIISource struct {
	*BIIViewer
}

func NewBIIViewer(reader *karmem.Reader, offset uint32) (v *BIIViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*BIIViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	v = (*BIIViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func NewBIISource(reader *karmem.Reader, offset uint32) (v BIISource) {
	return BIISource{NewBIIViewer(reader, offset)}
}

func newBIISourceByViewer(viewer *BIIViewer) (v BIISource) {
	return BIISource{viewer}
}

func (x *BIIViewer) size() uint32 {
	return 8
}

func (x *BIIViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BIIViewer) Value() (v int64) {
	return *(*int64)(unsafe.Add(unsafe.Pointer(x), 0))
}

func (x *BIISource) SetValue(v int64) {
	*(*int64)(unsafe.Add(unsafe.Pointer(x.BIIViewer), 0)) = v
}

type BIIRViewer [13]byte

type BIIRSource struct {
	*BIIRViewer
}

func NewBIIRViewer(reader *karmem.Reader, offset uint32) (v *BIIRViewer) {
	if !reader.IsValidOffset(offset, 4) {
		return (*BIIRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	v = (*BIIRViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*BIIRViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	return v
}

func NewBIIRSource(reader *karmem.Reader, offset uint32) (v BIIRSource) {
	return BIIRSource{NewBIIRViewer(reader, offset)}
}

func newBIIRSourceByViewer(viewer *BIIRViewer) (v BIIRSource) {
	return BIIRSource{viewer}
}

func (x *BIIRViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(x))
}

func (x *BIIRViewer) KarmemReader() *karmem.Reader {
	return karmem.NewReader(x[:])
}

func (x *BIIRViewer) isDataNil() bool {
	return !*(*bool)(unsafe.Add(unsafe.Pointer(x), 4))
}

func (x *BIIRViewer) Data() (v *BIIViewer) {
	if 4+8 > x.size() {
		return (*BIIViewer)(unsafe.Pointer(&_Null_case_refone[0]))
	}
	return (*BIIViewer)(unsafe.Add(unsafe.Pointer(x), 4+1))
}
