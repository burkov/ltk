#!/bin/bash
set -e

echo ":: Installing apt-upgrade service"

cp systemd/apt-upgrade/apt-upgrade.service /etc/systemd/system
cp systemd/apt-upgrade/apt-upgrade.timer   /etc/systemd/system
cp systemd/apt-upgrade/apt-upgrade-service /usr/local/bin

systemctl enable apt-upgrade.timer
systemctl start apt-upgrade.timer
systemctl disable apt-daily-upgrade.timer
systemctl stop apt-daily-upgrade.timer
systemctl disable apt-daily.timer
systemctl stop apt-daily.timer
