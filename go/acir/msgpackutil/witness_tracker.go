package msgpackutil

// witnessTracker records the highest witness index observed during decode.
type witnessTracker struct {
	maxWitness  uint32
	witnessSeen bool
}

// ObserveWitness updates the high-water mark of witness indices seen so far.
func (t *witnessTracker) ObserveWitness(idx uint32) {
	if !t.witnessSeen || idx > t.maxWitness {
		t.maxWitness = idx
		t.witnessSeen = true
	}
}

// MaxWitness returns the highest witness index observed since the last reset.
// `ok` is false if no witnesses were decoded.
func (t *witnessTracker) MaxWitness() (uint32, bool) {
	return t.maxWitness, t.witnessSeen
}

// ResetWitnessTracker zeros the witness high-water mark.
func (t *witnessTracker) ResetWitnessTracker() {
	t.maxWitness = 0
	t.witnessSeen = false
}
