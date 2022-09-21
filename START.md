After cloning, you will have all files in a directory.

### Run the commands to update and upgrade the system
- sudo apt update
- sudo apt upgrade

### Install go lang
- Download Go: wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz
- Clean the old go folder: sudo rm -rf /usr/local/go
- Extract Go: sudo tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz
- Add go to path: nano $HOME/.profile, then add the line at bottom export PATH=$PATH:/usr/local/go/bin
- Reload environment variables 'source $HOME/.profile'
- Test that go 1.19 is installed by typing go version

### Get SpineChain
- CD to the location you want the code, and type https://github.com/spinechain/dtp.git
- CD into the newly created folder called dtp
- Run with go run .


### Start SpineChain
- Install the dependencies with 'go get spinedtp'
- Run build.sh. It takes very long to compile the first time. The screen will be blank for long.
- At the end, you will have binaries in the build folder. Run these.
