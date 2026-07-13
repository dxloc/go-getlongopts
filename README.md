# Go-getlongopts

This package helps parse the command-line argument in the short opt/long opt style.

## Installation

```Shell
go get -u github.com/dxloc/go-getlongopts
```

 Import the package to your code:

```Go
import "github.com/dxloc/go-getlongopts"
```

## Usage

First define some variable to hold the parsed value.

```Go
var a = ""
var b = ""
```

Next define the option list for the parser.

```Go
opts := []getlongopts.LongOption{
 {
  Long:    "long-a",
  Short:   "a",
  ArgType: getlongopts.ArgTypeFile,
  SetFn: func(value string) {
   a = value
  },
  Description: "Option (-a|--long-a), file argument",
 },
 {
  Long:    "long-b",
  ArgType: getlongopts.ArgTypeDir,
  SetFn: func(value string) {
   b = value
  },
  Description: "Option (--long-b), directory argument",
 },
}
```

Initialize and parse the arguments.

```Go
p := getlongopts.NewParser(opts)
p.Parse(os.Args)
```

See complete example usage in `parser_test.go`.

## License

This package is released under the [Unlicense](https://unlicense.org/).

For more information, please refer to [https://unlicense.org/](https://unlicense.org/).

## Contributing

Fell free to [open an issue](https://github.com/dxloc/go-getlongopts/issues) or [submit a pull request](https://github.com/dxloc/go-getlongopts/pull/new).
