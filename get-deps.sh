#!/usr/bin/env bash

TOOLCHAIN_VER="5.3.0-r07"
OLD_TOOLCHAIN_VER="5.2.1-r06"

if [[ -d vicos-sdk/dist/$OLD_TOOLCHAIN_VER/ ]]; then
  echo "Removing old $OLD_TOOLCHAIN_VER toolchain"
  rm -rf vicos-sdk/dist/$OLD_TOOLCHAIN_VER/
fi

if [[ ! -d vicos-sdk/dist/$TOOLCHAIN_VER/ ]]; then
  echo Getting deps...
  mkdir -p vicos-sdk/dist/$TOOLCHAIN_VER/
  cd vicos-sdk/dist/$TOOLCHAIN_VER/
  wget https://github.com/os-vector/wire-os-externals/releases/download/$TOOLCHAIN_VER/vicos-sdk_"$TOOLCHAIN_VER"_amd64-linux.tar.gz
  tar -xzvf vicos-sdk_"$TOOLCHAIN_VER"_amd64-linux.tar.gz
  rm vicos-sdk_"$TOOLCHAIN_VER"_amd64-linux.tar.gz
else
  echo You already have the $TOOLCHAIN_VER toolchain installed. Exiting...
  exit
fi
