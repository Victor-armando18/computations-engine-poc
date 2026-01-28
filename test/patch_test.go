package test

import (
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/assert"
)

func TestPatchDecode(t *testing.T) {
	patch := []byte(`[
	  { "op": "replace", "path": "/items/0/quantity", "value": 3 }
	]`)

	_, err := jsonpatch.DecodePatch(patch)
	assert.NoError(t, err)
}
