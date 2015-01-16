onchange [![Build Status](https://travis-ci.org/minitools/onchange.svg?branch=master)](https://travis-ci.org/minitools/onchange)
========

Monitors changes to files in local directory and executes a command on each change event.

The utility supports:
* Filtering names by pattern, with ```-name <pattern>```

Examples
--------
For example, to rebuild a Go project every time a local file has changed:

```$ onchange -name "*.go" go build```

To upload a binary to a remote machine, every time it's built:

```$ onchange -name some_app scp some_app $(target)```

To Do
----
* Support monitoring of subfolders (```-r``` option)

Questions, ideas?
-----------------
Contact me at Tom Paoletti <tpaoletti_JUSTDROP@users.sf.net>
