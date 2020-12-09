package database

import (
	"bookshelf/config"
	"bookshelf/model"
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Repository struct {
	Engine *xorm.Engine
}

func NewRepository(c *config.DatabaseConfig) (*Repository, error) {

	dbConn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True",
		c.User, c.Password, c.Host, c.Port, c.Database)

	engine, err := xorm.NewEngine("mysql", dbConn)
	if err != nil {
		return nil, err
	}

	err = engine.Sync2(new(model.Book))
	if err != nil {
		return nil, err
	}

	engine.SetMaxOpenConns(10)
	//engine.SetConnMaxLifetime(time.Second * 60)

	return &Repository{
		Engine: engine,
	}, nil
}

func (r *Repository) Transaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		session := r.Engine.NewSession()
		defer session.Close()

		req := ctx.Request()
		ctx.SetRequest(req.WithContext(context.WithValue(req.Context(), "DB", session)))

		if err := session.Begin(); err != nil {
			log.Println(err)
		}
		if err := next(ctx); err != nil {
			session.Rollback()
			return err
		}
		if ctx.Response().Status >= 500 {
			session.Rollback()
			return nil
		}
		if err := session.Commit(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func (r *Repository) GetDBConn(ctx context.Context) *xorm.Session {
	v := ctx.Value("DB")
	if v == nil {
		return nil
	}
	if session, ok := v.(*xorm.Session); ok {
		return session
	}
	return nil
}
