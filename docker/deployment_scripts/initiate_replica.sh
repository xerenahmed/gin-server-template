#!/bin/bash

until mongo --host mongo
do
    sleep 2
done

mongo --host mongo <<EOF
rs.initiate(
  {
    _id : 'rs0',
    members: [
      { _id : 0, host : "mongo:27017" },
      { _id : 1, host : "mongodb1:27017" },
      { _id : 2, host : "mongodb2:27017" }
    ]
  }
)
EOF

