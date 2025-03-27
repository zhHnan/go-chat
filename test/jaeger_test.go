package test

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"testing"
)

// Test_jaeger 测试 Jaeger Tracer 的初始化和使用
// 该测试函数配置 Jaeger Tracer，并创建一个或多个 Span 以模拟跟踪请求
func Test_jaeger(t *testing.T) {
	// 创建一个 Jaeger 配置
	// 配置包括服务名称、采样器配置和报告器配置
	cfg := jaegercfg.Configuration{
		ServiceName: "my-service",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: fmt.Sprintf("http://%s/api/traces", "192.168.232.12:14268"),
		},
	}

	// 创建 Jaeger Tracer
	// 如果初始化失败，记录错误并返回
	Jaeger, err := cfg.InitGlobalTracer("client_test", jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		t.Log(err)
		return
	}
	defer Jaeger.Close()

	// 创建一个 Jaeger Span
	// 这里创建的是一个父 Span，它将作为其他 Span 的父级
	tracer := opentracing.GlobalTracer()
	parentSpan := tracer.StartSpan("A")
	defer parentSpan.Finish()

	// 调用函数 B，并传入 tracer 和 parentSpan
	// 函数 B 将创建一个子 Span，并将其链接到 parentSpan
	B(tracer, parentSpan)
}

// B 创建一个子 Span，并将其作为父 Span 的子级
// 该函数演示了如何在分布式系统中创建和链接 Span
func B(tracer opentracing.Tracer, parentSpan opentracing.Span) {
	// 创建一个子 Span
	// 使用父 Span 的上下文来建立父子关系
	childSpan := tracer.StartSpan("B", opentracing.ChildOf(parentSpan.Context()))
	defer childSpan.Finish()
}
