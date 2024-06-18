package tracer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TraceConfig struct {
	ServiceName string            // 服务名称
	Ratio       float64           // 采样率
	Exporter    string            // tracer接收者（otlp）
	Endpoint    string            // 远程输出日志路径(Tempo)
	Params      map[string]string // 预留参数字段（有些认证信息可以通过该字段放到header中）
	Path        string            // 本地输出日志路径(当trace接收暂不支持时，默认输出到本地)
}

// InitTracer 设置全局trace
func InitTracer(ctx context.Context, cnf *TraceConfig) (err error) {
	var spanExporter tracesdk.SpanExporter
	switch strings.ToLower(cnf.Exporter) {
	case "otlp":
		spanExporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cnf.Endpoint),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithHeaders(cnf.Params),
		)
		if err != nil {
			return err
		}
	case "file":
		writer, err := filetWriter(cnf.Path)
		if err != nil {
			return err
		}
		spanExporter, err = stdouttrace.New(stdouttrace.WithWriter(writer))
		if err != nil {
			return err
		}
	default:
		spanExporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
		if err != nil {
			return err
		}
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(cnf.Ratio))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(spanExporter),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(cnf.ServiceName),
			attribute.String("exporter", cnf.Exporter),
			attribute.Float64("ratio", cnf.Ratio),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func filetWriter(logPath string) (io.Writer, error) {
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s.log", logPath, "%Y%m%d"),
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		slog.Error("初始化链路日志出错:%v", err)
		return nil, err
	}
	return writer, nil
}
