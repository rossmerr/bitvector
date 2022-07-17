package bitvector

// The BitVector manages a compact array of bit values, which are represented as bytes.
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

func (v *BitVector) Data() []byte {
	return v.data
}

func (v *BitVector) Length() int {
	return v.length
}

// Performs the bitwise AND operation on the elements in the current BitVector against the corresponding elements in the specified BitVector.
func (v *BitVector) And(value *BitVector) *BitVector {
	return v
}

// Gets the value of the bit at a specific position in the BitVector.
func (v *BitVector) Get(index int) bool {
	return false
}

// Inverts all the bit values in the current BitVector, so that elements set to true are changed to false, and elements set to false are changed to true.
func (v *BitVector) Not() *BitVector {
	return v
}

// Performs the bitwise OR operation on the elements in the current BitVector against the corresponding elements in the specified BitVector.
func (v *BitVector) Or(value *BitVector) *BitVector {
	return v
}

// Sets the bit at a specific position in the BitVector to the specified value.
func (v *BitVector) Set(index int, value bool) {

}

// Sets all bits in the BitVector to the specified value.
func (v *BitVector) SetAll(value bool) {

}

// Performs the bitwise eXclusive OR operation on the elements in the current BitVector against the corresponding elements in the specified BitVector.
func (v *BitVector) Xor(value *BitVector) *BitVector {
	return v

}
