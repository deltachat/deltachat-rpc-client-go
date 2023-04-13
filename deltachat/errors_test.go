package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Parallel()
	var err error

	err = &BotRunningErr{}
	assert.NotEmpty(t, err.Error())

	err = &RpcRunningErr{}
	assert.NotEmpty(t, err.Error())
}
