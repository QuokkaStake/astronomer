package templates

import (
	"bytes"
	"html/template"
	"main/pkg/utils"
	"main/templates"

	"github.com/rs/zerolog"
)

type TelegramTemplatesManager struct {
	Templates map[string]*template.Template
	Logger    zerolog.Logger
}

func NewTelegramTemplatesManager(
	logger *zerolog.Logger,
) *TelegramTemplatesManager {
	return &TelegramTemplatesManager{
		Templates: map[string]*template.Template{},
		Logger:    logger.With().Str("component", "telegram_templates_manager").Logger(),
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
		"FormatDuration": utils.FormatDuration,
		"FormatPercent":  utils.FormatPercent,
		"FormatFloat":    utils.FormatFloat,
		"FormatSince":    utils.FormatSince,
	}).ParseFS(templates.TemplatesFs, "telegram/"+filename)
	if err != nil {
		return nil, err
	}

	m.Templates[templateName] = t

	return t, nil
}
