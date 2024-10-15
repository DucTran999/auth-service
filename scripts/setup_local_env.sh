#!/bin/bash
yellow() {
   echo -e "\033[33m$1\033[0m"
}

ENV_FILE_PATH=".env.local"
if [ ! -e "$ENV_FILE_PATH" ]; then
    cp .env.example .env.local
fi

PKG_INFRA=./infra
read -p "Do you fill config in .env.local file? (y/n): " answer
if [[ "$answer" == "y" || "$answer" == "Y" ]]; then
    docker-compose -f $PKG_INFRA/docker-compose.yml --env-file .env.local up -d
else
    yellow "hint: Fill the .env.local then run 'make localenv' again."
fi