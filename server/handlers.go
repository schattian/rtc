package server

import (
	"context"
	"encoding/json"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/sebach1/git-crud/config"
	"github.com/sebach1/git-crud/git"
	"github.com/sebach1/git-crud/internal/name"
	"github.com/valyala/fasthttp"
)

// Router is a handler which redirects to handlers
func Router(reqCtx *fasthttp.RequestCtx) {
	db := databaseHandler(reqCtx)
	switch string(reqCtx.Path()) {
	case "/add":
		addHandler(context.TODO(), db, reqCtx)
	case "/rm":
		rmHandler(context.TODO(), db, reqCtx)
	default:
		reqCtx.NotFound()
	}
}

func databaseHandler(reqCtx *fasthttp.RequestCtx) *sqlx.DB {
	dataSource := os.Getenv("db")
	if dataSource == "" {
		dataSource = config.DefaultDBSrc
	}
	db, err := sqlx.Open("postgres", dataSource)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadGateway)
	}
	db.MapperFunc(name.ToSnakeCase)
	return db
}

func decoderHandler(reqCtx *fasthttp.RequestCtx, validator func(body *reqBody) error) *reqBody {
	body := &reqBody{}
	rawBody := reqCtx.PostBody()
	err := json.Unmarshal(rawBody, body)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusPreconditionFailed)
	}
	err = validator(body)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusUnprocessableEntity)
	}

	return body
}

func addHandler(ctx context.Context, db *sqlx.DB, reqCtx *fasthttp.RequestCtx) {
	body := decoderHandler(reqCtx, validateAdd)
	err := git.Add(ctx, db,
		body.Entity, body.Table, body.Column, body.Branch, body.Value, body.Type, body.Opts)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}

func rmHandler(ctx context.Context, db *sqlx.DB, reqCtx *fasthttp.RequestCtx) {
	body := decoderHandler(reqCtx, validateRm)
	err := git.Rm(ctx, db,
		body.Entity, body.Table, body.Column, body.Branch, body.Value, body.Type, body.Opts)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}
