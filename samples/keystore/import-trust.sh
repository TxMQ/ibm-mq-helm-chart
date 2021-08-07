#!/bin/bash -x

outdir=$1
capath=$2

stem=key
keydbpath="$outdir/$stem.kdb"

label="ca"

/opt/mqm/bin/runmqckm -cert -add -db $keydbpath -stashed -label $label -file $capath -format ascii
