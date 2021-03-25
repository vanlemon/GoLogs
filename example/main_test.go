package main

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoLogs/example/config"
)

func init() {
	config.InitConfig("./conf/logs_example_dev.json")
	logs.InitDefaultLogger(config.ConfigInstance.LogConfig)
}

func BenchmarkGateWay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GateWay()
	}
}

func BenchmarkCtxInfo(b *testing.B) {
	logid := time.Now().Format("20060102150405")
	for i := 0; i < b.N; i++ {
		logs.CtxInfo(context.WithValue(context.Background(), logs.LOGID_KEY, logid), "1")
	}
}

func BenchmarkCtxWarn(b *testing.B) {
	logid := time.Now().Format("20060102150405")
	for i := 0; i < b.N; i++ {
		logs.CtxWarn(context.WithValue(context.Background(), logs.LOGID_KEY, logid), "1")
	}
}

func BenchmarkLogrus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logrus.Info("1")
	}
}

func TestGateWay(t *testing.T) {
	GateWay()
}