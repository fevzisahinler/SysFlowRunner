package couchbase

import (
	"fmt"
	"sysflowrunner/models"

	"github.com/couchbase/gocb/v2"
)

type CouchbaseClient struct {
	Cluster *gocb.Cluster
}

func NewCouchbaseClient(connectionString, username, password string) *CouchbaseClient {
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(fmt.Errorf("Error connecting to Couchbase: %v", err))
	}
	return &CouchbaseClient{Cluster: cluster}
}

func (c *CouchbaseClient) GetDocumentIDs(bucketName string) ([]string, error) {
	query := fmt.Sprintf("SELECT META().id FROM `%s`", bucketName)
	rows, err := c.Cluster.Query(query, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var row map[string]interface{}
		if err := rows.Row(&row); err != nil {
			return nil, err
		}
		id, ok := row["id"].(string)
		if !ok {
			continue
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (c *CouchbaseClient) GetSysBeaconByID(bucketName, id string) (*models.SysBeacon, error) {
	query := fmt.Sprintf("SELECT rabbitmq.hostname AS `rabbitmq.host`, rabbitmq.passwd AS `rabbitmq.passwd`, rabbitmq.username AS `rabbitmq.username`, rabbitmq.queue AS `rabbitmq.queue`, rabbitmq.port AS `rabbitmq.port` FROM `%s` WHERE META().id = $1", bucketName)
	rows, err := c.Cluster.Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{id}})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result struct {
		RabbitmqHostname string `json:"rabbitmq.host"`
		RabbitmqPassword string `json:"rabbitmq.passwd"`
		RabbitmqUsername string `json:"rabbitmq.username"`
		RabbitmqQueue    string `json:"rabbitmq.queue"`
		RabbitmqPort     string `json:"rabbitmq.port"`
	}
	if rows.Next() {
		if err := rows.Row(&result); err != nil {
			return nil, err
		}

		return &models.SysBeacon{
			RabbitmqHostname: result.RabbitmqHostname,
			RabbitmqPassword: result.RabbitmqPassword,
			RabbitmqUsername: result.RabbitmqUsername,
			RabbitmqQueue:    result.RabbitmqQueue,
			RabbitmqPort:     result.RabbitmqPort,
		}, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Document not found")
}
