package log

/**
* created by mengqi on 2023/11/14
 */

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/uptrace/opentelemetry-go-extra/otelutil"
)

const numAttr = 5

var (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")
	logTemplateKey = attribute.Key("log.template")
)

var (
	std  = New(NewOptions())
	lock sync.Mutex
)

// Logger 封装了zap.Logger
type Logger struct {
	*zap.Logger
	skipCaller *zap.Logger

	withTraceID bool

	minLevel         zapcore.Level
	errorStatusLevel zapcore.Level

	caller     bool
	stackTrace bool

	extraFields []zap.Field
}

//func New(logger *zap.Logger, opts ...Option) *Logger {
//	l := &Logger{
//		Logger:     logger,
//		skipCaller: logger.WithOptions(zap.AddCallerSkip(1)),
//
//		minLevel:         zap.WarnLevel,
//		errorStatusLevel: zap.ErrorLevel,
//		caller:           true,
//	}
//	for _, opt := range opts {
//		opt(l)
//	}
//	return l
//}

func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	var extraFields []zap.Field
	// zap.New side effect is extracting fields from .WithOptions(zap.Fields(...))
	zap.New(&fieldExtractorCore{extraFields: &extraFields}, opts...)
	clone := *l
	clone.Logger = l.Logger.WithOptions(opts...)
	clone.skipCaller = l.skipCaller.WithOptions(opts...)
	clone.extraFields = append(clone.extraFields, extraFields...)
	return &clone
}

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: l.Logger.Sugar(),
		skipCaller:    l.skipCaller.Sugar(),
		l:             l,
	}
}

func (l *Logger) Clone(opts ...Option) *Logger {
	clone := *l
	for _, opt := range opts {
		opt(&clone)
	}
	return &clone
}

func (l *Logger) Ctx(ctx context.Context) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: ctx,
		l:   l,
	}
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.DebugLevel, msg, fields)
	l.skipCaller.Debug(msg, fields...)
}

func (l *Logger) DebugfContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.DebugLevel, format, []zapcore.Field{})
	l.skipCaller.Debug(msg, fields...)
}

func (l *Logger) DebugwContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.DebugLevel, format, []zapcore.Field{})
	l.skipCaller.Debug(msg, fields...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.InfoLevel, msg, fields)
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) InfofContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.InfoLevel, msg, []zapcore.Field{})
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.WarnLevel, msg, fields)
	l.skipCaller.Warn(msg, fields...)
}

func (l *Logger) WarnfContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.WarnLevel, msg, []zapcore.Field{})
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.ErrorLevel, msg, fields)
	l.skipCaller.Error(msg, fields...)
}

func (l *Logger) ErrorfContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.ErrorLevel, msg, []zapcore.Field{})
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) Flush() {
	_ = l.Logger.Sync()
}

func (l *Logger) DPanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.DPanicLevel, msg, fields)
	l.skipCaller.DPanic(msg, fields...)
}

func (l *Logger) DPanicfContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.DPanicLevel, msg, []zapcore.Field{})
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.PanicLevel, msg, fields)
	l.skipCaller.Panic(msg, fields...)
}

func (l *Logger) PanicfContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.PanicLevel, msg, []zapcore.Field{})
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.FatalLevel, msg, fields)
	l.skipCaller.Fatal(msg, fields...)
}

func (l *Logger) FatalfContext(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := l.logFields(ctx, zap.FatalLevel, msg, []zapcore.Field{})
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) logFields(ctx context.Context, lvl zapcore.Level, msg string, fields []zapcore.Field) []zapcore.Field {
	if lvl < l.minLevel {
		return fields
	}

	switch ctx.(type) {
	case *gin.Context:
		requestID, _ := ctx.Value(KeyRequestID).(string)
		username, _ := ctx.Value(KeyUsername).(string)
		if requestID == "" {
			fields = append(fields, zap.String(KeyRequestID, requestID))
		}
		if username == "" {
			fields = append(fields, zap.String(KeyUsername, username))
		}
		ctx = ctx.(*gin.Context).Request.Context()
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return fields
	}

	attrs := make([]attribute.KeyValue, 0, numAttr+len(fields)+len(l.extraFields))

	for _, f := range fields {
		if f.Type == zapcore.NamespaceType {
			// should this be a prefix?
			continue
		}
		attrs = appendField(attrs, f)
	}

	for _, f := range l.extraFields {
		if f.Type == zapcore.NamespaceType {
			// should this be a prefix?
			continue
		}
		attrs = appendField(attrs, f)
	}

	l.log(span, lvl, msg, attrs)

	if l.withTraceID {
		traceID := span.SpanContext().TraceID().String()
		fields = append(fields, zap.String("trace_id", traceID))
	}

	return fields
}

