package fake

import (
	"testing"

	"github.com/samborkent/fake"
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

func (f *fakeBack) Get() {
	if f == nil || f.t == nil {
		panic("fake: " + fake.ErrFakeNotInitialized.Error())
	}
	const methodName = "Get"
	if f.On == nil || f.On.get == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrMethodNotInitialized.Error())
	}
	index := f.On.getCounter
	if index+1 > len(f.On.get) {
		f.t.Fatalf("fake: '%s.%s': expectations missing: called '%d' time(s), '%d' expectation(s) registered", backName, methodName, index+1, len(f.On.get))
	}
	f.On.getCounter++
	return
}

func (f *backExpect) Get() *backGet {
	if f == nil {
		return nil
	}
	f.get = append(f.get, &backGet{})
	return f.get[len(f.get)-1]
}

type backGet struct{}
