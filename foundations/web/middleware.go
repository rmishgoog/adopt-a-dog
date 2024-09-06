package web

// MidHandler is a type that represents a middleware handler, which is a function that takes a handler and returns a handler.
// This is a common pattern in Go for creating middleware chains & run code before and after another handler.
type MidHandler func(Handler) Handler

func wrapMiddleware(mw []MidHandler, hanlder Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := mw[i]
		if mwFunc != nil {
			hanlder = mwFunc(hanlder)
		}
	}
	return hanlder
}
