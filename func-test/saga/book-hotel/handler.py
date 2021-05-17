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

    hotel_bookings = db.hotel_bookings
    hotel={"trip_id": r["trip_id"], 
              "hotel": r["hotel"],
              "check_in": r["check_in"],
              "check_out": r["check_out"]}
    res = hotel_bookings.insert_one(hotel)

    # followers.remove()
    # print("book-hotel")


    # a = []

    # for post in hotel_bookings.find():
    #     a.append(post)

    return "Record inserted in hotel_bookings: {}".format(res.inserted_id)