package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler ...
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler ...
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yamlMap, err := newPathMap(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(yamlMap, fallback), nil
}

type pathMap struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func newPathMap(b []byte) (map[string]string, error) {
	var mapping []pathMap
	result := make(map[string]string)

	err := yaml.Unmarshal(b, &mapping)
	if err != nil {
		return make(map[string]string), err
	}
	for _, m := range mapping {
		result[m.Path] = m.URL
	}
	return result, nil
}
