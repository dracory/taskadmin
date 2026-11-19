package shared

import (
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// Layout is the interface a host-provided page layout must implement.
// It is stateful — each request gets a fresh instance via FuncLayout.
type Layout interface {
	SetTitle(title string)
	SetScriptURLs(scripts []string)
	SetScripts(scripts []string)
	SetStyleURLs(styles []string)
	SetStyles(styles []string)
	SetBody(string)
	Render(w http.ResponseWriter, r *http.Request) string
}

const defaultFavicon = `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`

// defaultLayout wraps *hb.HtmlWebpage to implement the Layout interface.
type defaultLayout struct {
	webpage *hb.HtmlWebpage
}

// SetTitle sets the webpage title
func (l *defaultLayout) SetTitle(title string) { l.webpage.SetTitle(title) }

// SetScriptURLs sets the webpage script URLs
func (l *defaultLayout) SetScriptURLs(scripts []string) { l.webpage.ScriptURLs(scripts) }

// SetScripts sets the webpage inline scripts
func (l *defaultLayout) SetScripts(scripts []string) { l.webpage.AddScripts(scripts) }

// SetStyleURLs sets the webpage style URLs
func (l *defaultLayout) SetStyleURLs(styles []string) { l.webpage.StyleURLs(styles) }

// SetStyles sets the webpage inline styles
func (l *defaultLayout) SetStyles(styles []string) { l.webpage.AddStyles(styles) }

// SetBody sets the webpage body HTML
func (l *defaultLayout) SetBody(html string) { l.webpage.AddChild(hb.NewHTML(html)) }

// Render returns the full webpage HTML
func (l *defaultLayout) Render(_ http.ResponseWriter, _ *http.Request) string {
	return l.webpage.ToHTML()
}

// DefaultWebpage provides a built-in default webpage layout with Bootstrap,
// Vue, jQuery, and SweetAlert2. Used when no FuncLayout is provided.
//
// SECURITY WARNING: The content parameter is rendered as raw HTML without
// sanitization. Ensure content is trusted or properly sanitized before
// passing to this function to avoid XSS vulnerabilities.
func DefaultWebpage(title, content string) Layout {
	webpage := hb.NewWebpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(defaultFavicon)

	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_3_8(),
	})
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_3_8(),
		cdn.Jquery_3_7_1(),
		cdn.VueJs_3(),
		cdn.Sweetalert2_11(),
	})
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}
	.form-select {
		display: block;
		width: 100%;
		padding: .375rem 2.25rem .375rem .75rem;
		font-size: 1rem;
		font-weight: 400;
		line-height: 1.5;
		color: #212529;
		background-color: #fff;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='none' stroke='%23343a40' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M2 5l6 6 6-6'/%3e%3c/svg%3e");
		background-repeat: no-repeat;
		background-position: right .75rem center;
		background-size: 16px 12px;
		border: 1px solid #ced4da;
		border-radius: .25rem;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
	}`)
	webpage.AddChild(hb.NewHTML(content))
	return &defaultLayout{webpage: webpage}
}