func (l *Logger) log(span trace.Span, lvl zapcore.Level, msg string, attrs []attribute.KeyValue) {
	attrs = append(attrs, logSeverityKey.String(levelString(lvl)))
	attrs = append(attrs, logMessageKey.String(msg))

	if l.caller {
		if fn, file, line, ok := runtimeCaller(4); ok {
			if fn != "" {
				attrs = append(attrs, semconv.CodeFunctionKey.String(fn))
			}
			if file != "" {
				attrs = append(attrs, semconv.CodeFilepathKey.String(file))
				attrs = append(attrs, semconv.CodeLineNumberKey.Int(line))
			}
		}
	}

	if l.stackTrace {
		stackTrace := make([]byte, 2048)
		n := runtime.Stack(stackTrace, false)
		attrs = append(attrs, semconv.ExceptionStacktraceKey.String(string(stackTrace[0:n])))
	}

	span.AddEvent("log", trace.WithAttributes(attrs...))

	if lvl >= l.errorStatusLevel {
		span.SetStatus(codes.Error, msg)
	}
}

func runtimeCaller(skip int) (fn, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame.Function, frame.File, frame.Line, frame.PC != 0
}

//------------------------------------------------------------------------------

// LoggerWithCtx 对Logger的封装，加入了context
type LoggerWithCtx struct {
	ctx context.Context
	l   *Logger
}

func (l LoggerWithCtx) Context() context.Context {
	return l.ctx
}

func (l LoggerWithCtx) Logger() *Logger {
	return l.l
}

func (l LoggerWithCtx) ZapLogger() *zap.Logger {
	return l.l.Logger
}

// Sugar 返回zap的sugared logger
func (l LoggerWithCtx) Sugar() SugaredLoggerWithCtx {
	return SugaredLoggerWithCtx{
		ctx: l.ctx,
		s:   l.l.Sugar(),
	}
}

func (l LoggerWithCtx) WithOptions(opts ...zap.Option) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: l.ctx,
		l:   l.l.WithOptions(opts...),
	}
}

func (l LoggerWithCtx) Clone() LoggerWithCtx {
	return LoggerWithCtx{
		ctx: l.ctx,
		l:   l.l.Clone(),
	}
}

func (l LoggerWithCtx) Debug(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.DebugLevel, msg, fields)
	l.l.skipCaller.Debug(msg, fields...)
}

func (l LoggerWithCtx) Info(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.InfoLevel, msg, fields)
	l.l.skipCaller.Info(msg, fields...)
}

func (l LoggerWithCtx) Warn(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.WarnLevel, msg, fields)
	l.l.skipCaller.Warn(msg, fields...)
}

func (l LoggerWithCtx) Error(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.ErrorLevel, msg, fields)
	l.l.skipCaller.Error(msg, fields...)
}

func (l LoggerWithCtx) DPanic(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.DPanicLevel, msg, fields)
	l.l.skipCaller.DPanic(msg, fields...)
}

func (l LoggerWithCtx) Panic(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.PanicLevel, msg, fields)
	l.l.skipCaller.Panic(msg, fields...)
}

func (l LoggerWithCtx) Fatal(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.FatalLevel, msg, fields)
	l.l.skipCaller.Fatal(msg, fields...)
}

//------------------------------------------------------------------------------

type SugaredLogger struct {
	*zap.SugaredLogger
	skipCaller *zap.SugaredLogger

	l *Logger
}

func (s *SugaredLogger) Desugar() *Logger {
	return s.l
}

func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: s.SugaredLogger.With(args...),
		l:             s.l,
	}
}

func (s *SugaredLogger) Ctx(ctx context.Context) SugaredLoggerWithCtx {
	return SugaredLoggerWithCtx{
		ctx: ctx,
		s:   s,
	}
}

func (s *SugaredLogger) DebugfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.DebugLevel, template, args)
	s.Debugf(template, args...)
}

func (s *SugaredLogger) InfofContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.InfoLevel, template, args)
	s.Infof(template, args...)
}

func (s *SugaredLogger) WarnfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.WarnLevel, template, args)
	s.Warnf(template, args...)
}

