package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var errorLogger *zap.SugaredLogger

func init() {
	var coreArr []zapcore.Core
	//encoder := GetEncoder()
	// 实现两个判断日志等级的interface
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.WarnLevel
	//})

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	//infoWriter := getWriter(config.LogInfoFileName)
	//errorWriter := getWriter(config.LogWarnFileName)
	consoleEncoder := GetConsoleEncoder()
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	coreArr = append(coreArr, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), errorLevel))
	caller := zap.AddCaller()
	// 构造日志
	log := zap.New(zapcore.NewTee(coreArr...), caller, zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.WarnLevel))
	errorLogger = log.Sugar()
}

func getWriter(filename string) io.Writer {
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+"/%Y-%m-%d.log",
		rotatelogs.WithMaxAge(time.Hour*24*60),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}

// Recover http请求未捕获异常和捕获异常的调用深度不一样，在+2层可获取异常行号
func Recover(skip int, args ...interface{}) {
	var coreArr []zapcore.Core
	//encoder := GetEncoder()
	// 实现两个判断日志等级的interface
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.WarnLevel
	//})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	//infoWriter := getWriter(config.LogInfoFileName)
	//errorWriter := getWriter(config.LogWarnFileName)
	consoleEncoder := GetConsoleEncoder()
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	//日志开发环境只打印到console,生产环境打印到log文件
	coreArr = append(coreArr, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), errorLevel))
	caller := zap.AddCaller()
	// 构造日志
	log := zap.New(zapcore.NewTee(coreArr...), caller, zap.AddCallerSkip(skip), zap.AddStacktrace(zapcore.WarnLevel))
	errorLogger = log.Sugar()
	errorLogger.Error(args...)
}

// RespFail http请求请求失败后调用深度不一样，在+1层可获取异常行号
func RespFail(args ...interface{}) {
	var coreArr []zapcore.Core
	//encoder := GetEncoder()
	// 实现两个判断日志等级的interface
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.WarnLevel
	//})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	//infoWriter := getWriter(config.LogInfoFileName)
	//errorWriter := getWriter(config.LogWarnFileName)
	consoleEncoder := GetConsoleEncoder()
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	coreArr = append(coreArr, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), errorLevel))
	caller := zap.AddCaller()
	// 构造日志
	log := zap.New(zapcore.NewTee(coreArr...), caller, zap.AddCallerSkip(2), zap.AddStacktrace(zapcore.WarnLevel))
	errorLogger = log.Sugar()
	errorLogger.Error(args...)
}

func GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller_line",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    customLevelEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
}
func GetConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller_line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    cEncodeLevel,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	})
}

// cEncodeLevel 自定义日志级别显示+颜色
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s := _levelToColor[Level(level)].Add("[" + level.CapitalString() + "]")
	enc.AppendString(s)
	enc.AppendString("[" + time.Now().Format("2006-01-02 15:04:05") + "]")
}

// log文件中输出，有颜色会乱码
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
	enc.AppendString("[" + time.Now().Format("2006-01-02 15:04:05") + "]")
}

// 给level设置颜色
var (
	_levelToColor = map[Level]Color{
		DebugLevel:  Magenta,
		InfoLevel:   Green,
		WarnLevel:   Yellow,
		ErrorLevel:  Red,
		DPanicLevel: Red,
		PanicLevel:  Red,
		FatalLevel:  Red,
	}
)
