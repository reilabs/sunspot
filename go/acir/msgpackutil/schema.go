package msgpackutil

import "fmt"

// Field is the key passed to a decode callback.
type Field struct {
	Tag  int
	Name string
}

// Renders renders a Field for diagnostic messages.
func (f Field) String() string {
	if f.Name != "" {
		return fmt.Sprintf("tag=%d name=%q", f.Tag, f.Name)
	}
	return fmt.Sprintf("tag=%d", f.Tag)
}

// Schema maps serde field/variant names to their int tags so that
// string-keyed wire formats  can dispatch through the same tag-switch as MsgpackTagged.
type Schema struct {
	nameToTag map[string]int
}

func NewSchema(nameToTag map[string]int) Schema {
	return Schema{nameToTag: nameToTag}
}

// TagFor resolves a serde field/variant name to its int tag. Returns
// (-1, false) if the schema has no entry or is the zero value.
func (s Schema) TagFor(name string) (int, bool) {
	if s.nameToTag == nil {
		return -1, false
	}
	t, ok := s.nameToTag[name]
	return t, ok
}
