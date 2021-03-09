from datetime import datetime

def digital_root (n):
    ap = 0
    n = abs(int(n))
    while n >= 10:
        n = sum(int(digit) for digit in str(n))
        ap += 1
    return ap, n
 
def test():
    start = datetime.now()
    persistance, root = digital_root(79563)
    print(f' 79563 has additive persistance {persistance} and digital root {root}.')
    end = datetime.now()
    delta = end - start
    print(f'{delta.total_seconds():.3f}s')

test()
