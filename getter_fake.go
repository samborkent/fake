package fake

import (
	"testing"

	"github.com/samborkent/fake/external"
)

var _ Getter = &fakeGetter{}

func NewFakeGetter(t *testing.T) *fakeGetter {
	return &fakeGetter{t: t, On: &getterExpect{get: make([]*getterGet, 0), put: make([]*getterPut, 0), update: make([]*getterUpdate, 0), getExternal: make([]*getterGetExternal, 0)}}
}

type fakeGetter struct {
	t  *testing.T
	On *getterExpect
}
type getterExpect struct {
	get                []*getterGet
	getCounter         int
	put                []*getterPut
	putCounter         int
	update             []*getterUpdate
	updateCounter      int
	getExternal        []*getterGetExternal
	getExternalCounter int
}

const getterName = "Getter"

type getterGet struct {
	id      string
	returns *getterGetReturn
}
type getterGetReturn struct {
	object Object
	err    error
}

func (f *getterGet) Return(object Object, err error)

type getterPut struct {
	object external.External
}
type getterUpdate struct {
	object *external.External
}
type getterGetExternal struct {
	externalID int
	returns    *getterGetExternalReturn
}
type getterGetExternalReturn struct {
	object *external.External
	err    error
}

func (f *getterGetExternal) Return(object *external.External, err error)
