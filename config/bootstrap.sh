#!/bin/bash -e

sudo add-apt-repository ppa:sergio-br2/vbam-trunk
sudo apt-get update
sudo apt-get install vbam-gtk

sudo apt-get install -y golang
sudo apt-get install -y xvfb
sudo apt-get install -y avconv
sudo apt-get install -y libav-tools
sudo apt-get install -y alsa-utils
sudo apt-get install -y xdotool

