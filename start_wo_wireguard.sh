#!/bin/bash
set -e

WLAN_IFACE="wlan0"

# 1. Detect active Ethernet / WAN interface
ETH_IFACE=$(ip route show default | awk '{print $5}' | head -n1)

if [ -z "$ETH_IFACE" ]; then
    echo "[!] Error: No default gateway/interface found. Are you connected via Ethernet?"
    exit 1
fi

echo "[*] Active WAN interface detected: $ETH_IFACE"

# 2. Stop firewalld to prevent DNS/DHCP drops
echo "[+] Stopping firewalld..."
sudo systemctl stop firewalld 2>/dev/null || true

# 3. Enable kernel IP forwarding (both IPv4 and IPv6)
echo "[+] Enabling kernel packet forwarding for IPv4 and IPv6..."
sudo sysctl -w net.ipv4.ip_forward=1 >/dev/null
sudo sysctl -w net.ipv6.conf.all.forwarding=1 >/dev/null

# 4. Start Hotspot on wlan0
echo "[+] Starting Hotspot on $WLAN_IFACE..."
# Disconnect previous session cleanly first
nmcli device disconnect "$WLAN_IFACE" 2>/dev/null || true
nmcli device wifi hotspot ifname "$WLAN_IFACE" ssid "GatewayTestbed" password "GatewayPassword123"

# 5. Apply forwarding and NAT masquerading
echo "[+] Setting up iptables NAT and forwarding rules..."
# Remove potential duplicate rules before appending
sudo iptables -D FORWARD -i "$WLAN_IFACE" -o "$ETH_IFACE" -j ACCEPT 2>/dev/null || true
sudo iptables -D FORWARD -i "$ETH_IFACE" -o "$WLAN_IFACE" -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT 2>/dev/null || true
sudo iptables -t nat -D POSTROUTING -o "$ETH_IFACE" -j MASQUERADE 2>/dev/null || true

# Append fresh rules
sudo iptables -A FORWARD -i "$WLAN_IFACE" -o "$ETH_IFACE" -j ACCEPT
sudo iptables -A FORWARD -i "$ETH_IFACE" -o "$WLAN_IFACE" -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -t nat -A POSTROUTING -o "$ETH_IFACE" -j MASQUERADE

echo "[✓] Hotspot active! Connected devices on $WLAN_IFACE route through $ETH_IFACE."
echo "[*] You can now start the Go engine and Python QoS controller."