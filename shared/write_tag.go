package shared

import (
	"net/http"

	"github.com/dracory/hb"
)

// WriteTag writes an hb tag to the response with an HTML content type.
func WriteTag(w http.ResponseWriter, tag hb.TagInterface) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	// #nosec G705 -- html is rendered by the page builder, not user input
	_, _ = w.Write([]byte(tag.ToHTML()))
}
