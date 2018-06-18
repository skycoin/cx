import time

def factorial(n):
    if n == 0: 
        return 1 
    else: 
        return n * factorial(n - 1) 

def test():
    start = time.time()
    
    print "Factorial 10: {}".format(factorial(10))

    end = time.time()
    print (end - start)

test()