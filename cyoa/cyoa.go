package cyoa

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const defaultTemplate = "template/template.html"

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

type HandlerOption func(h *handler)

func WithTemplate(t string) HandlerOption {
	return func(h *handler) {
		tpl := template.Must(template.ParseFiles(t))
		h.t = tpl
	}
}

type handler struct {
	s Story
	t *template.Template
}

// NewHandler creates and returns a handler type value
func NewHandler(s Story, opts ...HandlerOption) handler {
	h := handler{s: s}
	for _, o := range opts {
		o(&h)
	}
	if h.t == nil {
		h.t = template.Must(template.ParseFiles(defaultTemplate))
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; p {
	case "/":
		if err := h.t.Execute(w, h.s["intro"]); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println(err)
		}
	case "/favicon.ico":
		w.WriteHeader(http.StatusFound)
	default:
		p = strings.TrimLeft(p, "/")
		arc, ok := h.s[p]
		if ok {
			if err := h.t.Execute(w, arc); err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				log.Println(err)
			}
		} else {
			http.Error(w, "Invalid arc", http.StatusFound)
		}
	}
}
