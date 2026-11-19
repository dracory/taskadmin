package shared

import "github.com/dracory/hb"

// ErrorAlert renders a Bootstrap danger alert for a store/config error.
func ErrorAlert(message string) string {
	return hb.Div().
		Class("alert alert-danger").
		HTML(message).
		ToHTML()
}
