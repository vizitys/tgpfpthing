# tgpfpthing

Go project for eventually reordering Telegram profile pictures

## Usage

Get your app id and api hash from [telegram.org](https://core.telegram.org/api/obtaining_api_id). Create a `.env` file with the following content:

```env
APP_ID=(your app id)
API_HASH=(your api hash)
PHONE_NUMBER=(your phone number in international format without the +)
```

Install the dependencies with `go mod tidy`. Run the program with `go run tgpfpthing.go`. (or use the vscode launch configuration)
