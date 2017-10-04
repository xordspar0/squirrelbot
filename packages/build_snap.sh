#!/bin/sh
set -xe

version="$1"

# Download youtube-dl.
curl -L https://yt-dl.org/downloads/latest/youtube-dl -o bin/youtube-dl
chmod a+rx bin/youtube-dl

# Prepare snapcraft.yaml.
cd $(dirname $0)
cp snap/snapcraft.yaml.template snap/snapcraft.yaml

sed -i "s/{{version}}/$version/" snap/snapcraft.yaml

if [ $version = 'devel' ]; then
  sed -i "s/{{grade}}/devel/" snap/snapcraft.yaml
else
  sed -i "s/{{grade}}/stable/" snap/snapcraft.yaml
fi

snapcraft clean
snapcraft
