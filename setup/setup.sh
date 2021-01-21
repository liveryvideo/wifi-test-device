#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

apt update
apt upgrade

# The software managing traffic control runs on golang.
apt install golang-go

# In order to work as an access point, the Raspberry Pi needs to have the hostapd access point software package installed.
apt install hostapd

# Enable the wireless access point service and set it to start when your Raspberry Pi boots.
systemctl unmask hostapd
systemctl enable hostapd

# In order to provide network management services (DNS, DHCP) to wireless clients, the Raspberry Pi needs to have the dnsmasq software package installed.
apt install dnsmasq

# Finally, install netfilter-persistent and its plugin iptables-persistent. This utilty helps by saving firewall rules and restoring them when the Raspberry Pi boots.
DEBIAN_FRONTEND=noninteractive apt install -y netfilter-persistent iptables-persistent

# Copy over dhcpcd.conf
cp -i dhcpcd.conf /etc/dhcpcd.conf

# Copy over routed-ap.conf
cp -i routed-ap.conf /etc/sysctl.d/routed-ap.conf

# Add firewall rule.
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

# Save firewall rule.
netfilter-persistent save

# Copy over dnsmasq.conf
cp -i dnsmasq.conf /etc/dnsmasq.conf

# Ensure WiFi radio isn't blocked.
rfkill unblock wlan

# Copy over hostapd.conf.
cp -i hostapd.conf /etc/hostapd/hostapd.conf

# Install iproute2 for traffic control.
apt install iproute2

echo " "
echo "Installation complete, please check for any errors and resolve them."
echo "If there are no errors, please reboot the system and all should be setup to go."
echo " "
echo "Default network name: Test-Device-01"
echo "Default password: LiveryTest"
echo "These can be changed in: '/etc/hostapd/hostapd.conf'"
