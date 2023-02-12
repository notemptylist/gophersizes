package cyoa

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const tmpl = "template/template.html"

type Story map[string]StoryArc
type ArcOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type StoryArc struct {
	Title   string      `json:"title,omitempty"`
	Story   []string    `json:"story,omitempty"`
	Options []ArcOption `json:"options,omitempty"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles(tmpl))
}

// ParseFile decodes the JSON file and returns a Story
func ParseFile(fname string) (Story, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)
	var story Story
	if err = d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s Story
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; p {
	case "/":
		if err := tpl.Execute(w, h.s["intro"]); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println(err)
		}
	case "/favicon.ico":
		w.WriteHeader(http.StatusFound)
	default:
		p = strings.TrimLeft(p, "/")
		arc, ok := h.s[p]
		if ok {
			if err := tpl.Execute(w, arc); err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				log.Println(err)
			}
		} else {
			http.Error(w, "Invalid arc", http.StatusFound)
		}
	}
}
