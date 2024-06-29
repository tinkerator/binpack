# binpack - a package for 2D rectangle bin-packing

## Overview

This package implements a very simple algorithm to pack a set of
rectangles (`binpack.Tile`s) into a larger rectangle
(`binpack.Board`).

The package includes a unit test that demonstrates bin-packing in
action and generates an ASCII art board (5 wide, 6 high):
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
    binpack_test.go:34: 3: [4] = binpack.BBox{LL:binpack.Point{X:4, Y:3}, TR:binpack.Point{X:5, Y:4}} (rotated=false)
    binpack_test.go:34: 4: [6] = binpack.BBox{LL:binpack.Point{X:4, Y:4}, TR:binpack.Point{X:5, Y:5}} (rotated=false)
    binpack_test.go:34: 5: [7] = binpack.BBox{LL:binpack.Point{X:3, Y:5}, TR:binpack.Point{X:5, Y:6}} (rotated=true)
    binpack_test.go:34: 6: [1] = binpack.BBox{LL:binpack.Point{X:4, Y:0}, TR:binpack.Point{X:5, Y:2}} (rotated=false)
    binpack_test.go:34: 7: [3] = binpack.BBox{LL:binpack.Point{X:4, Y:2}, TR:binpack.Point{X:5, Y:3}} (rotated=false)
    binpack_test.go:42: 00077
    binpack_test.go:42: 00056
    binpack_test.go:42: 00054
    binpack_test.go:42: 22223
    binpack_test.go:42: 22221
    binpack_test.go:42: 22221
--- PASS: TestPack (0.00s)
PASS
ok  	zappem.net/pub/graphics/binpack	0.003s
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
