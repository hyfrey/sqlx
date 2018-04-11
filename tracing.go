package sqlx

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type tracerx struct {
	tracer     opentracing.Tracer
	driverName string
}

func (tx tracerx) startSpanFromContext(ctx context.Context, query string, args ...interface{}) (context.Context, opentracing.Span) {
	if tx.tracer == nil {
		return ctx, nil
	}

	var pspanCtx opentracing.SpanContext
	if pspan := opentracing.SpanFromContext(ctx); pspan != nil {
		pspanCtx = pspan.Context()
	}

	span := tx.tracer.StartSpan(firstWord(query), opentracing.ChildOf(pspanCtx))
	ext.Component.Set(span, "sqlx")
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, tx.driverName)
	span.LogFields(log.Object("args", args))

	return opentracing.ContextWithSpan(ctx, span), span
}

func (tx tracerx) finishSpan(span opentracing.Span, err error) {
	if tx.tracer == nil {
		return
	}

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("event", "error"),
			log.String("message", err.Error()),
		)
	}
	span.Finish()
}

// extract first word of query
func firstWord(query string) string {
	for i, c := range query {
		if c == ' ' {
			return strings.ToLower(query[0:i])
		}
	}
	return strings.ToLower(query)
}
