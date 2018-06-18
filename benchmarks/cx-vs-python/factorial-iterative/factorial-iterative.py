import time

def factorial(n):
    resultado = 1

    for i in range(n):
        resultado *= (i + 1)
    
    return resultado

def test():
    start = time.time()
    
    print "Factorial 10: {}".format(factorial(10))

    end = time.time()
    print (end - start)

test()