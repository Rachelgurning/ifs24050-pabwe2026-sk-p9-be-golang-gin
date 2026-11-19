package ai

type GenerateRequest struct {
	Topic                 string `json:"topic"`
	Category              string `json:"category"`
	Tone                  string `json:"tone"`
	IncludeWebSearch      *bool  `json:"include_web_search"`
	IncludeWebSearchCamel *bool  `json:"includeWebSearch"`
}

type Source struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type GeneratedPost struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Excerpt  string `json:"excerpt"`
	Content  string `json:"content"`
}

type GenerateResult struct {
	Title         string        `json:"title"`
	Category      string        `json:"category"`
	Excerpt       string        `json:"excerpt"`
	Content       string        `json:"content"`
	Post          GeneratedPost `json:"post"`
	Sources       []Source      `json:"sources"`
	SearchSources []Source      `json:"search_sources"`
	WebSearchUsed bool          `json:"web_search_used"`
	Generator     string        `json:"generator"`
}
