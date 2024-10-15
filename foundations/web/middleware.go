package web

// MidHandler is a type that represents a middleware handler, which is a function that takes a handler and returns a handler.
// This is a common pattern in Go for creating middleware chains & run code before and after another handler.
type MidHandler func(Handler) Handler

func wrapMiddleware(mw []MidHandler, hanlder Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := mw[i]
		if mwFunc != nil {
			hanlder = mwFunc(hanlder) // Equate this to m(profile.Authenticate)
			// You get back h (hanlder) with a hdl temp handler wrapping around a call to profile.Authenticate by virtue of closure
			// and subsequent call to real app middleware which will also invoke the hdl() causing profile.Authenticate to execure.
			// So, when h() swished, all exec went-in, hdl constructed w/ call to profile.Authenticate & app middlware invoke.
			// During app middlware invoke, hdl() swished causing profile.Authenticate to run.
		}
	}
	return hanlder
}

// m := func(handler web.Handler) web.Handler {
// 	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 		hdl := func(ctx context.Context) error {
// 			return handler(ctx, w, r)
// 		}
// 		return middleware.Authenticate(ctx, log, client, r.Header.Get("authorization"), hdl)
// 	}
// 	return h
// }
// return m
