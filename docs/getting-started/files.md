# Files

Telegram supports specifying files in many different formats. In order to
accommodate them all, there are multiple structs and type aliases required.

All of these types implement the `RequestFileData` interface.

| Type         | Description                                                               |
| ------------ | ------------------------------------------------------------------------- |
| `FilePath`   | A local path to a file                                                    |
| `FileID`     | Existing file ID on Telegram's servers                                    |
| `FileURL`    | URL to file, must be served with expected MIME type                       |
| `FileReader` | Use an `io.Reader` to provide a file. Lazily read to save memory.         |
| `FileBytes`  | `[]byte` containing file data. Prefer to use `FileReader` to save memory. |

## `FilePath`

A path to a local file.

```go
file := tgbotapi.FilePath("tests/image.jpg")
```

## `FileID`

An ID previously uploaded to Telegram. IDs may only be reused by the same bot
that received them. Additionally, thumbnail IDs cannot be reused.

```go
file := tgbotapi.FileID("AgACAgIAAxkDAALesF8dCjAAAa_…")
```

## `FileURL`

A URL to an existing resource. It must be served with a correct MIME type to
work as expected.

```go
file := tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg")
```

## `FileReader`

Use an `io.Reader` to provide file contents as needed. Requires a filename for
the virtual file.

```go
var reader io.Reader

file := tgbotapi.FileReader{
    Name: "image.jpg",
    Reader: reader,
}
```

## `FileBytes`

Use a `[]byte` to provide file contents. Generally try to avoid this as it
results in high memory usage. Also requires a filename for the virtual file.

```go
var data []byte

file := tgbotapi.FileBytes{
    Name: "image.jpg",
    Bytes: data,
}
```
