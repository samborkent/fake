package fakes

import (
	"context"

	"github.com/samborkent/fake/external"
)

const (
	ID         = "testID"
	ExternalID = 1
)

func HandleStuff(ctx context.Context, getter Getter) (Object, external.External, error) {
	object, err := getter.Get(ctx, ID)
	if err != nil {
		return Object{}, external.External{}, err
	}

	ext, err := getter.GetExternal(ctx, ExternalID)
	if err != nil {
		return Object{}, external.External{}, err
	}

	return object, *ext, nil
}
