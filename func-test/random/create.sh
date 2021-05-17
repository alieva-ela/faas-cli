#!/bin/bash
export OPENFAAS_URL=http://127.0.0.1:31112
max=50
printf '{
  "StartFunction": "random-python-0",
  "States":
  	{\n' > rand.json

for (( i = 0; i <= 49; i++ ))
do
name="random-python-$i"
faas-cli new --lang python $name
cp random-python/handler.py  $name
yml="$name.yml"
printf "version: 1.0
provider:
  name: openfaas
  gateway: http://127.0.0.1:31112
functions:
  $name:
    lang: python
    handler: ./$name
    image: 12111999/$name:latest\n" > $yml

faas up -f $yml --gateway=$GATEWAY

printf "	    \"random-python-$i\": {
	      \"Type\": \"Task\",
	      \"ResultPath\": \"random-python-$i\",
	      \"Next\": \"random-python-$(($i +1))\"
		},\n" >> rand.json
done

name="random-python-$max"
faas-cli new --lang python $name
cp random-python/handler.py  $name
yml="$name.yml"
printf "version: 1.0
provider:
  name: openfaas
  gateway: http://127.0.0.1:31112
functions:
  $name:
    lang: python
    handler: ./$name
    image: 12111999/$name:latest\n" > $yml
faas up -f $yml --gateway=$GATEWAY


printf "	    \"random-python-$max\": {
	      \"Type\": \"Task\",
	      \"ResultPath\": \"random-python-$max\",
	      \"End\": true
		}
		}}" >> rand.json

