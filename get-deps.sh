#!/usr/bin/env bash

if [[ ! -d ~/.anki/vicos-sdk/dist/5.3.0-r07 ]]; then
  echo Getting deps...
  mkdir ~/.anki/vicos-sdk/dist/5.3.0-r07
  cd ~/.anki/vicos-sdk/dist/5.3.0-r07
  wget https://froggitti.net/5.3.0-r07.tar.gz
  gunzip 5.3.0-r07.tar.gz
  tar -xvf 5.3.0-r07.tar
else
  echo You already have the 5.3.0-r07 toolchain installed. Exiting...
  exit
fi
