package store

import "database/sql"

type Store struct {
	User UserRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		User: &UserStore{db: db},
	}
}