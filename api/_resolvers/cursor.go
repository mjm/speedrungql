package resolvers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
)

type Cursor string

func (Cursor) ImplementsGraphQLType(name string) bool {
	return name == "Cursor"
}

func (c Cursor) MarshalJSON() ([]byte, error) {
	s := base64.URLEncoding.EncodeToString([]byte(c))
	return json.Marshal(s)
}

func (c *Cursor) UnmarshalGraphQL(input interface{}) error {
	if s, ok := input.(string); ok {
		decoded, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			return err
		}

		*c = Cursor(decoded)
		return nil
	}

	return fmt.Errorf("cursor was not a string")
}

func (c *Cursor) GetOffset() (int, error) {
	return strconv.Atoi(string(*c))
}
