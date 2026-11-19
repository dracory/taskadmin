package shared

import (
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// LayoutOptions are the assets a page can inject into the layout
type LayoutOptions struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}

// Layout is the default layout renderer. It builds a complete HTML page
// using hb.NewWebpage() with Bootstrap + Vue CDN, matching the
// shopadmin pattern. If a FuncLayout is provided in AdminOptions,
// it takes precedence over this default.
func Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options LayoutOptions) string {
	_ = r // kept for layout signature compatibility
	_ = w
	return webpageComplete(webpageTitle, webpageHtml, options).ToHTML()
}

const defaultFavicon = `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`

// webpageComplete builds the default webpage template
func webpageComplete(webpageTitle, webpageHtml string, options LayoutOptions) hb.TagInterface {
	webpage := hb.NewWebpage()
	webpage.SetTitle(webpageTitle)
	webpage.SetFavicon(defaultFavicon)

	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_3_8(),
		cdn.BootstrapIconsCss_1_13_1(),
	})
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_3_8(),
	})
	webpage.AddStyleURLs(options.StyleURLs)
	webpage.AddScriptURLs(options.ScriptURLs)
	webpage.AddStyles(options.Styles)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}`)
	webpage.AddChild(hb.NewHTML(webpageHtml))
	webpage.AddScripts(options.Scripts)

	return webpage
}
