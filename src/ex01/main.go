package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sql.DB

type Article struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

// func blogMainPage(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {

// 	}
// }

func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {}

func newArticle(w http.ResponseWriter, r *http.Request) {}

func createArticle(w http.ResponseWriter, r *http.Request) {}

func articleCtx(next http.Handler) http.Handler {
	return nil
}

func getArticle(w http.ResponseWriter, r *http.Request) {}

func main() {
	router := chi.NewRouter()
	/*used when server panics, recover server and response 500 internal server error to user*/
	router.Use(middleware.Recoverer)
	var err error
	db, err = connectToDB()
	catch(err)
	// http.HandleFunc("/", blogMainPage)
	// http.ListenAndServe(":8888", nil)
	router.Use(ChangeMethod)
	router.Get("/", getAllArticles)
	router.Route("/articles", func(r chi.Router) {
		r.Get("/", newArticle)
		r.Post("/", createArticle)
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(articleCtx)
			r.Get("/", getArticle)
		})
	})
}
