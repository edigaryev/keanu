# Keanu

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
