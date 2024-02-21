package tlog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"time"
)

// var envMode = configs.EnvModeTypeDevelop
// var developLogger *zap.Logger
var productLogger *zap.Logger

func InitLog() {
	initProductLog("./logs/flygoose.log")
}

func Info(msg string, fields ...zap.Field) {
	productLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	productLogger.Warn(msg, fields...)

}

func Error2(msg string, err error) {
	if err != nil && err != gorm.ErrRecordNotFound {
		productLogger.Error(msg, zap.Error(err))
	}
}

func Error(msg string, fields ...zap.Field) {
	productLogger.Error(msg, fields...)
}

//func initDevelopLog(logFileName string) {
//	// 日志分割
//	hook := lumberjack.Logger{
//		Filename:   logFileName, // 日志文件路径，默认 os.TempDir()
//		MaxSize:    100,         // 每个日志文件保存10M，默认 100M
//		MaxBackups: 30,          // 保留30个备份，默认不限
//		MaxAge:     7,           // 保留5天，默认不限
//		Compress:   false,       // 是否压缩，默认不压缩
//	}
//
//	write := zapcore.AddSync(&hook)
//	encoderConfig := zapcore.EncoderConfig{
//		TimeKey:       "time",
//		LevelKey:      "level",
//		NameKey:       "logger",
//		CallerKey:     "linenum",
//		MessageKey:    "msg",
//		StacktraceKey: "stacktrace",
//		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
//			enc.AppendString(t.Format("2006-01-02 15:04:05"))
//		},
//		LineEnding:     zapcore.DefaultLineEnding,
//		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
//		EncodeDuration: zapcore.SecondsDurationEncoder, //
//		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
//		EncodeName:     zapcore.FullNameEncoder,
//	}
//	// 设置日志级别
//	// debug 可以打印出 info debug warn
//	// info  级别可以打印 warn info
//	// warn  只能打印 warn
//	// debug->info->warn->error
//	atomicLevel := zap.NewAtomicLevel()
//	atomicLevel.SetLevel(zap.InfoLevel)
//	core := zapcore.NewCore(
//		// zapcore.NewConsoleEncoder(encoderConfig),
//		zapcore.NewConsoleEncoder(encoderConfig),
//		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(write)),
//		zap.InfoLevel,
//	)
//	// 开启开发模式，堆栈跟踪
//	caller := zap.AddCaller()
//
//	stack := zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel))
//
//	// 开启文件及行号
//	development := zap.Development()
//	// 设置初始化字段,如：添加一个服务器名称
//	//filed := zap.Fields(zap.String("serviceName", "serviceName"))
//	// 构造日志
//	developLogger = zap.New(core, caller, stack, development)
//}

func initProductLog(logFileName string) {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   logFileName, // 日志文件路径，默认 os.TempDir()
		MaxSize:    100,         // 每个日志文件保存10M，默认 100M
		MaxBackups: 30,          // 保留30个备份，默认不限
		MaxAge:     7,           // 保留5天，默认不限
		Compress:   false,       // 是否压缩，默认不压缩
	}

	write := zapcore.AddSync(&hook)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "linenum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(write)),
		zap.InfoLevel,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	stack := zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel))

	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	productLogger = zap.New(core, caller, stack, development)
}
