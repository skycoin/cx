from datetime import datetime

def fib(n):
    if n < 2:
        return n
    return (fib(n-2) + fib(n-1))

def test():
    start = datetime.now()
    print(fib(30))
    end = datetime.now()
    delta = end - start

    print(f'Fib Recursive time elapsed (in seconds):{delta.total_seconds():.3f}s')
test()