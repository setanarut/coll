[![GoDoc](https://godoc.org/github.com/setanarut/coll?status.svg)](https://pkg.go.dev/github.com/setanarut/coll)

# coll - 2d collision library for Go

Features

* Collisions only - no gravity, rigid body handling, or complex solvers
* Is data-oriented and functional

## Conventions

"Sweep" tests indicate at least 1 of the objects is moving. 
The number indicates how many objects are moving. e.g., `box-box-sweep2` means we are comparing 2 aabbs, both of which are moving.
"Overlap" tests don't take movement into account, and this is a static check to see if the 2 entities overlap.
plural forms imply a collection. e.g., `BoxSegmentSweep1Indexed()` checks one box segment against a set of line segments.
If there is more than one collision, the closest collision is set in the `hitInfo` argument.

## Available collision checks

### Box-Box overlap

![Box-Box-overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-overlap.png)

```go
BoxBoxOverlap(boxA, boxB *AABB, hitInfo *HitInfo) bool
```

### Box-Box contain

```go
// BoxBoxContain returns true if a fully contains b.
BoxBoxContain(a, b *AABB) bool
```

### Box-Box sweep 1

![Box-Box sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-sweep1.png)

```go
BoxBoxSweep1(staticBoxA, boxB *AABB, boxBVel v.Vec, hitInfo *HitInfo) bool 
```

### Box-Box sweep 2

![Box-Box sweep 2](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-sweep2.png)

```go
BoxBoxSweep2(boxA, boxB *AABB, boxAVel, boxBVel v.Vec, hitInfo *HitInfo) bool
```

### Box-Segment sweep 1

![Box-Segment sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segment-sweep1.png)

```go
BoxSegmentSweep1(line *Segment, box *AABB, delta v.Vec, hitInfo *HitInfo) bool
```

### Box-Segment sweep 1 indexed

![Box-Segment sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segments-sweep1-indexed.png)

```go
BoxSegmentSweep1Indexed(lines []*Segment, aabb *AABB, delta v.Vec, hitInfo *HitInfo)  (index int)
```

### Box-Point overlap

![Box-Point overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-point-overlap.png)

```go
BoxPointOverlap(box *AABB, point v.Vec, hitInfo *HitInfo) bool
```

### Box-Segment overlap

![Box-Segment overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segment-overlap.png)

```go
BoxSegmentOverlap(box *AABB, start, delta, padding v.Vec, hitInfo *HitInfo) bool
```

### Line-Circle overlap

![alt text](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/ray-circle-overlap.png)

```go
LineCircleOverlap(raySeg *Segment, circ *Circle, overlapSeg *Segment) bool
```

### Segment-Circle overlap

![alt text](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/segment-circle-overlap.png)

```go
SegmentCircleOverlap(seg *Segment, c *Circle) []v.Vec
```

### Box-Circle sweep 2

```go
BoxCircleSweep2(box *AABB, circle *Circle, boxVel, circleVel v.Vec) bool
```

### Box-Tilemap sweep

```go
(c *TileCollider) Collide(box AABB, delta v.Vec, onCollide TileCollisionCallback) v.Vec
```

### Ray-Tilemap DDA

```go
RayTilemapDDA(pos, dir v.Vec, length float64, tm [][]uint8, cellSize float64, h *HitInfo) (bool, image.Point) 
```

## Examples (Ebitengine)

1. Clone this repository
2. In the terminal, change to the examples directory `cd examples`
3. Run a folder with `go run ./foldername`. Example: `go run ./box_box_sweep1 `

## Credits

Most of these collision checks were adapted from existing open source repos:

* [github.com/mreinstein/collision-2d](https://github.com/mreinstein/collision-2d)
* [youtube.com/watch?v=NbSee-XM7WA](https://youtube.com/watch?v=NbSee-XM7WA) - ray-tilemap (RayTilemapDDA)
* [jonathanwhiting.com/tutorial/collision](https://jonathanwhiting.com/tutorial/collision) - box-tilemap (TileCollider)