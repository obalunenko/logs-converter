// logs-converter-cli is a command-line application that allow to parse files with different log formats and according
// on their basis insert MongoDB documents with a monotonous structure.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/oleg-balunenko/logs-converter/internal/config"
	"github.com/oleg-balunenko/logs-converter/internal/converter"
	"github.com/oleg-balunenko/logs-converter/internal/db"
	"github.com/oleg-balunenko/logs-converter/internal/models"
)

func main() {
	versionInfo()

	cfg, errLoadCfg := config.LoadConfig("config.toml")
	if errLoadCfg != nil {
		log.Fatalf("Failed to load config: %v \nExiting", errLoadCfg)
	}

	dbc, err := db.Connect(db.StorageTypeMongo, db.Params{
		URL:        cfg.DBURL,
		DB:         cfg.DBName,
		Collection: cfg.MongoCollection,
		Username:   cfg.DBUsername,
		Password:   cfg.DBPassword,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if cfg.DropDB {
		if err := dbc.Drop(); err != nil {
			log.Fatal(err)
		}
	}

	resChan := make(chan *models.LogModel)
	errorsChan := make(chan error)

	wg := &sync.WaitGroup{}
	startJobs(cfg.GetFilesList(), cfg.FilesMustExist, cfg.FollowFiles, wg, resChan, errorsChan)

	stop := make(chan struct{})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	go func() {
		wg.Wait()
		stop <- struct{}{}
	}()

	process(dbc, resChan, signals, errorsChan, stop)
}

func process(dbc db.Repository, resChan <-chan *models.LogModel, signals <-chan os.Signal,
	errorsChan <-chan error, stopChan <-chan struct{}) {
	var (
		storedModelsCnt, failedToStoreCnt, totalRecCnt uint64
	)

	defer func() {
		dbc.Close()
		executionSummary(totalRecCnt, storedModelsCnt, failedToStoreCnt)
	}()

	for {
		select {
		case <-signals:
			log.Infof("Got UNIX signal, shutting down")
			return
		case data := <-resChan:
			totalRecCnt++
			log.Debugf("Received model: %+v", data)
			log.Infof("Current amount of received models is: [%d]", totalRecCnt)

			id, errStore := dbc.Store(data)
			if errStore != nil {
				log.Errorf("Failed to store model...: %v", errStore)
				failedToStoreCnt++
			} else {
				log.Debugf("Successfully stored model[id: %s] [%+v].", id, data)
				storedModelsCnt++
				log.Infof("Current amount of stored models: %d", storedModelsCnt)
			}
		case err := <-errorsChan:
			if err != nil {
				log.Errorf("Receive error: %v", err)
			}
		case <-stopChan:
			log.Printf("stop received")
			return
		}
	}
}
func startJobs(files map[string]string, filesmustExist bool, followFiles bool, wg *sync.WaitGroup,
	resChan chan *models.LogModel, errorsChan chan error) {
	for l, format := range files {
		wg.Add(1)

		go converter.Start(l, format, filesmustExist, followFiles, resChan, errorsChan, wg)
	}
}

func executionSummary(received uint64, stored uint64, failed uint64) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)

	_, err := fmt.Fprintf(w, "Execution statistics:\n"+
		"Total models received\tStored in DBName\tFailed to store in DBName\n"+
		"%d\t%d\t%d", received, stored, failed)
	if err != nil {
		log.Errorf("failed to print execution summary: %v", err)
	}

	if err := w.Flush(); err != nil {
		log.Errorf("failed to flush statistic writer: %v", err)
	}
}
