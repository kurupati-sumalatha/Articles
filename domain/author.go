package domain

import "context"

// Author representing the Author data struct
type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (Author, error)
}
