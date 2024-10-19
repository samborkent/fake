package fake

import (
	"context"

	"github.com/samborkent/fake/external"
)

func (f *getterGet) Return(object Object, err error) {
	if f == nil {
		return
	}

	f.returns = &getterGetReturn{
		object: object,
		err:    err,
	}
}

func (f *fakeGetter) GetExternal(ctx context.Context, externalID int) (*external.External, error) {
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
		f.t.Fatalf(errArgumentMismatch, getterName, methodName, "externalID", externalID, f.On.getExternal[index].externalID)
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
	externalExternal *external.External
	err              error
}

func (f *getterGetExternal) Return(externalExternal *external.External, err error) {
	if f == nil {
		return
	}

	f.returns = &getterGetExternalReturn{
		externalExternal: externalExternal,
		err:              err,
	}
}
