package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rest-article/log"
	"rest-article/repo"
	"testing"
)

func NewMockArticleRepo(err error) repo.Repo {
	mockRepo := &repo.ArticleRepoMock{
		Err: err,
	}
	return mockRepo
}

func TestGetArticleFunction(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestGetArticleFunction"),
	}

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestGetArticleFunctionBadId(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestGetArticleFunction"),
	}

	req := httptest.NewRequest(http.MethodGet, "/articles/s", nil)
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, respBody.Error, ResponseError("provided id is not a number"))
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

func TestGetArticleFunctionNoId(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestGetArticleFunction"),
	}

	req := httptest.NewRequest(http.MethodGet, "/articles/", nil)
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, respBody.Success, false)
	assert.Equal(t, resp.Code, http.StatusNotFound)
}

func TestPostArticleFunction(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestPostArticleFunction"),
	}

	reqBody := PostArticleRequest{Article{
		Id:    "1",
		Title: "test article",
		Date:  "2020-02-01",
		Body:  "test art",
		Tags:  []string{"science", "math"},
	}}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewReader(body))
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.Code, http.StatusOK)
}

func TestPostArticleFunctionBadId(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestPostArticleFunction"),
	}

	reqBody := PostArticleRequest{Article{
		Id:    "",
		Title: "test article",
		Date:  "2020-02-01",
		Body:  "test art",
		Tags:  []string{"science", "math"},
	}}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewReader(body))
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, respBody.Error, ResponseError("no id provided"))
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

func TestGetTagsFunction(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestPostArticleFunction"),
	}

	tag := "science"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tag/%s/20200201", tag), nil)
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody TagSummaryResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, respBody.Tag, tag)
	assert.Equal(t, respBody.Tag, "science")
	assert.Equal(t, resp.Code, http.StatusOK)
}

func TestGetTagsFunctionBadDate(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestPostArticleFunction"),
	}

	tag := "science"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tag/%s/20200lkk201", tag), nil)
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, respBody.Success, false)
	assert.Equal(t, respBody.Error, ResponseError("bad date format provided"))
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}

func TestGetTagsFunctionBadTag(t *testing.T) {
	app := &App{
		Database: nil,
		ctx:      nil,
		repo:     NewMockArticleRepo(nil),
		Router:   mux.NewRouter(),
		logger:   log.NewLogger().WithField("test", "TestPostArticleFunction"),
	}

	tag := ""
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tag/%s/20200201", tag), nil)
	resp := httptest.NewRecorder()

	app.SetupRouter()
	app.Router.ServeHTTP(resp, req)

	var respBody ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&respBody)
	t.Log(respBody)

	assert.Equal(t, false, respBody.Success)
}
