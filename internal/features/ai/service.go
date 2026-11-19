package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/internal/features/post"
	"app/internal/response"
)

const systemPrompt = "Anda adalah penulis blog profesional berbahasa Indonesia. Balas HANYA dengan satu objek JSON valid berisi kunci \"title\", \"category\", \"excerpt\", dan \"content\". Nilai \"content\" berformat Markdown, gunakan subjudul level 2 (##) dan level 3 (###), JANGAN menulis judul level 1. Jika ada hasil penelusuran web, gunakan sebagai fakta terkini dan beri sitasi [n] sesuai nomor sumber. \"excerpt\" berisi 1-2 kalimat pemikat."

type Service struct {
	search  Searcher
	llm     LLM
	timeout time.Duration
}

func NewService(search Searcher, llm LLM) *Service {
	return &Service{search: search, llm: llm, timeout: 55 * time.Second}
}

func (s *Service) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
	topic := strings.TrimSpace(req.Topic)
	if len([]rune(topic)) < 3 {
		return nil, response.NewError(http.StatusBadRequest, "Topik wajib diisi (minimal 3 karakter)")
	}
	category, ok := post.NormalizeCategory(req.Category)
	if !ok {
		category = "Umum"
	}
	include := true
	if req.IncludeWebSearch != nil {
		include = *req.IncludeWebSearch
	} else if req.IncludeWebSearchCamel != nil {
		include = *req.IncludeWebSearchCamel
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sources := []Source{}
	if include && s.search != nil {
		if res, err := s.search.Search(ctx, topic, 5); err == nil {
			sources = res
		}
	}
	generator := "openai"
	gp, err := s.callLLM(ctx, topic, category, req.Tone, sources)
	if err != nil {
		generator = "template"
		gp = fallbackArticle(topic, category, sources)
	}
	if strings.TrimSpace(req.Category) == "" {
		if c, ok := post.NormalizeCategory(gp.Category); ok {
			category = c
		}
	}
	gp.Category = category
	gp.Content = appendSources(stripH1(gp.Content), sources)
	if gp.Title == "" {
		gp.Title = topic
	}
	if strings.TrimSpace(gp.Excerpt) == "" {
		gp.Excerpt = post.MakeExcerpt(gp.Content)
	}
	return &GenerateResult{
		Title: gp.Title, Category: gp.Category, Excerpt: gp.Excerpt, Content: gp.Content,
		Post: *gp, Sources: sources, SearchSources: sources,
		WebSearchUsed: include && len(sources) > 0, Generator: generator,
	}, nil
}

func (s *Service) callLLM(ctx context.Context, topic, category, tone string, sources []Source) (*GeneratedPost, error) {
	if s.llm == nil {
		return nil, fmt.Errorf("LLM tidak tersedia")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Topik: %s\nKategori: %s\n", topic, category)
	if tone != "" {
		fmt.Fprintf(&b, "Gaya bahasa: %s\n", tone)
	}
	if len(sources) > 0 {
		b.WriteString("\nHasil penelusuran mesin pencari:\n")
		for i, src := range sources {
			fmt.Fprintf(&b, "[%d] %s (%s)\n%s\n", i+1, src.Title, src.URL, src.Snippet)
		}
	}
	b.WriteString("\nTulis artikel blog 600-900 kata. Kategori wajib salah satu dari: Teknologi, Pendidikan, Gaya Hidup, Bisnis, Opini, Umum.")
	raw, err := s.llm.Chat(ctx, systemPrompt, b.String())
	if err != nil {
		return nil, err
	}
	return ParseGenerated(raw)
}

func ParseGenerated(raw string) (*GeneratedPost, error) {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end <= start {
		return nil, fmt.Errorf("format JSON tidak ditemukan")
	}
	var gp GeneratedPost
	if err := json.Unmarshal([]byte(raw[start:end+1]), &gp); err != nil {
		return nil, err
	}
	gp.Title = strings.TrimSpace(gp.Title)
	gp.Content = strings.TrimSpace(gp.Content)
	if gp.Content == "" {
		return nil, fmt.Errorf("konten AI kosong")
	}
	return &gp, nil
}

func stripH1(content string) string {
	lines := strings.Split(content, "\n")
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "# ") {
			return strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
		}
		break
	}
	return content
}

func appendSources(content string, sources []Source) string {
	if len(sources) == 0 || strings.Contains(content, sources[0].URL) {
		return content
	}
	var b strings.Builder
	b.WriteString(strings.TrimSpace(content))
	b.WriteString("\n\n## Sumber Rujukan\n\n")
	for i, src := range sources {
		fmt.Fprintf(&b, "%d. [%s](%s)\n", i+1, strings.NewReplacer("[", "(", "]", ")").Replace(src.Title), src.URL)
	}
	return b.String()
}

func fallbackArticle(topic, category string, sources []Source) *GeneratedPost {
	var b strings.Builder
	fmt.Fprintf(&b, "## Pendahuluan\n\n%s menjadi salah satu topik yang semakin banyak dibicarakan. Artikel ini merangkum konsep utama, perkembangan terkini, serta langkah praktis yang dapat langsung Anda terapkan.\n\n", topic)
	if len(sources) > 0 {
		b.WriteString("## Temuan Terkini dari Penelusuran Web\n\n")
		for i, src := range sources {
			text := src.Snippet
			if text == "" {
				text = src.Title
			}
			fmt.Fprintf(&b, "- %s [%d]\n", text, i+1)
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "## Mengapa %s Penting?\n\nMemahami %s membantu kita mengambil keputusan yang lebih tepat, meningkatkan produktivitas, dan siap menghadapi perubahan.\n\n", topic, topic)
	b.WriteString("## Langkah Praktis\n\n1. Pelajari konsep dasarnya dari sumber terpercaya.\n2. Terapkan dalam proyek atau aktivitas kecil sehari-hari.\n3. Evaluasi hasilnya secara berkala dan lakukan perbaikan.\n\n")
	fmt.Fprintf(&b, "## Kesimpulan\n\n%s bukan sekadar tren. Dengan pemahaman yang tepat dan praktik yang konsisten, manfaatnya dapat dirasakan dalam jangka panjang.", topic)
	return &GeneratedPost{
		Title:    "Memahami " + topic + ": Konsep, Tren Terkini, dan Langkah Praktis",
		Category: category,
		Excerpt:  fmt.Sprintf("Ulasan ringkas tentang %s: konsep kunci, tren terkini, dan langkah praktis yang dapat langsung diterapkan.", topic),
		Content:  b.String(),
	}
}
