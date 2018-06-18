import time


def ack1(M, N):
    return (N + 1) if M == 0 else (
        ack1(M-1, 1) if N == 0 else ack1(M-1, ack1(M, N-1)))


def test():
    start = time.time()
    print "(0, 5)"
    print "Result: {value}".format(value=ack1(0, 5))

    end = time.time()
    print (end - start)


test()
