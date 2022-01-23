package tester

import (
	"client/internal/client"
	"client/internal/config"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	startScriptPath = "start_db.sh"
)

type Tester struct {
	logger *zap.Logger
	client *client.Client

	storageSystemProcess *os.Process
	totalNumOfNodes      int

	report reportInfo
}

type reportInfo struct {
	timeStarted      time.Time
	insertedToNodes  []nodeOperation
	readFromNodes    []nodeOperation
	deletedFromNodes []nodeOperation
	allInsertedData  string
}
type nodeOperation struct {
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
		t.stop()
		t.logger.Panic("storage system failure!", zap.Error(err))
		return
	}

	t.insertAllKeys(inputData)

	t.randomlyRead(inputData)

	t.randomlyDelete(inputData)

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
	process, err := os.StartProcess(startScriptPath, []string{startScriptPath, fmt.Sprint(numberOfNodes)}, &attr)
	if err != nil {
		return err
	}

	t.logger.Info("starting the system, number of nodes: " + fmt.Sprint(numberOfNodes))
	t.storageSystemProcess = process
	t.report.timeStarted = time.Now()
	t.totalNumOfNodes = numberOfNodes

	t.client.Wait()

	return t.client.ConfigureSystem(numberOfNodes, config.GetCollectionName())
}

func (t *Tester) insertAllKeys(inputData map[string]interface{}) {
	for key, val := range inputData {
		t.logger.Debug("test data", zap.String("key", key), zap.Any("value", val))
		insertTime := time.Now()

		node, err := t.client.Insert(key, val)

		if err != nil {

			t.logger.Error("failed to insert: "+err.Error(), zap.String("key", key))

			continue
		}

		insertedInfo := nodeOperation{
			NodeName:  node,
			Key:       key,
			Timestamp: insertTime,
		}

		t.report.insertedToNodes = append(t.report.insertedToNodes, insertedInfo)

	}

	data, err := t.client.GetAllData()
	if err != nil {
		t.logger.Error("failed to get all inserted data: " + err.Error())
	}
	t.report.allInsertedData = data

}

func (t *Tester) randomlyRead(inputData map[string]interface{}) {

	keys := reflect.ValueOf(inputData).MapKeys()
	rand.Seed(time.Now().UnixNano())

	for {

		// break when desired number of nodes is accessed
		if len(t.report.readFromNodes) >= t.getMinNumOfAccessedNodes() {
			break
		}

		index := rand.Intn(len(keys))

		key, ok := keys[index].Interface().(string)
		if !ok {
			t.logger.Warn("invalid key, can't cast to string", zap.String("key", key))
			continue
		}

		value, nodeName, err := t.client.Get(key)
		if err != nil {
			t.logger.Warn("failed to get the element: "+key, zap.Error(err))
			continue
		}

		info := nodeOperation{
			NodeName:  nodeName,
			Key:       key,
			Timestamp: time.Now(),
		}
		t.report.readFromNodes = append(t.report.readFromNodes, info)

		t.logger.Info("value read", zap.String("key", key), zap.Any("value", value))
	}

}

func (t *Tester) randomlyDelete(inputData map[string]interface{}) {

	keys := reflect.ValueOf(inputData).MapKeys()
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 3; i++ {

		index := rand.Intn(len(keys))

		key, ok := keys[index].Interface().(string)
		if !ok {
			t.logger.Warn("invalid key, can't cast to string", zap.String("key", key))
			continue
		}

		nodeName, err := t.client.Delete(key)
		if err != nil {
			t.logger.Warn("failed to delete the element: "+key, zap.Error(err))
		}

		info := nodeOperation{
			NodeName:  nodeName,
			Key:       key,
			Timestamp: time.Now(),
		}
		t.report.deletedFromNodes = append(t.report.deletedFromNodes, info)
	}
}

func (t *Tester) stop() {
	pwd, err := os.Getwd()
	if err != nil {
		var attr = os.ProcAttr{
			Env: os.Environ(),
			Dir: pwd + "/..",
		}

		_, err := os.StartProcess("docker-compose", []string{"down"}, &attr)
		if err != nil {
			t.logger.Error("failed to execute docker-compose down", zap.Error(err))
		}

	}

	if err := t.storageSystemProcess.Kill(); err != nil {
		t.logger.Error("failed to kill the storage system", zap.Error(err))
	}
}

func (t Tester) getMinNumOfAccessedNodes() int {
	desiredNumber := 5
	if t.totalNumOfNodes < desiredNumber {
		return t.totalNumOfNodes
	}

	return desiredNumber
}

func (t Tester) GenerateReport() string {
	var report string
	report += "# Start\n"
	report += fmt.Sprintf("System started at %s.\n", t.report.timeStarted.Format(time.RFC3339))
	report += fmt.Sprintf("Created %d nodes.\n", t.totalNumOfNodes)
	report += fmt.Sprintf("Saving to collection: %s.\n", config.GetCollectionName())
	report += "\n"

	report += "# Inserting the elements\n"
	for i, info := range t.report.insertedToNodes {
		report += fmt.Sprintf("%d) key: %s, node: %s, timestamp: %s\n", i, info.Key, info.NodeName, info.Timestamp.Format(time.RFC3339))
	}
	report += "\n"

	report += "# Randomly reading the elements\n"
	for i, info := range t.report.readFromNodes {
		report += fmt.Sprintf("%d) key: %s, node: %s, timestamp: %s\n", i, info.Key, info.NodeName, info.Timestamp.Format(time.RFC3339))
	}
	report += "\n"

	report += "# Randomly deleting the elements\n"
	for i, info := range t.report.deletedFromNodes {
		report += fmt.Sprintf("%d) key: %s, node: %s, timestamp: %s\n", i, info.Key, info.NodeName, info.Timestamp.Format(time.RFC3339))
	}
	report += "\n"

	report += "# All the inserted elements"
	report += t.report.allInsertedData
	report += "\n"
	report += "\n"

	return report
}
