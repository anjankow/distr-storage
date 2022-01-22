package tester

import (
	"client/internal/client"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Tester struct {
	logger *zap.Logger
	client *client.Client

	nodesInserted []inserted
}

type inserted struct {
	NodeName  string
	Key       string
	Timestamp time.Time
}

func NewTester(logger *zap.Logger, clnt *client.Client) Tester {
	return Tester{
		logger: logger,
		client: clnt,
	}
}

func (t *Tester) Run(inputData map[string]interface{}) {

	for key, val := range inputData {
		t.logger.Debug("test data", zap.String("key", key), zap.Any("value", val))
		insertTime := time.Now()

		node, err := t.client.Insert(key, val)

		if err != nil {

			t.logger.Warn("failed to insert: "+err.Error(), zap.String("key", key))

			continue
		}

		insertedInfo := inserted{
			NodeName:  node,
			Key:       key,
			Timestamp: insertTime,
		}

		t.nodesInserted = append(t.nodesInserted, insertedInfo)

	}

}

func (t Tester) GenerateReport() string {
	var report string
	report += "# Inserting the elements\n"
	for i, info := range t.nodesInserted {
		report += fmt.Sprintf("%d) key: %s, node: %s, timestamp: %s\n", i, info.Key, info.NodeName, info.Timestamp.Format(time.RFC3339))
	}

	return report
}
