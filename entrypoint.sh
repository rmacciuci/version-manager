#!/bin/sh
set -eu

usage() {
  echo "Uso: $0 --minor=<true|false> --major=<true|false> --patch=<true|false>"
  exit 1
}

git_setup () {
  git config --global user.email "actions@github.com"
  git config --global user.name "Version Manager"
}

INCREMENT_MAJOR=false
INCREASE_MINOR=false
INCREASE_PATCH=false

while [ $# -gt 0 ]; do
  case $1 in
    --minor=*) INCREASE_MINOR="${1#*=}" ;;
    --major=*) INCREASE_MAJOR="${1#*=}" ;;
    --patch=*) INCREASE_PATCH="${1#*=}" ;;
    * ) usage ;;
  esac
  shift
done

echo "---------- Version Manager ----------"
echo "Starting Version Manager"
echo "-------------------------------------"

echo "1. Setting up git"
git_setup

echo "2. Updating repository tags"
git fetch origin --tags --quiet

last_tag=$(git describe --tags `git rev-list --tags --max-count=1`)
echo "Last tag: $last_tag"

version=$(echo $last_tag | sed 's/v//')

major=$(echo $version | cut -d. -f1)
minor=$(echo $version | cut -d. -f2)
patch=$(echo $version | cut -d. -f3)

if ! $INCREASE_MINOR && ! $INCREASE_MAJOR && ! $INCREASE_PATCH; then
  echo "No version increase selected. You must select at least one."
  usage
fi

if [ "$INCREMENT_MAJOR" = true ]; then
  echo "Increasing major version"
  major=$((major + 1))
  minor=0
  patch=0
fi

if [ "$INCREASE_MINOR" = true ]; then
  echo "Increasing minor version"
  minor=$((minor + 1))
  patch=0
fi

if [ "$INCREASE_PATCH" = true ]; then
  echo "Increasing patch version"
  patch=$((patch + 1))
fi

next_tag="v$major.$minor.$patch"
echo "New tag: $next_tag"

git pull

echo "3. Creating new tag"
git tag $next_tag

echo "4. Pushing new tag"
git push origin $next_tag

echo "5. Done"
echo $next_tag > VERSION
echo "-------------------------------------"
