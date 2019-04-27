fpick
===

pick files you want

[![CircleCI](https://circleci.com/gh/sawadashota/fpick/tree/master.svg?style=svg)](https://circleci.com/gh/sawadashota/fpick/tree/master)
[![GoDoc](https://godoc.org/github.com/sawadashota/fpick?status.svg)](https://godoc.org/github.com/sawadashota/fpick)
[![Go Report Card](https://goreportcard.com/badge/github.com/sawadashota/fpick)](https://goreportcard.com/report/github.com/sawadashota/fpick)
[![GolangCI](https://golangci.com/badges/github.com/sawadashota/fpick.svg)](https://golangci.com/r/github.com/sawadashota/fpick)


Usage
---

Extract filename match

```
$ fpick -s /path/to/repos/root -d ~/Downloads/envrc.bk -n .envrc
```

Regex filename match

```
$ fpick -s /path/to/repos/root -d ~/Downloads/envrc.bk -r "^\.envrc$"
```

Installation
---

```bash
brew tap sawadashota/homebrew-cheers
brew install fpick
```

License
---

MIT