#!/usr/bin/env python3
import random
import string

def rand(n:int):
    """Generate a random string of n length using lowercase and digits."""
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=n))
