package resolvers

import (
  "errors"
  "strconv"
)

type Cursor string

func (Cursor) ImplementsGraphQLType(name string) bool {
  return name == "Cursor"
}

func (c *Cursor) UnmarshalGraphQL(input interface{}) error {
  var err error
  switch input := input.(type) {
  case string:
    *c = Cursor(input)
  default:
    err = errors.New("wrong type")
  }
  return err
}

func (c Cursor) MarshalJSON() ([]byte, error) {
  return strconv.AppendQuote(nil, string(c)), nil
}
