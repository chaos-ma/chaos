package trace

/**
* created by mengqi on 2023/11/13
 */

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"sync"

	"github.com/chaos-ma/chaos/log"
)

/*
初始化不同的export的设置
*/

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
)

var (
	//set ,struct 空结构体不占内存， zerobase
	agents = make(map[string]struct{})
	lock   sync.Mutex
)

func InitAgent(o Options) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := agents[o.Endpoint]
	if ok {
		return
	}
	err := startAgent(o)
	if err != nil {
		return
	}
	agents[o.Endpoint] = struct{}{}
}

func startAgent(o Options) error {
	var sexp trace.SpanExporter
	var err error

	opts := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(o.Sampler))),
		trace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(o.Name))),
	}

	if len(o.Endpoint) > 0 {
		switch o.Batcher {
		case kindJaeger:
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
			if err != nil {
				return err
			}
		case kindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err
			}
		}
		opts = append(opts, trace.WithBatcher(sexp))
	}

	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))
	return nil
}
