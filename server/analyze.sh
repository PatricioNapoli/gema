#!/bin/zsh

sonar-scanner \
  -Dsonar.projectKey=gema \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://sonar.hq.localhost \
  -Dsonar.login=$1 \
