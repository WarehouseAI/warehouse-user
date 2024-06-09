package domain

import (
	"github.com/warehouse/user-service/internal/repository/models"
	"github.com/warehouse/user-service/internal/warehousepb"
)

var (
	Roles = []Role{RoleAdmin, RoleUser}
)

type Role int64

const (
	RoleAdmin Role = iota
	RoleUser
)

type (
	User struct {
		Id        string
		Firstname string
		Lastname  string
		Username  string
		Email     string
		Hash      string
		Role      Role
		Verified  bool
		CreatedAt int64
		UpdatedAt int64
	}
)

func (u User) ToProto() *warehousepb.User {
	return &warehousepb.User{
		UserId:    u.Id,
		Role:      int64(u.Role),
		Username:  u.Username,
		Firstname: u.Firstname,
		Verified:  u.Verified,
		Email:     u.Email,
	}
}

func (u User) ToModel() models.User {
	return models.User{
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Username:  u.Username,
		Hash:      u.Hash,
		Role:      int(u.Role),
	}
}

func (User) FromModel(u models.User) User {
	return User{
		Id:        u.Id.String(),
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Username:  u.Username,
		Email:     u.Email,
		Hash:      u.Hash,
		Role:      Role(u.Role),
		Verified:  u.Verified,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
