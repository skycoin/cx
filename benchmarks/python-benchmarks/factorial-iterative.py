from datetime import datetime

def factorial(n):
    resultado = 1

    for i in range(n):
        resultado *= (i + 1)
    
    return resultado

def test():
    start = datetime.now()
    print(f'Factorial 10: {factorial(10)}')
    end = datetime.now()

    delta = end - start
    print(f'{delta.total_seconds():.3f}s')

test()