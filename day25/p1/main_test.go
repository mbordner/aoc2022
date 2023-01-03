package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Pow(t *testing.T) {
	assert.Equal(t, 25, pow(5, 2))
	assert.Equal(t, 1, pow(5, 0))
	assert.Equal(t, 5, pow(5, 1))
	assert.Equal(t, 125, pow(5, 3))
	assert.Equal(t, 625, pow(5, 4))
}

func Test_Snafu(t *testing.T) {

	testCases := []struct {
		snafu string
		dec   int
	}{
		{
			snafu: `1=-0-2`,
			dec:   1747,
		},
		{
			snafu: `12111`,
			dec:   906,
		},
		{
			snafu: `2=0=`,
			dec:   198,
		},
		{
			snafu: `21`,
			dec:   11,
		},
		{
			snafu: `2=01`,
			dec:   201,
		},
		{
			snafu: `111`,
			dec:   31,
		},
		{
			snafu: `20012`,
			dec:   1257,
		},
		{
			snafu: `112`,
			dec:   32,
		},
		{
			snafu: `1=-1=`,
			dec:   353,
		},
		{
			snafu: `1-12`,
			dec:   107,
		},
		{
			snafu: `12`,
			dec:   7,
		},
		{
			snafu: `1=`,
			dec:   3,
		},
		{
			snafu: `122`,
			dec:   37,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.snafu, func(t *testing.T) {
			v := snafuToDec(tc.snafu)
			assert.Equal(t, tc.dec, v)
			s := decToSnafu(v)
			assert.Equal(t, tc.snafu, s)
		})
	}

}
