package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/oleg-balunenko/logs-converter/config"
	"github.com/oleg-balunenko/logs-converter/converter"
	"github.com/oleg-balunenko/logs-converter/mongo"
	log "github.com/sirupsen/logrus"
)

func main() {

	cfg, errLoadCfg := config.LoadConfig("config.toml")
	if errLoadCfg != nil {
		log.Fatalf("Failed to load config: %v \nExiting", errLoadCfg)
	}

	db := mongo.NewConnection(cfg)

	if cfg.DropDB {
		db.DropDatabase()

	}

	resChan := make(chan *converter.LogModel)
	for l, format := range cfg.LogsFilesList {

		go converter.Start(l, format, resChan)

	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)
	var storedModelsCnt, failedToStoreCnt, totalRecCnt uint64

	for {

		select {
		case <-signals:
			log.Infof("Got UNIX signal, shutting down")
			db.CloseConnection()

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)
			_, err := fmt.Fprintf(w, "Execution statistics:\n"+
				"Total models received\tStored in DB\tFailed to store in DB\n"+
				"%d\t%d\t%d", totalRecCnt, storedModelsCnt, failedToStoreCnt)
			if err != nil {
				log.Fatalf("Failed to print statistic: %v", err)
			}
			//fmt.Fprintln(w)
			if err := w.Flush(); err != nil {
				log.Fatalf("Failed to flush statistic writer: %v", err)
			}

			return

		case data := <-resChan:

			totalRecCnt++
			log.Debugf("Received model: %+v", data)
			log.Infof("Current amount of received models is: [%d]", totalRecCnt)
			errStore := db.StoreModel(data)
			if errStore != nil {
				log.Errorf("Failed to store model...: %v", errStore)
				failedToStoreCnt++
			} else {
				log.Debugf("Successfully stored model [%+v].", data)
				storedModelsCnt++
				log.Infof("Current amount of stored models: %d", storedModelsCnt)

			}

		}
	}

}
