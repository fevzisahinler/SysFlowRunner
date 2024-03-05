package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
	"sysflowrunner/pkg/couchbase"
	"sysflowrunner/pkg/influxdb"
	"sysflowrunner/pkg/rabbitmq"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
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
		log.Fatalf("Error fetching document IDs: %v", err)
	}

	var wg sync.WaitGroup

	for _, id := range docIds {
		beacon, err := couchbaseClient.GetSysBeaconByID(couchbaseBucket, id)
		if err != nil {
			log.Printf("Error fetching SysBeacon by ID: %v", err)
			continue
		}

		rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s",
			beacon.RabbitmqUsername, beacon.RabbitmqPassword, beacon.RabbitmqHostname, beacon.RabbitmqPort)

		rmqClient, err := rabbitmq.NewRabbitMQClient(rabbitMQURL)
		if err != nil {
			log.Printf("Error connecting to RabbitMQ: Host: %s Port: %s Password: %s %v", beacon.RabbitmqHostname, beacon.RabbitmqPort, beacon.RabbitmqPassword, err)
			continue
		}

		go func(queueName, rabbitmqHostname string) {
			defer wg.Done()
			msgChan, err := rmqClient.ConsumeQueue(queueName)
			if err != nil {
				log.Printf("Error consuming queue: %v", err)
				return
			}
			for d := range msgChan {
				var data map[string]interface{}
				fmt.Printf("Raw Data: %s\n", d)
				if err := json.Unmarshal(d, &data); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}
				log.Printf("Received data: %+v, %v", d, time.Now())

				if err := influxDBClient.WriteData(bucketName, rabbitmqHostname, data); err != nil {
					log.Printf("Error writing data to InfluxDB: %v", err)
				}
			}
		}(beacon.RabbitmqQueue, beacon.RabbitmqHostname)
		wg.Add(1)
	}

	wg.Wait()
}
