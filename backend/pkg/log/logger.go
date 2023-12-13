package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"changeme/backend/pkg/configs"
)

func InitLog() {
	conf := configs.Default()
	level := Str2LogLevel(conf.Log.Level)
	l = NewLogger(
		SetLogFileDir(conf.Log.LogFileDir),
		SetLogPrefix(conf.App.Name),
		SetDevelopment(conf.Log.Env == "dev"),
		SetDebugFileName("debug.log"),
		SetInfoFileName("info.log"),
		SetWarnFileName("warn.log"),
		SetErrorFileName("error.log"),
		SetMaxAge(conf.Log.MaxAge),
		SetMaxBackups(conf.Log.MaxBackup),
		SetMaxSize(conf.Log.MaxFileSize),
		SetLevel(level))
}

type Options struct {
	LogFileDir    string //文件保存目录
	LogPrefix     string //日志文件前缀
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string

	Level       zapcore.Level //日志等级
	MaxSize     int           //日志文件小大（M）
	MaxBackups  int           // 最多存在多少个备份文件
	MaxAge      int           //保存的最大天数
	Development bool          //是否是开发模式
	zap.Config
}

type ModOptions func(options *Options)

var (
	l                              *Logger
	loggerOnce                     sync.Once
	sp                             = string(filepath.Separator) // 路径分隔符
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer          // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout)    // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

func Default() *Logger {
	return l
}

func ReplaceDefault(log *Logger) {
	l = log
}

func Debug(msg string, fields ...Field) {
	l.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	l.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	l.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	l.Error(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	l.DPanic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	l.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	l.Fatal(msg, fields...)
}

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel
	WarnLevel   Level = zap.WarnLevel
	ErrorLevel  Level = zap.ErrorLevel
	DPanicLevel Level = zap.DPanicLevel
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel
	DebugLevel  Level = zap.DebugLevel
)

func Str2LogLevel(s string) Level {
	switch strings.ToLower(s) {
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "debug":
		return DebugLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

type Field = zap.Field

type Logger struct {
	*zap.Logger
	aLevel  *zap.AtomicLevel
	Opts    *Options
	zapConf zap.Config
}

func NewLogger(modOpts ...ModOptions) *Logger {
	loggerOnce.Do(func() {
		l = &Logger{}
		l.Opts = &Options{
			LogFileDir:    "",
			LogPrefix:     "def",
			ErrorFileName: "error",
			WarnFileName:  "warn",
			InfoFileName:  "info",
			DebugFileName: "debug",
			Level:         zapcore.DebugLevel,
			MaxSize:       128,
			MaxBackups:    10,
			MaxAge:        30,
		}

		// 修改 opts
		for _, modOpt := range modOpts {
			modOpt(l.Opts)
		}

		if l.Opts.LogFileDir == "" {
			l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
			l.Opts.LogFileDir = fmt.Sprintf("%slogs%s", sp, sp)
		}
		if l.Opts.Development {
			l.zapConf = zap.NewDevelopmentConfig()
			l.zapConf.EncoderConfig.EncodeTime = timeEncoder
		} else {
			l.zapConf = zap.NewProductionConfig()
			l.zapConf.EncoderConfig.EncodeTime = timeUnixNano
		}
		if len(l.Opts.OutputPaths) == 0 {
			l.zapConf.OutputPaths = []string{"stdout"}
		}
		if len(l.Opts.ErrorOutputPaths) == 0 {
			l.zapConf.ErrorOutputPaths = []string{"stderr"}
		}
		l.zapConf.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		l.zapConf.Level.SetLevel(l.Opts.Level)
		l.init()
	})
	return l
}

func (l *Logger) init() {
	l.setSyncs()
	var err error
	l.Logger, err = l.zapConf.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.Logger.Sync()
}

func (l *Logger) setSyncs() {
	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.LogPrefix + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxBackups: l.Opts.MaxBackups,
			MaxAge:     l.Opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	return
}

func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConf.EncoderConfig)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConf.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConf.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConf.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConf.Level.Level() > -1
	})
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),
	}
	if l.Opts.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func SetMaxSize(MaxSize int) ModOptions {
	return func(option *Options) {
		option.MaxSize = MaxSize
	}
}
func SetMaxBackups(MaxBackups int) ModOptions {
	return func(option *Options) {
		option.MaxBackups = MaxBackups
	}
}
func SetMaxAge(MaxAge int) ModOptions {
	return func(option *Options) {
		option.MaxAge = MaxAge
	}
}

func SetLogFileDir(LogFileDir string) ModOptions {
	return func(option *Options) {
		option.LogFileDir = LogFileDir
	}
}

func SetLogPrefix(logPrefix string) ModOptions {
	return func(option *Options) {
		option.LogPrefix = logPrefix
	}
}

func SetLevel(Level zapcore.Level) ModOptions {
	return func(option *Options) {
		option.Level = Level
	}
}
func SetErrorFileName(ErrorFileName string) ModOptions {
	return func(option *Options) {
		option.ErrorFileName = ErrorFileName
	}
}
func SetWarnFileName(WarnFileName string) ModOptions {
	return func(option *Options) {
		option.WarnFileName = WarnFileName
	}
}

func SetInfoFileName(InfoFileName string) ModOptions {
	return func(option *Options) {
		option.InfoFileName = InfoFileName
	}
}
func SetDebugFileName(DebugFileName string) ModOptions {
	return func(option *Options) {
		option.DebugFileName = DebugFileName
	}
}
func SetDevelopment(Development bool) ModOptions {
	return func(option *Options) {
		option.Development = Development
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}

func (l *Logger) SetLevel(level Level) {
	if l.aLevel != nil {
		l.aLevel.SetLevel(level)
	}
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.Logger.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.Logger.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.Logger.Fatal(msg, fields...)
}

var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString
	Float64    = zap.Float64
	Float64p   = zap.Float64p
	Float32    = zap.Float32
	Float32p   = zap.Float32p
	Int        = zap.Int
	Intp       = zap.Intp
	Int64      = zap.Int64
	Int64p     = zap.Int64p
	Int32      = zap.Int32
	Int32p     = zap.Int32p
	Int16      = zap.Int16
	Int16p     = zap.Int16p
	Int8       = zap.Int8
	Int8p      = zap.Int8p
	String     = zap.String
	Stringp    = zap.Stringp
	Uint       = zap.Uint
	Uintp      = zap.Uintp
	Uint64     = zap.Uint64
	Uint64p    = zap.Uint64p
	Uint32     = zap.Uint32
	Uint32p    = zap.Uint32p
	Uint16     = zap.Uint16
	Uint16p    = zap.Uint16p
	Uint8      = zap.Uint8
	Uint8p     = zap.Uint8p
	Duration   = zap.Duration
	Durationp  = zap.Durationp
	Time       = zap.Time
	Any        = zap.Any
	Err        = zap.Error
)
