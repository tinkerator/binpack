# binpack - a package for 2D rectangle bin-packing

## Overview

This [package](https://zappem.net/pub/graphics/binpack/) implements a
very simple algorithm to pack a set of rectangles (`binpack.Tile`s)
into a larger rectangle (`binpack.Board`).

The package includes a unit test that demonstrates bin-packing in
action and generates an ASCII art boards (5 wide, 6 high):

```
$ go test -v
=== RUN   TestPack
    binpack_test.go:22: tile[0] = binpack.Tile{Dx:3, Dy:3}
    binpack_test.go:22: tile[1] = binpack.Tile{Dx:1, Dy:2}
    binpack_test.go:22: tile[2] = binpack.Tile{Dx:4, Dy:3}
    binpack_test.go:22: tile[3] = binpack.Tile{Dx:1, Dy:1}
    binpack_test.go:22: tile[4] = binpack.Tile{Dx:1, Dy:1}
    binpack_test.go:22: tile[5] = binpack.Tile{Dx:1, Dy:2}
    binpack_test.go:22: tile[6] = binpack.Tile{Dx:1, Dy:1}
    binpack_test.go:22: tile[7] = binpack.Tile{Dx:1, Dy:2}
    binpack_test.go:34: 0: [2] = binpack.BBox{LL:binpack.Point{X:0, Y:0}, TR:binpack.Point{X:4, Y:3}} (rotated=false)
    binpack_test.go:34: 1: [0] = binpack.BBox{LL:binpack.Point{X:0, Y:3}, TR:binpack.Point{X:3, Y:6}} (rotated=false)
    binpack_test.go:34: 2: [5] = binpack.BBox{LL:binpack.Point{X:3, Y:3}, TR:binpack.Point{X:4, Y:5}} (rotated=false)
    binpack_test.go:34: 3: [3] = binpack.BBox{LL:binpack.Point{X:4, Y:2}, TR:binpack.Point{X:5, Y:3}} (rotated=false)
    binpack_test.go:34: 4: [4] = binpack.BBox{LL:binpack.Point{X:4, Y:3}, TR:binpack.Point{X:5, Y:4}} (rotated=false)
    binpack_test.go:34: 5: [6] = binpack.BBox{LL:binpack.Point{X:4, Y:4}, TR:binpack.Point{X:5, Y:5}} (rotated=false)
    binpack_test.go:34: 6: [7] = binpack.BBox{LL:binpack.Point{X:3, Y:5}, TR:binpack.Point{X:5, Y:6}} (rotated=true)
    binpack_test.go:34: 7: [1] = binpack.BBox{LL:binpack.Point{X:4, Y:0}, TR:binpack.Point{X:5, Y:2}} (rotated=false)
    binpack_test.go:42: 00077
    binpack_test.go:42: 00056
    binpack_test.go:42: 00054
    binpack_test.go:42: 22223
    binpack_test.go:42: 22221
    binpack_test.go:42: 22221
--- PASS: TestPack (0.00s)
=== RUN   TestFullPack
    binpack_test.go:53: tile[0] = binpack.Tile{Dx:2, Dy:1}
    binpack_test.go:53: tile[1] = binpack.Tile{Dx:1, Dy:2}
    binpack_test.go:65: 0: [1] = binpack.BBox{LL:binpack.Point{X:0, Y:0}, TR:binpack.Point{X:1, Y:2}} (rotated=false)
    binpack_test.go:65: 1: [0] = binpack.BBox{LL:binpack.Point{X:0, Y:2}, TR:binpack.Point{X:2, Y:3}} (rotated=false)
    binpack_test.go:65: 2: [0] = binpack.BBox{LL:binpack.Point{X:0, Y:3}, TR:binpack.Point{X:2, Y:4}} (rotated=false)
    binpack_test.go:65: 3: [1] = binpack.BBox{LL:binpack.Point{X:0, Y:4}, TR:binpack.Point{X:1, Y:6}} (rotated=false)
    binpack_test.go:65: 4: [0] = binpack.BBox{LL:binpack.Point{X:1, Y:4}, TR:binpack.Point{X:3, Y:5}} (rotated=false)
    binpack_test.go:65: 5: [0] = binpack.BBox{LL:binpack.Point{X:1, Y:5}, TR:binpack.Point{X:3, Y:6}} (rotated=false)
    binpack_test.go:65: 6: [1] = binpack.BBox{LL:binpack.Point{X:3, Y:5}, TR:binpack.Point{X:5, Y:6}} (rotated=true)
    binpack_test.go:65: 7: [1] = binpack.BBox{LL:binpack.Point{X:3, Y:4}, TR:binpack.Point{X:5, Y:5}} (rotated=true)
    binpack_test.go:65: 8: [1] = binpack.BBox{LL:binpack.Point{X:2, Y:3}, TR:binpack.Point{X:4, Y:4}} (rotated=true)
    binpack_test.go:65: 9: [0] = binpack.BBox{LL:binpack.Point{X:4, Y:2}, TR:binpack.Point{X:5, Y:4}} (rotated=true)
    binpack_test.go:65: 10: [1] = binpack.BBox{LL:binpack.Point{X:2, Y:2}, TR:binpack.Point{X:4, Y:3}} (rotated=true)
    binpack_test.go:65: 11: [0] = binpack.BBox{LL:binpack.Point{X:1, Y:0}, TR:binpack.Point{X:3, Y:1}} (rotated=false)
    binpack_test.go:65: 12: [0] = binpack.BBox{LL:binpack.Point{X:1, Y:1}, TR:binpack.Point{X:3, Y:2}} (rotated=false)
    binpack_test.go:65: 13: [1] = binpack.BBox{LL:binpack.Point{X:3, Y:1}, TR:binpack.Point{X:5, Y:2}} (rotated=true)
    binpack_test.go:65: 14: [1] = binpack.BBox{LL:binpack.Point{X:3, Y:0}, TR:binpack.Point{X:5, Y:1}} (rotated=true)
    binpack_test.go:73: 10011
    binpack_test.go:73: 10011
    binpack_test.go:73: 00110
    binpack_test.go:73: 00110
    binpack_test.go:73: 10011
    binpack_test.go:73: 10011
--- PASS: TestFullPack (0.00s)
PASS
ok      zappem.net/pub/graphics/binpack 0.002s
```

## API

The API provided by this package is avalble using `go doc
zappem.net/pub/graphics/binpack`. It can also be browsed on the
[go.dev](http://go.dev) website:
[package zappem.net/pub/graphics/binpack](https://pkg.go.dev/zappem.net/pub/graphics/binpack).

## License info

The `binpack` package is distributed with the same BSD 3-clause
license as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs

Use the [github `binpack` bug
tracker](https://github.com/tinkerator/binpack/issues).
