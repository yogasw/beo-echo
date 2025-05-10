#!/bin/bash
# Simple script to run the Mockoon Control Panel Go backend

echo "Starting Mockoon Control Panel Go backend..."
cd "$(dirname "$0")"
go run main.go "$@"
