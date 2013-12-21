# tablewriter

Write tab-delimited data in a formatted ascii table, complete with word-wrapping
and padding

Turns:

```
aaa aaa aaa aaa aaa aaa\tbbb bbb bbb bbb bbb bbb\tccc ccc ccc ccc ccc ccc
ddd ddd ddd ddd ddd ddd\teee eee eee eee eee eee\tfff fff fff fff fff fff
```

Into:

```
aaa aaa aaa    bbb bbb bbb    ccc ccc ccc
aaa aaa aaa    bbb bbb bbb    ccc ccc ccc

ddd ddd ddd    eee eee eee    fff fff fff
ddd ddd ddd    eee eee eee    fff fff fff
```

## API

Check out the [api][api] for exact usage, but the [example][example] will
probably be more helpful for getting started.

## Usage

```
go get github.com/mediocregopher/tablewriter
```

or if you're using [goat][goat]

```yaml
    - loc: https://github.com/mediocregopher/tablewriter.git
      type: git
      ref: v0.0.0
      loc: github.com/mediocregopher/tablewriter
```

[api]: http://godoc.org/github.com/mediocregopher/tablewriter
[example]: /example/example.go
[goat]: https://github.com/mediocregopher/goat
