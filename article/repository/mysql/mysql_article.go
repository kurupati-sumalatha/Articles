package mysql

import (
	"art/article/repository"
	"art/domain"
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewMysqlArticleRepository(conn *sql.DB) domain.ArticleRepository {
	return &mysqlArticleRepository{conn}
}

func (m *mysqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Article, 0)
	for rows.Next() {
		t := domain.Article{}
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = domain.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `SELECT id,title,content, author_id
  						FROM article`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (m *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	query := `SELECT id,title,content, author_id
  						FROM article WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlArticleRepository) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	query := `SELECT id,title,content, author_id
  						FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	query := `INSERT  article SET title=? , content=? , author_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}
