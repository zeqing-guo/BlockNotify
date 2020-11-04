#!/bin/bash

export $(less .env | xargs)
./bin/blocknotify