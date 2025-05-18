#!/bin/bash
# Simple script to run the Beo Echo Go backend

echo "Starting BeoEcho Go backend..."
cd "$(dirname "$0")"
go run main.go "$@"
