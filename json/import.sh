#!/bin/bash

for i in `seq 1 17`;
do
  nfl schedule import -f json/schedule-2017-wk${i}.json -d test.db
done
