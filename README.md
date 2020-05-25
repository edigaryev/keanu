# Keanu

[![Cirrus CI](https://api.cirrus-ci.com/github/edigaryev/keanu.svg)](https://cirrus-ci.com/github/edigaryev/keanu)
[![GoDoc](https://godoc.org/github.com/edigaryev/keanu?status.svg)](https://godoc.org/github.com/edigaryev/keanu/preprocessor?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/edigaryev/keanu)](https://goreportcard.com/report/github.com/edigaryev/keanu)
[![Codecov](https://codecov.io/gh/edigaryev/keanu/branch/master/graph/badge.svg)](https://codecov.io/gh/edigaryev/keanu)
[![LGTM alerts](https://img.shields.io/lgtm/alerts/g/edigaryev/keanu.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/edigaryev/keanu/alerts)

Preprocesses YAML files with `matrix` modifiers, semantically similar to [what Cirrus CI does](https://cirrus-ci.org/guide/writing-tasks/#matrix-modification).

## Installing

```
go get github.com/edigaryev/keanu
```

Make sure that the `$GOPATH/bin` directory is in `PATH` (see [article in the Go wiki](https://github.com/golang/go/wiki/SettingGOPATH) for more details).

## Using

The most simplest invocation produces output to `stdout`:

```
keanu .cirrus.yml
```

To write the output to a file, specify the file path as the second argument:

```
keanu .cirrus.yml .cirrus-processed.yml
```
