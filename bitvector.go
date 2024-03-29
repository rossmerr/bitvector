package bitvector

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

const bitsPerInt32 = 32

type BitVector struct {
	array   []uint32
	length  int
	version int
}

// Allocates space to hold the length of bit. All of the values in the BitVector are set to false.
func NewBitVector(length int) *BitVector {
	return NewBitVectorOfLength(length, false)
}

// Allocates space to hold the length of bit. All of the values in the BitVector are set to defaultBit.
func NewBitVectorOfLength(length int, defaultBit bool) *BitVector {
	arrayLength, err := getArrayLength(length, bitsPerInt32)
	if err != nil {
		panic(err)
	}
	array := make([]uint32, arrayLength)

	fillValue := uint32(0)
	if defaultBit {
		fillValue = 0xffffffff
	}

	for i := 0; i < arrayLength; i++ {
		array[i] = fillValue
	}

	return &BitVector{
		array:   array,
		length:  length,
		version: 0,
	}
}

// Allocates space to hold the values from the booleans.
func NewBitVectorFromBool(values []bool) *BitVector {
	arrayLength, err := getArrayLength(len(values), bitsPerInt32)
	if err != nil {
		panic(err)
	}
	array := make([]uint32, arrayLength)
	for i, value := range values {
		if value {
			array[i/bitsPerInt32] |= (1 << (i % bitsPerInt32))
		} else {
			array[i/bitsPerInt32] &= ^(1 << (i % bitsPerInt32))
		}
	}

	return &BitVector{
		array:   array,
		length:  len(values),
		version: 0,
	}
}

// Allocates a new BitVector with the same length and bit values as vector.
func NewBitVectorFromVector(vector BitVector) *BitVector {
	array := make([]uint32, len(vector.array))

	copy(array, vector.array)

	return &BitVector{
		array:   array,
		length:  vector.length,
		version: 0,
	}
}

// Allocates a new BitVector padded with the same length and values as the vector but left shifted by the padding.
func NewBitVectorFromVectorPadStart(vector *BitVector, padding int) *BitVector {
	length, err := getArrayLength(vector.Length()+padding, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	array := make([]uint32, length)

	index, err := getArrayLength(padding+1, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	index--

	offset := padding % bitsPerInt32

	arrayLength, err := getArrayLength(vector.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := index; i < length; i++ {
		for y := 0; y < arrayLength; y++ {
			array[i] = (vector.array[y] << offset)
		}
	}

	return &BitVector{
		array:   array,
		length:  vector.length + padding,
		version: 0,
	}
}

// Returns the bit value at position index.
func (s BitVector) Get(index int) bool {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}

	return (s.array[index/bitsPerInt32] & (1 << (index % bitsPerInt32))) != 0
}

// Sets the bit value at position index to value.
func (s BitVector) Set(index int, bit bool) {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}

	if bit {
		s.array[index/bitsPerInt32] |= (1 << (index % bitsPerInt32))
	} else {
		s.array[index/bitsPerInt32] &= ^(1 << (index % bitsPerInt32))
	}

	s.version++
}

// Sets all the bit values to value.
func (s *BitVector) SetAll(bit bool) {
	fillValue := uint32(0)
	if bit {
		fillValue = 0xffffffff
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] = fillValue
	}

	s.version++
}

func (s *BitVector) Copy(vector *BitVector, indexStart, indexEnd int) {
	if indexStart < 0 {
		panic("indexStart must be non negative number")
	}

	if indexEnd > s.Length() {
		panic("indexEnd must be equal to or less than bitvector")
	}

	if vector.Length() < (s.Length()-indexStart)-indexEnd {
		panic("invalid vector length is to small")
	}

	var err error

	arrayEnd := 0
	if indexEnd > 0 {
		arrayEnd, err = getArrayLength(indexEnd+1, bitsPerInt32)
		if err != nil {
			panic(err)
		}
		arrayEnd--
	}

	arrayStart := 0
	if indexStart > 0 {
		arrayStart, err = getArrayLength(indexStart+1, bitsPerInt32)
		if err != nil {
			panic(err)
		}
		arrayStart--
	}

	index := 0
	offset := indexStart % bitsPerInt32

	for i := arrayStart; i < arrayEnd; i++ {
		vector.array[index] = (s.array[i] >> offset) ^ (s.array[i+1] << (bitsPerInt32 - offset))
		index++
	}

	vector.array[index] = s.array[arrayEnd] >> offset
	vector.version++
}

func (s *BitVector) Length() int {
	return s.length
}

func (s *BitVector) Resize(length int) {
	if length < 0 {
		panic(fmt.Errorf("need non-negative number"))
	}

	arrayLength, err := getArrayLength(length, bitsPerInt32)

	if err != nil {
		panic(err)
	}

	if length != s.length {
		newarray := make([]uint32, arrayLength)
		if len(s.array) != 0 {
			copy(newarray, s.array[:arrayLength])
		}
		s.array = newarray
	}

	if length > s.length {
		last, err := getArrayLength(s.length, bitsPerInt32)
		if err != nil {
			panic(err)
		}
		last--

		bits := s.length % bitsPerInt32
		if bits > 0 {
			s.array[last] &= (1 << bits) - 1
		}
	}

	s.length = length
	s.version++
}

