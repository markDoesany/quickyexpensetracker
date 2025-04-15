package templates

type Button struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload,omitempty"`
	URL     string `json:"url,omitempty"`
}

type Template struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle,omitempty"`
	ImageURL string   `json:"image_url"`
	Buttons  []Button `json:"buttons"`
}
