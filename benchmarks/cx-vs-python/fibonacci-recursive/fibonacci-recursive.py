import time

def fib(n):
    if n < 2:
        return n
    return (fib(n-2) + fib(n-1))

def test():
    start = time.time()
    print fib(30)
    end = time.time()
    print (end - start)

print "Fib Recursive time elapsed (in seconds):"
test()