func getArrayLength(n int, div int) (int, error) {
	if div < 0 {
		return 0, fmt.Errorf("div arg must be greater than 0")
	}
	if n > 0 {
		return ((n - 1) / div) + 1, nil
	}
	return 0, nil
}

// Rank counts the number of true or false (depending on what the bit is set to)
// in the bitvector but not including the offset
func (s *BitVector) Rank(bit bool, offset int) int {
	rank := 0

	iterator := s.EnumerateFromOffset(0, offset)
	for iterator.HasNext() {
		v, _ := iterator.Next()

		if v == bit {
			rank++
		}
	}

	return rank
}

// find the offset of true or false (depending on what the bit is set to) from the rank
// (number of times the bit occurs)
func (s *BitVector) Select(bit bool, rank int) int {
	offset := -1
	match := -1
	iterator := s.EnumerateFromOffset(0, s.Length())

	for iterator.HasNext() {
		v, index := iterator.Next()

		if v == bit {
			offset++
			match = index
		}
		if offset == rank {
			break
		}
	}

	return match
}

func (s *BitVector) Concat(vectors []*BitVector) *BitVector {

	length := s.Length()
	for _, v := range vectors {
		length += v.Length()
	}

	vector := NewBitVector(length)

	iterator := s.Enumerate()

	for iterator.HasNext() {
		value, i := iterator.Next()

		vector.Set(i, value)
	}

	index := s.Length()
	for i, v := range vectors {
		index += i

		iterator := v.Enumerate()

		for iterator.HasNext() {
			value, i := iterator.Next()

			index += i
			vector.Set(index, value)
		}

	}
	return vector
}

func (s *BitVector) TrueBits() int {
	output := 0

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		output += bits.OnesCount32(s.array[i])
	}

	return output
}

// ANDed with vector.
func (s *BitVector) And(vector *BitVector) {
	if vector == nil {
		panic(fmt.Errorf("vector is null"))
	}

	if s.Length() != vector.Length() {
		panic(fmt.Errorf("vector length is different"))
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] &= vector.array[i]
	}

	s.version++
}

// ORed with vector.
func (s *BitVector) Or(vector *BitVector) {
	if vector == nil {
		panic(fmt.Errorf("vector is null"))
	}

	if s.Length() != vector.Length() {
		panic(fmt.Errorf("vector length is different"))
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] |= vector.array[i]
	}

	s.version++
}

// XORed with vector.
func (s *BitVector) Xor(vector *BitVector) {
	if vector == nil {
		panic(fmt.Errorf("vector is null"))
	}

	if s.Length() != vector.Length() {
		panic(fmt.Errorf("vector length is different"))
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] ^= vector.array[i]
	}

	s.version++
}

// Inverts all the bit values. On/true bit values are converted to off/false. Off/false bit values are turned on/true.
func (s *BitVector) Not() {
	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] = ^s.array[i]
	}

	s.version++
}

func (s BitVector) String() string {
	str := []string{}
	iterator := s.Enumerate()

	for iterator.HasNext() {
		value, _ := iterator.Next()

		str = append(str, strconv.FormatBool(value))
	}
	return fmt.Sprintf("{ %s }\n", strings.Join(str, ", "))
}

func (s *BitVector) Enumerate() *BitVectorIterator {
	return NewBitVectorIteratorWithOffset(s, 0, s.Length())
}

func (s *BitVector) EnumerateFromOffset(indexStart, indexEnd int) *BitVectorIterator {
	return NewBitVectorIteratorWithOffset(s, indexStart, indexEnd)
}

type BitVectorIterator struct {
	vector         *BitVector
	version        int
	indexStart     int
	indexEnd       int
	currentElement bool
}

func NewBitVectorIteratorWithOffset(vector *BitVector, indexStart, indexEnd int) *BitVectorIterator {
	if indexStart > vector.Length() {
		panic("indexStart grater or equal to length")
	}
	if indexStart < 0 {
		panic("indexStart must be non negative number")
	}

	if indexStart > indexEnd {
		panic("indexEnd must be greater then indexStart")
	}

	if vector.Length()-indexStart < 0 {
		panic("invalid indexStart length")
	}

	if indexEnd > vector.Length() {
		panic("indexEnd must be greater then vector length")
	}

	return &BitVectorIterator{
		vector:     vector,
		indexStart: indexStart,
		indexEnd:   indexEnd,
		version:    vector.version,
	}
}

func (s *BitVectorIterator) Reset() {
	if s.version != s.vector.version {
		panic("version failed")
	}
	s.indexStart = 0
}

func (s *BitVectorIterator) HasNext() bool {
	return s.indexStart < s.indexEnd
}

func (s *BitVectorIterator) Next() (bool, int) {
	if s.version != s.vector.version {
		panic("version failed")
	}

	if s.indexStart < s.vector.Length() {
		index := s.indexStart
		currentElement := s.vector.Get(s.indexStart)
		s.currentElement = currentElement
		s.indexStart++
		return currentElement, index
	}

	s.indexStart = s.vector.Length()

	return false, s.indexStart
}
