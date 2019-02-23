#!/bin/sh

mkdir build/phase1

echo "Building Phase 1: Playfair Cipher & Steganography..."

echo "Zipping Phase 2..."
cd build/phase2
zip -r ../../build/phase1/phase1.zip * .flag
cd ../..

echo "Hiding Phase 2 .zip inside image file..."
cat phase1/us.jpeg build/phase1/phase1.zip >> build/phase1/puzzle.jpeg
rm build/phase1/phase1.zip

