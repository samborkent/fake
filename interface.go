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
	GetExternal(ctx context.Context, externalID int) (external.External, error)
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
