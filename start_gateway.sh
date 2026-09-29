#!/bin/bash
set -e

echo "[+] Starting Hotspot on wlan0..."
nmcli device wifi hotspot ifname wlan0 ssid GatewayTestbed password "GatewayPassword123"

echo "[+] Bringing up WireGuard tunnel..."
sudo wg-quick up wg0

echo "[+] Applying routing and iptables rules..."
sudo sysctl -w net.ipv4.ip_forward=1
sudo iptables -A FORWARD -i wlan0 -o wg0 -j ACCEPT
sudo iptables -A FORWARD -i wg0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
sudo iptables -t nat -A POSTROUTING -o wg0 -j MASQUERADE
sudo iptables -t mangle -A FORWARD -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu

echo "[✓] Testbed network ready. Run ./run.sh to start capturing."
