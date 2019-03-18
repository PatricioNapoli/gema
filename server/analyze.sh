#!/bin/zsh

sonar-scanner \
  -Dsonar.projectKey=gema \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://localhost:9005 \
  -Dsonar.login=b82a411c7523afb58abc2eb5a0ef71248ef5262b