package speedrungql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) newLoader() *dataloader.Loader {
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var wg sync.WaitGroup
		results := make([]*dataloader.Result, len(keys))

		for i, key := range keys {
			wg.Add(1)
			go func(i int, key dataloader.Key) {
				defer wg.Done()

				u := key.String()
				if !strings.HasPrefix(u, c.BaseURL) {
					u = c.BaseURL + u
				}

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
				if err != nil {
					results[i] = &dataloader.Result{Error: err}
					return
				}

				res, err := c.HTTPClient.Do(req)
				if err != nil {
					results[i] = &dataloader.Result{Error: err}
					return
				}
				defer res.Body.Close()

				if res.StatusCode > 299 {
					results[i] = &dataloader.Result{
						Error: fmt.Errorf("Unexpected status code for key %q: %d", key, res.StatusCode),
					}
					return
				}

				var resp EnvelopeResponse
				if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
					results[i] = &dataloader.Result{Error: err}
					return
				}

				results[i] = &dataloader.Result{Data: &resp}
			}(i, key)
		}

		wg.Wait()
		return results
	})
}

func (c *Client) loadItem(ctx context.Context, path string, result interface{}) error {
	res, err := c.loader.Load(ctx, dataloader.StringKey(path))()
	if err != nil {
		return err
	}

	data := res.(*EnvelopeResponse).Data
	return json.Unmarshal(data, result)
}
