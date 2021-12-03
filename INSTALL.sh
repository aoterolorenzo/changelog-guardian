#!/bin/sh

LAST_VERSION=$(curl -sf https://gitlab.com/aoterocom/changelog-guardian/-/raw/main/VERSION)

case "$LAST_VERSION" in
  "")
      echo "Unable to found latest version. Exiting..."
      exit
    ;;
  *)
      echo "Found latest version $LAST_VERSION"
    ;;
esac

STATUS_CODE=$(curl -L -I HEAD https://gitlab.com/aoterocom/changelog-guardian/-/releases/v"$LAST_VERSION"/downloads/changelog-guardian_"$LAST_VERSION"_"$(uname -s)"_"$(uname -m)".tar.gz --silent | head -n 1)

case "$STATUS_CODE" in
  *"302"*)
    mkdir /tmp/changelog-guardian
      cd /tmp/changelog-guardian || exit

      echo "Downloading..."
      curl -L https://gitlab.com/aoterocom/changelog-guardian/-/releases/v"$LAST_VERSION"/downloads/changelog-guardian_"$LAST_VERSION"_"$(uname -s)"_"$(uname -m)".tar.gz --output changelog-guardian_"$LAST_VERSION"_"$(uname -s)"_"$(uname -m)".tar.gz
      echo "Decompressing..."
      tar -zxvf changelog-guardian_"$LAST_VERSION"_"$(uname -s)"_"$(uname -m)".tar.gz > /dev/null 2>&1
      echo "Installing to /usr/local/bin..."
      cp -nf changelog-guardian /usr/local/bin
      cd /tmp || exit
      rm -Rf /tmp/changelog-guardian
      echo "Success."
    ;;
  *)
      echo "Package not found for Version: $LAST_VERSION, OS: $(uname -s), ARCH: $(uname -m)"
      ;;
esac