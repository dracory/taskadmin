package shared

import (
	"net/http"
	urlpkg "net/url"
)

// URL builds a controller URL for the current endpoint
func URL(r *http.Request, controller string, params map[string]string) string {
	endpoint := r.URL.Path

	if params == nil {
		params = map[string]string{}
	}

	params["controller"] = controller

	return endpoint + Query(params)
}

// Query encodes the given parameters as a query string
func Query(queryData map[string]string) string {
	if len(queryData) == 0 {
		return ""
	}

	v := urlpkg.Values{}
	for key, value := range queryData {
		v.Set(key, value)
	}
	return "?" + v.Encode()
}

// HTTPBuildQuery encodes url.Values into a query string
func HTTPBuildQuery(queryData urlpkg.Values) string {
	return queryData.Encode()
}
