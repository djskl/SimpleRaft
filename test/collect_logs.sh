#!/usr/bin/env bash

servers=("10.0.138.151" "10.0.138.152" "10.0.138.153" "10.0.138.155" "10.0.138.158")

for idx in {0..4}
do
    ip=${servers[$idx]}
    scp ${ip}:/raftserver/db.dat /tmp/raft/${idx}.dat
done