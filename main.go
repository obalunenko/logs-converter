package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/tabwriter"

	"github.com/oleg-balunenko/logs-converter/config"
	"github.com/oleg-balunenko/logs-converter/converter"
	"github.com/oleg-balunenko/logs-converter/db"
	"github.com/oleg-balunenko/logs-converter/model"
	log "github.com/sirupsen/logrus"
)

var (
	version string
	build   string
	commit  string
)

func main() {
	fmt.Printf("Version info: %s:%s", version, build)
	fmt.Printf("commit: %s ", commit)

	cfg, errLoadCfg := config.LoadConfig("config.toml")
	if errLoadCfg != nil {
		log.Fatalf("Failed to load config: %v \nExiting", errLoadCfg)
	}

	dbc, err := db.Connect(db.Mongo, cfg.MongoURL, cfg.MongoDB, cfg.MongoCollection, cfg.MongoUsername, cfg.MongoPassword)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := dbc.Drop(cfg.DropDB); err != nil {
		log.Fatal(err)
	}

	resChan := make(chan *model.LogModel)
	errorsChan := make(chan error)

	wg := &sync.WaitGroup{}
	startJobs(cfg.LogsFilesList, wg, resChan, errorsChan)

	signals := make(chan os.Signal, 1)
	stop := make(chan struct{})
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	go func() {
		wg.Wait()
		stop <- struct{}{}

	}()

	var storedModelsCnt, failedToStoreCnt, totalRecCnt uint64
	defer executionSummary(totalRecCnt, storedModelsCnt, failedToStoreCnt)

	for {

		select {
		case <-signals:
			log.Infof("Got UNIX signal, shutting down")
			dbc.Close()

			close(resChan)
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
		case errors := <-errorsChan:
			log.Errorf("Receive error: %v", errors)

		case <-stop:
			log.Printf("stop received")
			close(resChan)
			dbc.Close()

			return

		}

	}

}

func startJobs(files map[string]string, group *sync.WaitGroup, resChan chan *model.LogModel, errorsChan chan error) {
	for l, format := range files {
		group.Add(1)
		go converter.Start(l, format, true, false, resChan, errorsChan, group)

	}

}

func executionSummary(received uint64, stored uint64, failed uint64) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)

	_, err := fmt.Fprintf(w, "Execution statistics:\n"+
		"Total models received\tStored in MongoDB\tFailed to store in MongoDB\n"+
		"%d\t%d\t%d", received, stored, failed)
	if err != nil {
		log.Fatalf("failed to print execution summary: %v", err)
	}

	if err := w.Flush(); err != nil {
		log.Fatalf("failed to flush statistic writer: %v", err)
	}

}
