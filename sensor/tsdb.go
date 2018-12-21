package main

import (
	"errors"
	"log"
	"os"

	"github.com/influxdata/influxdb/client/v2"
)

type tsdbData interface {
	Tags() map[string]string
	Fields() map[string]interface{}
	Measurement() string
}

type database struct {
	cli      client.Client
	queue    chan tsdbData
	database string
}

func (db *database) add(d tsdbData) error {
	db.queue <- d
	return nil
}

func (db *database) init() error {
	var err error
	db.cli, err = client.NewHTTPClient(
		client.HTTPConfig{
			Addr:     os.Getenv("SENSOR_INFLUXDB_HOST"),
			Username: os.Getenv("SENSOR_INFLUXDB_USERNAME"),
			Password: os.Getenv("SENSOR_INFLUXDB_PASSWORD"),
		},
	)
	if err != nil {
		return err
	}

	db.database = os.Getenv("SENSOR_INFLUXDB_DATABASE")
	if db.database == "" {
		return errors.New("no database specified")
	}

	db.queue = make(chan tsdbData, 1024)

	return nil
}

func (db *database) run() {
	for {
		d, ok := <-db.queue
		if !ok {
			log.Println("stop database")
			return
		}

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  db.database,
			Precision: "s",
		})
		if err != nil {
			log.Println(err)
			continue
		}

		pt, err := client.NewPoint(d.Measurement(), d.Tags(), d.Fields())
		if err != nil {
			log.Println(err)
			continue
		}
		bp.AddPoint(pt)
		err = db.cli.Write(bp)
		if err != nil {
			log.Println(err)
			continue
		}
	}

}
func (db *database) close() error {
	close(db.queue)
	return db.cli.Close()
}
