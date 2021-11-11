#!/bin/sh

set -ex

DOCKER_BUILDKIT=1 docker-compose --project-name isac --file ./deploy/compose.yml up --build
