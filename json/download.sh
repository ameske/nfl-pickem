#!/bin/bash

for i in `seq 1 17`;
do
  nfl schedule download -y 2017 -w ${i} > schedule-2017-wk${i}.json
done
