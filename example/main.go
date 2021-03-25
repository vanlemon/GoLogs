package main

import (
	"context"
	"time"

	"lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoLogs/example/config"
	"lmf.mortal.com/GoLogs/util"
)

// 初始化配置，启动日志服务
func InitConfig() {
	config.InitConfig(util.GetExecPath() + "/../conf/logs_example_dev.json")
	logs.InitDefaultLogger(config.ConfigInstance.LogConfig)
}

func main() {
	InitConfig()
	for i := 0; i < 10; i++ {
		go GateWay()
	}
	time.Sleep(time.Second * 10)
}

// 网关入口接口
func GateWay() {
	ctx := logs.NewCtxWithLogId()
	ret := call1(ctx)
	logs.CtxInfo(ctx, "resp: %+v", ret)
}

// 第一个调用
func call1(ctx context.Context) string {
	logs.CtxInfo(ctx, "a log in call1: %s", "call1")
	logs.CtxDebug(ctx, "debug in call1: %s", "call1")
	return call2(ctx)
}

// 第二个调用
func call2(ctx context.Context) string {
	logs.CtxInfo(ctx, "a log in call2: %s", "call2")
	logs.CtxWarn(ctx, "warning in call2: %s", "call2")
	return call3(ctx)
}

// 第三个调用
func call3(ctx context.Context) string {
	logs.CtxError(ctx, "a error in call3: %s", "call3")
	// 可预知的错误使用 Fatal，不可预知的错误使用 panic
	//logs.CtxFatal(ctx, "a error in call3: %s", "call3")
	//panic(fmt.Sprintf("a error in call3: %s", "call3"))
	return "Success in call3"
}
