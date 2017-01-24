package trace

type Tracer interface {
	NewSpan() Span
}

type Span interface {
	SetLabel(key string, value interface{})
	Finish()
}
