package speedrun

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) ListGenres(ctx context.Context, opts ...FetchOption) ([]*Genre, *PageInfo, error) {
	var resp GenresResponse
	if err := c.fetch(ctx, "/genres", &resp, opts...); err != nil {
		return nil, nil, err
	}
	return resp.Data, resp.Pagination, nil
}

func (c *Client) GetGenre(ctx context.Context, genreID string) (*Genre, error) {
	var genre Genre
	if err := c.loadItem(ctx, c.genreKey(genreID), &genre); err != nil {
		return nil, err
	}
	return &genre, nil
}

func (c *Client) GetGenres(ctx context.Context, ids []string) ([]*Genre, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(c.genreKey(id)))
	}

	ress, errs := c.loader.LoadMany(ctx, keys)()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var genres []*Genre
	for _, res := range ress {
		var genre Genre
		if err := json.Unmarshal(res.(*EnvelopeResponse).Data, &genre); err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}
	return genres, nil
}

func (c *Client) genreKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/genres/%s", c.BaseURL, id)
}
