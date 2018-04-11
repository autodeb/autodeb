package app

import (
	"bytes"
	"html/template"
	"path/filepath"
)

// RenderTemplate renders a template
func (app *App) RenderTemplate(templateName string, data interface{}) (string, error) {
	templatePath := filepath.Join(app.config.TemplatesDirectory, templateName)

	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", err
	}

	output := string(buf.Bytes())

	return output, nil
}
