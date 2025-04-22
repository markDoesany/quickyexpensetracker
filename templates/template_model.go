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

type AttachmentPayload struct {
	TemplateType string        `json:"template_type"`
	Elements     []interface{} `json:"elements"`
}

type Attachment struct {
	Type    string            `json:"type"`
	Payload AttachmentPayload `json:"payload"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Attachment Attachment `json:"attachment"`
}

type RequestPayload = struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
}
