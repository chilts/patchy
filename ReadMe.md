# patchy : Easy patching of your Postgres database

## Overview [![GoDoc](https://godoc.org/github.com/chilts/patchy?status.svg)](https://godoc.org/github.com/chilts/patchy) [![Build Status](https://travis-ci.org/chilts/patchy.svg?branch=master)](https://travis-ci.org/chilts/patchy)

Patchy aims to provide you both a binary you can use to patch your database, or a library you can use when starting your app to make sure your database is patched to the level required by your app.

Patches are written in simple files ("forward-00000-00001.sql" and "reverse-00001-00000.sql") and patchy will run all patches required to bring your database up to the level required.


## Install

```
go get github.com/chilts/patchy
```

## Example

```
level, err := patchy.Patch(3)

```

## Author

ToDo.

## License

Apache 2.0.
