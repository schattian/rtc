package server

import (
	"encoding/json"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/sebach1/rtc/config"
	"github.com/sebach1/rtc/git"
	"github.com/sebach1/rtc/internal/name"
	"github.com/valyala/fasthttp"
)

// Router is a handler which redirects to handlers
func Router(reqCtx *fasthttp.RequestCtx) {
	db := databaseHandler(reqCtx)
	switch string(reqCtx.Path()) {
	case "/add":
		addHandler(reqCtx, db)
	case "/rm":
		rmHandler(reqCtx, db)
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

func addHandler(reqCtx *fasthttp.RequestCtx, db *sqlx.DB) {
	body := decoderHandler(reqCtx, validateAdd)
	err := git.Add(reqCtx, db,
		body.Entity, body.Table, body.Column, body.Branch, body.Value, body.Type, body.Opts,
	)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}

func rmHandler(reqCtx *fasthttp.RequestCtx, db *sqlx.DB) {
	body := decoderHandler(reqCtx, validateRm)
	err := git.Rm(reqCtx, db,
		body.Entity, body.Table, body.Column, body.Branch, body.Value, body.Type, body.Opts,
	)
	if err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}
