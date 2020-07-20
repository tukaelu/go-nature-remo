package cloud

import (
	"context"
	"net/url"
)

// User represents user data.
type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

// Users provides interface of /users end-point.
type Users interface {
	GetMe(ctx context.Context) (*User, error)
	UpdateMe(ctx context.Context, u *User) (*User, error)
}

type users struct {
	cli *Client
}

// GetMe provides implementation of GET /users/me
// https://swagger.nature.global/#/default/get_1_users_me
func (api *users) GetMe(ctx context.Context) (*User, error) {
	var u *User
	if err := api.cli.Get(ctx, "users/me", nil, u); err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateMe provides implementation of POST /users/me
// https://swagger.nature.global/#/default/post_1_users_me
func (api *users) UpdateMe(ctx context.Context, me *User) (*User, error) {
	p := url.Values{}
	p.Set("nickname", me.Nickname)

	var u *User
	if err := api.cli.Post(ctx, "users/me", p, u); err != nil {
		return nil, err
	}
	return u, nil
}
