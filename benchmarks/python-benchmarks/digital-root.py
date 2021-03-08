import time

def digital_root (n):
    ap = 0
    n = abs(int(n))
    while n >= 10:
        n = sum(int(digit) for digit in str(n))
        ap += 1
    return ap, n
 
def test():
    start = time.time()
    persistance, root = digital_root(79563)
    print(f' 79563 has additive persistance {persistance} and digital root {root}.')
    end = time.time()
    print (end - start)

test()
