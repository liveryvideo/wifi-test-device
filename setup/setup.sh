#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

apt update
apt upgrade -y

apt install -y git dhcpcd

# The software managing traffic control runs on golang.
wget https://go.dev/dl/go1.20.2.linux-arm64.tar.gz

# Unpack package to go's directory.
echo "Unpacking go.."
tar -C /usr/local -xzf go1.20.2.linux-arm64.tar.gz

# Remove package after unzipping.
echo "Removing package.."
rm go1.20.2.linux-arm64.tar.gz

echo "Adding go bin to PATH.."
echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile.d/gopath.sh
echo "export GOPATH=/home/pi/go" >> /etc/profile.d/gopath.sh
source /etc/profile.d/gopath.sh

# In order to work as an access point, the Raspberry Pi needs to have the hostapd access point software package installed.
apt install -y hostapd

# Enable the wireless access point service and set it to start when your Raspberry Pi boots.
systemctl unmask hostapd
systemctl enable hostapd

# In order to provide network management services (DNS, DHCP) to wireless clients, the Raspberry Pi needs to have the dnsmasq software package installed.
apt install -y dnsmasq

# Finally, install netfilter-persistent and its plugin iptables-persistent. This utilty helps by saving firewall rules and restoring them when the Raspberry Pi boots.
DEBIAN_FRONTEND=noninteractive apt install -y netfilter-persistent iptables-persistent

# Copy over dhcpcd.conf
cp dhcpcd.conf /etc/dhcpcd.conf

# Copy over routed-ap.conf
cp routed-ap.conf /etc/sysctl.d/routed-ap.conf

# Copy over dnsmasq.conf
cp dnsmasq.conf /etc/dnsmasq.conf

# Ensure WiFi radio isn't blocked.
rfkill unblock wlan

# Copy over hostapd.conf.
cp hostapd.conf /etc/hostapd/hostapd.conf
# hostapd.conf override to let it wait for the wifi adapter to be fully started
mkdir -p /etc/systemd/system/hostapd.service.d
cp hostapd.service.d.override /etc/systemd/system/hostapd.service.d/override.conf

# Install iproute2 for traffic control.
apt install -y iproute2

# Disable WiFi power management
cp wifipwr.service /etc/systemd/system/
systemctl start wifipwr
systemctl enable wifipwr

echo "Configuring service.."
cp testdevice.service /etc/systemd/system/testdevice.service

cd ..

echo "Copying to src.."
mkdir -p /home/pi/go/src/bitbucket.org/exmachina/wifi-test-device
cp -r . /home/pi/go/src/bitbucket.org/exmachina/wifi-test-device

cd /home/pi/go/src/bitbucket.org/exmachina/wifi-test-device

echo "Getting dependencies.."
go get

echo "Building main.go.."
go build -o wifi-test-device

echo "Starting service.."
systemctl start testdevice
systemctl enable testdevice

echo "Adding firewall rules.."
# Add firewall rule.
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
# Save firewall rule.
netfilter-persistent save

echo "Reloading daemon.."
systemctl daemon-reload

echo " "
echo "Installation complete, please check for any errors and resolve them."
echo "If there are no errors, please reboot the system and all should be setup to go."
echo " "
echo "Default network name: Test-Device-01"
echo "Default password: LiveryTest"
echo "These can be changed in: '/etc/hostapd/hostapd.conf'"
