package option

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	t.Parallel()
	var options map[string]Option[string]
	data := []byte(`{"key":"value"}`)

	assert.Nil(t, json.Unmarshal(data, &options))
	assert.Equal(t, "value", options["key"].Unwrap())

	bytes, err := json.Marshal(&options)
	assert.Nil(t, err)
	assert.Equal(t, data, bytes)
}
