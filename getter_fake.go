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

func (f *getterExpect) Get(id string) *getterGet {
	if f == nil {
		return nil
	}
	f.get = append(f.get, &getterGet{id: id})
	return f.get[len(f.get)-1]
}

type getterGet struct {
	id      string
	returns *getterGetReturn
}
type getterGetReturn struct {
	object Object
	err    error
}

func (f *getterGet) Return(object Object, err error) {
	if f == nil {
		return
	}
	f.returns = &getterGetReturn{object: object, err: err}
}

func (f *getterExpect) Put(object external.External) *getterPut {
	if f == nil {
		return nil
	}
	f.put = append(f.put, &getterPut{object: object})
	return f.put[len(f.put)-1]
}

type getterPut struct {
	object external.External
}

func (f *getterExpect) Update(object *external.External) *getterUpdate {
	if f == nil {
		return nil
	}
	f.update = append(f.update, &getterUpdate{object: object})
	return f.update[len(f.update)-1]
}

type getterUpdate struct {
	object *external.External
}

func (f *getterExpect) GetExternal(externalID int) *getterGetExternal {
	if f == nil {
		return nil
	}
	f.getExternal = append(f.getExternal, &getterGetExternal{externalID: externalID})
	return f.getExternal[len(f.getExternal)-1]
}

type getterGetExternal struct {
	externalID int
	returns    *getterGetExternalReturn
}
type getterGetExternalReturn struct {
	object *external.External
	err    error
}

func (f *getterGetExternal) Return(object *external.External, err error) {
	if f == nil {
		return
	}
	f.returns = &getterGetExternalReturn{object: object, err: err}
}
