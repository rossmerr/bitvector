package bitvector

type BitVector struct {
	data   []byte
	length int
}

func NewBitVector(data []byte) *BitVector {
	return &BitVector{
		data:   data,
		length: len(data),
	}
}

func (v *BitVector) Length() int {
	return v.length
}
