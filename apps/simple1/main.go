package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

const simple1Version = "0.1.0"

type Response struct {
	Messages []string `json:"messages"`
}

func getIndex(c echo.Context) error {
	_, span := tracer.Start(c.Request().Context(), "getIndex")
	defer span.End()
	messages := []string{
		generateSimple1Message(),
	}

	return c.JSON(http.StatusOK, Response{Messages: messages})
}

func generateSimple1Message() string {
	return fmt.Sprintf("Hello from Simple1(v%s)!", simple1Version)
}

var tracer = otel.Tracer("k8s-in-the-house.com/simple1")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	exporter, err := newJaegerExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}
	tp, err := newTraceProvider(ctx, exporter)
	if err != nil {
		log.Fatal(err)
	}
	tracer = tp.Tracer("k8s-in-the-house.com/simple1")
	defer func() {
		if err := tp.ForceFlush(ctx); err != nil {
			log.Printf("Error flush tracer provider: %v", err)
		}
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	r := echo.New()
	r.Use(otelecho.Middleware("k8s-in-the-house.com/simple1"))

	r.GET("/", getIndex)
	go func() {
		if err = r.Start(":12345"); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			if err := r.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func newJaegerExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	conn, err := grpc.DialContext(
		ctx, "trace-collector-collector.default.svc.cluster.local:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func newTraceProvider(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("k8s-in-the-house.com/simple1"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tp.Tracer("k8s-in-the-house.com/simple1")
	return tp, nil
}
