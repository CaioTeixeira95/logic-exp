package repositories

import (
	"database/sql"
)

type DefaultRepository struct {
	db *sql.DB
}

type DefaultRepositoryOption func(r *DefaultRepository)

func NewRepository(options ...DefaultRepositoryOption) *DefaultRepository {
	r := &DefaultRepository{}

	for _, option := range options {
		option(r)
	}

	return r
}

func WithDatabaseOption(db *sql.DB) DefaultRepositoryOption {
	return func(r *DefaultRepository) {
		r.db = db
	}
}

var _ ExpressionRepository = (*DefaultRepository)(nil)
