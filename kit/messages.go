package kit

import "context"

// Messages are the few texts the kit renders itself (mostly aria-labels). An app sets them once per
// request with WithMessages; a component's own props still override them per call. The kit has no
// translation files: an app with several languages picks the Messages for the request's language.
type Messages struct {
	Close    string
	Pages    string
	PrevPage string
	NextPage string
}

var (
	English = Messages{Close: "Close", Pages: "Pages", PrevPage: "Previous page", NextPage: "Next page"}
	German  = Messages{Close: "Schließen", Pages: "Seiten", PrevPage: "Vorige Seite", NextPage: "Nächste Seite"}
)

type messagesKey struct{}

// WithMessages returns a context whose components render m.
func WithMessages(ctx context.Context, m Messages) context.Context {
	return context.WithValue(ctx, messagesKey{}, m)
}

// T returns the context's Messages, English by default.
func T(ctx context.Context) Messages {
	if m, ok := ctx.Value(messagesKey{}).(Messages); ok {
		return m
	}
	return English
}

// Or returns s, or fallback when s is empty.
func Or(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}
