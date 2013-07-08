package controllers

import (
	"github.com/robfig/revel"
	"code.google.com/p/go.talks/pkg/present"
	"html/template"
	"bytes"
	"path"
	"path/filepath"
		"os"
	)


type Blog struct {
	*revel.Controller
	docs     []*Doc
	tags     []string
	docPaths map[string]*Doc
	docTags  map[string][]*Doc
	template struct {
		home, index, article, doc *template.Template
	}
	atomFeed []byte // pre-rendered Atom feed
	jsonFeed []byte // pre-rendered JSON feed
}

type Doc struct {
	*present.Doc
	Permalink    string
	Path         string
	Related      []*Doc
	Newer, Older *Doc
	HTML         template.HTML // rendered article
}

func (c Blog) Index() revel.Result {

	return c.Render()
}
func (c Blog) Show(id string) revel.Result {
	var err error
	p := present.Template().Funcs(funcMap)
	tmpth := path.Join(revel.ViewsPath, "Blog")
	c.template.doc, err = p.ParseFiles(filepath.Join(tmpth, "doc.tmpl"))
	if err != nil {
		revel.ERROR.Println(err)
	}


	articlePath := path.Join(revel.AppPath, "content")
	f, err := os.Open(articlePath + "/" + id + ".article")
		if err != nil {
			revel.ERROR.Println(err)
		}
		defer f.Close()
	d, err := present.Parse(f, id, 0)
		if err != nil {
			revel.ERROR.Println(err)
		}

	html := new(bytes.Buffer)
		err = d.Render(html, c.template.doc)
		if err != nil {
			revel.ERROR.Println(err)
		}
	revel.ERROR.Println(d)

			doc :=  &Doc{
			Doc:       d,
			Path:      id,
			Permalink: "http://www.gopheracademy.com/blog/" + id,
			HTML:      template.HTML(html.String()),
		}

		c.RenderArgs["Doc"] = doc
	return c.Render()
}



// authors returns a comma-separated list of author names.
func authors(authors []present.Author) string {
	var b bytes.Buffer
	last := len(authors) - 1
	for i, a := range authors {
		if i > 0 {
			if i == last {
				b.WriteString(" and ")
			} else {
				b.WriteString(", ")
			}
		}
		b.WriteString(authorName(a))
	}
	return b.String()
}

// authorName returns the first line of the Author text: the author's name.
func authorName(a present.Author) string {
	el := a.TextElem()
	if len(el) == 0 {
		return ""
	}
	text, ok := el[0].(present.Text)
	if !ok || len(text.Lines) == 0 {
		return ""
	}
	return text.Lines[0]
}

var funcMap = template.FuncMap{
	"sectioned": sectioned,
	"authors":   authors,
}