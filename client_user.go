package speedrungql

import (
	"context"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) newUserLoader() *dataloader.Loader {
	return c.newLoader(func(key dataloader.Key) string {
		return "/users/" + key.String()
	})
}

func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := c.loadItem(ctx, c.userLoader, userID, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
