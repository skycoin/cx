from datetime import datetime
import csv

def digital_root (n):
    ap = 0
    n = abs(int(n))
    while n >= 10:
        n = sum(int(digit) for digit in str(n))
        ap += 1
    result = [ap, n]
    return result
 
def test():
    start = datetime.now()
    persistance, root = digital_root(79563)
    # print(f' 79563 has additive persistance {persistance} and digital root {root}.')
    end = datetime.now()
    delta = end - start
    return f'{delta.total_seconds():.3f}s'

test_result = test()

with open('results.csv', mode='w') as csv_file:
    fieldnames = ['Language', 'Test Name', 'Input', 'Time']
    writer = csv.DictWriter(csv_file, fieldnames=fieldnames)
    writer.writeheader()
    writer.writerow({'Language': 'python', 'Test Name': 'digital-root', 'Input': '79563','Time': test_result})

with open('results.csv','r', newline='') as file:
    reader = csv.reader(file, delimiter=',', quoting=csv.QUOTE_NONE)
    for row in reader:
        print(f'{row[0]:<15}  {row[1]:<15} {row[2]:<15} {row[3]:<15}')