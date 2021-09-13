# GoLox

[![Go](https://github.com/goropikari/golox/actions/workflows/go.yml/badge.svg)](https://github.com/goropikari/golox/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goropikari/golox)](https://goreportcard.com/report/github.com/goropikari/golox)

This is a go implementation of jlox (tree-walk interpreter) by [munificent/craftinginterpreters](https://github.com/munificent/craftinginterpreters).

Related implementation: [goropikari/tlps](https://github.com/goropikari/tlps) Lox, off-side rule version

```bash
git clone https://github.com/goropikari/golox
cd golox
docker build -t golox .
docker run -it golox # launch REPL
```

# Todo

- [x] escape sequence
- [x] import another file
  - [ ] detect circular import
- [ ] support varargs
- [ ] support IO
