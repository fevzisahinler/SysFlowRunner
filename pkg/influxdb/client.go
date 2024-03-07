package influxdb

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

type InfluxDBClient struct {
	Client influxdb2.Client
	Org    string
}

func NewInfluxDBClient(url, token, org string) *InfluxDBClient {
	client := influxdb2.NewClient(url, token)
	return &InfluxDBClient{Client: client, Org: org}
}

func (client *InfluxDBClient) WriteData(bucketName, hostname string, data map[string]interface{}) error {
	writeAPI := client.Client.WriteAPIBlocking(client.Org, bucketName)

	cpuUsed := data["used_cpu"]
	cpuTotal := data["total_cpu"]
	memoryUsed := data["used_memory"]
	memoryTotal := data["total_memory"]

	var ts time.Time
	tsFloat := data["timestamp"].(float64)
	timestamp := int64(tsFloat)
	ts = time.Unix(timestamp, 0)

	p := influxdb2.NewPoint(
		"system_metrics",
		map[string]string{"host": hostname},
		map[string]interface{}{
			"used_cpu":     cpuUsed,
			"total_cpu":    cpuTotal,
			"used_memory":  memoryUsed,
			"total_memory": memoryTotal,
		},
		ts,
	)

	return writeAPI.WritePoint(context.Background(), p)
}
