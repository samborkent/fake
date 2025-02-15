package test_test

import (
	"context"
	"testing"

	"github.com/samborkent/check"
	"github.com/samborkent/fake/external"
	fakes "github.com/samborkent/fake/gen"
)

func TestFake(t *testing.T) {
	t.Parallel()

	t.Run("missing expectations", func(t *testing.T) {
		t.Parallel()

		_, _, _ = fakes.HandleStuff(t.Context(), fakes.NewFakeGetter(t))
	})

	t.Run("nil context", func(t *testing.T) {
		t.Parallel()

		fakeGetter := fakes.NewFakeGetter(t)

		fakeGetter.On.Get("")

		_, _, _ = fakes.HandleStuff(nil, fakeGetter) //nolint:go-staticcheck
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()

		fakeGetter := fakes.NewFakeGetter(t)

		fakeGetter.On.Get("")

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, _, err := fakes.HandleStuff(ctx, fakeGetter)
		check.ErrorNil(t, err)
	})

	t.Run("argument mismatch", func(t *testing.T) {
		t.Parallel()

		fakeGetter := fakes.NewFakeGetter(t)

		fakeGetter.On.Get("")

		_, _, err := fakes.HandleStuff(t.Context(), fakeGetter)
		check.ErrorNil(t, err)
	})

	t.Run("no return value for get", func(t *testing.T) {
		t.Parallel()

		fakeGetter := fakes.NewFakeGetter(t)

		fakeGetter.On.Get("testID")

		_, _, err := fakes.HandleStuff(t.Context(), fakeGetter)
		check.ErrorNil(t, err)
	})

	t.Run("no return value for get external", func(t *testing.T) {
		t.Parallel()

		fakeGetter := fakes.NewFakeGetter(t)

		fakeGetter.On.Get("testID").Return(fakes.Object{}, nil)
		fakeGetter.On.GetExternal(1)

		_, _, err := fakes.HandleStuff(t.Context(), fakeGetter)
		check.ErrorNil(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		fakeGetter := fakes.NewFakeGetter(t)

		fakeGetter.On.Get("testID").Return(fakes.Object{}, nil)
		fakeGetter.On.GetExternal(1).Return(&external.External{}, nil)

		_, _, err := fakes.HandleStuff(t.Context(), fakeGetter)
		check.ErrorNil(t, err)
	})
}
