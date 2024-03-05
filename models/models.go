package models

type SysBeacon struct {
	RabbitmqHostname string
	RabbitmqPassword string
	RabbitmqUsername string
	RabbitmqQueue    string
	InfluxBucket     string
	RabbitmqPort     string
}
