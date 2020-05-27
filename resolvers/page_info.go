package resolvers

type PageInfo struct {
  StartCursor     *Cursor
  EndCursor       *Cursor
  HasNextPage     bool
  HasPreviousPage bool
}
