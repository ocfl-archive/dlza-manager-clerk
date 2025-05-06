#!/bin/bash
docker build --rm -f ./Dockerfile -t gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-miniresolver/dlza-manager-clerk:latest .
docker tag gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-miniresolver/dlza-manager-clerk:latest registry.localhost:5001/gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-miniresolver/dlza-manager-clerk
docker push registry.localhost:5001/gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-miniresolver/dlza-manager-clerk
