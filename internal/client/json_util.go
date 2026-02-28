package client

import "encoding/json"

// OneOrMany handles the Cisco CUPI API quirk where a collection containing
// exactly one item is returned as a JSON object instead of a single-element
// array. When the result set has two or more items the normal JSON array form
// is used, so this type tries array first and falls back to a single object.
type OneOrMany[T any] []T

func (o *OneOrMany[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*o = nil
		return nil
	}
	var arr []T
	if err := json.Unmarshal(data, &arr); err == nil {
		*o = arr
		return nil
	}
	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*o = []T{obj}
	return nil
}
