# coll

There are many Go collision libraries for 2d. None satisifed all of these criteria:

* Collisions only - no gravity, rigid body handling, or complex solvers
* Is data-oriented and functional
* Consistent API interface

## Conventions

All collision checking functions return a bool indicating if there was a collision.
They also accept an optional `hitInfo` argument, which gets filled in if there is an actual collision.  

"Sweep" tests indicate at least 1 of the objects is moving.

The number indicates how many objects are moving. e.g., `aabb-aabb-sweep2` means we are comparing 2 aabbs, both of which are moving.

"Overlap" tests don't take movement into account, and this is a static check to see if the 2 entities overlap.

plural forms imply a collection. e.g., `segments-segment-ovelap` checks one line segment against a set of line segments. If there is more than one collision, the closest collision is set in the `hitInfo` argument.

"indexed" tests are the same as their non-indexed forms, except they take in an array of segment indices to use. These are nice in that you can avoid having to build large arrays of line segments every frame, if you have things like dynamic line segments (platforms) or have a spatial culling algorithm that selects line segments to include.

## Available collision checks (work in progress)

- [x] aabb-aabb overlap
- [x] aabb-aabb contain
- [x] aabb-aabb sweep 1
- [x] aabb-aabb sweep 2
- [x] aabb-segment sweep
- [ ] aabb-segments sweep-indexed
- [ ] aabb-point overlap
- [x] aabb-segment overlap
- [ ] ray-plane-distance
- [ ] ray-sphere overlap
- [ ] segment-sphere overlap
- [ ] segment-normal
- [ ] segment-point-overlap
- [ ] segment-segment-overlap
- [ ] segments-segment-overlap
- [ ] segments-segment-overlap-indexed
- [x] aabb-circle sweep
- [ ] segments-circle-sweep 1
- [ ] segments-circle-sweep-1-indexed
- [ ] circle-circle-overlap
- [ ] circle-circle-sweep2
- [ ] cone-point-overlap
- [ ] triangle-point-overlap
- [x] ray-tilemap
- [x] aabb-tilemap

## Credits

Most of these collision checks were adapted from existing open source repo:

* [github.com/mreinstein/collision-2d](https://github.com/mreinstein/collision-2d)
* [youtube.com/watch?v=NbSee-XM7WA](https://youtube.com/watch?v=NbSee-XM7WA) - ray-tilemap (RaycastDDA)
* https://jonathanwhiting.com/tutorial/collision - aabb-tilemap (TileCollider)