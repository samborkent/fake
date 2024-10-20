package fake

import (
	"context"
	"fmt"

	"github.com/samborkent/fake/external"
)

// Object is an object.
type Object struct{}

type Back interface{ Get() }

type Getter interface {
	Get(ctx context.Context, id string) (Object, error)
	Put(object external.External)
	Update(object *external.External)
	GetExternal(ctx context.Context, externalID int) (object *external.External, err error)
}

func fooBar() {
	if true {
		return
	}

	return
}

func checkItOut(
	hello string,
) error {
	return fmt.Errorf("%s", hello)
}
