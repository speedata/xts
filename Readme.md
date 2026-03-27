# XTS - XML typesetting system

XTS turns XML data into PDF. You write a layout in XML, point it at your data, and get a fully typeset document -- no GUI, no manual intervention. Think product catalogs, price lists, data sheets, or anything where content changes but the design stays the same.

Under the hood, XTS uses [boxes and glue](https://github.com/boxesandglue/boxesandglue), a Go library that implements TeX's typesetting algorithms. If you know the [speedata Publisher](https://github.com/speedata/publisher/), XTS is its next-generation successor.

## Getting started

Grab the latest release from the [releases page](https://github.com/speedata/xts/releases/latest), unzip, add `bin/` to your PATH, and you're good to go:

```
xts new hello
cd hello
xts
open xts.pdf
```

That's it -- your first PDF from XML.

## Documentation

- **[Manual and reference](https://doc.speedata.de/xts/)** -- everything from "Hello World" to advanced layouts
- **[Examples](https://github.com/speedata/xts-examples)** -- complete, runnable projects you can learn from

## Building from source

```
git clone https://github.com/speedata/xts.git
cd xts
rake build
```

Needs Go 1.21+ and Ruby/rake (or just run the `go build` command from the Rakefile directly).

## License

BSD -- see [License.md](License.md)

## Contact

Patrick Gundlach, gundlach@speedata.de
