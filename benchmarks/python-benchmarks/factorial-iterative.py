import time

def factorial(n):
    resultado = 1

    for i in range(n):
        resultado *= (i + 1)
    
    return resultado

def test():
    start = time.time()
    
    print(f'Factorial 10: {factorial(10)}')

    end = time.time()
    print (end - start)

test()