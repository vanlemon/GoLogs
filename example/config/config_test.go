package config

import (
	"log"
	"testing"
)

func TestConfigInit(t *testing.T) {
	InitConfig("../conf/logs_example_dev.json")
	log.Printf("%+v\n", ConfigInstance)
	log.Printf("%+v\n", ConfigJson)
}
