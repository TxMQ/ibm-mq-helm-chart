#!/bin/bash -x

password="password"
stem=key

outdir=$1

keydbpath="$outdir/$stem.kdb"
rdbpath="$outdir/$stem.rdb"
sthpath="$outdir/$stem.sth"

/opt/mqm/bin/runmqckm -keydb -create -db $keydbpath -pw $password -type cms -stash

chmod g+rw,o+r $keydbpath $rdbpath $sthpath

