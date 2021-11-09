#!/bin/sh

# Install docker-compose
# https://docs.docker.com/compose/install/
echo "installing docker-compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version

# For bash completion.
# https://docs.docker.com/compose/completion/

echo "installing auto-completion for docker-compose..."
sudo curl \
    -L https://raw.githubusercontent.com/docker/compose/1.29.2/contrib/completion/bash/docker-compose \
    -o /etc/bash_completion.d/docker-compose

# Add user to docker usergroup
echo "Add user to docker usergroup"
sudo usermod -aG docker ${USER}
