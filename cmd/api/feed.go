package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
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

	ctx := r.Context()

	feed, err := app.store.Queries.GetPostsWithMetaData(ctx, 1)
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
