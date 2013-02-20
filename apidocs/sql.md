Note to self on SQL:

We'll use:

http://godoc.org/code.google.com/p/go-sqlite/go1/sqlite3

But we must not use a single connection - just have each part of reservoir open its own connection, operate on it, and close it afterwards.

So for pings, the pings listener would have its own connection to the db.
So would the job logger.
Etc.