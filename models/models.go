package models

type BeaconInfo struct {
	Hostname         string
	RabbitmqHostname string
	RabbitmqQueue    string
}

type SysBeacon struct {
	RabbitmqHostname string
	RabbitmqPassword string
	RabbitmqUsername string
	RabbitmqQueue    string
	InfluxBucket     string
	RabbitmqPort     string
}
