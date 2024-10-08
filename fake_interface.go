package fake

import (
	"context"
	"testing"

	"github.com/samborkent/fake/external"
)

var (
	errFakeNotInitialized   = "fake: not initialized"
	errMethodNotInitialized = "fake: '%s.%s': method not initialized"
	errContextNil           = "fake: '%s.%s': nil context"
	errContextCancel        = "fake: '%s.%s': %s"
	errExpectationsMissing  = "fake: '%s.%s': expectations missing: called '%d' time(s), '%d' expectation(s) registered"
	errMismatch             = "fake: '%s.%s': '%s' mismatch: got '%+v', want '%+v'"
	errReturnMissing        = "fake: '%s.%s': missing return value(s) for expectation '%d'"
)

const getterName = "Getter"

type fakeGetter struct {
	t  *testing.T
	On *getterExpect
}

type getterExpect struct {
	get        []*getterGet
	getCounter int

	getExternal        []*getterGetExternal
	getExternalCounter int
}

func NewFakeGetter(t *testing.T) *fakeGetter {
	return &fakeGetter{
		t: t,
		On: &getterExpect{
			get:         make([]*getterGet, 0),
			getExternal: make([]*getterGetExternal, 0),
		},
	}
}

func (f *fakeGetter) Get(ctx context.Context, id string) (Object, error) {
	if f == nil || f.t == nil {
		panic(errFakeNotInitialized)
	}

	const methodName = "Get"

	if f.On == nil || f.On.get == nil {
		f.t.Fatalf(errMethodNotInitialized, getterName, methodName)
	}

	if ctx == nil {
		f.t.Fatalf(errContextNil, getterName, methodName)
	}

	if err := ctx.Err(); err != nil {
		f.t.Fatalf(errContextCancel, getterName, methodName, err.Error())
	}

	index := f.On.getCounter

	if index+1 > len(f.On.get) {
		f.t.Fatalf(errExpectationsMissing, getterName, methodName, index+1, len(f.On.get))
	}

	if id != f.On.get[index].id {
		f.t.Fatalf(errMismatch, getterName, methodName, "id", id, f.On.get[index].id)
	}

	if f.On.get[index].returns == nil {
		f.t.Fatalf(errReturnMissing, getterName, methodName, index+1)
	}

	f.On.getCounter++

	return f.On.get[index].returns.object, f.On.get[index].returns.err
}

func (e *getterExpect) Get(id string) *getterGet {
	if e == nil {
		return nil
	}

	e.get = append(e.get, &getterGet{
		id: id,
	})

	return e.get[len(e.get)-1]
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

	f.returns = &getterGetReturn{
		object: object,
		err:    err,
	}
}

func (f *fakeGetter) GetExternal(ctx context.Context, externalID int) (external.External, error) {
	if f == nil || f.t == nil {
		panic(errFakeNotInitialized)
	}

	const methodName = "GetExternal"

	if f.On == nil || f.On.getExternal == nil {
		f.t.Fatalf(errMethodNotInitialized, getterName, methodName)
	}

	if ctx == nil {
		f.t.Fatalf(errContextNil, getterName, methodName)
	}

	if err := ctx.Err(); err != nil {
		f.t.Fatalf(errContextCancel, getterName, methodName, err.Error())
	}

	index := f.On.getExternalCounter

	if index+1 > len(f.On.getExternal) {
		f.t.Fatalf(errExpectationsMissing, getterName, methodName, index+1, len(f.On.getExternal))
	}

	if externalID != f.On.getExternal[index].externalID {
		f.t.Fatalf(errMismatch, getterName, methodName, "externalID", externalID, f.On.getExternal[index].externalID)
	}

	if f.On.getExternal[index].returns == nil {
		f.t.Fatalf(errReturnMissing, getterName, methodName, index+1)
	}

	f.On.getExternalCounter++

	return f.On.getExternal[index].returns.externalExternal,
		f.On.getExternal[index].returns.err
}

func (e *getterExpect) GetExternal(externalID int) *getterGetExternal {
	if e == nil {
		return nil
	}

	e.getExternal = append(e.getExternal, &getterGetExternal{
		externalID: externalID,
	})

	return e.getExternal[len(e.getExternal)-1]
}

type getterGetExternal struct {
	externalID int
	returns    *getterGetExternalReturn
}

type getterGetExternalReturn struct {
	externalExternal external.External
	err              error
}

func (f *getterGetExternal) Return(externalExternal external.External, err error) {
	if f == nil {
		return
	}

	f.returns = &getterGetExternalReturn{
		externalExternal: externalExternal,
		err:              err,
	}
}
