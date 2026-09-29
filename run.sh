#!/bin/bash

# 1. Open Python listener first so the ZeroMQ port is ready
gnome-terminal -- bash -c "python3 script.py; exec bash"

# 2. Brief sleep to allow Python socket binding
sleep 1

# 3. Run Go capture engine with sudo for raw packet socket privileges
gnome-terminal -- bash -c "sudo go run cmd/main/main.go; exec bash"