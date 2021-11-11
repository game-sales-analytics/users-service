#!/bin/sh

set -ex

docker-compose --project-name isac --file ./deploy/compose.yml down --rmi local
