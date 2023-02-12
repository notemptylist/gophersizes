package cyoa

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
