package module

import (
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Utility function during module invocation.
func unwrapError(rerr error) string {
	var gqlErr *gqlerror.Error
	if errors.As(rerr, &gqlErr) {
		return gqlErr.Message
	}
	
	return rerr.Error()
}