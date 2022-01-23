package bucket

import (
	"encoding/json"
	"time"
)

type Document struct {
	ID        string          `bson:"_id"`
	Timestamp time.Time       `bson:"timestamp"`
	Content   json.RawMessage `bson:"content"`
}
