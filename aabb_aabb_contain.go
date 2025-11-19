package coll

// AABBAABBContain returns true if a fully contains b (b is fully inside of the bounds of a).
func AABBAABBContain(a, b *AABB) bool {
	if b.Left() < a.Left() {
		return false
	}
	if b.Right() > a.Right() {
		return false
	}
	if b.Top() < a.Top() {
		return false
	}
	if b.Bottom() > a.Bottom() {
		return false
	}
	return true
}
