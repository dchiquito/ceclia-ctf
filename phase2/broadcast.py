import Crypto
from Crypto.PublicKey import RSA
from Crypto import Random
from Crypto.Util import number
import os
import sys

def encrypt(message):
    key = RSA.generate(1024, e=7)
    return (number.long_to_bytes(key.n), key.encrypt(message, 32)[0])

def writeFile(contents, filename):
    f = open(filename, "w")
    f.write(contents)
    f.close()

def broadcast(message, rootDir, copies):
    if not os.path.exists(rootDir):
        os.makedirs(rootDir)
    for i in range(0, copies):
        print "\tEncrypting message " + str(i) + "..."
        pubkeyFile = os.path.join(rootDir, "pubkey" + str(i))
        cipherFile = os.path.join(rootDir, "cipher" + str(i))
        n, cipher = encrypt(message)
        print "Cipher: " + str(number.bytes_to_long(cipher))
        writeFile(n, pubkeyFile)
        writeFile(cipher, cipherFile)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print "You must specify a path to the broadcast and a message to broadcast, like so:"
        print "\tpython broadcast.py broadcast/ ThisIsMyMessage"
        sys.exit(1)
    broadcast(sys.argv[2], sys.argv[1], 10)

