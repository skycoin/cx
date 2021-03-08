import time

def factorial(n):
    if n == 0: 
        return 1 
    else: 
        return n * factorial(n - 1) 

def test():
    start = time.time()
    
    print(f'Factorial 10: {factorial(10)}')

    end = time.time()
    print (end - start)

test()