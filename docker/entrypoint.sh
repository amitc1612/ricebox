#!/bin/bash
mkdir -p /tmp/runtime-rice
chmod 700 /tmp/runtime-rice

# Start VNC server for visual preview
wayvnc 0.0.0.0 5900 &
VNC_PID=$!

# Trap shutdown
trap "kill $VNC_PID; exit" SIGTERM SIGINT

# Launch Hyprland
Hyprland &
HYPR_PID=$!

echo "Sandbox ready! VNC on port 5900"
wait $HYPR_PID