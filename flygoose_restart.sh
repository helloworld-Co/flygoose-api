#!/bin/bash
source /etc/profile &> /dev/null
flag="flygoose"
flag2="admin"
port=29090
port2=29091

case $1 in 
   start)
        if [ "$2" == "online" -a "$3" == "$flag" ]; then
            su -s /bin/bash www -c "nohup /apps/$flag/$flag -env app-prod.yaml >> /tmp/$flag.log &" 
        elif [ "$2" == "test" -a "$3" == "$flag" ]; then
            su -s /bin/bash www -c "nohup /apps/$flag/$flag -env app-dev.yaml >> /tmp/$flag.log &" 
        elif [ "$2" == "online" -a "$3" == "$flag2" ]; then
            su -s /bin/bash www -c "nohup /apps/$flag2/$flag2 -env app-prod.yaml>> /tmp//$flag-$flag2.log &" 
        elif [ "$2" == "test" -a "$3" == "$flag2" ]; then
            su -s /bin/bash www -c "nohup /apps/$flag2/$flag2 -env app-dev.yaml>> /tmp/$flag-$flag2.log &" 
        fi
        sleep 3
        ;;
   stop)
        if [ "$3" == "$flag2" ]; then
            port=$port2
        fi
        #(ps aux | grep '/apps/top_server/top_server' | grep -v grep) && (ps aux | grep '/apps/top_server/top_server' | grep -v grep | awk '{print $2}' | xargs kill) 
        netstat -nulpt|grep -w ":$port"|grep -v grep|awk '{print $NF}'|awk -F'/' '{print $1}' |xargs kill
        sleep 3
        #(ps aux | grep '/apps/top_server/top_server' | grep -v grep) && (ps aux | grep '/apps/top_server/top_server' | grep -v grep | awk '{print $2}' | xargs kill -9) 
        (netstat -nulpt|grep -w ":$port"|grep -v grep|awk '{print $NF}'|awk -F'/' '{print $1}' ) && (netstat -nulpt|grep -w ":$port"|grep -v grep|awk '{print $NF}'|awk -F'/' '{print $1}' | xargs kill -9)
        echo 0
        ;;
esac
chmod 755 /tmp/*.log
