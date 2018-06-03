package htmltemplate

import (
	"html/template"
	"sync"
)

//templateCache is a thread-safe template cache
type templateCache struct {
	m     map[string]*template.Template
	mutex sync.RWMutex
}

func newTemplateCache() *templateCache {
	templateCache := templateCache{
		m: make(map[string]*template.Template),
	}
	return &templateCache
}

func (cache *templateCache) Clear() {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	cache.m = make(map[string]*template.Template)
}

func (cache *templateCache) Load(key string) (*template.Template, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	value, ok := cache.m[key]

	return value, ok
}

func (cache *templateCache) Store(key string, value *template.Template) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.m[key] = value
}
