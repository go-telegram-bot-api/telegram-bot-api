# Golang Telegram bot using the Bot API
A simple Golang bot for the Telegram Bot API

Really simple bot for interacting with the Telegram Bot API, not nearly done yet. Expect frequent breaking changes!

All methods have been added, and all features should be available.
If you want a feature that hasn't been added yet, open an issue and I'll see what I can do.

There's a few plugins in here, named as `plugin_*.go`.

## Getting started

After installing all the dependencies, run

```
go build
./telegram-bot-api -newbot
```

Fill in any asked information, enable whatever plugins, etc.

## Plugins

All plugins implement the `Plugin` interface.

```go
type Plugin interface {
	GetName() string
	GetCommands() []string
	GetHelpText() []string
	GotCommand(string, Message, []string)
	Setup()
}
```

`GetName` should return the plugin's name. This must be unique!

`GetCommands` should return a slice of strings, each command should look like `/help`, it must have the forward slash!

`GetHelpText` should return a slice of strings with each command and usage. You many include any number of items in here.

`GotCommand` is called when a command is executed for this plugin, the parameters are the command name, the Message struct, and a list of arguments passed to the command. The original text is available in the Message struct.

`Setup` is called when the bot first starts, if it needs any configuration, ask here.

To add your plugin, you must edit a line of code and then run the `go build` again.

```go
// current version
plugins = []Plugin{&HelpPlugin{}, &ManagePlugin{}}

// add your own plugins
plugins = []Plugin{&HelpPlugin{}, &FAPlugin{}, &ManagePlugin{}}
```
