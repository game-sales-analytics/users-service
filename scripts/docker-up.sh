#!/bin/sh

set -ex

DOCKER_BUILDKIT=1 docker-compose --project-name gsa --file ./deploy/compose.yml up --build
