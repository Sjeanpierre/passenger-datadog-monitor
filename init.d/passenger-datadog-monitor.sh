#!/bin/env sh
#
# passenger-datadog-monitor

# Source function library.
. /etc/init.d/functions

RETVAL=0
prog="passenger-datadog-monitor"
PIDFILE=/var/run/$prog.pid


start() {
        printf "%s\n" "Starting $prog: "
        PID=`nohup /usr/bin/passenger-datadog-monitor > /tmp/passenger-datadog-monitor.log 2>&1 & echo $!`
        RETVAL=$?
        if [ $RETVAL -eq 0 ]; then
          echo "$PID" > $PIDFILE
          printf "%s\n" "Ok"
        fi
        return $RETVAL
}

stop() {
        echo -n "Shutting down $prog: "
        PID=`cat $PIDFILE`
        if [ -f $PIDFILE ]; then
            kill -9 $PID
            printf "%s\n" "Ok"
            rm -f $PIDFILE
        else
            printf "%s\n" "pidfile not found"
        fi
}

status() {
        echo -n "Checking $prog status: "
        if [ -f $PIDFILE ]; then
            PID=`cat $PIDFILE`
            if [ -z "`ps axf | grep ${PID} | grep -v grep`" ]; then
                printf "%s\n" $prog" dead but pidfile exists"
            else
                echo $prog" Running"
            fi
        else
            printf "%s\n" $prog" not running"
        fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: $prog {start|stop|status|restart}"
        exit 1
        ;;
esac
exit $RETVAL
