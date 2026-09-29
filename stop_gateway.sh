#!/bin/bash

echo "[-] Bringing down WireGuard..."
sudo wg-quick down wg0 2>/dev/null || true

echo "[-] Stopping Hotspot..."
nmcli connection down Hotspot 2>/dev/null || true

echo "[-] Clearing forwarding rules..."
sudo iptables -D FORWARD -i wlan0 -o wg0 -j ACCEPT 2>/dev/null || true
sudo iptables -D FORWARD -i wg0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT 2>/dev/null || true
sudo iptables -t nat -D POSTROUTING -o wg0 -j MASQUERADE 2>/dev/null || true
sudo iptables -t mangle -D FORWARD -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu 2>/dev/null || true

echo "[✓] Network restored to normal."
