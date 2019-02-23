from Crypto.Util import number
import os
import sys

def chinese_remainder(a, n):
    sum = 0
    prod = reduce(lambda a, b: a*b, n)
    for n_i, a_i in zip(n, a):
        p = prod / n_i
        sum += a_i * number.inverse(p, n_i) * p
    return sum % prod

def searchForRoot(n, r):
    a = 2
    b = n/r
    while a < b:
        c = (a+b)/2
        p = pow(c,r)
        if p < n:
            a = c
        elif p > n:
            b = c
        else:
            return c
    return a

def readFile(filename):
    f = open(filename, "rb")
    contents = f.read()
    f.close()
    return number.bytes_to_long(contents)

def solve(rootDir, e):
    a = []
    n = []
    for i in range(0, e):
        a.append(readFile(os.path.join(rootDir, "cipher" + str(i))))
        n.append(readFile(os.path.join(rootDir, "pubkey" + str(i))))
    crt = chinese_remainder(a, n)
    m = searchForRoot(crt, e)
    return number.long_to_bytes(m)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print "You must specify a path to the broadcast, like so:"
        print "\tpython solution.py broadcast/"
        sys.exit(1)
    
    print solve(sys.argv[1], 7)



