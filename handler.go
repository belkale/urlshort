package urlshort

import (
	yaml "gopkg.in/yaml.v2"
	"html"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := html.EscapeString(r.URL.Path)
		if val, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, val, http.StatusMovedPermanently)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

type pathMap struct {
	Path string
	URL  string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathMaps []pathMap
	err := yaml.Unmarshal(yml, &pathMaps)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := html.EscapeString(r.URL.Path)
		for _, p := range pathMaps {
			if p.Path == path {
				http.Redirect(w, r, p.URL, http.StatusMovedPermanently)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}
