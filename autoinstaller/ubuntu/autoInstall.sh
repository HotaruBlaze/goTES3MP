#!/bin/bash

# Prereq
echo -e "\e[32m[goTES3MP Installer]: Getting Prerequisites\e[0m"
sudo apt update && sudo apt install jq curl -y
mkdir -p $HOME/goTES3MP && cd $HOME/goTES3MP/
cd $HOME/goTES3MP
printf '%s' 'ergo,tmp' | xargs -d, mkdir

# Install goTES3MP: Binary
# echo "Done Install goTES3MP and ergo IRC"
echo -e "\e[32m[goTES3MP Installer]: Installing goTES3MP\e[0m"
Latest_goTES3MP_Version=$(curl https://api.github.com/repos/HotaruBlaze/goTES3MP/releases -s | jq -r .[].tag_name | grep '^v[0-9]\.[0-9]*\.[0-9]*$' | sort -nr | head -n1)
Latest_goTES3MP_DownloadURL=$(curl -sL https://api.github.com/repos/HotaruBlaze/goTES3MP/releases/tags/$Latest_goTES3MP_Version | jq -r ".assets[] | select(.name | contains(\"goTES3MP-Linux\")) | .browser_download_url")
curl -o $HOME/goTES3MP/goTES3MP-Linux -J -L $Latest_goTES3MP_DownloadURL
chmod +x $HOME/goTES3MP/goTES3MP-Linux
if [ ! -f $HOME/goTES3MP/config.yaml ]
then
    cd $HOME/goTES3MP && ./goTES3MP-Linux
    sed -i -E "s/(\enableinteractiveconsole: true\b)/enableinteractiveconsole: false/" config.yaml
fi

# install goTES3MP: Service
echo -e "\e[32m[goTES3MP Installer]: Installing goTES3MP service\e[0m"
curl -o $HOME/goTES3MP/tmp/gotes3mp.service -J -L "https://raw.githubusercontent.com/HotaruBlaze/goTES3MP/main/autoinstaller/ubuntu/gotes3mp.service"
sed -i -E "s/(\bUSERNAME\b)/$USER/" $HOME/goTES3MP/tmp/gotes3mp.service
sudo mv $HOME/goTES3MP/tmp/gotes3mp.service /etc/systemd/system/gotes3mp.service

# install ergo: Binary
echo -e "\e[32m[goTES3MP Installer]: Installing ergo IRC Server\e[0m"
curl -o $HOME/goTES3MP/ergo/ergo-2.11.1-ergo-tes3mp-linux-x86_64.tar.gz -J -L "https://github.com/HotaruBlaze/ergo-tes3mp/releases/download/v2.11.1/ergo-2.11.1-ergo-tes3mp-linux-x86_64.tar.gz"
cd $HOME/goTES3MP/ergo
tar -xzf *.tar.gz
mv ergo-2.11.1-ergo-tes3mp-linux-x86_64/* .
rm *.tar.gz
rm -Rf ergo-2.11.1-ergo-tes3mp-linux-x86_64/
chmod +x ./ergo
cp default.yaml ircd.yaml
./ergo mkcerts

# install ergo: Service
echo -e "\e[32m[goTES3MP Installer]: Installing ergo service file\e[0m"
curl -o $HOME/goTES3MP/tmp/gotes3mp-ergo.service -J -L "https://raw.githubusercontent.com/HotaruBlaze/goTES3MP/main/autoinstaller/ubuntu/gotes3mp-ergo.service"
sed -i -E "s/(\bUSERNAME\b)/$USER/" $HOME/goTES3MP/tmp/gotes3mp-ergo.service
sudo mv $HOME/goTES3MP/tmp/gotes3mp-ergo.service /etc/systemd/system/gotes3mp-ergo.service

# Enable Services and turn them on.
echo -e "\e[32m[goTES3MP Installer]: Enabling ergo and goTES3MP service files\e[0m"
printf '%s' 'gotes3mp,gotes3mp-ergo' | xargs -d, sudo systemctl enable 
printf '%s' 'gotes3mp,gotes3mp-ergo' | xargs -d, sudo systemctl start 
sleep 1
printf '%s' 'gotes3mp,gotes3mp-ergo' | xargs -td, sudo systemctl is-active
echo -e "\e[32m[goTES3MP Installer]: Done Install goTES3MP and ergo IRC\e[0m"
