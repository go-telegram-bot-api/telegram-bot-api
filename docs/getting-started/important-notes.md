# Important Notes

The Telegram Bot API has a few potentially unanticipated behaviors. Here are a
few of them. If any behavior was surprising to you, please feel free to open a
pull request!

## Callback Queries

- Every callback query must be answered, even if there is nothing to display to
  the user. Failure to do so will show a loading icon on the keyboard until the
  operation times out.

## ChatMemberUpdated

- In order to receive `chat_member` and updates, you must explicitly add it to
  your `allowed_updates` when getting updates or setting your webhook.

## Entities use UTF16

- When extracting text entities using offsets and lengths, characters can appear
  to be in incorrect positions. This is because Telegram uses UTF16 lengths
  while Golang uses UTF8. It's possible to convert between the two, see
  [issue #231][issue-231] for more details.

[issue-231]: https://github.com/go-telegram-bot-api/telegram-bot-api/issues/231

## GetUpdatesChan

- This method is very basic and likely unsuitable for production use. Consider
  creating your own implementation instead, as it's very simple to replicate.
- This method only allows your bot to process one update at a time. You can
  spawn goroutines to handle updates concurrently or switch to webhooks instead.
  Webhooks are suggested for high traffic bots.

## Nil Updates

- At most one of the fields in an `Update` will be set to a non-nil value. When
  evaluating updates, you must make sure you check that the field is not nil
  before trying to access any of it's fields.

## User and Chat ID size

- These types require up to 52 significant bits to store correctly, making a
  64-bit integer type required in most languages. They are already `int64` types
  in this library, but make sure you use correct types when saving them to a
  database or passing them to another language.
