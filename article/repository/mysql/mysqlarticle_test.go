package mysql_test

import (
	articleMysqlRepo "art/article/repository/mysql"
	"art/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockArticles := []domain.Article{
		{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: domain.Author{ID: 1},
		},
		{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: domain.Author{ID: 1},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id"}).
		AddRow(mockArticles[0].ID, mockArticles[0].Title, mockArticles[0].Content,
			mockArticles[0].Author.ID).
		AddRow(mockArticles[1].ID, mockArticles[1].Title, mockArticles[1].Content,
			mockArticles[1].Author.ID)

	query := "SELECT id,title,content, author_id FROM article"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := articleMysqlRepo.NewMysqlArticleRepository(db)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), "", num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM article WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	num := int64(5)
	anArticle, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}

func TestStore(t *testing.T) {

	ar := &domain.Article{
		Title:   "Judul",
		Content: "Content",
		Author: domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT  article SET title=\\? , content=\\? , author_id=\\?, updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}
