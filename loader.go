package speedrungql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) newLoader(pathFn func(dataloader.Key) string) *dataloader.Loader {
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var wg sync.WaitGroup
		results := make([]*dataloader.Result, len(keys))

		for i, key := range keys {
			wg.Add(1)
			go func(i int, key dataloader.Key) {
				defer wg.Done()

				res, err := http.Get(c.BaseURL + pathFn(key))
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

func (c *Client) loadItem(ctx context.Context, loader *dataloader.Loader, id string, result interface{}) error {
	res, err := loader.Load(ctx, dataloader.StringKey(id))()
	if err != nil {
		return err
	}

	data := res.(*EnvelopeResponse).Data
	return json.Unmarshal(data, result)
}
