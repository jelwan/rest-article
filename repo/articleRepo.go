package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"rest-article/database/model"
	"rest-article/field"
	"rest-article/log"
)

type ArticleRepo struct {
	ctx    context.Context
	db     *sql.DB
	logger *logrus.Entry
}

func NewArticleRepo(ctx context.Context, db *sql.DB) *ArticleRepo {
	return &ArticleRepo{
		ctx:    ctx,
		db:     db,
		logger: log.NewLogger().WithContext(ctx).WithField("module", "repo"),
	}
}

func (articleRepo *ArticleRepo) GetArticleByID(id string) (*model.Article, []*model.Tag, error) {

	statement, err := articleRepo.db.PrepareContext(articleRepo.ctx,
		"Select `id`, `title`, `date`, `body` FROM `svc-article`.articles where id = ?")
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetArticleByID", "PrepareContext")).
			Errorf("statement creation failed because: %v", err)
		return nil, nil, err
	}

	var article model.Article
	err = statement.QueryRow(id).Scan(&article.Id, &article.Title, &article.Date, &article.Body)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetArticleByID", "QueryRow")).
			Errorf("failed to query article id %s because %v", id, err)
		return nil, nil, err
	}

	tagIDs, err := articleRepo.getArticleTagsByArticleID(articleRepo.ctx, id)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetArticleByID", "getArticleTagsByArticleID")).
			Errorf("failed to query article id %s because %v", id, err)
		return nil, nil, err
	}

	var tags []*model.Tag
	for _, id := range tagIDs {
		tag, err := articleRepo.getTagById(articleRepo.ctx, id)
		if err != nil {
			articleRepo.logger.
				WithFields(field.ErrorFields("GetArticleByID", "getTagById")).
				Errorf("error getting tag id %d because %v", id, err)
			return nil, nil, err
		}
		tags = append(tags, tag)
	}

	return &article, tags, nil
}

func (articleRepo *ArticleRepo) CountTagForDateName(name, date string) (int, error) {

	countStmt, err := articleRepo.db.PrepareContext(articleRepo.ctx,
		"SELECT count(tags.id) as tag_count "+
			"FROM tags "+
			"INNER JOIN article_tags on tags.id = article_tags.tag_id "+
			"INNER JOIN articles on article_tags.article_id = articles.id "+
			"WHERE tags.tag_title = ? AND articles.date = ?")

	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("countTagForDateName", "PrepareContext")).
			Errorf("statement creation failed because: %v", err)
		return -1, err
	}

	var count int
	err = countStmt.QueryRow(name, date).Scan(&count)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetArticleByID", "QueryRow")).
			Errorf("failed to get count of tag name %s on date %s because %v", name, date, err)
		return -1, err
	}

	return count, nil
}

func (articleRepo *ArticleRepo) GetRelatedTagForDateAndName(name, date string) ([]string, error) {

	relatedTagsStmt, err := articleRepo.db.PrepareContext(articleRepo.ctx,
		"SELECT DISTINCT tags.tag_title "+
			"FROM tags "+
			"INNER JOIN article_tags on tags.id = article_tags.tag_id "+
			"INNER JOIN articles on article_tags.article_id = articles.id "+
			"WHERE tag_title != ? AND articles.date = ?")

	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetRelatedTagForDate", "PrepareContext")).
			Errorf("statement creation failed because: %v", err)
		return nil, err
	}

	rows, err := relatedTagsStmt.Query(name, date)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetRelatedTagForDate", "Query")).
			Errorf("statement query failed because: %v", err)
		return nil, err
	}
	defer rows.Close()

	var relatedTags []string
	for rows.Next() {
		var relatedTag string
		err := rows.Scan(&relatedTag)
		if err != nil {
			articleRepo.logger.
				WithFields(field.ErrorFields("GetRelatedTagForDate", "Scan")).
				Errorf("failed to get list of related tags on date %s because %v", date, err)
			return nil, err
		}
		relatedTags = append(relatedTags, relatedTag)
	}

	return relatedTags, nil
}

