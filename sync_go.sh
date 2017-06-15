#!/usr/bin/env bash

servers=("10.0.138.151" "10.0.138.152" "10.0.138.153" "10.0.138.155")

for ip in ${servers[@]}
do
    rsync -avP /root/gcodes/src/SimpleRaft ${ip}:/opt/golang/src/
    #ssh root@${ip} "mkdir /opt/golang/src"
    #ssh root@${ip} "mkdir /opt/golang/bin"
    #ssh root@${ip} "mkdir /opt/golang/pkg"
done