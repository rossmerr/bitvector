package bitvector_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/rossmerr/bitvector"
)

func TestBitVector_Get_Set(t *testing.T) {
	tests := []struct {
		name   string
		values []bool
		length int
		index  int
		want   bool
	}{
		{
			name:   "Get Set",
			values: []bool{true, false, true, true, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVector(len(tt.values))
			for i := 0; i < len(tt.values); i++ {
				s.Set(i, tt.values[i])
			}

			for i := 0; i < len(tt.values); i++ {
				if got := s.Get(i); got != tt.values[i] {
					t.Errorf("BitVector.Get() = %v, want %v", got, tt.values[i])
				}
			}
		})
	}
}

func TestBitVector_SetAll(t *testing.T) {

	tests := []struct {
		name         string
		length       int
		defaultValue bool
	}{
		{
			name:         "SetAll",
			length:       5,
			defaultValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorOfLength(tt.length, tt.defaultValue)
			s.SetAll(!tt.defaultValue)

			for i := 0; i < tt.length; i++ {
				got := s.Get(i)
				if got != !tt.defaultValue {
					t.Errorf("BitVector.Get() = %v, want %v", got, !tt.defaultValue)
				}
			}

			s.SetAll(tt.defaultValue)
			for i := 0; i < tt.length; i++ {
				got := s.Get(i)
				if got != tt.defaultValue {
					t.Errorf("BitVector.Get() = %v, want %v", got, tt.defaultValue)
				}
			}
		})
	}
}

func TestBitVector_Enumerate(t *testing.T) {

	tests := []struct {
		name   string
		values []bool
		want   *bitvector.BitVectorIterator
	}{
		{
			name:   "Enumerate",
			values: []bool{true, false, true, true, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorFromBool(tt.values)

			iterator := s.Enumerate()
			counter := 0
			for iterator.HasNext() {
				value, _ := iterator.Next()

				current := s.Get(counter)
				if value != current {
					t.Errorf("BitVector.Get() = %v, want %v", value, current)
				}
				counter++
			}

			if counter != len(tt.values) {
				t.Errorf("counter = %v, want %v", counter, len(tt.values))
			}

		})
	}
}

func TestBitVector_Resize(t *testing.T) {

	tests := []struct {
		name         string
		length       int
		newLength    int
		defaultValue bool
	}{
		{
			name:         "Resize",
			length:       5,
			newLength:    10,
			defaultValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorOfLength(tt.length, tt.defaultValue)

			s.Resize(tt.newLength)
			if s.Length() != tt.newLength {
				t.Errorf("BitVector.Get() = %v, want %v", s.Length(), tt.newLength)
			}

			for i := 0; i < int(math.Min(float64(tt.length), float64(s.Length()))); i++ {
				got := s.Get(i)
				if !got {
					t.Errorf("BitVector.Get() = %v, want %v", got, true)
				}
			}

			for i := tt.length; i < tt.newLength; i++ {
				got := s.Get(i)
				if got {
					t.Errorf("BitVector.Get() = %v, want %v", got, false)
				}

			}

			s.Resize(0)
			if s.Length() != 0 {
				t.Errorf("BitVector.Get() newLength= %v, want %v", s.Length(), 0)
			}

			s.Resize(tt.newLength)
			if s.Length() != tt.newLength {
				t.Errorf("BitVector.Get() newLength= %v, want %v", s.Length(), tt.newLength)
			}

			first := s.Get(0)
			if first {
				t.Errorf("BitVector.Get() first= %v, want %v", first, false)
			}
			last := s.Get(tt.newLength - 1)
			if last {
				t.Errorf("BitVector.Get() last = %v, want %v", last, false)
			}
		})
	}
}

func TestBitVector_Copy(t *testing.T) {
	tests := []struct {
		name       string
		values     []bool
		length     int
		indexStart int
		indexEnd   int
		want       []bool
	}{
		{
			name: "Copy",
			values: []bool{
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
				true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
			},
			want: []bool{
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
				true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
			},
			length:     86,
			indexStart: 10,
			indexEnd:   86,
		},
		{
			name: "Copy over array bounds",
			values: []bool{
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
				true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
			},
			want: []bool{
				true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
			},
			length:     32,
			indexStart: 32,
			indexEnd:   32,
		},
		{
			name: "Copy over array bounds 2",
			values: []bool{
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
				true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
			},
			want: []bool{
				true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
			},
			length:     28,
			indexStart: 34,
			indexEnd:   34,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVector(len(tt.values))

			for i := 0; i < len(tt.values); i++ {
				s.Set(i, tt.values[i])
			}
			copy := bitvector.NewBitVector(tt.length)

			s.Copy(copy, tt.indexStart, tt.indexEnd)

			fmt.Println(s.Length() - 10)
			for i := 0; i < len(tt.want); i++ {
				if got := copy.Get(i); got != tt.want[i] {
					t.Errorf("BitVector.Get(%v) = %v, want %v", i, got, tt.want[i])
				}
			}
		})
	}
}

