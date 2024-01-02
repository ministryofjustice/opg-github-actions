import random
import string

def rand(n:int):
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=n))
