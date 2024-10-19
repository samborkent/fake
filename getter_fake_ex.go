package fake

import (
	"context"

	"github.com/samborkent/fake/external"
)

func (f *fakeGetter) Get(ctx context.Context, id string) (Object Object, err error) {
	if f == nil || f.t == nil {
		panic(errFakeNotInitialized)
	}

	const methodName = "Get"

	if f.On == nil || f.On.get == nil {
		f.t.Fatalf(errMethodNotInitialized, getterName, methodName)
	}

	index := f.On.getCounter

	if index+1 > len(f.On.get) {
		f.t.Fatalf(errExpectationsMissing, getterName, methodName, index+1, len(f.On.get))
	}

	if ctx == nil {
		f.t.Fatalf(errContextNil, getterName, methodName)
	}

	if err := ctx.Err(); err != nil {
		f.t.Fatalf(errContextCancel, getterName, methodName, err.Error())
	}

	if id != f.On.get[index].id {
		f.t.Fatalf(errArgumentMismatch, getterName, methodName, "id", id, f.On.get[index].id)
	}

	f.On.getCounter++

	return f.On.get[index].returns.object, f.On.get[index].returns.err
}

func (f *fakeGetter) Put(object external.External) {
	if f == nil || f.t == nil {
		panic(errFakeNotInitialized)
	}

	const methodName = "Put"

	if f.On == nil || f.On.put == nil {
		f.t.Fatalf(errMethodNotInitialized, getterName, methodName)
	}

	index := f.On.putCounter

	if index+1 > len(f.On.put) {
		f.t.Fatalf(errExpectationsMissing, getterName, methodName, index+1, len(f.On.put))
	}

	if object != f.On.put[index].object {
		f.t.Fatalf(errArgumentMismatch, getterName, methodName, "object", object, f.On.put[index].object)
	}

	f.On.putCounter++
}

func (f *fakeGetter) Update(object *external.External) {
	if f == nil || f.t == nil {
		panic(errFakeNotInitialized)
	}

	const methodName = "Update"

	if f.On == nil || f.On.update == nil {
		f.t.Fatalf(errMethodNotInitialized, getterName, methodName)
	}

	index := f.On.updateCounter

	if index+1 > len(f.On.update) {
		f.t.Fatalf(errExpectationsMissing, getterName, methodName, index+1, len(f.On.update))
	}

	if object != f.On.update[index].object {
		f.t.Fatalf(errArgumentMismatch, getterName, methodName, "object", object, f.On.update[index].object)
	}

	f.On.updateCounter++
}

func (f *fakeGetter) GetExternal(ctx context.Context, externalID int) (object *external.External, err error) {
	if f == nil || f.t == nil {
		panic(errFakeNotInitialized)
	}

	const methodName = "GetExternal"

	if f.On == nil || f.On.getExternal == nil {
		f.t.Fatalf(errMethodNotInitialized, getterName, methodName)
	}

	index := f.On.getExternalCounter

	if index+1 > len(f.On.getExternal) {
		f.t.Fatalf(errExpectationsMissing, getterName, methodName, index+1, len(f.On.getExternal))
	}

	if ctx == nil {
		f.t.Fatalf(errContextNil, getterName, methodName)
	}

	if err := ctx.Err(); err != nil {
		f.t.Fatalf(errContextCancel, getterName, methodName, err.Error())
	}

	if externalID != f.On.getExternal[index].externalID {
		f.t.Fatalf(errArgumentMismatch, getterName, methodName, "externalID", externalID, f.On.getExternal[index].externalID)
	}

	f.On.getExternalCounter++

	return f.On.getExternal[index].returns.object, f.On.getExternal[index].returns.err
}
