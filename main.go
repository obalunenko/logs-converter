package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gitlab.com/oleg.balunenko/logs-converter/config"
	"gitlab.com/oleg.balunenko/logs-converter/models"
	"gitlab.com/oleg.balunenko/logs-converter/mongo"
)

func main() {
	cfg := config.LoadConfig()
	if len(cfg.LogsFilesList) == 0 {
		log.Fatalf("No log files provided: [%+v], Exiting", cfg.LogsFilesList)
	}
	dbCollection := mongo.Connect(cfg)

	if cfg.DropDB {
		if errDrop := dbCollection.DropCollection(); errDrop != nil {
			log.Fatalf("Failed to drop the collection [%+v.%+v]:%v", dbCollection, dbCollection.Database, errDrop)
		}

	}

	resChan := make(chan *models.LogModel)
	for l, format := range cfg.LogsFilesList {

		go startConverting(l, format, resChan)

	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)
	var storedModelsCnt, failedToStoreCnt int
	for {

		select {
		case <-signals:
			log.Infof("Got UNIX signal, shutting down")
			mongo.CloseConnection(dbCollection)
			log.Infof("Total stored logs in DB: [%d]", storedModelsCnt)
			log.Infof("Total failed to store logs in DB: [%d]", failedToStoreCnt)
			return

		case data := <-resChan:
			log.Debugf("Received model: %+v", data)
			errStore := mongo.StoreModel(data, dbCollection)
			if errStore != nil {
				log.Errorf("Failed to store model...: %v", errStore)
				failedToStoreCnt++
			} else {
				storedModelsCnt++
			}

		}
	}

}
