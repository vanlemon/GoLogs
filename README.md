# GoLogs

使用示例：[example](https://github.com/lilinxi/GoLogs/blob/master/example/main.go)

敏感数据已用`something`脱敏。

---

## 功能说明


1. 日志功能

- [基于logrus 配置](https://mojotv.cn/2018/12/27/golang-logrus-tutorial)

|功能|可选/必选|依赖|
|---|---|---|
|日志输出到命令行|必选|github.com/sirupsen/logrus|
|日志输出到滚动文件|可选|github.com/lestrrat-go/file-rotatelogs & github.com/rifflock/lfshook|
|日志输出到远程服务器|可选|github.com/gemnasium/logrus-graylog-hook|
|日志邮件报警|可选|github.com/zbindenren/logrus_mail|

2. Context 功能

- 生成唯一 LogId（24）
    - YYYYMMDDHHmmss:   14
    - ns/1000:          6
    - RandOct:          4

---

## 接口说明

```go
// 初始化日志服务
func InitDefaultLogger(config LogConfig) {

// 输出服务日志
func CtxFatal(ctx context.Context, format string, v ...interface{}) {
func CtxError(ctx context.Context, format string, v ...interface{}) {
func CtxWarn(ctx context.Context, format string, v ...interface{}) {
func CtxInfo(ctx context.Context, format string, v ...interface{}) {
func CtxDebug(ctx context.Context, format string, v ...interface{}) {

// 创建带 logid 的 Context
func CtxWithLogId(ctx context.Context) context.Context {
	
// Context 获取 LogId
func CtxGetLogId(ctx context.Context) string {
	
// 创建带 logid 的 Context
func NewCtxWithLogId() context.Context {

// 日志服务需要传入 Context
var SysCtx = context.WithValue(NewCtxWithLogId(), ENV_KEY, Sys) // 系统启动 Context，只有一个
func DevCtx() context.Context { // 创建开发环境 Context
func ProdCtx() context.Context { // 创建运行环境 Context
func TestCtx() context.Context { // 创建压测环境 Context
func Ctx(env string) context.Context { // 根据环境标识创建 Context（dev，prod，test）
```

---

## 日志服务器配置

- 日志服务配置

```go
type LogConfig struct { // 日志系统配置
	Env           string `json:"env"`             // 日志系统的运行环境
	LogDir        string `json:"log_dir"`         // 日志输出的文件夹
	LogFileName   string `json:"log_file_name"`   // 日志输出的文件名，默认按天分割
	LogServerIp   string `json:"log_server_ip"`   // 日志服务器 IP 地址
	LogServerPort int    `json:"log_server_port"` // 日志服务器端口
	MailBot       `json:"mail_bot"`               // 邮件告警机器人配置
}

type MailBot struct { // 邮件告警机器人配置
	Name              string   `json:"name"`                 // 机器人名字
	SmtpServerIp      string   `json:"smtp_server_ip"`       // Smtp 服务器 IP 地址
	SmtpServerPort    int      `json:"smtp_server_port"`     // Smtp 服务器端口
	FromMailAddress   string   `json:"from_mail_address"`    // 机器人邮箱地址
	ToMailAddressList []string `json:"to_mail_address_list"` // 收件人邮箱地址列表
	UserName          string   `json:"username"`             // 机器人邮箱用户名
	Password          string   `json:"password"`             // 机器人邮箱密码
	Enable            bool     `json:"enable"`               // 是否启用告警
}
```

- [graylog 配置](http://docs.graylog.org/en/2.3/pages/installation/docker.html)
    - [docker compose](https://yeasy.gitbooks.io/docker_practice/compose/usage.html)
    - http://127.0.0.1:9000(只能在配置中更改)
    - admin/admin(密码盐值只能在配置中更改)
    - update time zone from UTC to PRC(时区只能在配置中更改)
    - **配置输入，和输入端口**：Inputs -> GELF UDP -> Launch new input

---

## 日志中间件实例

```go
func LogsMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		// step1: 根据运行环境初始化 ctx，每次请求会分配唯一的 logid
		ctx := logs.Ctx(config.ConfigInstance.Env)
		gctx.Set(cconst.CtxKey, ctx) // 设置 ctx

		// step2: 打印请求日志
		host := gctx.Request.Host     // 请求主机
		url := gctx.Request.URL       // 请求 url
		method := gctx.Request.Method // 请求接口
		reqTime := time.Now()            // 请求时间
		logs.CtxInfo(ctx, "[Access] %s \t %s \t %s \t %s", reqTime.Format("2006-01-02 15:04:05"), host, url, method)

		// step3: 执行服务
		gctx.Next()

		// step4: 打印返回日志
		respTime := time.Now()            // 返回时间
		costTime := respTime.Sub(reqTime) // 耗时
		logs.CtxInfo(ctx, "[Access] %s \t %s \t %s \t %s \t %+v \t %+v", respTime.Format("2006-01-02 15:04:05"), host, url, method, gctx.Writer.Status(), costTime)
	}
}
```
    
