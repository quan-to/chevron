#!/usr/bin/env bash


RESOLUTIONS="16 32 48 128 256 512 1024"

mkdir -p multires
rm -f multires/*

for res in ${RESOLUTIONS}
do
    echo "Rendering image ${res}x${res} px"
    inkscape -z -e multires/logo-${res}.png -w ${res} -h ${res} AgentUI.svg 2>/dev/null
done

echo "Generating Linux Icon"
cp multires/logo-512.png icon.png

echo "Generating Apple Icons"
png2icns icon.icns multires/logo-*.png

echo "Generating Windows Icons"
# Broken
#convert AgentUI.svg -bordercolor white -border 0 \
#      \( -clone 0 -resize 16x16 \) \
#      \( -clone 0 -resize 32x32 \) \
#      \( -clone 0 -resize 48x48 \) \
#      \( -clone 0 -resize 57x57 \) \
#      \( -clone 0 -resize 64x64 \) \
#      \( -clone 0 -resize 72x72 \) \
#      \( -clone 0 -resize 110x110 \) \
#      \( -clone 0 -resize 114x114 \) \
#      \( -clone 0 -resize 120x120 \) \
#      \( -clone 0 -resize 128x128 \) \
#      \( -clone 0 -resize 144x144 \) \
#      \( -clone 0 -resize 152x152 \) \
#      \( -clone 0 -resize 192x192 \) \
#      \( -clone 0 -resize 256x256 \) \
#      \( -clone 0 -resize 512x512 \) \
#      -delete 0 -alpha off -colors 256 logo.ico
echo "  The current generator is broken for windows, please use https://convertio.co/pt/svg-ico/"

rm -fr multires
