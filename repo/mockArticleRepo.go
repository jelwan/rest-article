package repo

import (
	"context"
	"fmt"
	"rest-article/database/model"
	"time"
)

type ArticleRepoMock struct {
	Err error
}

func (mr *ArticleRepoMock) NewMockArticleRepo(err error) Repo {
	mockRepo := &ArticleRepoMock{
		Err: err,
	}
	return mockRepo
}

func (mr *ArticleRepoMock) GetArticleByID(id string) (*model.Article, []*model.Tag, error) {

	if mr.Err != nil {
		return nil, nil, mr.Err
	}

	article := &model.Article{
		Id:    1,
		Title: "test article",
		Date:  time.Now(),
		Body:  "test article",
	}

	tags := []*model.Tag{
		{Id: 1, Name: "test"},
		{Id: 2, Name: "test2"},
	}

	return article, tags, nil
}

func (mr *ArticleRepoMock) CountTagForDateName(name, date string) (int, error) {

	if mr.Err != nil {
		return -1, mr.Err
	}

	return 999, nil
}

func (mr *ArticleRepoMock) GetRelatedTagForDateAndName(name, date string) ([]string, error) {

	if mr.Err != nil {
		return nil, mr.Err
	}

	return []string{"test", "test2"}, nil
}

func (mr *ArticleRepoMock) GetArticleIDForDateAndTag(name, date string) ([]string, error) {
	if mr.Err != nil {
		return nil, mr.Err
	}
	return []string{"1", "2", "3", "4"}, nil
}

func (mr *ArticleRepoMock) CreateArticle(article model.Article, tags []string) (*model.Article, []*model.Tag, error) {

	if mr.Err != nil {
		return nil, nil, mr.Err
	}

	var tagItems []*model.Tag
	for i, tag := range tags {
		tagItems = append(tagItems, &model.Tag{Id: i, Name: tag})
	}

	return &article, tagItems, nil
}

func (mr *ArticleRepoMock) getTagsByName(ctx context.Context, tagNames []string) ([]*model.Tag, error) {
	if mr.Err != nil {
		return nil, mr.Err
	}

	var tagItems []*model.Tag
	for i, tag := range tagNames {
		tagItems = append(tagItems, &model.Tag{Id: i, Name: tag})
	}

	return tagItems, nil
}

func (mr *ArticleRepoMock) getTagById(ctx context.Context, id int) (*model.Tag, error) {

	if mr.Err != nil {
		return nil, mr.Err
	}

	return &model.Tag{
			Id:   id,
			Name: fmt.Sprintf("test-%d", id)},
		nil
}

func (mr *ArticleRepoMock) getArticleTagsByArticleID(ctx context.Context, articleID string) ([]int, error) {

	if mr.Err != nil {
		return nil, mr.Err
	}

	return []int{1, 2, 3, 4, 5}, nil

}

func (mr *ArticleRepoMock) insertTag(ctx context.Context, tagName string) (int, error) {
	if mr.Err != nil {
		return -1, mr.Err
	}

	return 1, nil
}

func (mr *ArticleRepoMock) insertArticle(ctx context.Context, article model.Article) error {
	if mr.Err != nil {
		return mr.Err
	}

	return nil
}

func (mr *ArticleRepoMock) insertArticleTags(ctx context.Context, articleID int, tagIDs []int) error {
	if mr.Err != nil {
		return mr.Err
	}

	return nil
}
