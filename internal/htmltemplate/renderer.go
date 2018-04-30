package htmltemplate

import (
	"bytes"
	"html/template"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

//Renderer renders html templates
type Renderer struct {
	fs filesystem.FS
}

//NewRenderer created a new renderer
func NewRenderer(fs filesystem.FS) *Renderer {
	r := Renderer{
		fs: fs,
	}
	return &r
}

//RenderTemplate renders an html template with the given data
func (renderer *Renderer) RenderTemplate(templateName string, data interface{}) (string, error) {
	b, err := filesystem.ReadFile(renderer.fs, templateName)
	if err != nil {
		return "", err
	}

	str := string(b)

	tmpl := template.New(templateName)

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
