[![GoDoc](https://godoc.org/github.com/setanarut/coll?status.svg)](https://pkg.go.dev/github.com/setanarut/coll)

# coll - 2d collision library for Go

Features

* Collisions only - no gravity, rigid body handling, or complex solvers
* Is data-oriented and functional

## Conventions

"Sweep" tests indicate at least 1 of the objects is moving. 
The number indicates how many objects are moving. e.g., `aabb-aabb-sweep2` means we are comparing 2 aabbs, both of which are moving.
"Overlap" tests don't take movement into account, and this is a static check to see if the 2 entities overlap.
plural forms imply a collection. e.g., `segments-segment-ovelap` checks one line segment against a set of line segments.
If there is more than one collision, the closest collision is set in the `hitInfo` argument.

## Available collision checks

### AABB-AABB-overlap

![AABB-AABB-overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-overlap.png)

```go
AABBAABBOverlap(boxA, boxB *AABB, hitInfo *HitInfo) bool
```

### AABB-AABB-contain

```go
// AABBAABBContain returns true if a fully contains b.
AABBAABBContain(a, b *AABB) bool
```

### AABB-AABB sweep 1

![AABB-AABB sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-sweep1.png)

```go
AABBAABBSweep1(staticBoxA, boxB *AABB, boxBVel v.Vec, hitInfo *HitInfo) bool 
```

### AABB-AABB sweep 2

![AABB-AABB sweep 2](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-sweep2.png)

```go
AABBAABBSweep2(boxA, boxB *AABB, boxAVel, boxBVel v.Vec, hitInfo *HitInfo) bool
```

### AABB-Segment sweep 1

![AABB-Segment sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segment-sweep1.png)

```go
AABBSegmentSweep1(line *Segment, box *AABB, delta v.Vec, hitInfo *HitInfo) bool
```

### AABB-Segment sweep 1 indexed

![AABB-Segment sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segments-sweep1-indexed.png)

```go
AABBSegmentSweep1Indexed(lines []*Segment, aabb *AABB, delta v.Vec, hitInfo *HitInfo)  (index int)
```

### AABB-Point overlap

![AABB-Point overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-point-overlap.png)

```go
AABBPointOverlap(box *AABB, point v.Vec, hitInfo *HitInfo) bool
```

### AABB-Segment overlap

![AABB-Segment overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segment-overlap.png)

```go
AABBSegmentOverlap(box *AABB, start, delta, padding v.Vec, hitInfo *HitInfo) bool
```

### Ray-Circle overlap

![alt text](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/ray-sphere-overlap.png)

```go
RayCircleOverlap(raySeg *Segment, circ *Circle, overlapSeg *Segment) bool
```

### Segment-Circle overlap

![alt text](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/segment-sphere-overlap.png)

```go
SegmentCircleOverlap(seg *Segment, c *Circle) []v.Vec
```

### AABB-Circle sweep 2

```go
AABBCircleSweep2(box *AABB, circle *Circle, boxVel, circleVel v.Vec) bool
```

### AABB-Tilemap sweep

```go
(c *TileCollider) Collide(box AABB, delta v.Vec, onCollide TileCollisionCallback) v.Vec
```

### Ray-Tilemap overlap (DDA)

```go
RaycastDDA(pos, dir v.Vec, length float64, tm [][]uint8, cellSize float64, h *HitInfo) (bool, image.Point) 
```

## Examples (Ebitengine)

1. Clone this repository
2. In the terminal, change to the examples directory `cd examples`
3. Run a folder with `go run ./foldername`. Example: `go run ./aabb_aabb_sweep1 `

## Credits

Most of these collision checks were adapted from existing open source repos:

* [github.com/mreinstein/collision-2d](https://github.com/mreinstein/collision-2d)
* [youtube.com/watch?v=NbSee-XM7WA](https://youtube.com/watch?v=NbSee-XM7WA) - ray-tilemap (RaycastDDA)
* [jonathanwhiting.com/tutorial/collision](https://jonathanwhiting.com/tutorial/collision) - aabb-tilemap (TileCollider)