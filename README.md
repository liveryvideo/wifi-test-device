# Information
The wifi test device is a tool used to simulate networking problems such a packet loss, duplication or corruption.
Combined with the ability to limit bandwidth and increase latency.

This allows us to test if the livery player can recover from the myriad of problems that can occure while watching a stream.


# Specifications
This tool is designed to work on a raspberry pi, with access to at least etho0 and wlan0.
I.E: it needs to accept an ethernet input for internet access and a wifi chip to send out a wifi network for devices to connect to.

The Front-end runs on plain html, css and javascript.
The Front-end communicates to the Back-end through a REST api.
The Back-end runs a go http server. This server serves both the static Front-end files and manages the REST api.


# Manual Installation
Please follow the instruction below to setup the pi to be able to run GO applications and to serve as an access point.

Golang installation tutorial:               https://golang.org/doc/install
Raspberry pi access point tutorial:         https://www.raspberrypi.org/documentation/configuration/wireless/access-point-routed.md


Next, check if `tc` (traffic control) is installed by typing the command `tc` in a terminal.
If the command is unknown install it through:
`apt-get install iproute` or if prompted; `apt-get install iproute2`

Traffic control is what allows us to set bandwidth, latency and package manipulation on the network.

After having successfully setup the pi through the tutorials above, you can now `go run main.go` in the project. (With super user access)

Additionally you can `go build main.go` and natively execute the resulting image.

The server should now be running and can be accessed on your `localhost` on port 80.

# Automatic installation
Make sure you have rebooted the system after the first kernel update to ensure ip-tables functions correctly.

Install the nececary packages:

    sudo apt install git

Clone the repository.

Within the setup directory is a bash script called `setup.sh`.
Run this using `sudo ./setup.sh` to install and configure the WiFi access point.

Besides installing the necessary programs this script overwrites the following files:
`/etc/dnsmasq.conf`
`/etc/dhcpcd.conf`
`/etc/sysctl.d/routed-ap.conf`
`/etc/hostapd/hostapd.conf`

After the installation has completed, reboot the system and the network should show up with the default network name; unless changed in `/etc/hostapd/hostapd.conf`.


# Troubleshooting
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


If ./setup.sh is an unknown command use these modifiers to fix this.
`chmod 777 setup.sh`