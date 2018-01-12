#!/usr/bin/env bash
cd bin
for((i=1;i<=2;i++));do

httpPort=$(expr $i \+ 8000);
rpcPort=$(expr $i \+ 9000);

rm -rf config && ./riff start -name node$i -rpc :$rpcPort -http :$httpPort -join 192.168.3.2:9001,192.168.3.2:9002,192.168.3.2:9003 &
#sleep 1

done


#rm -rf config && go run *.go start -name node5 -rpc :8634 -http :8614 -join 192.168.1.220:8630,192.168.1.220:8631,192.168.1.220:8632
