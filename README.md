
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# **🐾 Puggy**  
Introducing Puggy, a simple webapp router w/ CORS!

## **🦴 What Makes Puggy Special?**
-	Lightweight! 
-   No dependencies!
-   Supports CORS out of the box!
-   URL path variables support `/articles/{article-slug}`

## **🤔 Why Was Puggy Made?**

I really liked [gorilla/mux](https://github.com/gorilla/mux) but it's now archived and though I considered stepping up to maintain it, it wouldn't be fair to maintain only the parts I know and used the most, so I made a new thing that accomplished all the things I loved [gorilla/mux](https://github.com/gorilla/mux) for, and hopefully will help other people as well.

## **🎉 Getting Started**

First, `go get github.com/donuts-are-good/puggy` and import it at the top of your project:

```
import "github.com/donuts-are-good/puggy"
```

Next, create a new router by calling `NewRouter` with an array of CORS domains:

```
corsDomains := []string{
    "puggy.local", 
    "pugtastic.com",
    }

r := puggy.NewRouter(corsDomains)
```

Then, add your routes using the `AddRoute` method:


```
r.AddRoute("GET", "/", handlerRoot) 
r.AddRoute("POST", "/users", handlerUsers)
```

Finally, pass your router instance to the `http.ListenAndServe` method to start your server:

```
pugPort := ":8080"
http.ListenAndServe(pugPort, r)
```

And that's it! You now have a working webapp using Puggy as your router.

## **💡 Usage**

Here are some examples to get you started with Puggy:

### **Simple route and handler example**

```
// handler
func handlerRoot(w http.ResponseWriter, req *http.Request) { 	
    w.Write([]byte("Hello, World!")) 
}

// route
r.AddRoute("GET", "/", handlerRoot)
```

### **Route with path variable example**

```
// handler
func handlerUsers(w http.ResponseWriter, req *http.Request) {
    userID := req.Context().Value("userID").(string) 	
    w.Write([]byte("Hello, User " + userID)) 
}

// route
r.AddRoute("GET", "/users/{userID}", handlerUsers)
```

### **POST request with JSON payload example**


```
// handler
func handlerUserJSON(w http.ResponseWriter, req *http.Request) { 	
    var user struct { 		
        Name string `json:"name"` 		
        Age  int    `json:"age"` 	
    }  	
    if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest) 		
        return
    }  	
    w.Write([]byte("Hello, " + user.Name)) 
}

// route
r.AddRoute("POST", "/users", handlerUserJSON)
```



## **🚀 Full CRUD example**
```
package main

import (
	"encoding/json"
	"net/http"

    "github.com/donuts-are-good/puggy"
)

var books = []struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}{
	{ID: "1", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"},
	{ID: "2", Title: "To Kill a Mockingbird", Author: "Harper Lee"},
	{ID: "3", Title: "Pride and Prejudice", Author: "Jane Austen"},
}

func main() {
	r := puggy.NewRouter([]string{
        "puggy.local", 
        "pugtaculous.local",
        })
	
	// Get all books
	r.AddRoute("GET", "/books", getBooks)
	
	// Get a single book by ID
	r.AddRoute("GET", "/books/{bookID}", getBook)
	
	// Create a new book
	r.AddRoute("POST", "/books", createBook)
	
	// Update a book by ID
	r.AddRoute("PUT", "/books/{bookID}", updateBook)
	
	// Delete a book by ID
	r.AddRoute("DELETE", "/books/{bookID}", deleteBook)
	
	http.ListenAndServe(":8000", r)
}

func getBooks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, req *http.Request) {
	bookID := req.Context().Value("bookID").(string)
	for _, book := range books {
		if book.ID == bookID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, req)
}

func createBook(w http.ResponseWriter, req *http.Request) {
	var book struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(req.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	books = append(books, book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, req *http.Request) {
	bookID := req.Context().Value("bookID").(string)
	var newBook struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(req.Body).Decode(&newBook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, book := range books {
		if book.ID == bookID {
			books[i] = newBook
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newBook)
			return
		}
	}
	http.NotFound(w, req)
}

func deleteBook(w http.ResponseWriter, req *http.Request) {
	bookID := req.Context().Value("bookID").(string)
	for i, book := range books {
		if book.ID == bookID {
			books = append(books[:i], books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, req)
}


```



## Greetz

the Dozens, code-cartel, offtopic-gophers, the garrison, and the monster beverage company.

## License

this code uses the MIT license, not that anybody cares. If you don't know, then don't sweat it.

made with ☕ by 🍩 😋 donuts-are-good


😆👏 Thanks
