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

    flight_bookings = db.flight_bookings
    flight={"trip_id": r["trip_id"], 
              "depart": r["depart"],
              "depart_at": r["depart_at"],
              "arrive": r["arrive"],
              "arrive_at": r["arrive_at"]}
    res = flight_bookings.insert_one(flight)

    # followers.remove()
     # print("book-flight")


    # a = []

    # for post in hotel_bookings.find():
    #     a.append(post)

    return "Record inserted in flight_bookings: {}".format(res.inserted_id)

    