#!/bin/bash

BUCKET_NUM=$1
if [ -z $BUCKET_NUM ]; then
    echo "Number of buckets not specified, defaults to 3"
    BUCKET_NUM=3
fi

# prepare the docker compose file, it will be started in the last step
cp docker-compose.template docker-compose.yml

# read the services template
TEMPLATE=$(cat services.template)

# for each requested bucket, create the service definition in docker-compose file
for i in $(seq 1 $BUCKET_NUM)
do
    BUCKET_IDX=$(expr $i - 1)
    
    # prepare the services template for this bucket index
    SERVICES="$(echo "${TEMPLATE}" | sed "s/#/$BUCKET_IDX/g")"
    echo "${SERVICES}"

    echo -e "${SERVICES}" >> docker-compose.yml
    echo -e "\n" >> docker-compose.yml

done

cat docker-compose.yml