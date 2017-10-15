# patchy : Easy patching of your Postgres database

## Overview [![GoDoc](https://godoc.org/github.com/chilts/patchy?status.svg)](https://godoc.org/github.com/chilts/patchy) [![Build Status](https://travis-ci.org/chilts/patchy.svg?branch=master)](https://travis-ci.org/chilts/patchy)

Patchy aims to provide you both a binary you can use to patch your database, or a library you can use when starting
your app to make sure your database is patched to the level required by your app. You could also run this as a separate
program as part of your deployment workflow.

Patches are written in simple files ("0001-forward.sql" and "0001-reverse.sql") and patchy will run all patches
required to bring your database up to the level required.

## Install

```
go get github.com/chilts/patchy
```

## Example

Open the database, set a couple of options and call `patchy.Patch`:

```
db, err := sql.Open("postgres", "...")

level := 17
opts := patchy.Options{
	Dir:           "schema",
	PropertyTable: "property",
}
newLevel, err := patchy.Patch(db, level, &opts)
if err != nil {
	// `newLevel` will be the level we managed to patch to
    // but we didn't make it to the patch requested.
}
// `newLevel` will be 17 and `err` will be `nil`
```

## Best Practices ##

Remember that sometimes you will patch your database prior to updating your deployed code. This will happen either
because you're deploying a canary for your new release and therefore it needs the updated schema, or because you're
just stopping your (single) webserver, upgrading, and patching the database when the new webserver restarts.

Ultimately this means that you should make sure that if you have two programs in production that require v4 or v5 of
the database schema, that you make your schema changes backwards compatible and can be used by two successive versions
of the webserver.

To do this, you may have to perform changes over two or (possibly) three schema patches. In reality this isn't a big
problem but it does mean you have to forward plan a little bit. When I was working at Mozilla working on Firefox
Accounts, we also made sure each program checked the patch version of the database to make sure it was either at the
version it expected, or the version+1 of the version expected. It quit if the database didn't satisfy the level
expected.

We also noticed some drawbacks to this kind of database patching, however in my opinion, this is a great procedure for
small to medium sized projects. Also, we used MySql at Mozilla which was also a pain because it didn't have
transactional DDL like Postgres does. Either way, I think we could have gotten around some of our difficulties by
separating our patches out into separate releases. (See the
[mysql-patcher](https://www.npmjs.com/package/mysql-patcher) project and it's use in the
[fxa-auth-db-mysql](https://github.com/mozilla/fxa-auth-db-mysql/blob/e0088495b7a8a56af0a9d3b823bab801a7745c3f/bin/db_patcher.js)
project. There is also some notes on
[stored procedures here](https://github.com/mozilla/fxa-auth-db-mysql/blob/526ee73bfb9ea0c77c92d73bfcc20d6975fa453b/lib/db/schema/README.md#stored-procedures-and-future-patches)
which might be of interest.

## Author

By [Andrew Chilton](https://chilts.org/), [@twitter](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com/), [@AppsAttic](https://twitter.com/AppsAttic).

## License

Apache 2.0.

(Ends)