func (s *SugaredLogger) ErrorfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.ErrorLevel, template, args)
	s.Errorf(template, args...)
}

func (s *SugaredLogger) DPanicfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.DPanicLevel, template, args)
	s.DPanicf(template, args...)
}

func (s *SugaredLogger) PanicfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.PanicLevel, template, args)
	s.Panicf(template, args...)
}

func (s *SugaredLogger) FatalfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.FatalLevel, template, args)
	s.Fatalf(template, args...)
}

func (s *SugaredLogger) logArgs(
	ctx context.Context, lvl zapcore.Level, template string, args []interface{},
) {
	if lvl < s.l.minLevel {
		return
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	attrs := make([]attribute.KeyValue, 0, numAttr+1)
	attrs = append(attrs, logTemplateKey.String(template))

	s.l.log(span, lvl, fmt.Sprintf(template, args...), attrs)
}

func (s *SugaredLogger) DebugwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.DebugLevel, msg, keysAndValues)
	s.Debugw(msg, keysAndValues...)
}

func (s *SugaredLogger) InfowContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.InfoLevel, msg, keysAndValues)
	s.Infow(msg, keysAndValues...)
}

func (s *SugaredLogger) WarnwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.WarnLevel, msg, keysAndValues)
	s.Warnw(msg, keysAndValues...)
}

func (s *SugaredLogger) ErrorwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.ErrorLevel, msg, keysAndValues)
	s.Errorw(msg, keysAndValues...)
}

func (s *SugaredLogger) DPanicwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.DPanicLevel, msg, keysAndValues)
	s.DPanicw(msg, keysAndValues...)
}

func (s *SugaredLogger) PanicwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.PanicLevel, msg, keysAndValues)
	s.Panicw(msg, keysAndValues...)
}

func (s *SugaredLogger) FatalwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	keysAndValues = s.logKVs(ctx, zap.FatalLevel, msg, keysAndValues)
	s.Fatalw(msg, keysAndValues...)
}

func (s *SugaredLogger) logKVs(
	ctx context.Context, lvl zapcore.Level, msg string, kvs []interface{},
) []interface{} {
	if lvl < s.l.minLevel {
		return kvs
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return kvs
	}

	attrs := make([]attribute.KeyValue, 0, numAttr+len(kvs))

	for i := 0; i < len(kvs); i += 2 {
		if key, ok := kvs[i].(string); ok {
			attrs = append(attrs, otelutil.Attribute(key, kvs[i+1]))
		}
	}

	s.l.log(span, lvl, msg, attrs)

	if s.l.withTraceID {
		traceID := span.SpanContext().TraceID().String()
		kvs = append(kvs, "trace_id", traceID)
	}

	return kvs
}

//------------------------------------------------------------------------------

type SugaredLoggerWithCtx struct {
	ctx context.Context
	s   *SugaredLogger
}

func (s SugaredLoggerWithCtx) Desugar() *Logger {
	return s.s.Desugar()
}

func (s SugaredLoggerWithCtx) Debugf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.DebugLevel, template, args)
	s.s.skipCaller.Debugf(template, args...)
}

func (s SugaredLoggerWithCtx) Infof(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.InfoLevel, template, args)
	s.s.skipCaller.Infof(template, args...)
}

func (s SugaredLoggerWithCtx) Warnf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.WarnLevel, template, args)
	s.s.skipCaller.Warnf(template, args...)
}

func (s SugaredLoggerWithCtx) Errorf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.ErrorLevel, template, args)
	s.s.skipCaller.Errorf(template, args...)
}

func (s SugaredLoggerWithCtx) DPanicf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.DPanicLevel, template, args)
	s.s.skipCaller.DPanicf(template, args...)
}

func (s SugaredLoggerWithCtx) Panicf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.PanicLevel, template, args)
	s.s.skipCaller.Panicf(template, args...)
}

func (s SugaredLoggerWithCtx) Fatalf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.FatalLevel, template, args)
	s.s.skipCaller.Fatalf(template, args...)
}

func (s SugaredLoggerWithCtx) Debugw(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.DebugLevel, msg, keysAndValues)
	s.s.skipCaller.Debugw(msg, keysAndValues...)
}

func (s SugaredLoggerWithCtx) Infow(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.InfoLevel, msg, keysAndValues)
	s.s.skipCaller.Infow(msg, keysAndValues...)
}

