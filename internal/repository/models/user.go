package models

import "github.com/rs/xid"

type (
	User struct {
		Id        xid.ID `db:"id"`
		Firstname string `db:"firstname"`
		Lastname  string `db:"lastname"`
		Username  string `db:"username"`
		Hash      string `db:"hash"`
		Role      int    `db:"role"`
		Email     string `db:"email"`
		Verified  bool   `db:"verified"`
		CreatedAt int64  `db:"created_at"`
		UpdatedAt int64  `db:"updated_at"`
	}
)