func (articleRepo *ArticleRepo) GetArticleIDForDateAndTag(name, date string) ([]string, error) {

	taggedArticleStmt, err := articleRepo.db.PrepareContext(articleRepo.ctx,
		"SELECT articles.id "+
			"FROM tags "+
			"INNER JOIN article_tags on tags.id = article_tags.tag_id "+
			"INNER JOIN articles on article_tags.article_id = articles.id "+
			"WHERE tag_title = ? "+
			"AND articles.date = ? "+
			"LIMIT 10")

	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetArticleIDForDateAndTag", "PrepareContext")).
			Errorf("statement creation failed because: %v", err)
		return nil, err
	}

	rows, err := taggedArticleStmt.Query(name, date)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("GetArticleIDForDateAndTag", "Query")).
			Errorf("statement query failed because: %v", err)
		return nil, err
	}
	defer rows.Close()

	var taggedArticles []string
	for rows.Next() {
		var articleID string
		err := rows.Scan(&articleID)
		if err != nil {
			articleRepo.logger.
				WithFields(field.ErrorFields("GetArticleIDForDateAndTag", "Scan")).
				Errorf("failed to get list of tagged articles on date %s for tag %s because %v", date, name, err)
			return nil, err
		}
		taggedArticles = append(taggedArticles, articleID)
	}

	return taggedArticles, nil
}

func (articleRepo *ArticleRepo) GetTagSummary(tag string, date string) (*model.Article, error) {
	return nil, nil
}

func (articleRepo *ArticleRepo) CreateArticle(article model.Article, tags []string) (*model.Article, []*model.Tag, error) {

	tagItems, err := articleRepo.getTagsByName(articleRepo.ctx, tags)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("CreateArticle", "getTagsByName")).
			Infof("error getting tags because %v", err)
		return nil, nil, err
	}

	var tagMap = make(map[string]int)
	for _, tag := range tags {
		tagMap[tag] = -1
	}

	for _, tag := range tagItems {
		tagMap[tag.Name] = tag.Id
	}

	for tagName, id := range tagMap {
		if id == -1 {

			newTagId, err := articleRepo.insertTag(articleRepo.ctx, tagName)
			if err != nil {
				articleRepo.logger.
					WithFields(field.ErrorFields("CreateArticle", "insertTag")).
					Errorf("failed to insert new tag %s, because %v", tagName, err)
				return nil, nil, err
			}

			tagMap[tagName] = newTagId
			tag := &model.Tag{Id: newTagId, Name: tagName}
			tagItems = append(tagItems, tag)
		}
	}

	tagIdList := tagIDList(tagItems)

	err = articleRepo.insertArticle(articleRepo.ctx, article)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("CreateArticle", "insertArticle")).
			Errorf("failed to insert new article %s, because %v", article.Id, err)
		return nil, nil, err
	}

	err = articleRepo.insertArticleTags(articleRepo.ctx, article.Id, tagIdList)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("CreateArticle", "insertArticleTags")).
			Errorf("failed to insert article tag because %v", err)
		return nil, nil, err
	}

	return &article, tagItems, nil
}

func (articleRepo *ArticleRepo) getTagsByName(ctx context.Context, tagNames []string) ([]*model.Tag, error) {

	statement, err := articleRepo.db.PrepareContext(ctx,
		"SELECT `id`, `tag_title` FROM `svc-article`.tags WHERE `tag_title` = ?")
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("getTagsByName", "PrepareContext")).
			Errorf("statement creation failed because: %v", err)
		return nil, err
	}

	var tagItems []*model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		err = statement.QueryRow(name).Scan(&tag.Id, &tag.Name)
		if err == sql.ErrNoRows {
			articleRepo.logger.Infof("no tags found with names %s", name)
		} else if err != nil {
			articleRepo.logger.
				WithFields(field.ErrorFields("getTagsByName", "Scan")).
				Errorf("error selecting tags because: %v", err)
			return nil, err
		}
		tagItems = append(tagItems, &tag)
	}

	return tagItems, nil
}

func (articleRepo *ArticleRepo) getTagById(ctx context.Context, id int) (*model.Tag, error) {

	statement, err := articleRepo.db.PrepareContext(ctx,
		"SELECT `id`, `tag_title` FROM `svc-article`.tags WHERE id = ?")
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("getTagById", "PrepareContext")).
			Errorf("statement creation failed because: %v", err)
		return nil, err
	}

	var tag model.Tag

	err = statement.QueryRow(id).Scan(&tag.Id, &tag.Name)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("getTagById", "Scan")).
			Errorf("failed to query article id %s because %v", id, err)
		return nil, err
	}

	return &tag, nil
}

