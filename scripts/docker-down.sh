#!/bin/sh

set -ex

docker-compose --project-name gsa --file ./deploy/compose.yml down --rmi local
