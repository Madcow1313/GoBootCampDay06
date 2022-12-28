package main

import (
	"context"
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

// func ChangeMethod(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodPost {
// 			switch method := r.PostFormValue("_method"); method {
// 			case http.MethodPut:
// 				fallthrough
// 			case http.MethodPatch:
// 				fallthrough
// 			case http.MethodDelete:
// 				r.Method = method
// 			default:
// 			}
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := dbGetAllArticles(db)
	catch(err)

	t, _ := template.ParseFiles("base.html", "index.html")
	err = t.Execute(w, articles)
	catch(err)
}

func newArticle(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("base.html", "new.html")
	err := t.Execute(w, nil)
	catch(err)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	article := &Article{
		Title:   title,
		Content: template.HTML(content),
	}
	err := dbCreateArticle(article, db)
	catch(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func articleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		article, err := dbGetArticle(articleID, db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)
	t, _ := template.ParseFiles("base.html", "article.html")
	err := t.Execute(w, article)
	catch(err)
}

func main() {
	router := chi.NewRouter()
	/*used when server panics, recover server and response 500 internal server error to user*/
	router.Use(middleware.Recoverer)
	var err error
	db, err = connectToDB()
	catch(err)
	// http.HandleFunc("/", blogMainPage)
	// http.ListenAndServe(":8888", nil)
	//router.Use(ChangeMethod)
	router.Get("/", getAllArticles)
	router.Route("/articles", func(r chi.Router) {
		r.Get("/", newArticle)
		r.Post("/", createArticle)
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(articleCtx)
			r.Get("/", getArticle)
		})
	})

	err = http.ListenAndServe(":8080", router)
	catch(err)

}
