package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func BenchmarkGenLogId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genLogId()
	}
}

func TestGenLogId(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(genLogId())
	}
}

func TestSetReportCaller(t *testing.T) {
	logrus.SetReportCaller(true)

	logrus.Info("info msg")
}