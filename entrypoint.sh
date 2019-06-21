#!/bin/bash

ssh-keyscan -H github.com >> /root/.ssh/known_hosts > /dev/null 2>&1
exec "$(ssh-agent)"
ssh-add /root/.ssh/id_rsa

./ad server --config .ad.yaml

