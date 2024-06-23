package binpack

import "sort"

type Point struct {
	X, Y int
}

type Tile struct {
	Dx, Dy int
}

type BBox struct {
	LL, TR Point
}

type Board struct {
	Area    Tile
	Content []BBox
	Indices []int
}

// NewBoard defines a new area into which we can pack tiles.
func NewBoard(dx, dy int) *Board {
	return &Board{
		Area: Tile{Dx: dx, Dy: dy},
	}
}

// Space is a working structure to allocate space on a board.
type Space struct {
	// Free is the remaining largest free space of this empty tile
	// or the Left/Right child.
	Free BBox
	// Index is the numerical ID of the device placed within
	// the Box.  Index is considered valid only if Left or Right
	// are non-nil. Use the (*Space).Occupied() function.
	Index int
	Box   BBox

	// Left, Right are children.
	Left, Right *Space
}

// Occupied extracts the index and bounding box location of, s, if the
// bool return value is true. If false, the space is considered empty.
func (s *Space) Occupied() (int, BBox, bool) {
	return s.Index, s.Box, s.Left != s.Right
}

// Fits determines if tile fits in a bounding box. If the tile needs
// to be rotated 90 degrees, rotate is true. If no fit is possible ok
// is false.
func (bb BBox) Fits(tile Tile) (rotate, ok bool) {
	dx, dy := bb.TR.X-bb.LL.X, bb.TR.Y-bb.LL.Y
	if dx >= tile.Dx && dy >= tile.Dy {
		return false, true
	}
	if dy >= tile.Dx && dx >= tile.Dy {
		return true, true
	}
	return false, false
}

// Add inserts a bounding box of tile size and id into an available
// space.  A return value of true indicates space was allocated.
func (s *Space) Add(id int, tile Tile) bool {
	dx, dy := tile.Dx, tile.Dy
	if dx*dy > (s.Free.TR.X-s.Free.LL.X)*(s.Free.TR.Y-s.Free.LL.Y) {
		return false
	}
	left, right := s.Left, s.Right
	if left == right {
		if rotate, ok := s.Free.Fits(tile); !ok {
			return false
		} else if rotate {
			dx, dy = dy, dx
		}
		s.Index = id
		s.Box = BBox{
			LL: Point{s.Free.LL.X, s.Free.LL.Y},
			TR: Point{s.Free.LL.X + dx, s.Free.LL.Y + dy},
		}
		left = &Space{
			Free: BBox{
				LL: Point{s.Free.LL.X + dx, s.Free.LL.Y},
				TR: Point{s.Free.TR.X, s.Free.LL.Y + dy},
			},
		}
		right = &Space{
			Free: BBox{
				LL: Point{s.Free.LL.X, s.Free.LL.Y + dy},
				TR: Point{s.Free.TR.X, s.Free.TR.Y},
			},
		}
	} else if !(s.Left != nil && s.Left.Add(id, tile)) && !(right != nil && right.Add(id, tile)) {
		// no space for the tile after all.
		return false
	}
	dxL := left.Free.TR.X - left.Free.LL.X
	dyL := left.Free.TR.Y - left.Free.LL.Y
	dxR := right.Free.TR.X - right.Free.LL.X
	dyR := right.Free.TR.Y - right.Free.LL.Y
	// make sure that the s.Left has the smaller area.
	if dxL*dyL < dxR*dyR {
		s.Left = left
		s.Right = right
	} else {
		s.Left = right
		s.Right = left
	}
	s.Free = s.Right.Free
	return true
}

// build flattens the Space tree into a list of BBox values and their
// corresponding indices.
func (b *Board) build(root *Space) {
	if root == nil {
		return
	}
	if root.Left != root.Right {
		b.Content = append(b.Content, root.Box)
		b.Indices = append(b.Indices, root.Index)
	}
	b.build(root.Left)
	b.build(root.Right)
}

// Pack all of tiles into the board area. The returned value is the
// number of whole copies of all tiles that are placed out on the
// board.
func (b *Board) Pack(tiles []Tile) int {
	var indices []int
	for i := range tiles {
		indices = append(indices, i)
	}
	sort.Slice(indices, func(i, j int) bool {
		a := tiles[indices[i]]
		b := tiles[indices[j]]
		if a.Dy > b.Dy {
			return true
		} else if a.Dy < b.Dy {
			return false
		}
		if a.Dx > b.Dx {
			return true
		} else {
			return false
		}
	})
	root := &Space{
		Free: BBox{LL: Point{0, 0}, TR: Point{b.Area.Dx, b.Area.Dy}},
	}
	n := 0
	for some := true; some; {
		some = false
		all := true
		for _, id := range indices {
			if ok := root.Add(id, tiles[id]); ok {
				some = true
			} else {
				all = false
			}
		}
		if all {
			n++
		}
	}
	b.build(root)
	return n
}
