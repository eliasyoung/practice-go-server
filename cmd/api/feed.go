package main

import (
	"net/http"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetUserFeed godoc
//
//	@Summary		Get a user's feed
//	@Description	Get a user's feed
//	@Tags			Feed
//	@Accpet			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			offset	query		int		false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]db.GetPostsWithMetaDataRow
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := db.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   make([]string, 0),
		Search: "",
		Since:  "",
		Until:  "",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validatePagiationQuery(&fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	params := db.GetPostsWithMetaDataParams{
		UserID:  1,
		Column2: fq.Sort,
		Limit:   int32(fq.Limit),
		Offset:  int32(fq.Offset),
		Column5: pgtype.Text{
			String: fq.Search,
			Valid:  true,
		},
		Tags: fq.Tags,
	}

	if fq.Since != "" {
		ts, err := db.ParseToPgTimestamptz(fq.Since)
		if err == nil {
			params.CreatedAt = ts
		}
	}

	if fq.Until != "" {
		ts, err := db.ParseToPgTimestamptz(fq.Until)
		if err == nil {
			params.CreatedAt_2 = ts
		}
	}

	ctx := r.Context()

	feed, err := app.store.Queries.GetPostsWithMetaData(ctx, params)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if len(feed) > 0 {
		count := len(feed)

		for i := 0; i < count; i++ {
			feed[i].Tags = responseSliceFormater(feed[i].Tags)
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
