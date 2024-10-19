package fake

import (
	"context"
	"fmt"

	"github.com/samborkent/fake/external"
)

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

func (f *backExpect) Get() *backGet {
	if f == nil {
		return nil
	}
	f.get = append(f.get, &backGet{})
	return f.get[len(f.get)-1]
}

type backGet struct {
}
