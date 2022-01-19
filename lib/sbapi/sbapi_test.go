package sbapi

import (
	"encoding/json"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestPageWithRawMessage_UnmarshalJSON(t *testing.T) {
	jsonText := `
{
    "id": "57c7df23d25ef00f00100aa6",
    "title": "check",
    "image": "https://gyazo.com/b8b95342853b35152ded9341ee704fd7/raw",
    "descriptions": [
      "[https://gyazo.com/b8b95342853b35152ded9341ee704fd7.png]"
    ],
    "user": {
      "id": "566f8b954fb08e1100af5c5b"
    },
    "pin": 0,
    "views": 3677,
    "linked": 0,
    "commitId": "57c7df2ed25ef00f00100aad",
    "created": 1472713402,
    "updated": 1472716590,
    "accessed": 1641763683,
    "snapshotCreated": 1540619203,
    "pageRank": 0

    , "__extra__": true
  }
`
	var page Page

	err := json.Unmarshal([]byte(jsonText), &page)
	assert.NilError(t, err)

	assert.Assert(t, page.Title == "check")

	t.Run("unkown fields can be marshaled", func(t *testing.T) {
		encodedJSON, err := json.Marshal(page)
		assert.NilError(t, err)
		assert.Assert(t, is.Contains(string(encodedJSON), `"__extra__"`))
	})

	t.Run("can work fine without rawJSON", func(t *testing.T) {
		page.rawJSON = nil
		encodedJSON, err := json.Marshal(page)
		assert.NilError(t, err)

		var roundTrippedPage Page
		err = json.Unmarshal(encodedJSON, &roundTrippedPage)
		assert.NilError(t, err)

		roundTripJSON, err := json.Marshal(roundTrippedPage)
		assert.NilError(t, err)

		assert.Equal(t, string(encodedJSON), string(roundTripJSON))
	})
}
