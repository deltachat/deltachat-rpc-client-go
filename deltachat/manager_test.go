package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountManager_String(t *testing.T) {
	manager := server.AccountManager()
	assert.NotEmpty(t, manager.String())
}

func TestAccountManager_SelectedAccount(t *testing.T) {
	manager := server.AccountManager()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)
	selected, err := manager.SelectedAccount()
	assert.Nil(t, err)
	assert.Equal(t, acc.Id, selected.Id)

	_, err = manager.AddAccount()
	assert.Nil(t, err)
	selected, err = manager.SelectedAccount()
	assert.Nil(t, err)
	assert.NotEqual(t, acc.Id, selected.Id)

	err = acc.Select()
	assert.Nil(t, err)
	selected, err = manager.SelectedAccount()
	assert.Nil(t, err)
	assert.Equal(t, acc.Id, selected.Id)
}

func TestAccountManager_Accounts(t *testing.T) {
	manager := server.AccountManager()

	accounts, err := manager.Accounts()
	assert.Nil(t, err)
	count := len(accounts)

	_, err = manager.AddAccount()
	assert.Nil(t, err)

	accounts, err = manager.Accounts()
	assert.Nil(t, err)
	assert.Equal(t, len(accounts), count+1)
}

func TestAccountManager_Remove(t *testing.T) {
	manager := server.AccountManager()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	accounts, err := manager.Accounts()
	assert.Nil(t, err)
	assert.NotEmpty(t, accounts)
	count := len(accounts)

	assert.Nil(t, acc.Remove())

	accounts, err = manager.Accounts()
	assert.Nil(t, err)
	assert.Equal(t, len(accounts), count-1)
}

func TestAccountManager_StartIO(t *testing.T) {
	manager := server.AccountManager()
	assert.Nil(t, manager.StartIO())
}

func TestAccountManager_StopIO(t *testing.T) {
	manager := server.AccountManager()
	assert.Nil(t, manager.StopIO())
}

func TestAccountManager_MaybeNetwork(t *testing.T) {
	manager := server.AccountManager()
	assert.Nil(t, manager.MaybeNetwork())
}

func TestAccountManager_SystemInfo(t *testing.T) {
	manager := server.AccountManager()
	sysinfo, err := manager.SystemInfo()
	assert.Nil(t, err)
	assert.NotEmpty(t, sysinfo["deltachat_core_version"], "invalid deltachat_core_version")
}

func TestAccountManager_SetTranslations(t *testing.T) {
	manager := server.AccountManager()
	assert.Nil(t, manager.SetTranslations(map[uint]string{1: "test"}))
}
