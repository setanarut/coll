[![GoDoc](https://godoc.org/github.com/setanarut/coll?status.svg)](https://pkg.go.dev/github.com/setanarut/coll)

# coll - 2d collision library for Go

Features

* Collisions only - no gravity, rigid body handling, or complex solvers
* Is data-oriented and functional

## Conventions

"Sweep" tests indicate at least 1 of the objects is moving. 
The number indicates how many objects are moving. e.g., `box-box-sweep2` means we are comparing 2 aabbs, both of which are moving.
"Overlap" tests don't take movement into account, and this is a static check to see if the 2 entities overlap.
plural forms imply a collection. e.g., `BoxSegmentsSweep1Indexed()` checks one box segment against a set of line segments.
If there is more than one collision, the closest collision is set in the `h *Hit` argument.

## Visualization of some (but not all) functions

### Box-Box overlap

![Box-Box-overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-overlap.png)

### Box-Box sweep 1

![Box-Box sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-sweep1.png)

### Box-Box sweep 2

![Box-Box sweep 2](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-aabb-sweep2.png)

### Box-OrientedBox overlap

![abb-obb-overlap](https://github.com/user-attachments/assets/972ae63f-6ed0-4d10-a893-ab80f0f59f00)

### Box-Segment sweep 1

![Box-Segment sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segment-sweep1.png)

### Box-Segments sweep 1 indexed

![Box-Segments sweep 1](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segments-sweep1-indexed.png)

### Box-Segment overlap

![Box-Segment overlap](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/aabb-segment-overlap.png)

### Line-Circle overlap

![alt text](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/ray-circle-overlap.png)

### Segment-Circle overlap

![alt text](https://raw.githubusercontent.com/mreinstein/collision-2d/refs/heads/main/docs/segment-circle-overlap.png)

## Examples (Ebitengine)

1. Clone this repository
2. In the terminal, change to the examples directory `cd examples`
3. Run a folder with `go run ./foldername`. Example: `go run ./box_box_sweep1 `

## Credits

Most of these collision checks were adapted from existing open source repos:

* [github.com/mreinstein/collision-2d](https://github.com/mreinstein/collision-2d)
* [youtube.com/watch?v=NbSee-XM7WA](https://youtube.com/watch?v=NbSee-XM7WA) - ray-tilemap (RayTilemapDDA)
* [jonathanwhiting.com/tutorial/collision](https://jonathanwhiting.com/tutorial/collision) - box-tilemap (TileCollider)
