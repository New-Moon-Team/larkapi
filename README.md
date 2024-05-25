Deciding which workspace information to collect

Collecting workspace information

Filtering to most relevant information

# larkapi

This is a Go package for interacting with spreadsheets.

## Installation

First, you need to install the package. Run the following command:

```sh
go get github.com/yourusername/larkapi
```

Replace `yourusername` with your actual GitHub username.

## Usage

### Environment Variables

You need to set the following environment variables:

- `APP_ID`: Your application ID.
- `APP_SECRET`: Your application secret.

You can set these in your shell, or directly in your Go code:

```go
os.Setenv("APP_ID", "your-app-id")
os.Setenv("APP_SECRET", "your-app-secret")
```

### Working with Sheets

Here's an example of how to use the [`Sheet`](command:_github.copilot.openSymbolFromReferences?%5B%7B%22%24mid%22%3A1%2C%22path%22%3A%22%2FC%3A%2FUsers%2FToanT%2FDocuments%2Frepo%2Flarkapi%2Fsheets%2Fsheet.go%22%2C%22scheme%22%3A%22file%22%7D%2C%7B%22line%22%3A466%2C%22character%22%3A5%7D%5D "sheets/sheet.go") struct from the [`sheet.go`](command:_github.copilot.openSymbolInFile?%5B%22sheets%2Fsheet.go%22%2C%22sheet.go%22%5D "sheets/sheet.go") file:

```go
import "larkapi/sheets"

// Initialize a new sheet
sh := sheets.Sheet{...}

// Get all data from the sheet
var result interface{}
err := sh.GetAll(&result)
if err != nil {
    // handle error
}
```

### Working with Spreadsheets

Here's an example of how to use the `Spreadsheet` struct from the [`spreadsheet.go`](command:_github.copilot.openSymbolInFile?%5B%22sheets%2Fspreadsheet.go%22%2C%22spreadsheet.go%22%5D "sheets/spreadsheet.go") file:

```go
import "larkapi/sheets"

// Initialize a new spreadsheet
ss := sheets.Spreadsheet{...}

// Perform operations on the spreadsheet
...
```

## Examples

You can find more usage examples in the [`example`](example/) directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.