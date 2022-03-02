package waterjug

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaterJug(t *testing.T) {
	tcases := []struct {
		x, y, z int
		solved  bool
		steps   []Node
	}{
		{
			x: 2, y: 10, z: 4,
			solved: true,
			steps: []Node{
				{X: 0, Y: 0, Op: ""},
				{X: 2, Y: 0, Op: OpFillX},
				{X: 0, Y: 2, Op: OpTransferXY},
				{X: 2, Y: 2, Op: OpFillX},
				{X: 0, Y: 4, Op: OpTransferXY},
			},
		},
		{
			x: 10, y: 2, z: 4,
			solved: true,
			steps: []Node{
				{X: 0, Y: 0, Op: ""},
				{X: 0, Y: 2, Op: OpFillY},
				{X: 2, Y: 0, Op: OpTransferYX},
				{X: 2, Y: 2, Op: OpFillY},
				{X: 4, Y: 0, Op: OpTransferYX},
			},
		},
		{
			x: 2, y: 10, z: 14,
			solved: false,
		},
		{
			x: 0, y: 0, z: 0,
			solved: true,
			steps: []Node{
				{X: 0, Y: 0, Op: ""},
			},
		},
		{
			x: 10, y: 10, z: 0,
			solved: true,
			steps: []Node{
				{X: 0, Y: 0, Op: ""},
			},
		},
		{
			x: 10, y: 0, z: 3,
			solved: false,
		},
		{
			x: 0, y: 0, z: 10,
			solved: false,
		},
	}

	for i, tcase := range tcases {
		t.Run(fmt.Sprintf("case %d: X:%d, Y:%d, z:%d", i, tcase.x, tcase.y, tcase.z), func(t *testing.T) {
			steps, solved := SearchSolution(tcase.x, tcase.y, tcase.z)
			assert.Equal(t, tcase.solved, solved)
			assert.Equal(t, tcase.steps, steps)
		})
	}
}
