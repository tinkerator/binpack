package binpack

import (
	"fmt"
	"sort"
)

// Point is an integer coordinate.
type Point struct {
	X, Y int
}

// Tile is an positive integer sized rectangle.
type Tile struct {
	Dx, Dy int
}

// BBox holds the lower-left and top-right Points.
type BBox struct {
	LL, TR Point
}

// Board holds the current state of a board.
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

// dropMerges looks for a space that shares x coordinates with the
// specified space and touches a y coordinate. Having found the space,
// said space is dropped from the space tree.
func (s *Space) dropMerges(this BBox) (dropped bool, bb BBox) {
	if s == nil || s.Free == bb {
		// s has no empty space (bb is the zero value here).
		return
	}
	fr := s.Free
	var larger BBox
	found := false
	if fr.LL.X == this.LL.X && fr.TR.X == this.TR.X && (this.LL.Y == fr.TR.Y || this.TR.Y == fr.LL.Y) {
		found = true
	}

	dropped, bb = s.Left.dropMerges(this)
	if !dropped {
		dropped, bb = s.Right.dropMerges(this)
	}
	if !dropped {
		if found {
			s.Free = bb
			dropped = true
			bb = fr
		}
		return
	}

	var aL, aR int
	if s.Left != nil {
		aL = s.Left.Free.Area()
	}
	if s.Right != nil {
		aR = s.Right.Free.Area()
		larger = s.Right.Free
	}
	if aL > aR {
		larger = s.Left.Free
		s.Right, s.Left = s.Left, s.Right
	} else {
		s.Free = larger
	}

	return
}

// replace replaces a s.Free that matches older with newer.
func (s *Space) replace(older, newer BBox) bool {
	if s == nil || s.Free.TR.X == 0 {
		return false
	}
	if s.Free == older {
		s.Free = newer
		s.Right.replace(older, newer)
		return true
	}

	done := s.Left.replace(older, newer)
	if !done {
		if done = s.Right.replace(older, newer); !done {
			return false
		}
		s.Free = s.Right.Free
		return true
	}

	var aR int
	if s.Right != nil {
		aR = s.Right.Free.Area()
		s.Free = s.Right.Free
	}
	if s.Left.Free.Area() > aR {
		s.Free = s.Left.Free
		s.Left, s.Right = s.Right, s.Left
	}
	return true
}

// Area computes the area of the bounding box.
func (bb BBox) Area() int {
	return (bb.TR.X - bb.LL.X) * (bb.TR.Y - bb.LL.Y)
}

// add inserts a bounding box of tile size and id into an available
// space.  A return value of true indicates space was allocated. The
// fragment value returned is the fragment of space (further X)
// remaining after the addition.
func (s *Space) add(id int, tile Tile) (added bool, fragment BBox) {
	dx, dy := tile.Dx, tile.Dy
	if dx*dy > s.Free.Area() {
		return
	}
	left, right := s.Left, s.Right
	if left == right { // both nil
		if rotate, ok := s.Free.Fits(tile); !ok {
			return
		} else if rotate {
			dx, dy = dy, dx
		}
		s.Index = id
		s.Box = BBox{
			LL: Point{s.Free.LL.X, s.Free.LL.Y},
			TR: Point{s.Free.LL.X + dx, s.Free.LL.Y + dy},
		}
		if s.Free.LL.X+dx != s.Free.TR.X {
			left = &Space{
				Free: BBox{
					LL: Point{s.Free.LL.X + dx, s.Free.LL.Y},
					TR: Point{s.Free.TR.X, s.Free.LL.Y + dy},
				},
			}
			fragment = left.Free
		}
		if s.Free.LL.Y+dy != s.Free.TR.Y {
			right = &Space{
				Free: BBox{
					LL: Point{s.Free.LL.X, s.Free.LL.Y + dy},
					TR: Point{s.Free.TR.X, s.Free.TR.Y},
				},
			}
		}
		added = true
	} else {
		if left != nil {
			added, fragment = left.add(id, tile)
		}
		if !added && right != nil {
			added, fragment = right.add(id, tile)
		}
		if !added {
			// no space for the tile after all.
			return
		}
	}
	// Adjust the Free space to the larger of the children,
	// and make that the right child.
	var dxL, dyL, dxR, dyR int
	if left != nil {
		dxL = left.Free.TR.X - left.Free.LL.X
		dyL = left.Free.TR.Y - left.Free.LL.Y
	}
	if right != nil {
		dxR = right.Free.TR.X - right.Free.LL.X
		dyR = right.Free.TR.Y - right.Free.LL.Y
	}
	// make sure that the s.Left has the smaller area.
	if dxL*dyL < dxR*dyR {
		s.Left = left
		s.Right = right
	} else {
		s.Left = right
		s.Right = left
	}
	if s.Right != nil {
		s.Free = s.Right.Free
	} else {
		s.Free = BBox{}
	}
	return
}

// build flattens the Space tree into a list of BBox values and their
// corresponding indices.
func (b *Board) build(root *Space) {
	if root == nil {
		return
	}
	if root.Box.TR.X != 0 {
		b.Content = append(b.Content, root.Box)
		b.Indices = append(b.Indices, root.Index)
	}
	b.build(root.Left)
	b.build(root.Right)
}

// Add inserts a bounding box of tile size and id into an available
// space.  A return value of true indicates space was allocated.
func (s *Space) Add(id int, tile Tile) bool {
	ok, _ := s.add(id, tile)
	return ok
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
			if ok, fragment := root.add(id, tiles[id]); ok {
				if fragment.TR.X != 0 {
					dropped, bb := root.dropMerges(fragment)
					if dropped {
						resize := fragment
						if resize.LL.Y == bb.TR.Y {
							resize.LL = bb.LL
						} else {
							resize.TR = bb.TR
						}
						if !root.replace(fragment, resize) {
							panic(fmt.Sprintf("programming error: failed to replace fragment=%v with resize=%v", fragment, resize))
						}
					}
				}
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
