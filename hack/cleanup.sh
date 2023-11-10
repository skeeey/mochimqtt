# !/bin/bash

ps ax | grep "mochimqtt/bin/client" | awk '{print $1}' | xargs kill
ps ax | grep "mqtt-client/bin/client" | awk '{print $1}' | xargs kill
