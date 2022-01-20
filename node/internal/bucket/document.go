package bucket

import "encoding/json"

type Document struct {
	Content json.RawMessage `bson:"content"`
}
