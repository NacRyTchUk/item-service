package middleware

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func HandleRoutingError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	_ = HttpResponseModifier(ctx, w, nil)

	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}
