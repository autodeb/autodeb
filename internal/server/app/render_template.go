package app

import (
	"bytes"
	"html/template"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

// RenderTemplate renders a template
func (app *App) RenderTemplate(templateName string, data interface{}) (string, error) {
	b, err := filesystem.ReadFile(app.templatesFS, templateName)
	if err != nil {
		return "", err
	}

	str := string(b)

	tmpl := template.New(filepath.Base(templateName))

	if _, err := tmpl.Parse(str); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	output := string(buf.Bytes())

	return output, nil
}
