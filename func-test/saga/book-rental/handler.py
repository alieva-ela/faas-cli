import os
from pymongo import MongoClient
from urllib.parse import quote_plus
import json

def get_uri():
    password=""
    with open("/var/openfaas/secrets/mongo-db-password") as f:
        password = f.read()

    return "mongodb://%s:%s@%s" % (
    quote_plus("root"), quote_plus(password), os.getenv("mongo_host"))

def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """

    uri = get_uri()
    client = MongoClient(uri)

    db = client['openfaas']

    r = json.loads(req)

    rental_bookings = db.rental_bookings
    rental={"trip_id": r["trip_id"], 
              "rental": r["rental"],
              "rental_from": r["rental_from"],
              "rental_to": r["rental_to"]}
    res = rental_bookings.insert_one(rental)

    # followers.remove()
    # print("book-hotel")


    # a = []

    # for post in hotel_bookings.find():
    #     a.append(post)

    return "Record inserted in rental_bookings: {}".format(res.inserted_id)