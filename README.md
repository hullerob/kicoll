# KiColl

Kindle Collections

## About

This is quick hack to generate file `collections.json` on Kindle 3.  It works
for me, maybe it will work for you.

## Build

    GOOS=linux GOARCH=arm go build

## Install

Copy `kicoll` binary to your `$KINDLEDIR/kicoll`.

## Run

Put something like this to `launchpad/kicoll.ini`:

    [Actions]
    C R = !/mnt/us/kicoll/kicoll

Then you can run this with launchpad by pressing `shift c r` This will scan
your `documents` directory and overwrite `system/collections.json`.

As always, backup and backup and backup.

After regenerating `collections.json`, restart is needed.

## How to create collection

    cd documents/somedir
    echo 'My New Collection' > .collection.name
    echo 'Second Collection' >> .collection.name

This will put everything in directory and subdirectories to collections
in file. One collection per line, empty lines are ignored.
If there is other `.collection.name` deeper in directory tree,
that one takes precedence.

## Why

So I don't have to click and click and click in Calibre.  And no need to
install program on other computer, everything is on my Kindle.

## Support

Only works with PDF and calibre-MOBI (what I use).
