package controllers

import (
	"github.com/robfig/revel"
		"code.google.com/p/go.talks/pkg/present"
			"html/template"
	"bytes"
	"html"
	"strings"
	       "unicode"
        "unicode/utf8"
        "fmt"
              "net/url"
              "strconv"

)

func init() {
	revel.OnAppStart(InitDB)
	revel.TemplateFuncs["sectioned"] = sectioned
		revel.TemplateFuncs["elem"] = renderElem
		revel.TemplateFuncs["style"] = Style
        revel.TemplateFuncs["image"] = parseImage

}

// sectioned returns true if the provided Doc contains more than one section.
// This is used to control whether to display the table of contents and headings.
func sectioned(d *present.Doc) bool {
	return len(d.Sections) > 1
}

func renderElem(t *template.Template, e present.Elem) (template.HTML, error) {
        var data interface{} = e
        if s, ok := e.(present.Section); ok {
                data = struct {
                        present.Section
                        Template *template.Template
                }{s, t}
        }
        return execTemplate(t, e.TemplateName(), data)
}

// execTemplate is a helper to execute a template and return the output as a
// template.HTML value.
func execTemplate(t *template.Template, name string, data interface{}) (template.HTML, error) {
        b := new(bytes.Buffer)
        err := t.ExecuteTemplate(b, name, data)
        if err != nil {
                return "", err
        }
        return template.HTML(b.String()), nil
}

// Style returns s with HTML entities escaped and font indicators turned into
// HTML font tags.
func Style(s string) template.HTML {
        return template.HTML(font(html.EscapeString(s)))
}

// font returns s with font indicators turned into HTML font tags.
func font(s string) string {
        if strings.IndexAny(s, "[`_*") == -1 {
                return s
        }
        words := split(s)
        var b bytes.Buffer
Word:
        for w, word := range words {
                if len(word) < 2 {
                        continue Word
                }
                if link, _ := parseInlineLink(word); link != "" {
                        words[w] = link
                        continue Word
                }
                const punctuation = `.,;:()!?—–'"`
                const marker = "_*`"
                // Initial punctuation is OK but must be peeled off.
                first := strings.IndexAny(word, marker)
                if first == -1 {
                        continue Word
                }
                // Is the marker prefixed only by punctuation?
                for _, r := range word[:first] {
                        if !strings.ContainsRune(punctuation, r) {
                                continue Word
                        }
                }
                open, word := word[:first], word[first:]
                char := word[0] // ASCII is OK.
                close := ""
                switch char {
                default:
                        continue Word
                case '_':
                        open += "<i>"
                        close = "</i>"
                case '*':
                        open += "<b>"
                        close = "</b>"
                case '`':
                        open += "<code>"
                        close = "</code>"
                }
                // Terminal punctuation is OK but must be peeled off.
                last := strings.LastIndex(word, word[:1])
                if last == 0 {
                        continue Word
                }
                head, tail := word[:last+1], word[last+1:]
                for _, r := range tail {
                        if !strings.ContainsRune(punctuation, r) {
                                continue Word
                        }
                }
                b.Reset()
                b.WriteString(open)
                var wid int
                for i := 1; i < len(head)-1; i += wid {
                        var r rune
                        r, wid = utf8.DecodeRuneInString(head[i:])
                        if r != rune(char) {
                                // Ordinary character.
                                b.WriteRune(r)
                                continue
                        }
                        if head[i+1] != char {
                                // Inner char becomes space.
                                b.WriteRune(' ')
                                continue
                        }
                        // Doubled char becomes real char.
                        // Not worth worrying about "_x__".
                        b.WriteByte(char)
                        wid++ // Consumed two chars, both ASCII.
                }
                b.WriteString(close) // Write closing tag.
                b.WriteString(tail)  // Restore trailing punctuation.
                words[w] = b.String()
        }
        return strings.Join(words, "")
}

