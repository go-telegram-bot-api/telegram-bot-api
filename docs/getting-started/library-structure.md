# Library Structure

This library is generally broken into three components you need to understand.

## Configs

Configs are collections of fields related to a single request. For example, if
one wanted to use the `sendMessage` endpoint, you could use the `MessageConfig`
struct to configure the request. There is a one-to-one relationship between
Telegram endpoints and configs. They generally have the naming pattern of
removing the `send` prefix and they all end with the `Config` suffix. They
generally implement the `Chattable` interface. If they can send files, they
implement the `Fileable` interface.

## Helpers

Helpers are easier ways of constructing common Configs. Instead of having to
create a `MessageConfig` struct and remember to set the `ChatID` and `Text`,
you can use the `NewMessage` helper method. It takes the two required parameters
for the request to succeed. You can then set fields on the resulting
`MessageConfig` after it's creation. They are generally named the same as
method names except with `send` replaced with `New`.

## Methods

Methods are used to send Configs after they are constructed. Generally,
`Request` is the lowest level method you'll have to call. It accepts a
`Chattable` parameter and knows how to upload files if needed. It returns an
`APIResponse`, the most general return type from the Bot API. This method is
called for any endpoint that doesn't have a more specific return type. For
example, `setWebhook` only returns `true` or an error. Other methods may have
more specific return types. The `getFile` endpoint returns a `File`. Almost
every other method returns a `Message`, which you can use `Send` to obtain.

There's lower level methods such as `MakeRequest` which require an endpoint and
parameters instead of accepting configs. These are primarily used internally.
If you find yourself having to use them, please open an issue.
