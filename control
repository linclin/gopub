#!/bin/bash -e
export GOROOT=/usr/local/go
export PATH=$PATH:/usr/local/go/bin/ 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASENAME=`basename $DIR`
export GOPATH=$DIR
export GOBIN=$DIR/bin/ 

app=$BASENAME
conf=$DIR/src/conf/app.conf
pidfile=$DIR/$BASENAME.pid
logfile=$DIR/$BASENAME.log

 function check_pid() {
    if [ -f $pidfile ];then
        pid=`cat $pidfile`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function build() {
    gofmt -w $DIR/src/
    cd $DIR/src/ 
    go build -o $BASENAME
    if [ $? -ne 0 ]; then
        exit $?
    fi
}
function pack() {
    build
    cd  ..
    rm  -rf $BASENAME/src/logs/*
    cd  .. && tar zcvf $app.tar.gz $BASENAME/control $BASENAME/src/$app  $BASENAME/src/conf    $BASENAME/src/logs   $BASENAME/src/agent  $BASENAME/src/views $BASENAME/src/static  $BASENAME/src/favicon.ico
}
function start() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "$app now is running already, pid="
        cat $pidfile
        return 1
    fi

    if ! [ -f $conf ];then
        echo "Config file $conf doesn't exist, creating one." 
    fi
    cd $DIR/src/ 
    nohup  ./$BASENAME  >$logfile 2>&1 &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > $pidfile
        echo "$app started..., pid=$!"
    else
        echo "$app failed to start."
        return 1
    fi


}
function killall() {
    pid=`cat $pidfile`
    ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'|xargs kill -9 
    rm -f $pidfile
    echo "$app killed..., pid=$pid"
}
function stop() {
    #ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'|xargs kill -9
    pid=`cat $pidfile`
    kill $pid
    rm -f $pidfile
    echo "$app stoped..., pid=$pid"
}
function restart() {
    stop
    sleep 1
    start 
}
function reload() { 
    pid=`cat $pidfile`
    kill -HUP $pid
    sleep 1
    newpid=`ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'` 
    echo "$app reload..., pid=$newpid"
    echo $newpid > $pidfile 
}
function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo started
    else
        echo stoped
    fi
}
function run() {
   cd $DIR/src/
   ./$BASENAME -docker
   #go run main.go 
}
function rundocker() {
   cd $DIR/src/
   ./$BASENAME -docker
   #go run main.go
}
function init() {
   cd $DIR/src/
   ./$BASENAME -syncdb
   #go run main.go
}
function beerun() {
   cd $DIR/src/
   bee run 
}

function tailf() {
   tail -f $logfile
}
function docs() {
   cd $DIR/src/ 
   bee generate docs 
}

function sslkey() {
   cd $DIR/src/conf/ssl
   ###CA:
   #私钥文件
   openssl genrsa -out ca.key 2048
}  

 
function help() {
    echo "$0 build|start|stop|kill|restart|reload|run|rundocker|init|tail|docs|pack|beerun|sslkey"
}
if [ "$1" == "" ]; then
    help
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "pack" ];then
    pack
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "kill" ];then
    killall
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "reload" ];then
    reload
elif [ "$1" == "status" ];then
    status
elif [ "$1" == "run" ];then
    run
elif [ "$1" == "rundocker" ];then
    rundocker
elif [ "$1" == "init" ];then
    init
elif [ "$1" == "beerun" ];then
    beerun 
elif [ "$1" == "tail" ];then
    tailf
elif [ "$1" == "docs" ];then
    docs 
elif [ "$1" == "sslkey" ];then
    sslkey
else
    help
fi
