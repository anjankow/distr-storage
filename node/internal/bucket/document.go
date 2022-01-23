package bucket

import (
	"encoding/json"
	"time"
)

type Document struct {
	ID        string          `bson:"_id" json:"id"`
	Timestamp time.Time       `bson:"timestamp" json:"timestamp"`
	Content   json.RawMessage `bson:"content" json:"content"`
}
