package shared

import (
	"sunspot/go/acir/msgpackutil"
	"testing"
)

// ConsumeEnumTag opens the outer single-entry-fixmap of an enum-shaped
// fixture value, returning the variant tag.
func ConsumeEnumTag(t *testing.T, r *msgpackutil.Reader) int {
	n, err := r.ReadMapLen()
	if err != nil {
		t.Fatalf("expected enum fixmap: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected single-entry enum fixmap, got %d entries", n)
	}
	tag, err := r.ReadUint()
	if err != nil {
		t.Fatalf("failed to read enum tag: %v", err)
	}
	return int(tag)
}
