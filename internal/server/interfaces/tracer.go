//go:generate mockgen -destination=../mocks/tracer_mock.go -source=./tracer.go -package=mocks

package interfaces

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer interface for Tracer
type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

// Span interface for Span
type Span interface {
	End(options ...trace.SpanEndOption)
	AddEvent(name string, options ...trace.EventOption)
	IsRecording() bool
	RecordError(err error, options ...trace.EventOption)
	SpanContext() trace.SpanContext
	SetStatus(code codes.Code, description string)
	SetName(name string)
	SetAttributes(kv ...attribute.KeyValue)
	TracerProvider() trace.TracerProvider
}
