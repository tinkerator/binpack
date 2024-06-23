package binpack

import (
	"fmt"
	"strings"
	"testing"
)

func TestPack(t *testing.T) {
	b := NewBoard(5, 6)
	tiles := []Tile{
		{3, 3},
		{1, 2},
		{4, 3},
		{1, 1},
		{1, 1},
		{1, 2},
		{1, 1},
		{1, 2},
	}
	for i, tile := range tiles {
		t.Logf("tile[%d] = %#v", i, tile)
	}
	n := b.Pack(tiles)
	if n != 1 {
		t.Errorf("got n=%d, want=1", n)
	}
	dat := make([]string, b.Area.Dx*b.Area.Dy)
	for i := range dat {
		dat[i] = "."
	}
	for i, id := range b.Indices {
		dx := b.Content[i].TR.X - b.Content[i].LL.X
		t.Logf("%d: [%d] = %#v (rotated=%v)", i, id, b.Content[i], dx != tiles[id].Dx)
		for j := b.Content[i].LL.X; j < b.Content[i].TR.X; j++ {
			for k := b.Content[i].LL.Y; k < b.Content[i].TR.Y; k++ {
				dat[j+b.Area.Dx*k] = fmt.Sprintf("%X", id)
			}
		}
	}
	for k := 0; k < b.Area.Dy; k++ {
		t.Logf("%s", strings.Join(dat[b.Area.Dx*(b.Area.Dy-1-k):b.Area.Dx*(b.Area.Dy-k)], ""))
	}
}
