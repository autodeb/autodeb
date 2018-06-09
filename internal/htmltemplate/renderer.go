package htmltemplate

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

//Renderer renders html templates
type Renderer struct {
	fs           filesystem.FS
	cache        *templateCache
	cacheEnabled bool
	funcMap      template.FuncMap
}

//NewRenderer created a new renderer
func NewRenderer(fs filesystem.FS, funcMap template.FuncMap, cacheEnabled bool) *Renderer {
	r := Renderer{
		fs:           fs,
		cache:        newTemplateCache(),
		cacheEnabled: cacheEnabled,
		funcMap:      funcMap,
	}
	return &r
}

//RenderTemplate renders an html template with the given data
func (renderer *Renderer) RenderTemplate(data interface{}, filenames ...string) (string, error) {
	tmpl, err := renderer.getOrCreateTemplate(filenames...)
	if err != nil {
		return "", errors.WithMessage(err, "cannot get template")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	output := string(buf.Bytes())

	return output, nil
}

func (renderer *Renderer) getOrCreateTemplate(filenames ...string) (*template.Template, error) {
	templateName := strings.Join(filenames, "+")

	if renderer.cacheEnabled {
		tmpl, ok := renderer.cache.Load(templateName)
		if ok {
			return tmpl, nil
		}
	}

	tmpl, err := renderer.createTemplate(filenames...)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create template")
	}

	if renderer.cacheEnabled {
		renderer.cache.Store(templateName, tmpl)
	}

	return tmpl, nil
}

func (renderer *Renderer) createTemplate(filenames ...string) (*template.Template, error) {
	tmpl := template.New("")

	tmpl.Funcs(renderer.funcMap)

	for _, filename := range filenames {

		f, err := renderer.fs.Open(filename)
		if err != nil {
			return nil, errors.WithMessage(err, "could not open template")
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.WithMessage(err, "could not read template")
		}

		str := string(b)

		if _, err := tmpl.Parse(str); err != nil {
			return nil, errors.WithMessage(err, "could not parse template")
		}
	}

	return tmpl, nil
}
