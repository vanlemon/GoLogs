package logs

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	graylog "github.com/gemnasium/logrus-graylog-hook"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zbindenren/logrus_mail"

	"lmf.mortal.com/GoLogs/util"
)

/**************************************************************************/
/********************************** 常量 ***********************************/
/**************************************************************************/

const (
	LOGID_KEY = "K_LOGID" // 日志唯一标识
)

const (
	LEVEL_KEY = "K_LEVEL" // 日志级别
	DEBUG     = "Debu"
	INFO      = "Info"
	WARN      = "Warn"
	ERROR     = "Erro"
	FATAL     = "Fata"
)

const (
	ENV_KEY = "K_ENV" // 运行环境
	Sys     = "sys"   // 系统环境，系统启动时的 Context 运行环境
	Dev     = "dev"   // 开发环境
	Prod    = "prod"  // 运行环境
	Test    = "test"  // 压测环境
)

/**************************************************************************/
/********************************** 配置 ***********************************/
/**************************************************************************/

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

/**************************************************************************/
/******************************** 日志函数 **********************************/
/**************************************************************************/

// 初始化默认实例
func InitDefaultLogger(config LogConfig) {
	logger := DefaultLoggerInstance().GetLogrus() // 获取日志实例中需要配置的 logrus 引用

	CtxDebug(SysCtx, "[Bootstrap Logs]init logrus success")

	if config.Env == Dev || config.Env == Test { // 开发模式或压测模式
		logger.SetLevel(logrus.DebugLevel)
		CtxWarn(SysCtx, "[Bootstrap Logs]Logging in Debug Level")
	} else if config.Env == Prod { // 线上运行模式
		CtxInfo(SysCtx, "[Bootstrap Logs]Logging in Info Level")
		// 功能1：日志输出到滚动文件 rotatelogs hook
		if !util.HasNil(config.LogDir, config.LogFileName) { // 参数校检不为空
			// 创建滚动文件实例
			execPath := util.GetExecPath()
			CtxInfo(SysCtx, "[Bootstrap Logs]init logger execPath: %+v", execPath)
			rl, err := rotatelogs.New(
				strings.Join([]string{
					execPath, // 日志存储在执行目录同级别的 LogDir 目录下
					"..",
					config.LogDir,
					config.LogFileName + ".%Y%m%d%H", // %H%m
				}, string(os.PathSeparator)),                // 日志文件路径
				rotatelogs.WithLinkName(config.LogFileName), // 指向当前日志的软连接
				rotatelogs.WithRotationTime(time.Hour)) // 默认按小时分割
			if err != nil {
				CtxFatal(SysCtx, "[Bootstrap Logs]init rotatelogs error: %+v", err)
			}
			// 日志输出到文件 file hook
			// 默认 Info 级别以上日志输出到滚动文件
			logger.AddHook(lfshook.NewHook(
				lfshook.WriterMap{
					logrus.InfoLevel:  rl,
					logrus.WarnLevel:  rl,
					logrus.ErrorLevel: rl,
					logrus.FatalLevel: rl,
				},
				&logrus.TextFormatter{}, // 默认日志格式
			))
		} else {
			CtxWarn(SysCtx, "[Bootstrap Logs]no file hook")
		}
		// 功能2：日志输出到远程服务器 graylog hook
		if !util.HasNil(config.LogServerIp, config.LogServerPort) { // 参数校检不为空
			// TODO：是否正确连接服务器，无法检验
			//hook := graylog.NewAsyncGraylogHook(util.CombineIpAndPort(config.LogServerIp, config.LogServerPort), nil)
			// TODO：测试使用同步 hook 还是异步 hook
			hook := graylog.NewGraylogHook(util.CombineIpAndPort(config.LogServerIp, config.LogServerPort), nil)
			logger.AddHook(hook) // 添加远程服务器钩子
		} else {
			CtxWarn(SysCtx, "[Bootstrap Logs]no graylog hook")
		}
		// 功能3：日志邮件报警 mail hook: 默认发送 Error 级别以上的日志
		if !util.HasNil(
			config.MailBot.Name,
			config.MailBot.SmtpServerIp, config.MailBot.SmtpServerPort,
			config.MailBot.FromMailAddress,
			config.MailBot.UserName, config.MailBot.Password) && // 配置不为空
			!util.HasNil(config.MailBot.ToMailAddressList[:]) { // 列表不为空
			if config.MailBot.Enable {                          // 启用邮件报警
				for _, toMailAddress := range config.MailBot.ToMailAddressList { // 循环添加报警人
					mailHook, err := logrus_mail.NewMailAuthHook(
						config.MailBot.Name,
						config.MailBot.SmtpServerIp, config.MailBot.SmtpServerPort,
						config.MailBot.FromMailAddress, toMailAddress,
						config.MailBot.UserName, config.MailBot.Password)
					if err != nil {
						CtxError(SysCtx, "[Bootstrap Logs]connect to mail server err: %+v", err)
					}
					logger.AddHook(mailHook) // 添加邮件报警钩子
				}
			}

		} else {
			CtxWarn(SysCtx, "[Bootstrap Logs]no mail hook")
		}
	} else { // 未识别环境，报错
		CtxFatal(SysCtx, "[Bootstrap Logs]config env: %+v error", config.Env)
	}
}

