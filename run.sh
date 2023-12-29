#!/usr/bin/env bash

echo "Migrate"
./mongo-streamer migrate up

echo "Start server..."
./mongo-streamer
