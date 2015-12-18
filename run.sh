#!/bin/bash

export PRODUCTION=1

start() {
    cd /home/gyazo/git/Dyazo;
    rm main
    go build main.go
    exec nohup /home/gyazo/git/Dyazo/main > /tmp/gyazo.out 2>&1&
    echo $! > /home/gyazo/pids/dyazo.pid
    disown
}

stop() {
    kill `cat /home/gyazo/pids/dyazo.pid`
}


case $1 in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    stop
    start
    ;;
  *)
  echo "usage: run.sh {start|stop|restart}" ;;
esac
exit 0
