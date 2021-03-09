from datetime import datetime

def factorial(n):
    if n == 0: 
        return 1 
    else: 
        return n * factorial(n - 1) 

def test():
    start = datetime.now()
    
    print(f'Factorial 10: {factorial(10)}')

    end = datetime.now()
    delta = end - start

    print(f'{delta.total_seconds():.3f}s')

test()