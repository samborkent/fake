package fake

import (
	"testing"
)

var _ Back = &fakeBack{}

func NewFakeBack(t *testing.T) *fakeBack {
	return &fakeBack{
		t: t,
		On: &backExpect{
			get: make([]*backGet, 0),
		},
	}
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
		panic(errFakeNotInitialized)
	}

	const methodName = "Get"

	if f.On == nil || f.On.get == nil {
		f.t.Fatalf(errMethodNotInitialized, backName, methodName)
	}

	index := f.On.getCounter

	if index+1 > len(f.On.get) {
		f.t.Fatalf(errExpectationsMissing, backName, methodName, index+1, len(f.On.get))
	}

	f.On.getCounter++
}

func (e *backExpect) Get() *backGet {
	if e == nil {
		return nil
	}

	e.get = append(e.get, &backGet{})

	return e.get[len(e.get)-1]
}
