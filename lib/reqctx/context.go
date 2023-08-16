package reqctx

import (
	"net/http"
)

type String string

func GetContext(key String, r *http.Request) interface{} {
	parameter := r.Context().Value(key)
	// if parameter == nil {
	// 	panic(apierrors.NewClientError(errors.New("context parameter '" + key + "' missing")))
	// }
	return parameter
}
