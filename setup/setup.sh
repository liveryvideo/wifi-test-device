#!/bin/bash

# setup.sh
#
# Usage:
#   ./setup.sh [wifi_mode]
#
# Example:
#   ./setup.sh embedded   # uses embedded Wi-Fi chip
#   ./setup.sh external   # uses external USB Wi-Fi chip


if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

if [ $# -lt 1 ]; then
  echo "Usage: $0 [external|embedded]"
  exit 1
fi
WIFI_MODE="$1"
echo "Setting up the environment for '$WIFI_MODE' Wi-Fi"

if [ "$WIFI_MODE" = "external" ]; then
   echo "dtoverlay=disable-wifi" | tee -a /boot/firmware/config.txt
fi

# The software managing traffic control runs on golang.
wget https://go.dev/dl/go1.24.2.linux-arm64.tar.gz

# Unpack package to go's directory.
echo "Unpacking go.."
tar -C /usr/local -xzf go1.24.2.linux-arm64.tar.gz

# Remove package after unzipping.
echo "Removing package.."
rm go1.24.2.linux-arm64.tar.gz

echo "Adding go bin to PATH.."
echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile.d/gopath.sh
echo "export GOPATH=/home/pi/go" >> /etc/profile.d/gopath.sh
source /etc/profile.d/gopath.sh

# Use NetworkManager to setup an IP address for the wlan interface and configure the accesspoint
nmcli c add type wifi ifname wlan0 con-name rpi-ap autoconnect yes ssid "Test-Device-01"
nmcli c modify rpi-ap 802-11-wireless.mode ap \
                    802-11-wireless.band a \
                    802-11-wireless.channel 36 \
                    802-11-wireless-security.key-mgmt wpa-psk \
                    802-11-wireless-security.psk "LiveryTest"

nmcli c modify rpi-ap ipv4.addresses 192.168.4.1/24
nmcli c modify rpi-ap ipv4.method manual

# Disable ipv6 for now
nmcli c modify rpi-ap ipv6.method ignore
# Enable DHCP for the wifi clients
nmcli c modify rpi-ap ipv4.method shared

nmcli c up rpi-ap

# Copy over routed-ap.conf to enable IP forwarding
cp routed-ap.conf /etc/sysctl.d/routed-ap.conf


# Ensure WiFi radio isn't blocked.
rfkill unblock wlan

# Configure Wifi country
raspi-config nonint do_wifi_country NL

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

echo "Reloading daemon.."
systemctl daemon-reload

echo " "
echo "Installation complete, please check for any errors and resolve them."
echo "If there are no errors, please reboot the system and all should be setup to go."
echo " "
echo "Default network name: Test-Device-01"
echo "Default password: LiveryTest"
# echo "These can be changed in: '/etc/hostapd/hostapd.conf'"
