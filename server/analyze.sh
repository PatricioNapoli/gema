#!/bin/zsh

sonar-scanner \
  -Dsonar.projectKey=gema \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://localhost:9005 \
  -Dsonar.login=$SONAR_KEY