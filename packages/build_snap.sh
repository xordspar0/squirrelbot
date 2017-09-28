version="$1"

cd $(dirname $0)

cp snap/snapcraft.yaml.template snap/snapcraft.yaml

sed -i '' "s/{{version}}/$version/" packages/snap/snapcraft.yaml

if [ $version == 'devel' ]; then
  sed -i '' "s/{{grade}}/devel/" packages/snap/snapcraft.yaml
else
  sed -i '' "s/{{grade}}/stable/" packages/snap/snapcraft.yaml
fi

snapcraft
