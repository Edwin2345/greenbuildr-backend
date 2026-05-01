package gqlerrors

import (
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CodedError is implemented by domain errors that carry a machine-readable code.
type CodedError interface {
	error
	ErrorCode() string
	ErrorMessage() string
}

// ToGQLError translates a CodedError into a GraphQL error with extensions.code set.
// Non-coded errors are returned unchanged.
func ToGQLError(err error) error {
	var ce CodedError
	if errors.As(err, &ce) {
		return &gqlerror.Error{
			Message: ce.ErrorMessage(),
			Extensions: map[string]any{
				"code": ce.ErrorCode(),
			},
		}
	}
	return err
}
