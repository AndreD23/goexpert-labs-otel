package main

import (
	"context"
	"github.com/AndreD23/goexpert-labs-otel/serviceb/configs"
	"github.com/AndreD23/goexpert-labs-otel/serviceb/internal/handlers"
	"github.com/AndreD23/goexpert-labs-otel/serviceb/internal/weatherapi"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer()
	if err != nil {
		log.Fatal("Init Provider error: ", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	config := configs.NewConfig()
	weather := weatherapi.NewWeatherAPI(config.WeatherAPIKey)
	temperatureHandler := handlers.New(weather)

	r := chi.NewRouter()
	r.Get("/{zipCode}", temperatureHandler.GetTemperature)
	http.ListenAndServe(":8080", r)
}

func initTracer() (func(ctx context.Context) error, error) {
	traceExporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("serviceb"),
		)),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}
