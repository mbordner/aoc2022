package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Test_Pairs(t *testing.T) {

	testCases := []struct {
		left     string
		right    string
		expected bool
	}{
		{
			left:     `[1,1,3,1,1]`,
			right:    `[1,1,5,1,1]`,
			expected: true,
		},
		{
			left:     `[1,1,3,1,1]`,
			right:    `[1,1,3,1,1]`,
			expected: true,
		},
		{
			left:     `[1]`,
			right:    `[1]`,
			expected: true,
		},
		{
			left:     `[[1],[2,3,4]]`,
			right:    `[[1],4]`,
			expected: true,
		},
		{
			left:     `[9]`,
			right:    `[[8,7,6]]`,
			expected: false,
		},
		{
			left:     `[[4,4],4,4]`,
			right:    `[[4,4],4,4,4]`,
			expected: true,
		},
		{
			left:     `[7,7,7]`,
			right:    `[7,7,7,7]`,
			expected: true,
		},
		{
			left:     `[7,7,7,7]`,
			right:    `[7,7,7]`,
			expected: false,
		},
		{
			left:     `[]`,
			right:    `[3]`,
			expected: true,
		},
		{
			left:     `[3]`,
			right:    `[]`,
			expected: false,
		},
		{
			left:     `[[[]]]`,
			right:    `[[]]`,
			expected: false,
		},
		{
			left:     `[1,[2,[3,[4,[5,6,7]]]],8,9]`,
			right:    `[1,[2,[3,[4,[5,6,0]]]],8,9]`,
			expected: false,
		},
		{
			left:     `[[1],[2,3,4]]`,
			right:    `[[1],2]`,
			expected: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i+1), func(t *testing.T) {

			var pair Pair
			pair = make([]Packet, 2, 2)

			err := json.Unmarshal([]byte(tc.left), &(pair[0]))
			assert.Nil(t, err)

			err = json.Unmarshal([]byte(tc.right), &(pair[1]))
			assert.Nil(t, err)

			inOrder := sort.IsSorted(pair)
			assert.Equal(t, tc.expected, inOrder)

		})
	}
}
