# Changelog

## Unreleased

- return error in Rpc.Start() if the Rpc is already started
- bug fix: initialize Rpc.closed to true in NewRpcIO()

## v0.6.0

- bug fix in `Account.QrCode()`
- avoid panics when stopping Rpc and closing event channels
- bug fix: move `Connectivity()` from AccountManager to Account
- remove pasphrase argument from Account.ExportSelfKeys() and Account.ImportSelfKeys()
- fix bug: remove FreshMsgCount() from Account
- fix bug in Account.AddDeviceMsg()
- fix bug in Account.QueryChatListItems() and Account.QueryChatListEntries()
- `Chat.QrCode()` now returns `(string, string, error)` instead of `([2]string, err)`

## v0.5.0

- `Account.QrCode()` now returns `(string, string, error)` instead of `([2]string, err)`

## v0.4.0

- fix bug in Account.FreshMsgsInArrivalOrder()

## v0.3.0

- avoid deadlocks calling Bot.Stop() when the bot is already stopped
- fix bug: Bot.Run() never returning even if Bot.Stop() was called

## v0.2.0

- removed Bot.RunForever(), and now Bot.Run() doesn't require a channel as argument, added Bot.Stop()
  to stop processing events

## v0.1.0

- initial release
