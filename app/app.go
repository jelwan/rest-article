package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"rest-article/database/model"
	"rest-article/log"
	"rest-article/repo"
	"strconv"
	"time"
)

// Header type constants
const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type App struct {
	ctx      context.Context
	Router   *mux.Router
	Database *sql.DB
	repo     *repo.ArticleRepo
	logger   *logrus.Entry
}

type Article struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type PostArticleRequest struct {
	Article
}

type CreateArticleResponse struct {
	Id int `json:"id"`
}

type TagsResponse struct {
	Tag      string `json:"tag"`
	Count    int    `json:"count"`
	Articles struct {
		Id []string
	} `json:"articles"`
	RelatedTags struct {
		tags []string
	} `json:"related_tags"`
}

type ErrorResponse struct {
	ID      string        `json:"id,omitempty"`
	Success bool          `json:"success,omitempty"`
	Error   ResponseError `json:"error,omitempty"`
}

type ResponseError string

func NewApp(router *mux.Router, database *sql.DB, ctx context.Context) *App {
	return &App{
		Router:   router,
		Database: database,
		ctx:      ctx,
		repo:     repo.NewArticleRepo(ctx, database),
		logger:   log.NewLogger().WithContext(ctx).WithField("module", "app"),
	}
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/articles/{id}").
		HandlerFunc(app.getArticleFunction)

	app.Router.
		Methods("POST").
		Path("/articles").
		HandlerFunc(app.postArticleFunction)

	app.Router.
		Methods("GET").
		Path("/tag/{tagName}/{date}").
		HandlerFunc(app.getTagsFunction)

	app.Router.
		Methods("GET").
		Path("/ping").
		HandlerFunc(app.pingFunction)
}

func (app *App) getArticleFunction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the request path")
	}

	article, tags, err := app.repo.GetArticleByID(id)
	if err != nil {
		err = app.handleError(w, err.Error(), http.StatusInternalServerError)
		if err != nil {
			app.logger.Errorf("error sending error response because: %v", err)
		}
		return
	}

	var tagsList []string
	for _, tag := range tags {
		tagsList = append(tagsList, tag.Name)
	}

	response := Article{
		Id:    fmt.Sprintf("%d", article.Id),
		Title: article.Title,
		Date:  article.Date.Format("01-02-2006"),
		Body:  article.Body,
		Tags:  tagsList,
	}

	w.Header().Add(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (app *App) postArticleFunction(w http.ResponseWriter, r *http.Request) {

	var article Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		err = app.handleError(w, err.Error(), http.StatusInternalServerError)
		if err != nil {
			app.logger.Errorf("error decoding json post body because: %v", err)
		}
		return
	}

	date, _ := time.Parse("2006-01-02", article.Date)
	id, _ := strconv.Atoi(article.Id)
	articleModel := model.Article{
		Id:    id,
		Title: article.Title,
		Date:  date,
		Body:  article.Body,
	}

	articleRes, _, err := app.repo.CreateArticle(articleModel, article.Tags)
	if err != nil {
		err = app.handleError(w, err.Error(), http.StatusInternalServerError)
		if err != nil {
			app.logger.Errorf("error sending error response because: %v", err)
		}
		return
	}

	response := CreateArticleResponse{
		Id: articleRes.Id,
	}

	w.Header().Add(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (app *App) getTagsFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagName, ok := vars["tagName"]
	if !ok {
		log.Fatal("No ID in the request path")
	}

	date, ok := vars["date"]
	if !ok {
		log.Fatal("No ID in the request path")
	}

	dateStr, _ := time.Parse("20060102", date)

	tagCount, err := app.repo.CountTagForDateName(tagName, dateStr.String())
	if err != nil {
		err = app.handleError(w, err.Error(), http.StatusInternalServerError)
		if err != nil {
			app.logger.Errorf("error sending error response because: %v", err)
		}
		return
	}

	relatedTags, err := app.repo.GetRelatedTagForDateAndName(tagName, dateStr.String())
	if err != nil {
		err = app.handleError(w, err.Error(), http.StatusInternalServerError)
		if err != nil {
			app.logger.Errorf("error sending error response because: %v", err)
		}
		return
	}

	taggedArticles, err := app.repo.GetArticleIDForDateAndTag(tagName, dateStr.String())
	if err != nil {
		err = app.handleError(w, err.Error(), http.StatusInternalServerError)
		if err != nil {
			app.logger.Errorf("error sending error response because: %v", err)
		}
		return
	}

	type TagSummaryResponse struct {
		Tag         string   `json:"tag"`
		Count       int      `json:"count"`
		Articles    []string `json:"articles"`
		RelatedTags []string `json:"related_tag"`
	}

	response := TagSummaryResponse{
		Tag:         tagName,
		Count:       tagCount,
		Articles:    taggedArticles,
		RelatedTags: relatedTags,
	}

	w.Header().Add(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (app *App) pingFunction(w http.ResponseWriter, r *http.Request) {
	app.logger.Infof("ping req")
	w.WriteHeader(http.StatusOK)
}

func (app *App) handleError(w http.ResponseWriter, message string, statusCode int) error {

	msg := ResponseError(message)
	body, err := json.Marshal(ErrorResponse{
		Success: false,
		Error:   msg,
	})

	if err != nil {
		app.logger.Errorf("error in marshaling JSON success response because: %v", err)
		return err
	}

	w.Header().Add(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
		app.logger.Errorf("error writing response because: %v", err)
		return err
	}

	return nil
}