func TestBitVector_Not(t *testing.T) {

	tests := []struct {
		name   string
		values []bool
	}{
		{
			name:   "Not",
			values: []bool{true, false, true, true, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorFromBool(tt.values)

			s.Not()
			for i := 0; i < s.Length(); i++ {
				got := s.Get(i)
				if got == tt.values[i] {
					t.Errorf("BitVector.Get() = %v, want %v", got, !tt.values[i])
				}
			}

		})
	}
}

func TestBitVector_Xor(t *testing.T) {

	tests := []struct {
		name  string
		left  []bool
		right []bool
		want  []bool
	}{
		{
			name:  "Xor",
			left:  []bool{false, false, true, true},
			right: []bool{false, true, false, true},
			want:  []bool{false, true, true, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := bitvector.NewBitVectorFromBool(tt.left)
			right := bitvector.NewBitVectorFromBool(tt.right)

			left.Xor(right)
			for i := 0; i < left.Length(); i++ {
				got := left.Get(i)
				if got != tt.want[i] {
					t.Errorf("BitVector.Get() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}

func TestBitVector_Or(t *testing.T) {

	tests := []struct {
		name  string
		left  []bool
		right []bool
		want  []bool
	}{
		{
			name:  "Or",
			left:  []bool{false, false, true, true},
			right: []bool{false, true, false, true},
			want:  []bool{false, true, true, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := bitvector.NewBitVectorFromBool(tt.left)
			right := bitvector.NewBitVectorFromBool(tt.right)

			left.Or(right)
			for i := 0; i < left.Length(); i++ {
				got := left.Get(i)
				if got != tt.want[i] {
					t.Errorf("BitVector.Get() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}

func TestBitVector_And(t *testing.T) {
	tests := []struct {
		name  string
		left  []bool
		right []bool
		want  []bool
	}{
		{
			name:  "And",
			left:  []bool{false, false, true, true},
			right: []bool{false, true, false, true},
			want:  []bool{false, false, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := bitvector.NewBitVectorFromBool(tt.left)
			right := bitvector.NewBitVectorFromBool(tt.right)

			left.And(right)
			for i := 0; i < left.Length(); i++ {
				got := left.Get(i)
				if got != tt.want[i] {
					t.Errorf("BitVector.Get() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}

func TestNewBitVectorFromVectorPadStart(t *testing.T) {
	tests := []struct {
		name    string
		values  []bool
		padding int
		want    []bool
	}{
		{
			name:    "PadStart",
			values:  []bool{true, false, false, true, true},
			padding: 2,
			want:    []bool{false, false, true, false, false, true, true},
		},
		{
			name:    "PadStart",
			values:  []bool{true, false, false, true, true},
			padding: 32,
			want:    []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, true, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorFromBool(tt.values)
			result := bitvector.NewBitVectorFromVectorPadStart(s, tt.padding)
			for i := 0; i < len(tt.want); i++ {
				got := result.Get(i)
				if got != tt.want[i] {
					t.Errorf("BitVector.Get(%v) = %v, want %v", i, got, tt.want[i])
				}
			}
		})
	}
}

func TestBitVector_Rank(t *testing.T) {
	tests := []struct {
		name   string
		values []bool
		value  bool
		offset int
		want   int
	}{
		{
			name:   "rank",
			values: []bool{false, true, true, false},
			value:  true,
			offset: 2,
			want:   1,
		},
		{
			name:   "rank",
			values: []bool{false, true, true, false, true, false},
			value:  false,
			offset: 5,
			want:   2,
		},
		{
			name:   "rank",
			values: []bool{false, true, true, false, true, false, true, false, false},
			value:  true,
			offset: 8,
			want:   4,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorFromBool(tt.values)
			got := s.Rank(tt.value, tt.offset)
			if got != tt.want {
				t.Errorf("BitVector.Rank(%v) = %v, want %v", i, got, tt.want)
			}
		})
	}
}

func TestBitVector_Select(t *testing.T) {
	tests := []struct {
		name   string
		values []bool
		value  bool
		rank   int
		want   int
	}{
		{
			name:   "select",
			values: []bool{false, true, true, false},
			value:  true,
			rank:   0,
			want:   1,
		},
		{
			name:   "select",
			values: []bool{false, true, true, false, true, false, true, false, false, false},
			value:  false,
			rank:   3,
			want:   7,
		},
		{
			name:   "select",
			values: []bool{false, true, true, false, true, false, true, false, false, false, true, false, true},
			value:  true,
			rank:   5,
			want:   12,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := bitvector.NewBitVectorFromBool(tt.values)
			got := s.Select(tt.value, tt.rank)
			if got != tt.want {
				t.Errorf("BitVector.Select(%v) = %v, want %v", i, got, tt.want)
			}
		})
	}
}
