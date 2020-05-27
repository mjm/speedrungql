package resolvers

import (
  "context"
  "encoding/json"
  "net/http"

  "github.com/mjm/speedrungql"
)

type Resolvers struct {
  baseURL    string
  httpClient http.Client
}

func New(baseURL string) *Resolvers {
  return &Resolvers{
    baseURL:    baseURL,
    httpClient: http.Client{},
  }
}

func (r *Resolvers) Viewer() *Viewer {
  return &Viewer{r: r}
}

type Viewer struct {
  r *Resolvers
}

func (v *Viewer) Platforms(ctx context.Context, args struct{
  Order *struct{
    Field *string
    Direction *string
  }
  First *int32
  After *Cursor
}) (*PlatformConnection, error) {
  u := v.r.baseURL + "/platforms"
  res, err := v.r.httpClient.Get(u)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  var resp speedrungql.PlatformsResponse
  if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
    return nil, err
  }

  return &PlatformConnection{res: &resp}, nil
}

type PlatformConnection struct {
  res *speedrungql.PlatformsResponse
}

func (pc *PlatformConnection) Edges() []*PlatformEdge {
  var edges []*PlatformEdge
  for _, p := range pc.res.Data {
    edges = append(edges, &PlatformEdge{
      Node: &Platform{p},
    })
  }
  return edges
}

func (pc *PlatformConnection) Nodes() []*Platform {
  var nodes []*Platform
  for _, p := range pc.res.Data {
    nodes = append(nodes, &Platform{p})
  }
  return nodes
}

func (pc *PlatformConnection) PageInfo() PageInfo {
  return PageInfo{}
}

type PlatformEdge struct {
  Node *Platform
}

func (pe *PlatformEdge) Cursor() *Cursor {
  return nil
}

type Platform struct {
  speedrungql.Platform
}

