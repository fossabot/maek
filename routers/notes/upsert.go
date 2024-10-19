package notes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/karngyan/maek/routers/models"

	"github.com/karngyan/maek/domains/notes"
	"github.com/karngyan/maek/routers/base"
)

func Upsert(ctx *base.WebContext) {
	nuuid := strings.TrimSpace(ctx.Input.Param(":note_uuid"))
	if nuuid == "" {
		base.BadRequest(ctx, map[string]any{
			"note_uuid": "note_uuid is required",
		})
		return
	}

	var req models.Note

	if err := ctx.DecodeJSON(&req); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	rctx := ctx.Request.Context()

	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		base.InternalError(ctx, err)
		return
	}

	note, err := notes.UpsertNoteCtx(rctx, &notes.UpsertNoteRequest{
		Uuid:      nuuid,
		Content:   string(contentBytes),
		Favorite:  req.Favorite,
		Trashed:   req.Trashed,
		Created:   req.Created,
		Updated:   req.Updated,
		Workspace: ctx.Workspace,
		CreatedBy: ctx.User, // only used if it's an insert, otherwise ignored
		UpdatedBy: ctx.User,
	})
	if err != nil {
		base.InternalError(ctx, err)
		return
	}

	uiNote, err := modelForNote(note)
	if err != nil {
		base.InternalError(ctx, err)
		return
	}

	base.Respond(ctx, map[string]any{
		"note": uiNote,
	}, http.StatusOK)
}
