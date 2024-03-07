package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"sync"
	"sysflowrunner/logger" // logger paketini import et
	"sysflowrunner/pkg/couchbase"
	"sysflowrunner/pkg/influxdb"
	"sysflowrunner/pkg/rabbitmq"
	"time"
)

func main() {
	logger.Init()
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file: %v", err) // logger kullanarak loglama
	}

	couchbaseURL := os.Getenv("COUCHBASE_URL")
	couchbaseUsername := os.Getenv("COUCHBASE_USERNAME")
	couchbasePassword := os.Getenv("COUCHBASE_PASSWORD")
	couchbaseBucket := os.Getenv("COUCHBASE_BUCKET")
	influxDBURL := os.Getenv("INFLUXDB_URL")
	influxDBToken := os.Getenv("INFLUXDB_TOKEN")
	bucketName := os.Getenv("INFLUXDB_BUCKET_NAME")
	influxDBOrg := os.Getenv("INFLUXDB_ORG")

	couchbaseClient := couchbase.NewCouchbaseClient(couchbaseURL, couchbaseUsername, couchbasePassword)
	influxDBClient := influxdb.NewInfluxDBClient(influxDBURL, influxDBToken, influxDBOrg)

	docIds, err := couchbaseClient.GetDocumentIDs("sysbeacon")
	if err != nil {
		logger.Error("Error fetching document IDs: %v", err) // logger kullanarak loglama
	}

	var wg sync.WaitGroup

	for _, id := range docIds {
		beacon, err := couchbaseClient.GetSysBeaconByID(couchbaseBucket, id)
		if err != nil {
			logger.Error("Error fetching SysBeacon by ID: %v", err)
			continue
		}

		rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s",
			beacon.RabbitmqUsername, beacon.RabbitmqPassword, beacon.RabbitmqHostname, beacon.RabbitmqPort)

		rmqClient, err := rabbitmq.NewRabbitMQClient(rabbitMQURL)
		if err != nil {
			logger.Error("Error connecting to RabbitMQ: Host: %s Port: %s Password: %s %v", beacon.RabbitmqHostname, beacon.RabbitmqPort, beacon.RabbitmqPassword, err) // logger kullanarak loglama
			continue
		}

		go func(queueName, rabbitmqHostname string) {
			defer wg.Done()
			msgChan, err := rmqClient.ConsumeQueue(queueName)
			if err != nil {
				logger.Error("Error consuming queue: %v", err)
				return
			}
			for d := range msgChan {
				var data map[string]interface{}
				if err := json.Unmarshal(d, &data); err != nil {
					logger.Error("Error unmarshaling message: %v", err)
					continue
				}
				logger.Info("Received data: %+v, Time: %v", data, time.Now())

				if err := influxDBClient.WriteData(bucketName, rabbitmqHostname, data); err != nil {
					logger.Error("Error writing data to InfluxDB: %v", err)
				}
			}
		}(beacon.RabbitmqQueue, beacon.RabbitmqHostname)
		wg.Add(1)
	}

	wg.Wait()
}
