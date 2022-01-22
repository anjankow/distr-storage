package tester

import (
	"client/internal/client"
	"client/internal/config"
	"fmt"
	"os"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Tester struct {
	logger *zap.Logger
	client *client.Client

	storageSystemProcess *os.Process

	report reportInfo
}

type reportInfo struct {
	timeStarted     time.Time
	totalNumOfNodes int
	nodesInserted   []inserted
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
	if err := t.startStorageSystem(); err != nil {
		t.logger.Panic("storage system not started, failure!", zap.Error(err))
		return
	}

	t.insertAllKeys(inputData)

	t.randomlyRead()

	t.randomlyDelete()

	t.stop()

}

func (t *Tester) startStorageSystem() error {
	numberOfNodes := config.GetNumberOfNodes()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	outFile, ferr1 := os.Create("storage-system-stdout.log")
	errFile, ferr2 := os.Create("storage-system-stdout.log")

	if err := multierr.Append(ferr1, ferr2); err != nil {
		return err
	}

	var attr = os.ProcAttr{
		Dir: pwd + "/..",
		Env: os.Environ(),
		Files: []*os.File{
			nil,
			outFile,
			errFile,
		},
	}
	process, err := os.StartProcess("start_db.sh", []string{fmt.Sprint(numberOfNodes)}, &attr)
	if err != nil {
		return err
	}

	t.storageSystemProcess = process
	t.report.timeStarted = time.Now()
	t.report.totalNumOfNodes = numberOfNodes
	return nil
}

func (t *Tester) insertAllKeys(inputData map[string]interface{}) {
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

		t.report.nodesInserted = append(t.report.nodesInserted, insertedInfo)

	}
}

func (t *Tester) randomlyRead() {
	time.Sleep(8 * time.Second)
}

func (t *Tester) randomlyDelete() {

}

func (t *Tester) stop() {
	if err := t.storageSystemProcess.Kill(); err != nil {
		t.logger.Error("failed to kill the storage system", zap.Error(err))
	}
}

func (t Tester) GenerateReport() string {
	var report string
	report += "# Start\n"
	report += fmt.Sprintf("System started at %s\n", t.report.timeStarted.Format(time.RFC3339))
	report += fmt.Sprintf("Created %d nodes\n\n", t.report.totalNumOfNodes)

	report += "# Inserting the elements\n"
	for i, info := range t.report.nodesInserted {
		report += fmt.Sprintf("%d) key: %s, node: %s, timestamp: %s\n", i, info.Key, info.NodeName, info.Timestamp.Format(time.RFC3339))
	}

	return report
}
