#!/bin/sh

./build4.sh
./build3.sh
./build2.sh
./build1.sh

echo "Archiving build directory..."
cd build
zip -r ctf.zip phase1 phase2 phase3 phase4

