package chunk

import (
	"testing"

	"github.com/go-test/deep"
)

// TestMerge tests that merge is valid
func TestNumChunk(t *testing.T) {
	table := []struct{ l, n, result int }{
		{1, 1, 1},
		{800, 200, 4},
		{801, 200, 5},
		{1000, 200, 5},
		{1001, 200, 6},
	}

	for _, tbl := range table {
		result := NumChunk(tbl.l, tbl.n)

		if result != tbl.result {
			t.Errorf("fail. len=%d n=%d expected=%d got=%d\n", tbl.l, tbl.n, tbl.result, result)
		}
	}

}

// TestChunk tests Chunk
func TestChunk(t *testing.T) {
	table := []struct {
		l, n   int
		result [][2]int
	}{
		{100, 1, [][2]int{{0, 100}}},
		{800, 3, [][2]int{{0, 267}, {267, 534}, {534, 800}}},
	}

	for _, tbl := range table {
		result := Chunk(tbl.l, tbl.n)

		if diff := deep.Equal(result, tbl.result); diff != nil {
			for _, msg := range diff {
				t.Errorf("%s\n", msg)
			}
		}
	}
}
