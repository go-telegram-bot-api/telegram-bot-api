# Change Log

## v5.4.0

- Remove all methods that return `(APIResponse, error)`.
  - Use the `Request` method instead.
  - For more information, see [Library Structure][library-structure].
- Remove all `New*Upload` and `New*Share` methods, replace with `New*`.
  - Use different [file types][files] to specify if upload or share.
- Rename `UploadFile` to `UploadFiles`, accept `[]RequestFile` instead of a
  single fieldname and file.
- Fix methods returning `APIResponse` and errors to always use pointers.
- Update user IDs to `int64` because of Bot API changes.
- Add missing Bot API features.

[library-structure]: ./getting-started/library-structure.md#methods
[files]: ./getting-started/files.md
