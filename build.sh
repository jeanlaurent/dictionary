#!/bin/bash

GOOS=linux go build
docker build -ti dico
