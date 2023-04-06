package deltachat

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	t.Parallel()
	var timestamps []Timestamp
	data := []byte(`[1680779737]`)

	assert.Nil(t, json.Unmarshal(data, &timestamps))
	assert.Equal(t, Timestamp{time.Unix(1680779737, 0)}, timestamps[0])

	bytes, err := json.Marshal(timestamps)
	assert.Nil(t, err)
	assert.Equal(t, []byte("[1680779737]"), bytes)
}
