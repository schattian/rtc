package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	"github.com/sebach1/git-crud/config"

	"github.com/sebach1/git-crud/git"
	"github.com/valyala/fasthttp"
)

// Router is a handler which redirects to handlers
func Router(reqCtx *fasthttp.RequestCtx) {
	DB := databaseHandler(reqCtx)
	switch string(reqCtx.Path()) {
	case "/add":
		addHandler(context.TODO(), DB, reqCtx)
	case "/rm":
		rmHandler(context.TODO(), DB, reqCtx)
	default:
		reqCtx.NotFound()
	}
}

func databaseHandler(reqCtx *fasthttp.RequestCtx) *sql.DB {
	dataSource := os.Getenv("DB")
	if dataSource == "" {
		dataSource = config.DefaultDataSource
	}
	DB, err := sql.Open("postgres", dataSource)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadGateway)
	}
	return DB
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

func addHandler(ctx context.Context, DB *sql.DB, reqCtx *fasthttp.RequestCtx) {
	body := decoderHandler(reqCtx, validateAdd)
	err := git.Add(ctx, DB,
		body.Entity, body.Table, body.Column, body.Branch, body.Value, body.Type, body.Opts)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}

func rmHandler(ctx context.Context, DB *sql.DB, reqCtx *fasthttp.RequestCtx) {
	body := decoderHandler(reqCtx, validateRm)
	err := git.Rm(ctx, DB,
		body.Entity, body.Table, body.Column, body.Branch, body.Value, body.Type, body.Opts)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}
