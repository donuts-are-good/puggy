package puggy

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Route defines a single route
type Route struct {
	Method  string           `json:"method,omitempty"`
	Path    *regexp.Regexp   `json:"path,omitempty"`
	Handler http.HandlerFunc `json:"handler,omitempty"`
}

// Router holds routes and cors domains
type Router struct {
	Routes  []Route  `json:"routes,omitempty"`
	Domains []string `json:"domains,omitempty"`
}

// pathVarName satisfies a linter warning about collisiones in string[string] maps
type pathVarName string

// NewRouter makes a new router
func NewRouter(domains []string) *Router {
	return &Router{Domains: domains}
}

func (router *Router) AddRoute(method string, path string, handler http.HandlerFunc) {
	var regexPattern string
	if path == "/" {
		regexPattern = "^/$"
	} else {
		regexPattern = createRegexPattern(path)
	}
	route := Route{Method: method, Path: regexp.MustCompile(regexPattern), Handler: http.HandlerFunc(handler)}

	router.Routes = append(router.Routes, route)
}

// func createRegexPattern(path string) string {
// 	re := regexp.MustCompile(`{(\w+)}`)
// 	escapedPath := regexp.QuoteMeta(path)
// 	regexPattern := re.ReplaceAllStringFunc(escapedPath, func(m string) string {
// 		variableName := re.FindStringSubmatch(m)[1]
// 		return fmt.Sprintf(`(?P<%s>[^/]+)`, variableName)
// 	})
// 	return fmt.Sprintf("^%s$", regexPattern)
// }

func createRegexPattern(path string) string {
	re := regexp.MustCompile(`/{([^/]+)}`)
	escapedPath := regexp.QuoteMeta(path)
	regexPattern := re.ReplaceAllStringFunc(escapedPath, func(m string) string {
		variableName := re.FindStringSubmatch(m)[1]
		return fmt.Sprintf(`(?P<%s>[^/]+)`, variableName)
	})
	return fmt.Sprintf("^%s$", regexPattern)
}

// ServeHTTP matches the methods and paths with the right handler
func (router *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	// iterate over all the routes in the router
	for _, route := range router.Routes {

		// check if the request method matches the method of the current route
		if req.Method != route.Method {
			continue
		}

		// check if the request URL path matches the path of the current route
		vars := matchPath(route.Path, req.URL.Path)
		if vars == nil {
			continue
		}

		// set the CORS header for the response
		setCORS(writer, router.Domains)

		// if the request method is "OPTIONS", write an HTTP status code of 200 OK and return
		if req.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusOK)
			return
		}

		// add the URL path variables to the request context
		for pathKey, pathValue := range vars {
			req = req.WithContext(context.WithValue(req.Context(), pathKey, pathValue))
		}

		// call the handler function of the current route
		route.Handler(writer, req)
		return
	}

	// if no matching route is found, hit em with a 404
	http.NotFound(writer, req)

}

// setCORS defines allowed domains, methods and headers
func setCORS(writer http.ResponseWriter, domains []string) {

	//  allowed domains
	writer.Header().Set("Access-Control-Allow-Origin", strings.Join(domains, ","))

	//  allowed methods
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// allowed headers
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

}

// matchPath takes a regex and a string and returns a map of variables
// in the url
func matchPath(path *regexp.Regexp, urlPath string) map[pathVarName]string {

	// find all matches in the url
	matches := path.FindStringSubmatch(urlPath)

	// if there are no matches return
	if len(matches) == 0 {
		return nil
	}

	// put the vars in a map
	vars := map[pathVarName]string{}

	// loop through the matches
	for i, name := range path.SubexpNames() {

		// skip the first match and empty matches
		if i == 0 || name == "" {
			continue
		}

		// store the matched value in the vars map
		vars[pathVarName(name)] = matches[i]
	}

	// return the vars map
	return vars
}
