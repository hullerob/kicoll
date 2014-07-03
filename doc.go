// © 2014, Robert Hülle

/*

KiColl regenerates `collections.json` on Kindle 3.

Collection names are provided per-directory in special file `.collection.name`.
Collection name must be on single line, but more collections may be specified
on separate lines.

Access times in `collections.json` are preserved, but empty collections
are not written back. Be careful, running this command will completely
overwrite collection in your Kindle.

How it works:

File `collections.json` contains collection info in JSON format. Books
in collections are referred by some kind of hash/ID. For PDF files it is
sha1 hash of their full name, e.g.
	"*" + hex(sha1("/mnt/us/documents/foo.pdf"))
For MOBI files, situation is more complicated. Some files, usually created
by Calibre, contain UUID by which they are referenced, e.g.
	"#" + UUID + "^EBOK"
Some MOBI files use different types of ID. See `hash.go` for supported types
of IDs.

This command is intended to be run on jailbroken Kindle 3, through launchpad,
or by other means. But it can be run on desktop computer by specifying
different Kindle directory, e.g. path where Kindle is mounted. For purposes
of hashing filepath, this prefix is stripped and replaced by "/mnt/us".

*/
package main
