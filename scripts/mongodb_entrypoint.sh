#!/usr/bin/env bash
set -euo pipefail

MONGODB_USER="${MONGODB_USER?Please define MONGODB_USER}"
MONGODB_PASS="${MONGODB_PASS?Please define MONGODB_PASS}"
MONGODB_DB="${MONGODB_DB?Please define MONGODB_DB}"

# Create an ephemeral database directory if not mounted externally.
[ ! -d /data/db ] && mkdir -p /data/db

# Start `mongod` unauth'ed.
mongod --fork --logpath /tmp/mongo.log

# Wait for mongo to start
attempts=0
while [ "$attempts" -lt 60 ]
do
  nc -w 1 -z localhost 27017 && break
  attempts=$((attempts+1))
  >&2 echo "INFO: [${attempts}/60] waiting for mongodb to start..."
  sleep 0.25
done
# Create the MongoDB user in the MONGODB_DB database.
if ! mongo --quiet --eval 'db.getUsers()' "$MONGODB_DB"
then mongo --eval 'db.createUser({user: "'"$MONGODB_USER"'", pwd: "'"$MONGODB_PASS"'", roles: [{ role: "dbOwner", db: "'"$MONGODB_DB"'" }]})' "$MONGODB_DB"
fi

# Restart `mongod` with auth and logs forwarded to stdout
killall mongod
# Wait for mongo to stop
attempts=0
while [ "$attempts" -lt 60 ]
do
  nc -w 1 -z localhost 27017 || break
  attempts=$((attempts+1))
  >&2 echo "INFO: [${attempts}/60] waiting for mongodb to stop..."
  sleep 0.25
done
mongod --auth --bind_ip_all
