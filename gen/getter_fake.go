package fakes

import (
	"context"
	"reflect"
	"testing"

	"github.com/samborkent/fake"
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

func (f *fakeGetter) Get(ctx context.Context, id string) (Object Object, err error) {
	if f == nil || f.t == nil {
		panic("fake: " + fake.ErrFakeNotInitialized.Error())
	}
	const methodName = "Get"
	if f.On == nil || f.On.get == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrMethodNotInitialized.Error())
	}
	index := f.On.getCounter
	if index+1 > len(f.On.get) {
		f.t.Fatalf("fake: '%s.%s': %s: called '%d' time(s), '%d' expectation(s) registered", getterName, methodName, fake.ErrExpectationsMissing.Error(), index+1, len(f.On.get))
	}
	if ctx == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrContextNil.Error())
	}
	if err := ctx.Err(); err != nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, err.Error())
	}
	if !reflect.DeepEqual(id, f.On.get[index].id) {
		f.t.Fatalf("fake: '%s.%s': %s: '%s': got '%+v', want '%+v'", getterName, methodName, fake.ErrArgumentMismatch.Error(), "id", id, f.On.get[index].id)
	}
	if f.On.get[index].returns == nil {
		f.t.Fatalf("fake: '%s.%s': %s: '%d'", getterName, methodName, fake.ErrReturnMissing.Error(), index+1)
	}
	f.On.getCounter++
	return f.On.get[index].returns.object, f.On.get[index].returns.err
}
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
func (f *fakeGetter) Put(object external.External) {
	if f == nil || f.t == nil {
		panic("fake: " + fake.ErrFakeNotInitialized.Error())
	}
	const methodName = "Put"
	if f.On == nil || f.On.put == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrMethodNotInitialized.Error())
	}
	index := f.On.putCounter
	if index+1 > len(f.On.put) {
		f.t.Fatalf("fake: '%s.%s': %s: called '%d' time(s), '%d' expectation(s) registered", getterName, methodName, fake.ErrExpectationsMissing.Error(), index+1, len(f.On.put))
	}
	if !reflect.DeepEqual(object, f.On.put[index].object) {
		f.t.Fatalf("fake: '%s.%s': %s: '%s': got '%+v', want '%+v'", getterName, methodName, fake.ErrArgumentMismatch.Error(), "object", object, f.On.put[index].object)
	}
	f.On.putCounter++
	return
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

func (f *fakeGetter) Update(object *external.External) {
	if f == nil || f.t == nil {
		panic("fake: " + fake.ErrFakeNotInitialized.Error())
	}
	const methodName = "Update"
	if f.On == nil || f.On.update == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrMethodNotInitialized.Error())
	}
	index := f.On.updateCounter
	if index+1 > len(f.On.update) {
		f.t.Fatalf("fake: '%s.%s': %s: called '%d' time(s), '%d' expectation(s) registered", getterName, methodName, fake.ErrExpectationsMissing.Error(), index+1, len(f.On.update))
	}
	if !reflect.DeepEqual(object, f.On.update[index].object) {
		f.t.Fatalf("fake: '%s.%s': %s: '%s': got '%+v', want '%+v'", getterName, methodName, fake.ErrArgumentMismatch.Error(), "object", object, f.On.update[index].object)
	}
	f.On.updateCounter++
	return
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

func (f *fakeGetter) GetExternal(ctx context.Context, externalID int) (object *external.External, err error) {
	if f == nil || f.t == nil {
		panic("fake: " + fake.ErrFakeNotInitialized.Error())
	}
	const methodName = "GetExternal"
	if f.On == nil || f.On.getExternal == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrMethodNotInitialized.Error())
	}
	index := f.On.getExternalCounter
	if index+1 > len(f.On.getExternal) {
		f.t.Fatalf("fake: '%s.%s': %s: called '%d' time(s), '%d' expectation(s) registered", getterName, methodName, fake.ErrExpectationsMissing.Error(), index+1, len(f.On.getExternal))
	}
	if ctx == nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, fake.ErrContextNil.Error())
	}
	if err := ctx.Err(); err != nil {
		f.t.Fatalf("fake: '%s.%s': %s", getterName, methodName, err.Error())
	}
	if !reflect.DeepEqual(externalID, f.On.getExternal[index].externalID) {
		f.t.Fatalf("fake: '%s.%s': %s: '%s': got '%+v', want '%+v'", getterName, methodName, fake.ErrArgumentMismatch.Error(), "externalID", externalID, f.On.getExternal[index].externalID)
	}
	if f.On.getExternal[index].returns == nil {
		f.t.Fatalf("fake: '%s.%s': %s: '%d'", getterName, methodName, fake.ErrReturnMissing.Error(), index+1)
	}
	f.On.getExternalCounter++
	return f.On.getExternal[index].returns.object, f.On.getExternal[index].returns.err
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