func (s SugaredLoggerWithCtx) Warnw(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.WarnLevel, msg, keysAndValues)
	s.s.skipCaller.Warnw(msg, keysAndValues...)
}

func (s SugaredLoggerWithCtx) Errorw(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.ErrorLevel, msg, keysAndValues)
	s.s.skipCaller.Errorw(msg, keysAndValues...)
}

func (s SugaredLoggerWithCtx) DPanicw(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.DPanicLevel, msg, keysAndValues)
	s.s.skipCaller.DPanicw(msg, keysAndValues...)
}

func (s SugaredLoggerWithCtx) Panicw(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.PanicLevel, msg, keysAndValues)
	s.s.skipCaller.Panicw(msg, keysAndValues...)
}

func (s SugaredLoggerWithCtx) Fatalw(msg string, keysAndValues ...interface{}) {
	keysAndValues = s.s.logKVs(s.ctx, zap.FatalLevel, msg, keysAndValues)
	s.s.skipCaller.Fatalw(msg, keysAndValues...)
}

//------------------------------------------------------------------------------

func appendField(attrs []attribute.KeyValue, f zapcore.Field) []attribute.KeyValue {
	switch f.Type {
	case zapcore.BoolType:
		attr := attribute.Bool(f.Key, f.Integer == 1)
		return append(attrs, attr)

	case zapcore.Int8Type, zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type,
		zapcore.Uint32Type, zapcore.Uint8Type, zapcore.Uint16Type, zapcore.Uint64Type,
		zapcore.UintptrType:
		attr := attribute.Int64(f.Key, f.Integer)
		return append(attrs, attr)

	case zapcore.Float32Type, zapcore.Float64Type:
		attr := attribute.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))
		return append(attrs, attr)

	case zapcore.Complex64Type:
		s := strconv.FormatComplex(complex128(f.Interface.(complex64)), 'E', -1, 64)
		attr := attribute.String(f.Key, s)
		return append(attrs, attr)
	case zapcore.Complex128Type:
		s := strconv.FormatComplex(f.Interface.(complex128), 'E', -1, 128)
		attr := attribute.String(f.Key, s)
		return append(attrs, attr)

	case zapcore.StringType:
		attr := attribute.String(f.Key, f.String)
		return append(attrs, attr)
	case zapcore.BinaryType, zapcore.ByteStringType:
		attr := attribute.String(f.Key, string(f.Interface.([]byte)))
		return append(attrs, attr)
	case zapcore.StringerType:
		attr := attribute.String(f.Key, f.Interface.(fmt.Stringer).String())
		return append(attrs, attr)

	case zapcore.DurationType, zapcore.TimeType:
		attr := attribute.Int64(f.Key, f.Integer)
		return append(attrs, attr)
	case zapcore.TimeFullType:
		attr := attribute.Int64(f.Key, f.Interface.(time.Time).UnixNano())
		return append(attrs, attr)
	case zapcore.ErrorType:
		err := f.Interface.(error)
		typ := reflect.TypeOf(err).String()
		attrs = append(attrs, semconv.ExceptionTypeKey.String(typ))
		attrs = append(attrs, semconv.ExceptionMessageKey.String(err.Error()))
		return attrs
	case zapcore.ReflectType:
		attr := otelutil.Attribute(f.Key, f.Interface)
		return append(attrs, attr)
	case zapcore.SkipType:
		return attrs

	case zapcore.ArrayMarshalerType:
		var attr attribute.KeyValue
		arrayEncoder := &bufferArrayEncoder{
			stringsSlice: []string{},
		}
		err := f.Interface.(zapcore.ArrayMarshaler).MarshalLogArray(arrayEncoder)
		if err != nil {
			attr = attribute.String(f.Key+"_error", fmt.Sprintf("otelzap: unable to marshal array: %v", err))
		} else {
			attr = attribute.StringSlice(f.Key, arrayEncoder.stringsSlice)
		}
		return append(attrs, attr)

	case zapcore.ObjectMarshalerType:
		attr := attribute.String(f.Key+"_error", "otelzap: zapcore.ObjectMarshalerType is not implemented")
		return append(attrs, attr)

	default:
		attr := attribute.String(f.Key+"_error", fmt.Sprintf("otelzap: unknown field type: %v", f))
		return append(attrs, attr)
	}
}

func levelString(lvl zapcore.Level) string {
	if lvl == zapcore.DPanicLevel {
		return "PANIC"
	}
	return lvl.CapitalString()
}
