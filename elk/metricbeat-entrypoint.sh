#!/bin/sh

metricbeat setup -e -strict.perms=false -E setup.kibana.host=kibana:5601
metricbeat -e -strict.perms=false