func CtxFatal(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxFatal(ctx, format, v...)
}

func CtxError(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxError(ctx, format, v...)
}

func CtxWarn(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxWarn(ctx, format, v...)
}

func CtxInfo(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxInfo(ctx, format, v...)
}

func CtxDebug(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxDebug(ctx, format, v...)
}

/**************************************************************************/
/****************************** Context 函数 *******************************/
/**************************************************************************/

// 生成一个 24 位的 logId
func genLogId() string {
	t := time.Now()

	return strings.Join([]string{
		t.Format("20060102150405"),          // 14
		strconv.Itoa(t.Nanosecond() / 1000), // 6
		util.RandHexString(4),               // 4
	}, "")

	// 方法二
	//return fmt.Sprintf("%s%6d%s",
	//	t.Format("20060102150405"),
	//	t.Nanosecond()/1000,
	//	util.RandHexString(4))

	// 方法三
	//buf := bytes.Buffer{}
	//buf.WriteString(t.Format("20060102150405"))
	//buf.WriteString(strconv.Itoa(t.Nanosecond()/1000))
	//buf.WriteString(util.RandHexString(4))
	//return buf.String()
}

// Context 加 LogId
func CtxWithLogId(ctx context.Context) context.Context {
	if ctx == nil { // Context 为空，新建一个
		ctx = context.Background()
	}
	if logId := ctx.Value(LOGID_KEY); logId != nil { // 已有 LogId，直接返回
		return ctx
	}
	return context.WithValue(ctx, LOGID_KEY, genLogId()) // 创建新 LogId
}

// Context 获取 LogId
func CtxGetLogId(ctx context.Context) string {
	if logIdValue := ctx.Value(LOGID_KEY); logIdValue != nil { // 已有 LogId
		logId, ok := logIdValue.(string)
		if ok {
			return logId
		}
	}
	return "" // 暂无 LogId
}

// 新 Context 加 LogId
func NewCtxWithLogId() context.Context {
	return CtxWithLogId(nil)
}

// 创建 Context

var SysCtx = context.WithValue(NewCtxWithLogId(), ENV_KEY, Sys) // 系统启动 Context，只有一个

func DevCtx() context.Context { // 创建开发环境 Context
	return context.WithValue(NewCtxWithLogId(), ENV_KEY, Dev)
}

func ProdCtx() context.Context { // 创建运行环境 Context
	return context.WithValue(NewCtxWithLogId(), ENV_KEY, Prod)
}

func TestCtx() context.Context { // 创建压测环境 Context
	return context.WithValue(NewCtxWithLogId(), ENV_KEY, Test)
}

func Ctx(env string) context.Context { // 根据环境标识创建 Context（dev，prod，test）
	if env == Dev {
		return DevCtx()
	} else if env == Prod {
		return ProdCtx()
	} else if env == Test {
		return TestCtx()
	} else {
		CtxFatal(SysCtx, "[Bootstrap Logs]config env: %+v error", env)
		panic(fmt.Sprintf("config env: %+v error", env))
	}
}
