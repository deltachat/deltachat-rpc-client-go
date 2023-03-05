# Changelog

## v0.3.0

- avoid deadlocks calling Bot.Stop() when the bot is already stopped
- fix bug: Bot.Run() never returning even if Bot.Stop() was called

## v0.2.0

- removed Bot.RunForever(), and now Bot.Run() doesn't require a channel as argument, added Bot.Stop()
  to stop processing events

## v0.1.0

- initial release
