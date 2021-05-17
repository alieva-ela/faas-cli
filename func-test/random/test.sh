#!/bin/bash
s="time echo \"0\""
for i in {0..15};
do
   s+=" | faas-cli invoke random-python-$i"
done
#s+=" --gateway=$GATEWAY"
eval $s