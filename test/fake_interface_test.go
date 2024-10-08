package test_test

import (
	"context"
	"testing"

	"github.com/samborkent/check"
	"github.com/samborkent/fake"
	"github.com/samborkent/fake/external"
)

func TestFake(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		_, _, err := fake.HandleStuff(nil, fake.NewFakeGetter(t))
		check.ErrorNil(t, err)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := fake.HandleStuff(ctx, fake.NewFakeGetter(t))
		check.ErrorNil(t, err)
	})

	t.Run("missing expectations", func(t *testing.T) {
		_, _, err := fake.HandleStuff(context.TODO(), fake.NewFakeGetter(t))
		check.ErrorNil(t, err)
	})

	t.Run("no return value for get", func(t *testing.T) {
		fakeGetter := fake.NewFakeGetter(t)

		fakeGetter.On.Get("")

		_, _, err := fake.HandleStuff(context.TODO(), fakeGetter)
		check.ErrorNil(t, err)
	})

	t.Run("no return value for get external", func(t *testing.T) {
		fakeGetter := fake.NewFakeGetter(t)

		fakeGetter.On.Get("").Return(fake.Object{}, nil)
		fakeGetter.On.GetExternal(0)

		_, _, err := fake.HandleStuff(context.TODO(), fakeGetter)
		check.ErrorNil(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		fakeGetter := fake.NewFakeGetter(t)

		fakeGetter.On.Get("").Return(fake.Object{}, nil)
		fakeGetter.On.GetExternal(0).Return(external.External{}, nil)

		_, _, err := fake.HandleStuff(context.TODO(), fakeGetter)
		check.ErrorNil(t, err)
	})
}
