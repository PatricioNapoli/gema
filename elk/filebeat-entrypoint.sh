#!/bin/sh

filebeat setup -e -strict.perms=false -E setup.kibana.host=kibana:5601
filebeat -e -strict.perms=false