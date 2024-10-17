package notes_test

import (
	"fmt"
	"testing"

	approvals "github.com/customerio/go-approval-tests"
	"github.com/stretchr/testify/assert"

	"github.com/karngyan/maek/zarf/tests"
)

func TestUpsert(t *testing.T) {
	defer tests.CleanDBRows()

	cs := tests.NewClientStateWithUser(t)

	var testCases = []struct {
		name           string
		uuid           string
		content        map[string]any
		favorite       bool
		expectedStatus int
	}{
		{
			name: "valid note first time",
			uuid: "123",
			content: map[string]any{
				"dom": "some paragraph",
			},
			favorite:       true,
			expectedStatus: 200,
		},
		{
			name: "update created note",
			uuid: "123",
			content: map[string]any{
				"dom": "some more text",
			},
			favorite:       false,
			expectedStatus: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr, err := cs.Put(fmt.Sprintf("/v1/workspaces/%d/notes/%s", cs.Workspace.Id, tc.uuid), map[string]any{
				"content":  tc.content,
				"favorite": tc.favorite,
			})

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedStatus, rr.Code)

			approvals.VerifyJSONBytes(t, rr.Body.Bytes())
		})
	}
}
