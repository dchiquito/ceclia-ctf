#!/bin/sh

INTERNAL_IP=192.168.0.152
EXTERNAL_IP=98.26.35.179
PASSWORD=$1

sshpass -p $PASSWORD ssh admin@$INTERNAL_IP "rm -rf /usr/local/ceclia-ctf/ceclia-ctf-go-pi"
sshpass -p $PASSWORD scp build/phase3/ceclia-ctf-go-pi admin@$INTERNAL_IP:/usr/local/ceclia-ctf/ceclia-ctf-go-pi
