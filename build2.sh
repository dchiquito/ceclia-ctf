#!/bin/sh

mkdir -p build/phase2

echo "Building Phase 2: Hastad's Broadcast Attack..."

echo "Copying files over..."
cp phase2/README build/phase2/README
cp phase2/broadcast.py build/phase2/broadcast.py
cp phase2/.flag build/phase2/.flag
cp phase2/more-homework.pdf build/phase2/more-homework.pdf

cd build/phase2

echo "Generating encrypted files..."
python broadcast.py broadcast '<a href="http://127.0.0.1:9596/login"></a>'

echo "Zipping hint..."
zip are-you-sure.zip more-homework.pdf
zip hint.zip are-you-sure.zip

echo "Removing zip leftovers..."
rm more-homework.pdf
rm are-you-sure.zip