func (articleRepo *ArticleRepo) getArticleTagsByArticleID(ctx context.Context, articleID string) ([]int, error) {

	rows, err := articleRepo.db.QueryContext(ctx, "SELECT `tag_id` FROM `svc-article`.article_tags WHERE `article_id` = ?", articleID)
	if err != nil {
		articleRepo.logger.
			WithFields(field.ErrorFields("getArticleTagsByArticleID", "QueryContext")).
			Errorf("error selecting article tags because: %v", err)
		return nil, err
	}

	var tagIDs []int
	for rows.Next() {
		var tag int
		err := rows.Scan(&tag)
		if err != nil {
			if err == sql.ErrNoRows {
				articleRepo.logger.Infof("no article tags found for article id %d", articleID)
				return nil, fmt.Errorf("no article tags found found")
			} else {
				articleRepo.logger.
					WithFields(field.ErrorFields("getArticleTagsByArticleID", "Scan")).
					Errorf("error selecting article tags because: %v", err)
				return nil, err
			}
		}
		tagIDs = append(tagIDs, tag)
	}

	return tagIDs, nil
}

func (articleRepo *ArticleRepo) insertTag(ctx context.Context, tagName string) (int, error) {
	tx, err := articleRepo.db.Begin()
	if err != nil {
		articleRepo.logger.Errorf("failed to start transaction because: %v", err)
		return -1, err
	}
	defer tx.Rollback()

	insertTagStmt, err := tx.PrepareContext(ctx, "INSERT INTO  `svc-article`.tags(tag_title) VALUES(?)")
	if err != nil {
		articleRepo.logger.Errorf("tag statement creation failed because: %v", err)
		return -1, err
	}
	defer insertTagStmt.Close()

	result, err := insertTagStmt.Exec(tagName)
	if err != nil {
		articleRepo.logger.Errorf("error executing insert tag statement: %v", err)
		return -1, err
	}

	if count, err := result.RowsAffected(); count > 1 || err != nil {
		articleRepo.logger.Errorf("error too many rows affected by insert tag statement: %v", err)
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		articleRepo.logger.Errorf("error failed to commit transaction: %v", err)
		return -1, err
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		articleRepo.logger.Errorf("error retrieving tag ID: %v", err)
		return -1, err
	}

	return int(tagID), nil
}

func (articleRepo *ArticleRepo) insertArticle(ctx context.Context, article model.Article) error {

	tx, err := articleRepo.db.Begin()
	if err != nil {
		articleRepo.logger.Errorf("failed to start transaction because: %v", err)
		return err
	}
	defer tx.Rollback()

	insertArticleStmt, err := tx.PrepareContext(ctx,
		"INSERT INTO `svc-article`.articles(`id`, `title`, `date`, `body`) VALUES(?, ?, ?, ?)")
	if err != nil {
		articleRepo.logger.Errorf("article statement creation failed because: %v", err)
		return err
	}
	defer insertArticleStmt.Close()

	result, err := insertArticleStmt.Exec(article.Id, article.Title, article.Date, article.Body)
	if err != nil {
		articleRepo.logger.Errorf("error executing insert article statement: %v", err)
		return err
	}

	if count, err := result.RowsAffected(); count > 1 || err != nil {
		articleRepo.logger.Errorf("error too many rows affected by insert tag statement: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		articleRepo.logger.Errorf("error failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (articleRepo *ArticleRepo) insertArticleTags(ctx context.Context, articleID int, tagIDs []int) error {

	tx, err := articleRepo.db.Begin()
	if err != nil {
		articleRepo.logger.Errorf("failed to start transaction because: %v", err)
		return err
	}
	defer tx.Rollback()

	insertArticleTagStmt, err := tx.PrepareContext(ctx,
		"INSERT INTO `svc-article`.article_tags(article_id, tag_id) VALUES (?, ?)")
	if err != nil {
		articleRepo.logger.Errorf("tag statement creation failed because: %v", err)
		return err
	}
	defer insertArticleTagStmt.Close()

	for _, tagID := range tagIDs {
		_, err := insertArticleTagStmt.Exec(articleID, tagID)
		if err != nil {
			articleRepo.logger.Errorf("error executing insert article tags statement: %v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		articleRepo.logger.Errorf("error failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func tagIDList(tags []*model.Tag) []int {
	var list []int
	for _, tag := range tags {
		list = append(list, tag.Id)
	}
	return list
}
