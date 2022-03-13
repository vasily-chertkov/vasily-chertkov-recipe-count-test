package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {

	type input struct {
		vmIds []string
	}

	type expect struct {
		slice   *Slice
		vmBytes [][]byte
	}

	var tests = []struct {
		name   string
		input  input
		expect expect
	}{
		{
			"Different VMIds length",
			input{[]string{"vm-2000", "vm-DEAD", "vm-вмка", "vm-5abcd", "vm-100"}},
			expect{
				slice: &Slice{
					StringSlice: []string{"vm-100", "vm-2000", "vm-5abcd", "vm-DEAD", "vm-вмка"},
					pos: map[string][2]int{
						"vm-100":   {1, 10},
						"vm-2000":  {10, 20},
						"vm-5abcd": {20, 31},
						"vm-DEAD":  {31, 41},
						"vm-вмка":  {40, 54},
					},
					bytes: []byte(`["vm-100","vm-2000","vm-5abcd","vm-DEAD","vm-вмка"]`),
				},
				vmBytes: [][]byte{
					[]byte(`["vm-2000","vm-5abcd","vm-DEAD","vm-вмка"]`),
					[]byte(`["vm-100","vm-5abcd","vm-DEAD","vm-вмка"]`),
					[]byte(`["vm-100","vm-2000","vm-DEAD","vm-вмка"]`),
					[]byte(`["vm-100","vm-2000","vm-5abcd","vm-вмка"]`),
					[]byte(`["vm-100","vm-2000","vm-5abcd","vm-DEAD"]`),
				},
			},
		},
		{
			"Single VMId",
			input{[]string{"vm-2000"}},
			expect{
				slice: &Slice{
					StringSlice: []string{"vm-2000"},
					pos: map[string][2]int{
						"vm-2000": {1, 10},
					},
					bytes: []byte(`["vm-2000"]`),
				},
				vmBytes: [][]byte{
					[]byte(`[]`),
				},
			},
		},
	}

	for _, test := range tests {
		s := NewSlice(test.input.vmIds)
		assert.Equal(t, test.expect.slice, s, test.name+": Slice test")
		for idx, vmByte := range test.expect.vmBytes {
			b := s.Bytes(s.StringSlice[idx])
			assert.Equal(t, vmByte, b, test.name+": Bytes test")
		}
	}
}
