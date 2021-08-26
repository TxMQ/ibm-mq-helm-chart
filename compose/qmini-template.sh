#!/bin/bash

# not used
qmenv=$1

outdir=output/etc/mqm/qmini

mkdir -p $outdir

cat <<EOF > $outdir/qmini.yaml
qmini: |
  Log:
    LogBufferPages=0
EOF