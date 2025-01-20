package utils

import (
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func UnwrapError(rerr error) string {
	var gqlErr *gqlerror.Error
	if errors.As(rerr, &gqlErr) {
		return gqlErr.Message
	}
	
	return rerr.Error()
}