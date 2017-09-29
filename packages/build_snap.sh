#!/bin/sh
set -e

version="$1"

cd $(dirname $0)

cp snap/snapcraft.yaml.template snap/snapcraft.yaml

sed -i "s/{{version}}/$version/" snap/snapcraft.yaml

if [ $version = 'devel' ]; then
  sed -i "s/{{grade}}/devel/" snap/snapcraft.yaml
else
  sed -i "s/{{grade}}/stable/" snap/snapcraft.yaml
fi

snapcraft
