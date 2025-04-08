# Wi-Fi test device

The wifi test device is a tool used to simulate networking problems such a packet loss, duplication or corruption.
Combined with the ability to limit bandwidth and increase latency.

This allows us to test if the livery player can recover from the myriad of problems that can occur while watching a stream.

## Specifications

This tool is designed to work on a raspberry pi, with access to at least eth0 and wlan0.
I.E: it needs to accept an ethernet input for internet access and a wifi chip to send out a wifi network for devices to connect to.

The Front-end runs on plain html, css and javascript.
The Front-end communicates to the Back-end through a REST api.
The Back-end runs a go http server. This server serves both the static Front-end files and manages the REST api.

## Wi-Fi performance

The onboard wifi chip of the RapsberryPi is pretty bad in terms of range and speed.
Using a USB wifi dongle can help a lot.
This project expects a wifi dongle taht can do 5GHz, you need to change some of the config files if you want to use a 2.4GHz dongle.

If you want to use an external dongle, you can run the `setup.sh` script with the `external` parameter.
This will disable the internal wifi chip and use the external USB wifi adapter.

## Installation

Follow these steps to make sure the wifi-test-device is setup correctly.
This assumes you have installed `Raspberry PI OS Lite (64-bit)` on the Raspberry PI already and you are able to log into it via SSH.
We will overwrite files and install new packages so please don't run this on a dedicated Raspberry PI as you might loose previously made configuration changes.

- Run `sudo apt update`
- Run `sudo apt upgrade`
- Reboot the system (important!)
- Run `sudo apt install git`
- Run `git clone https://github.com/liveryvideo/wifi-test-device.git`
- Run `cd wifi-test-device/setup`
- Run `sudo ./setup.sh internal` as root
- Check the logs for any errors.
- Reboot

`setup.sh` Will setup your raspberry pi to act as a wifi access-point;
and a service which serves you a control panel for controlling and logging the network.

Besides installing the necessary programs this script overwrites the following files:
`/etc/dnsmasq.conf`
`/etc/dhcpcd.conf`
`/etc/sysctl.d/routed-ap.conf`
`/etc/hostapd/hostapd.conf`

After the installation has completed, reboot the system and the network should show up with the default network name; unless changed it in `/etc/hostapd/hostapd.conf`.

## Troubleshooting

If the server returns a 404 on the homepage but the api is accessible.
I.E: `localhost` returns 404 but `localhost/api/status` return the expected values.
This is cause by an invalid working directory. Make sure the directory in the terminal is set to the root of the application when launching it.
This usually happens when launching with sudo, as this references from your home directory.


If your wifi access-point isn't visible on your other devices, make sure you have set a password that is longer than 6 characters.
This error can only be found in the logs of hostapd. Also make sure you have rebooted either the service or the pi itself after you make changes to hostapd.


If you are able to connect to your wifi access-point but there is no internet connection there might be an issue with your ip-tables.
Run these two commands as sudo:
`iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE`
`netfilter-persistent save`
This usually happens if you didn't reboot the system after the first system update.

Alternatively a more hacky solution is to run `setup.sh` (again).

If ./setup.sh is an unknown command use these modifiers to fix this.
`chmod 777 setup.sh`
