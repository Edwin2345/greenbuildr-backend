package graph

import "github.com/greenbuildr/shared/gqlerrors"

func toGQLError(err error) error {
	return gqlerrors.ToGQLError(err)
}