// split is like strings.Fields but also returns the runs of spaces
// and treats inline links as distinct words.
func split(s string) []string {
        var (
                words = make([]string, 0, 10)
                start = 0
        )

        // appendWord appends the string s[start:end] to the words slice.
        // If the word contains the beginning of a link, the non-link portion
        // of the word and the entire link are appended as separate words,
        // and the start index is advanced to the end of the link.
        appendWord := func(end int) {
                if j := strings.Index(s[start:end], "[["); j > -1 {
                        if _, l := parseInlineLink(s[start+j:]); l > 0 {
                                // Append portion before link, if any.
                                if j > 0 {
                                        words = append(words, s[start:start+j])
                                }
                                // Append link itself.
                                words = append(words, s[start+j:start+j+l])
                                // Advance start index to end of link.
                                start = start + j + l
                                return
                        }
                }
                // No link; just add the word.
                words = append(words, s[start:end])
                start = end
        }

        wasSpace := false
        for i, r := range s {
                isSpace := unicode.IsSpace(r)
                if i > start && isSpace != wasSpace {
                        appendWord(i)
                }
                wasSpace = isSpace
        }
        for start < len(s) {
                appendWord(len(s))
        }
        return words
}

type Link struct {
        URL   *url.URL
        Label string
}

func (l Link) TemplateName() string { return "link" }

func parseLink(ctx *present.Context, fileName string, lineno int, text string) (present.Elem, error) {
        args := strings.Fields(text)
        url, err := url.Parse(args[1])
        if err != nil {
                return nil, err
        }
        label := ""
        if len(args) > 2 {
                label = strings.Join(args[2:], " ")
        } else {
                scheme := url.Scheme + "://"
                if url.Scheme == "mailto" {
                        scheme = "mailto:"
                }
                label = strings.Replace(url.String(), scheme, "", 1)
        }
        return Link{url, label}, nil
}

func renderLink(url, text string) string {
        text = font(text)
        if text == "" {
                text = url
        }
        return fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, url, text)
}

// parseInlineLink parses an inline link at the start of s, and returns
// a rendered HTML link and the total length of the raw inline link.
// If no inline link is present, it returns all zeroes.
func parseInlineLink(s string) (link string, length int) {
        if len(s) < 2 || s[:2] != "[[" {
                return
        }
        end := strings.Index(s, "]]")
        if end == -1 {
                return
        }
        urlEnd := strings.Index(s, "]")
        url := s[2:urlEnd]
        const badURLChars = `<>"{}|\^[] ` + "`" // per RFC2396 section 2.4.3
        if strings.ContainsAny(url, badURLChars) {
                return
        }
        if urlEnd == end {
                return renderLink(url, ""), end + 2
        }
        if s[urlEnd:urlEnd+2] != "][" {
                return
        }
        text := s[urlEnd+2 : end]
        return renderLink(url, text), end + 2
}

type Image struct {
        URL    string
        Width  int
        Height int
}

func (i Image) TemplateName() string { return "image" }

func parseImage(ctx *present.Context, fileName string, lineno int, text string) (present.Elem, error) {
        args := strings.Fields(text)
        img := Image{URL: args[1]}
        a, err := parseArgs(fileName, lineno, args[2:])
        if err != nil {
                return nil, err
        }
        switch len(a) {
        case 0:
                // no size parameters
        case 2:
                if v, ok := a[0].(int); ok {
                        img.Height = v
                }
                if v, ok := a[1].(int); ok {
                        img.Width = v
                }
        default:
                return nil, fmt.Errorf("incorrect image invocation: %q", text)
        }
        return img, nil
}
func parseArgs(name string, line int, args []string) (res []interface{}, err error) {
        res = make([]interface{}, len(args))
        for i, v := range args {
                if len(v) == 0 {
                        return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
                }
                switch v[0] {
                case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
                        n, err := strconv.Atoi(v)
                        if err != nil {
                                return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
                        }
                        res[i] = n
                case '/':
                        if len(v) < 2 || v[len(v)-1] != '/' {
                                return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
                        }
                        res[i] = v
                case '$':
                        res[i] = "$"
                default:
                        return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
                }
        }
        return
}
