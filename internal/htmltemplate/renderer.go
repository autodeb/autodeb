package htmltemplate

import (
	"bytes"
	"html/template"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

//Renderer renders html templates
type Renderer struct {
	fs    filesystem.FS
	cache *templateCache
}

//NewRenderer created a new renderer
func NewRenderer(fs filesystem.FS) *Renderer {
	r := Renderer{
		fs:    fs,
		cache: newTemplateCache(),
	}
	return &r
}

//RenderTemplate renders an html template with the given data
func (renderer *Renderer) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, err := renderer.getOrCreateTemplate(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	output := string(buf.Bytes())

	return output, nil
}

func (renderer *Renderer) getOrCreateTemplate(name string) (*template.Template, error) {
	tmpl, ok := renderer.cache.Load(name)
	if ok {
		return tmpl, nil
	}

	tmpl, err := renderer.createTemplate(name)
	if err != nil {
		return nil, err
	}

	renderer.cache.Store(name, tmpl)

	return tmpl, nil
}

func (renderer *Renderer) createTemplate(name string) (*template.Template, error) {
	b, err := filesystem.ReadFile(renderer.fs, name)
	if err != nil {
		return nil, err
	}

	str := string(b)

	tmpl := template.New(name)

	if _, err := tmpl.Parse(str); err != nil {
		return nil, err
	}

	return tmpl, nil
}
