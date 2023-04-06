package puggy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {

	// define test domains
	domains := []string{"example.com", "localhost:3000", "puggy.local"}

	// make a new router with the test domains
	router := NewRouter(domains)

	// test that the router we make has the right type
	if reflect.TypeOf(router) != reflect.TypeOf(&Router{}) {
		t.Errorf("NewRouter() returned wrong type: %T, want *Router", router)
	}

	// test against the domain in the slice
	if !reflect.DeepEqual(router.Domains, domains) {
		t.Errorf("NewRouter() returned wrong domains: got %v, want %v", router.Domains, domains)
	}
}

func TestAddRoute(t *testing.T) {

	// define test domains
	domains := []string{"example.com", "localhost:3000", "puggy.local"}

	// make a new router with the test domains
	router := NewRouter(domains)

	// add a new route to the router
	router.AddRoute("GET", "/users", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	// make a test server
	testPuggy := httptest.NewServer(router)

	// send a request to the new route using the test server
	res, err := http.Get(testPuggy.URL + "/users")
	if err != nil {
		t.Fatalf("http.Get() failed: %v", err)
	}
	defer res.Body.Close()

	// check for http 200 OK
	if res.StatusCode != http.StatusOK {
		t.Errorf("AddRoute() returned wrong status code: got %d, want %d", res.StatusCode, http.StatusOK)
	}

	// add a new route to the router with the new path format
	router.AddRoute("GET", "/books/{id}", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Book " + req.Context().Value(pathVarName("id")).(string)))
	})

	// check for errors in resp body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() failed: %v", err)
	}

	// if the body isn't what we wanted, fail
	if string(body) != "Hello, world!" {
		t.Errorf("AddRoute() returned wrong response body: got %s, want Hello, world!", string(body))
	}

}

func TestServeHTTP(t *testing.T) {

	// define test domains
	domains := []string{"example.com", "localhost:3000", "puggy.local"}

	// make a new router with the test domains
	router := NewRouter(domains)

	// add some routes to the router
	// get
	router.AddRoute("GET", "/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// add a new route with the new path format
	router.AddRoute("GET", "/books/{id}", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Book " + req.Context().Value(pathVarName("id")).(string)))
	})

	// post
	router.AddRoute("POST", "^/users$", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("User created"))
	})

	// put
	router.AddRoute("PUT", "^/users/(?P<id>[0-9]+)$", func(w http.ResponseWriter, req *http.Request) {
		id := req.Context().Value(pathVarName("id")).(string)
		w.Write([]byte("Updating user " + id))
	})

	// test get
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// test a GET request to /books/1
	req, err = http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "Book 1" {
		t.Errorf("ServeHTTP() failed to handle GET request to /books/{id}: got %q, want %q", w.Body.String(), "Book 1")
	}

	// test a GET request to /books
	req, err = http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("ServeHTTP() failed to return 404 Not Found for /books: got %d, want %d", w.Code, http.StatusNotFound)
	}

	// use the response recorder for testing
	w = httptest.NewRecorder()

	// serve the test server
	router.ServeHTTP(w, req)
	if w.Body.String() != "Hello, world!" {
		t.Errorf("ServeHTTP() failed to handle GET request to root URL: got %q, want %q", w.Body.String(), "Hello, world!")
	}

	// test POST to /users
	req, err = http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// use the response recorder for testing
	w = httptest.NewRecorder()

	// serve the test server
	router.ServeHTTP(w, req)
	if w.Body.String() != "User created" {
		t.Errorf("ServeHTTP() failed to handle POST request to /users: got %q, want %q", w.Body.String(), "User created")
	}

	// test a PUT request to /users/:id
	req, err = http.NewRequest("PUT", "/users/123", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// use the response recorder for testing
	w = httptest.NewRecorder()

	// serve the test server
	router.ServeHTTP(w, req)
	if w.Body.String() != "Updating user 123" {
		t.Errorf("ServeHTTP() failed to handle PUT request to /users/:id: got %q, want %q", w.Body.String(), "Updating user 123")
	}

	// test a request to a non-existent route
	req, err = http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// use the response recorder for testing
	w = httptest.NewRecorder()

	// serve the test server
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("ServeHTTP() failed to return 404 Not Found for non-existent route: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestSetCORS(t *testing.T) {

	// use the response recorder for testing
	w := httptest.NewRecorder()

	// define test domains
	domains := []string{"example.com", "localhost:3000", "puggy.local"}

	setCORS(w, domains)

	// test origin
	expectedOrigin := strings.Join(domains, ",")
	actualOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if actualOrigin != expectedOrigin {
		t.Errorf("setCORS() returned wrong Access-Control-Allow-Origin header: got %q, want %q", actualOrigin, expectedOrigin)
	}

	// test methods
	expectedMethods := "GET, POST, PUT, DELETE, OPTIONS"
	actualMethods := w.Header().Get("Access-Control-Allow-Methods")
	if actualMethods != expectedMethods {
		t.Errorf("setCORS() returned wrong Access-Control-Allow-Methods header: got %q, want %q", actualMethods, expectedMethods)
	}

	// test headers
	expectedHeaders := "Origin, Content-Type, Accept"
	actualHeaders := w.Header().Get("Access-Control-Allow-Headers")
	if actualHeaders != expectedHeaders {
		t.Errorf("setCORS() returned wrong Access-Control-Allow-Headers header: got %q, want %q", actualHeaders, expectedHeaders)
	}
}

func TestMatchPath(t *testing.T) {

	// define a test regex filter thing
	testRegex := regexp.MustCompile(`^/users/(?P<id>[0-9]+)/blog/(?P<blog_id>[0-9]+)$`)

	// define test paths
	testPaths := []string{
		"/users/123/blog/456",
		"/users/abc/blog/def",
		"/articles/123",
	}

	// test that we're matching vars correctly in the path
	expectedVars := []map[pathVarName]string{
		{"id": "123", "blog_id": "456"},
		nil,
		nil,
	}

	// range through the test paths and see what we get
	for i, path := range testPaths {
		actualVars := matchPath(testRegex, path)
		if !reflect.DeepEqual(actualVars, expectedVars[i]) {
			t.Errorf("matchPath() failed for %s: got %v, want %v", path, actualVars, expectedVars[i])
		}
	}

}
