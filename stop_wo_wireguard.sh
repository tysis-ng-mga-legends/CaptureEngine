#!/bin/bash

WLAN_IFACE="wlan0"
ETH_IFACE=$(ip route show default | awk '{print $5}' | head -n1)

echo "[-] Stopping Hotspot on $WLAN_IFACE..."
nmcli device disconnect "$WLAN_IFACE" 2>/dev/null || true
nmcli connection down "Hotspot" 2>/dev/null || true
nmcli connection down "GatewayTestbed" 2>/dev/null || true

if [ -n "$ETH_IFACE" ]; then
    echo "[-] Clearing iptables rules for $ETH_IFACE and $WLAN_IFACE..."
    sudo iptables -D FORWARD -i "$WLAN_IFACE" -o "$ETH_IFACE" -j ACCEPT 2>/dev/null || true
    sudo iptables -D FORWARD -i "$ETH_IFACE" -o "$WLAN_IFACE" -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT 2>/dev/null || true
    sudo iptables -t nat -D POSTROUTING -o "$ETH_IFACE" -j MASQUERADE 2>/dev/null || true
fi

echo "[-] Clearing any remaining QoS queue rules on $WLAN_IFACE..."
sudo tc qdisc del dev "$WLAN_IFACE" root 2>/dev/null || true

if [ -n "$ETH_IFACE" ]; then
    sudo tc qdisc del dev "$ETH_IFACE" root 2>/dev/null || true
fi

echo "[✓] Hotspot and routing rules stopped. Interfaces restored to defaults."