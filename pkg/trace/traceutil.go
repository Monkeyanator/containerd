package traceutil

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/golang/glog"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

// SpanContextFromBase64String takes string and returns decoded context from it
func SpanContextFromBase64String(stringEncodedContext string) (spanContext trace.SpanContext, err error) {

	decodedContextBytes, err := base64.StdEncoding.DecodeString(stringEncodedContext)
	if err != nil {
		return trace.SpanContext{}, err
	}

	spanContext, ok := propagation.FromBinary(decodedContextBytes)
	if !ok {
		return trace.SpanContext{}, errors.New("could not convert raw bytes to trace")
	}

	return spanContext, nil

}

// SpanContextToBase64String takes context and encodes it to a string
func SpanContextToBase64String(spanContext trace.SpanContext) string {

	rawContextBytes := propagation.Binary(spanContext)
	encodedContext := base64.StdEncoding.EncodeToString(rawContextBytes)

	return encodedContext
}

// DefaultExporter returns the default trace exporter for the project
// This is Stackdriver at the moment, but will be the OpenCensus agent
func DefaultExporter() (exporter trace.Exporter, err error) {

	glog.Errorf("default exporter created")

	// Create an register a OpenCensus
	// Stackdriver Trace exporter.
	// exporter, err = stackdriver.NewExporter(stackdriver.Options{})

	// Create the Zipkin exporter.
	localEndpoint, err := openzipkin.NewEndpoint("container-runtime", "192.168.1.5:5454")
	if err != nil {
		log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://35.193.38.26:9411/api/v2/spans")
	ze := zipkin.NewExporter(reporter, localEndpoint)

	return ze, err
}
