package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
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
		path := r.URL.Path
		longURL, exist := pathsToUrls[path]
		if !exist {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type YAMlData struct {
	Path string
	URL  string
}

func YAMLHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedYAMAL []YAMlData

	if err := yaml.Unmarshal([]byte(data), &parsedYAMAL); err != nil {
		return nil, err
	}

	pathsToUrls := map[string]string{}

	for _, YAMLEntrey := range parsedYAMAL {
		pathsToUrls[YAMLEntrey.Path] = YAMLEntrey.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJSON []YAMlData

	if err := json.Unmarshal([]byte(data), &parsedJSON); err != nil {
		return nil, err
	}

	pathsToUrls := map[string]string{}

	for _, JSONEntrey := range parsedJSON {
		pathsToUrls[JSONEntrey.Path] = JSONEntrey.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}
