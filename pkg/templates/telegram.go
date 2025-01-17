package templates

import (
	"bytes"
	"fmt"
	"html/template"
	timePkg "main/pkg/time"
	"main/pkg/types"
	"main/pkg/utils"
	"main/templates"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type TelegramTemplatesManager struct {
	Templates map[string]*template.Template
	Logger    zerolog.Logger
	Time      timePkg.Time
}

func NewTelegramTemplatesManager(
	logger *zerolog.Logger,
	time timePkg.Time,
) *TelegramTemplatesManager {
	return &TelegramTemplatesManager{
		Templates: map[string]*template.Template{},
		Logger:    logger.With().Str("component", "telegram_templates_manager").Logger(),
		Time:      time,
	}
}

func (m *TelegramTemplatesManager) Render(templateName string, data interface{}) (string, error) {
	templateToRender, err := m.GetTemplate(templateName)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("name", templateName).
			Msg("Error getting template")
		return "", err
	}

	var buffer bytes.Buffer
	if err := templateToRender.Execute(&buffer, data); err != nil {
		m.Logger.Error().
			Err(err).
			Str("name", templateName).
			Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), nil
}

func (m *TelegramTemplatesManager) GetTemplate(templateName string) (*template.Template, error) {
	if cachedTemplate, ok := m.Templates[templateName]; ok {
		m.Logger.Trace().Str("type", templateName).Msg("Using cached template")
		return cachedTemplate, nil
	}

	m.Logger.Trace().Str("type", templateName).Msg("Loading template")

	filename := templateName + ".html"

	t, err := template.New(filename).Funcs(template.FuncMap{
		"FormatDuration":   utils.FormatDuration,
		"FormatPercent":    utils.FormatPercent,
		"FormatPercentDec": utils.FormatPercentDec,
		"FormatFloat":      utils.FormatFloat,
		"FormatSince":      m.FormatSince,
		"FormatLink":       m.FormatLink,
		"FormatLinks":      m.FormatLinks,
		"SerializeAmount":  m.SerializeAmount,
	}).ParseFS(templates.TemplatesFs, "telegram/"+filename)
	if err != nil {
		return nil, err
	}

	m.Templates[templateName] = t

	return t, nil
}

func (m *TelegramTemplatesManager) FormatSince(sinceTime time.Time) string {
	return utils.FormatSince(m.Time.Since(sinceTime))
}

func (m *TelegramTemplatesManager) FormatLink(link types.Link) template.HTML {
	return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link.Href, link.Text))
}
func (m *TelegramTemplatesManager) FormatLinks(links []types.Link) template.HTML {
	linksConverted := utils.Map(links, func(link types.Link) string {
		return fmt.Sprintf("<a href='%s'>%s</a>", link.Href, link.Text)
	})

	return template.HTML(strings.Join(linksConverted, " "))
}

func (m *TelegramTemplatesManager) SerializeAmount(amount types.Amount) string {
	if amount.PriceUSD != nil {
		return fmt.Sprintf(
			"%s %s ($%s)",
			utils.FormatDec(amount.Amount),
			amount.Denom,
			utils.FormatDec(*amount.PriceUSD),
		)
	}

	return fmt.Sprintf(
		"%s %s",
		utils.FormatDec(amount.Amount),
		amount.Denom,
	)
}
