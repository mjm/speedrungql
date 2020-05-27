package speedrungql

type PlatformsResponse struct {
  Data []Platform `json:"data"`
}

type Platform struct {
  ID       string `json:"id"`
  Name     string `json:"name"`
  Released int32  `json:"released"`
}
