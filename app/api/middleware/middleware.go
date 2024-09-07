// This package provides the 'app' layer middleware which is protocol agnostic.

package middleware

import "context"

type Handler func(context.Context) error
