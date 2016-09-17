package viewbag

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"

	// "github.com/davecgh/go-spew/spew"
	"github.com/unrolled/render"
)

func New(w http.ResponseWriter, req *http.Request, db *runner.DB) *viewBag {
	viewBag := viewBag{}
	viewBag.Data = viewGlobals
	viewBag.renderer = r
	viewBag.w = w
	viewBag.db = db
	viewBag.db = db
	return &viewBag
}

func (viewBag *viewBag) Add(key string, value interface{}) {
	viewBag.Data[key] = value
	// spew.Dump(viewBag.data)
}

// func (viewBag *viewBag) AddStruct(key string, value interface{}) {
// 	viewBag.data[key] = func() interface{} {
// 		return value
// 	}
// }
func (viewBag *viewBag) LoadNavItems() {
	var navItems []*NavItem
	err := viewBag.db.
		Select("title", "slug").
		From("pages").
		QueryStructs(&navItems)
	if err != nil {
		panic(err)
	}
	viewBag.Add("NavItems", navItems)
}

func (viewBag *viewBag) Render(status int, templateName string) {
	// spew.Dump(viewBag.data)
	viewBag.renderer.HTML(viewBag.w, status, templateName, viewBag)
}

type viewBag struct {
	renderer *render.Render
	db       *runner.DB
	w        http.ResponseWriter
	res      *http.Request
	Data     map[string]interface{}
}

var viewGlobals = map[string]interface{}{
	"HeaderDate": time.Now(),
	"Copyright":  time.Now().Year(),
}

var templateFunctions = template.FuncMap{
	"javascript": javascriptTag,
	"sass":       sassTag,
	"stylesheet": stylesheetTag,
	"image":      imageTag,
	"imagepath":  imagePath,
	"content":    content,
	"htmlblock":  htmlblock,
	"navigation": navigation,
}

var r = render.New(render.Options{
	Layout:     "application",
	Extensions: []string{".html"},
	Funcs:      []template.FuncMap{templateFunctions},
})

func content(contents ...string) template.HTML {
	var str string
	for _, content := range contents {
		str += "<div class='standard'>" + content + "</standard>"
	}
	return template.HTML(str)
}

func javascriptTag(names ...string) template.HTML {
	var str string
	for _, name := range names {
		str += "<script src='assets/javascripts/" + name + ".js' type='text/javascript'></script>"
	}
	return template.HTML(str)
}

func sassTag(names ...string) template.HTML {
	var str string
	for _, name := range names {
		str += "<link rel='stylesheet' href='assets/sass/" + name + ".scss' type='text/css' media='screen' />\n"
	}
	return template.HTML(str)
}

func stylesheetTag(names ...string) template.HTML {
	var str string
	for _, name := range names {
		str += "<link rel='stylesheet' href='assets/stylesheets/" + name + ".css' type='text/css' media='screen'  />\n"
	}
	return template.HTML(str)
}

func imagePath(name string) string {
	return "assets/images/" + name
}

func imageTag(name string, class string) template.HTML {
	return template.HTML("<image src='" + imagePath(name) + " class='" + class + "' />")
}

func htmlblock(page *Page, code string) template.HTML {
	html := "<div class='textblock editable' "
	html += " data-textblock='page-" + strconv.FormatInt(page.PageID, 10) + "-" + code + "'"
	html += " data-placeholder='#{placeholder}'> "
	html += getHTMLFromTextblock(page, code)
	html += "</div>"
	return template.HTML(html)
}

func navigation(viewBag *viewBag) template.HTML {
	html := ""
	if viewBag.Data["NavItems"] != nil {
		navItems := viewBag.Data["NavItems"].([]*NavItem)
		html = "<nav class='main-nav closed'>"
		for _, navItem := range navItems {
			html += "<a href='/" + navItem.Slug + "'>" + navItem.Title + "</a>"
		}
		html += "</nav>"
	}
	return template.HTML(html)
}

type Page struct {
	PageID     int64        `db:"page_id"`
	Title      string       `db:"title"`
	Body       string       `db:"body"`
	Slug       string       `db:"slug"`
	Template   string       `db:"template"`
	CreatedAt  dat.NullTime `db:"created_at"`
	UpdatedAt  dat.NullTime `db:"updated_at"`
	Textblocks []*Textblock
}

type NavItem struct {
	Title string `db:"title"`
	Slug  string `db:"slug"`
}

func (navItem *NavItem) getURL() string {
	return ""
}

type Textblock struct {
	TextblockID int64        `db:"textblock_id"`
	Code        string       `db:"code"`
	Body        string       `db:"body"`
	CreatedAt   dat.NullTime `db:"created_at"`
	UpdatedAt   dat.NullTime `db:"updated_at"`
	PageID      int64        `db:"page_id"`
}

func getHTMLFromTextblock(page *Page, code string) string {
	var body string
	for _, tb := range page.Textblocks {
		if tb.Code == code {
			body = tb.Body
		}
	}
	return body
}
