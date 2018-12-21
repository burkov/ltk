#!/bin/bash
set -e

echo ":: Installing apt-upgrade service"

cp systemd/apt-upgrade/apt-upgrade.service /etc/systemd/system
cp systemd/apt-upgrade/apt-upgrade.timer   /etc/systemd/system
cp systemd/apt-upgrade/apt-upgrade-service /usr/local/bin

chmod +x /usr/local/bin/apt-upgrade-service

systemctl enable apt-upgrade.timer
systemctl start apt-upgrade.timer
systemctl disable apt-daily-upgrade.timer
systemctl stop apt-daily-upgrade.timer
systemctl disable apt-daily.timer
systemctl stop apt-daily.timer

################

echo ":: Installing off-on service"

cp systemd/off-on/off-on.service /etc/systemd/system
cp systemd/off-on/off-on.timer   /etc/systemd/system
cp systemd/off-on/off-on-service /usr/local/bin

chmod +x /usr/local/bin/off-on-service

systemctl enable off-on.timer
systemctl start off-on.timer