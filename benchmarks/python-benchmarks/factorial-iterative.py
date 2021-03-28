from datetime import datetime
import csv

def factorial(n):
    resultado = 1

    for i in range(n):
        resultado *= (i + 1)
    
    return resultado

def test():
    start = datetime.now()
    #print(f'Factorial 10: {factorial(10)}')
    end = datetime.now()

    delta = end - start
    return f'{delta.total_seconds():.3f}s'

result = test()

with open('results.csv', mode='w') as csv_file:
    fieldnames = ['Language', 'Test Name', 'Input', 'Output', 'Time']
    writer = csv.DictWriter(csv_file, fieldnames=fieldnames)
    writer.writeheader()
    writer.writerow({'Language': 'python', 'Test Name': 'factorial-iterative', 'Input': '10', 'Output': factorial(10), 'Time': result})

with open('results.csv','r', newline='') as file:
    reader = csv.reader(file, delimiter=',', quoting=csv.QUOTE_NONE)
    for row in reader:
        print(f'{row[0]:<15}  {row[1]:<20} {row[2]:<15} {row[3]:<15} {row[4]:<15}')