package fake

import "testing"

var _ Back = &fakeBack{}

func NewFakeBack(t *testing.T) *fakeBack {
	return &fakeBack{t: t, On: &backExpect{get: make([]*backGet, 0)}}
}

type fakeBack struct {
	t  *testing.T
	On *backExpect
}
type backExpect struct {
	get        []*backGet
	getCounter int
}

const backName = "Back"

type backGet struct{}
