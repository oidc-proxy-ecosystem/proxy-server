package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/plugin/navbar/files"
	"github.com/oidc-proxy-ecosystem/proxy-server/shared"
)

func setHead(doc *goquery.Document, str string) {
	doc.Find("head").AppendHtml(str)
}

func setLink(doc *goquery.Document, url string) {
	link := fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, url)
	setHead(doc, link)
}

func setScript(doc *goquery.Document, url string) {
	script := fmt.Sprintf(`<script src="%s"></script>`, url)
	setHead(doc, script)
}

func setBody(doc *goquery.Document, str string) {
	doc.Find("body").AppendHtml(str)
}

func setCSS(doc *goquery.Document) {
	if css, err := files.Templates.ReadFile("templates/css.html"); err != nil {
		panic(err)
	} else {
		setHead(doc, string(css))
	}
}

func setVuetify(doc *goquery.Document) {
	navBarjs, _ := files.Templates.ReadFile("templates/vuetify-bar.html")
	docBody, _ := doc.Find("body").Clone().Html()
	doc.Find("body").Empty()
	setBody(doc, string(navBarjs))
	setBody(doc, docBody)
}

func setNav(doc *goquery.Document) {
	navBarjs, _ := files.Templates.ReadFile("templates/navbarjs.html")
	// userImg, _ := files.Images.ReadFile("images/user.png")
	// userImgBase64 := base64.StdEncoding.EncodeToString(userImg)
	// t := template.Must(template.ParseFS(files.Templates, "templates/navbar.html"))
	// buf := new(bytes.Buffer)
	// t.Execute(buf, map[string]string{
	// 	"UserPicture": userImgBase64,
	// })
	// doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	docBody, _ := doc.Find("body").Clone().Html()
	doc.Find("body").Empty()
	setBody(doc, string(navBarjs))
	setBody(doc, docBody)
}

type navBar struct {
}

var _ shared.Response = (*navBar)(nil)

func (h *navBar) Modify(URL string, method string, header map[string][]string, body []byte) shared.ResponseResult {
	var resultBody []byte
	httpHeader := http.Header(header)
	if strings.HasPrefix(httpHeader.Get("Content-Type"), "text/html") {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
		setLink(doc, "https://unpkg.com/element-ui/lib/theme-chalk/index.css")
		setLink(doc, "https://fonts.googleapis.com/css?family=Roboto:100,300,400,500,700,900")
		setLink(doc, "https://cdn.jsdelivr.net/npm/@mdi/font@4.x/css/materialdesignicons.min.css")
		setLink(doc, "https://cdn.jsdelivr.net/npm/vuetify@2.x/dist/vuetify.min.css")
		setScript(doc, "https://unpkg.com/vue/dist/vue.js")
		// setScript(doc, "https://unpkg.com/element-ui/lib/index.js")
		setScript(doc, "https://cdn.jsdelivr.net/npm/vuetify@2.x/dist/vuetify.js")
		setCSS(doc)
		// setNav(doc)
		setVuetify(doc)
		v, _ := doc.Html()
		resultBody = []byte(v)
	} else {
		resultBody = body
	}
	return shared.ResponseResult{
		Header: header,
		Body:   resultBody,
	}
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				"response_modify": &internal.ResponsePlugin{Impl: &navBar{}},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
