package repository

import "context"

type ScopeRepository interface {
	GetAllScope(ctx context.Context, scopeNames []string)
}

// GetAllScope implements Repository.
func (r *dbRepository) GetAllScope(ctx context.Context, scopeNames []string) {
	panic("unimplemented")
}
