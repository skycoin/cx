from datetime import datetime

def ack1(M, N):
    return (N + 1) if M == 0 else (
        ack1(M-1, 1) if N == 0 else ack1(M-1, ack1(M, N-1)))


def test():
    start = datetime.now()
    print(0, 5)
    print(f'Result: {ack1(0, 5)}')
    end = datetime.now()
    delta = end - start

    print(f'{delta.total_seconds():.3f}s')


test()