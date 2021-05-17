from random import randint

def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """
    random_number = randint(-100000, 100000)
    #print("req:", req)
    #print("!!!!!!!!")
    #print("Result = ", req + str(random_number))
    res = int(req) + 1
    return str(res)
