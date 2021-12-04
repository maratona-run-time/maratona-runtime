#!/bin/bash
docker build . -t test 
docker run test
