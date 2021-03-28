from datetime import datetime
import csv 

def factorial(n):
    if n == 0: 
        return 1 
    else: 
        return n * factorial(n - 1) 

def test():
    start = datetime.now()
    
    # print(f'Factorial 10: {factorial(10)}')

    end = datetime.now()
    delta = end - start

    return f'{delta.total_seconds():.3f}s'

test_result = test()

with open('results.csv', mode='w') as csv_file:
    fieldnames = ['Language', 'Test Name', 'Input', 'Time']
    writer = csv.DictWriter(csv_file, fieldnames=fieldnames)
    writer.writeheader()
    writer.writerow({'Language': 'python', 'Test Name': 'factorial-recursive', 'Input': '10','Time': test_result})

with open('results.csv','r', newline='') as file:
    reader = csv.reader(file, delimiter=',', quoting=csv.QUOTE_NONE)
    for row in reader:
        print(f'{row[0]:<15}  {row[1]:<20} {row[2]:<15} {row[3]:<15}')