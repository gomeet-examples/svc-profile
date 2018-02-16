# Working with the sources

## Install

### Install requirements

svc-profile needs some requirements :

- [golang](https://golang.org/doc/install)
- [protobuf](https://github.com/google/protobuf)
- [git flow](https://danielkummer.github.io/git-flow-cheatsheet/)
- [docker](https://docs.docker.com/engine/installation/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [mysql](https://www.mysql.com/) or [mariaDB clone](https://mariadb.com/)

#### On Linux (Ubuntu Xenial)

```bash
sudo apt-get update
sudo apt-get install -y build-essential git software-properties-common python-software-properties
sudo add-apt-repository -y ppa:longsleep/golang-backports
sudo add-apt-repository -y ppa:maarten-fonville/protobuf
sudo apt-get update
sudo apt-get install -y golang-go protobuf-compiler git-flow

echo -e "export GOPATH=\$(go env GOPATH)\nexport PATH=\${PATH}:\${GOPATH}/bin" >> ~/.bashrc
source ~/.bashrc

# docker install cf. https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/
# 1. Add repository
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
# Response is something like this
#   pub   4096R/0EBFCD88 2017-02-22
#         Key fingerprint = 9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
#   uid                  Docker Release (CE deb) <docker@docker.com>
#   sub   4096R/F273FCD8 2017-02-22
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install docker-ce

# docker compose install https://docs.docker.com/compose/install/#install-compose
# 1. download the latest version
sudo curl \
     -L https://github.com/docker/compose/releases/download/1.17.0/docker-compose-`uname -s`-`uname -m` \
     -o /usr/local/bin/docker-compose
# 2. Apply executable permissions to the binary
sudo chmod +x /usr/local/bin/docker-compose

# install mariadb/mysql
sudo apt-get install mariadb-server

```

#### On MacOSX

```bash
brew install go
brew install git
brew install protobuf
brew install git-flow-avh

echo -e "export GOPATH=\$(go env GOPATH)\nexport PATH=\${PATH}:\${GOPATH}/bin" >> ~/.bashrc
echo "" >> ~/.bashrc
source ~/.bashrc

# for docker see https://docs.docker.com/docker-for-mac/install/
brew install mysql

```

#### On Windows

```txt
TODO
```

### Check version dependencies

```bash
$ protoc --version
libprotoc 3.3.0

$ go version
go version go1.8.1 ...snip...

$ echo $GOPATH
...snip...

$ docker --version
Docker version 17.06.2-ce, build cec0b72

$ docker-compse --version
docker-compose version 1.17.0, build ac53b73
```

### Install and build from source

```bash
# clone repository in $GOPATH
mkdir -p $GOPATH/src/github.com/gomeet-examples
cd $GOPATH/src/github.com/gomeet-examples
git clone https://github.com/gomeet-examples/svc-profile.git
cd svc-profile

# initalize git-flow needed for make release
git checkout master
git checkout develop
git flow init -d

# use make
make
_build/svc-profile
```

### Database initialization

- MySQL database

```bash
$ sudo mysql -u root -p
Enter password: ***********
MariaDB [(none)]> CREATE DATABASE svc_profile;
...
MariaDB [(none)]> CREATE DATABASE svc_profile_test;
...
MariaDB [(none)]> GRANT ALL PRIVILEGES ON svc_profile.* To '<USERNAME>'@'localhost' IDENTIFIED BY '<PASSWORD>';
...
MariaDB [(none)]> GRANT ALL PRIVILEGES ON svc_profile_test.* To '<USERNAME>'@'localhost' IDENTIFIED BY '<PASSWORD>';
...
```

### Make directives

- `make` - Installing all development tools and making specific platform binary. It's equivalent to `make build`
- `make build` - Installing all development tools and making specific platform binary
- `make clean` - Removing all generated files (all tools, all compiled files, generated from proto). It's equivalent to `make tools-clean package-clean proto-clean`
- `make docker-test` - Executing `make test` inside docker
- `make docker` - Building docker image
- `make docker-push` - Push the docker image to docker registry server - default registry is `docker.io` it can be overloaded via the environment variable `DOCKER_REGISTRY` like this `DOCKER_REGISTRY={{hostname}}:{{port}} make docker-push`
- `make install` - Performing a `go install` command
- `make package` - Building all packages (multi platform, and docker image)
- `make package-clean` - Clean up the builded packages
- `make package-proto` - Building the `_build/packaged/proto.tgz` file with dirstribluables protobuf files
- `make proto` - Generating files from proto
- `make proto-clean` - Clean up generated files from the proto file
- `make release` - Making a release (see below)
- `make start` - Building docker image and performing a `docker-compose up -d` command
- `make stop` - Performing a `docker-compose down` command
- `make test` - Runing tests
- `make tools` - Installing all development tools
- `make tools-sync` - Re-Syncronizes development tools
- `make tools-upgrade-gomeet` - Upgrading all gomeet's development tools [gomeet-tools-markdown-server](github.com/gomeet/gomeet-tools-markdown-server), [protoc-gen-gomeetfaker](github.com/gomeet/go-proto-gomeetfaker/protoc-gen-gomeetfaker), [gomeet & protoc-gen-gomeet-service](https://github.com/gomeet/gomeet)
- `make tools-upgrade` - Upgrading all development tools
- `make tools-clean` - Uninstall all development tools
- `make dep` - Executes the `dep ensure` command
- `make dep-prune` - Executes the `dep prune` command
- `make dep-update-gomeetexamples [individual svc name without svc- prefix|default all]` - Executes the `dep ensure -update github.com/gomeet-examples/svc-[individual svc name without svc- prefix|default all]`
- `make dep-update-gomeet-utils` - Executes the `dep ensure -update github.com/gomeet/gomeet`
- `make gomeet-regenerate-project` - regenerate the project with [gomeet](https://github.com/gomeet/gomeet) be careful this replaces files except for the protobuf file

#### Add a tool

Build tool chain:

```shell
make tools-sync
make tools
```

Add a tool dependency:

```shell
_tools/bin/retool add retool github.com/jteeuwen/go-bindata/go-bindata origin/master
```

Use a tool:

```shell
_tools/bin/go-bindata
```

Commit changes

#### Make a release

```bash
make release <Git flow option : start|finish> <Release version : major|minor|patch> [Release version metadata (optional)]
```

- Git flow option
  - The `start` option does not finish the `git flow release` so to finish the release and prepare `VERSION` file in `develop` branch do :
  ```bash
  git flow release finish "v$(cat VERSION)"
  # NB: _tools/bin/semver is compiled with "make tools"
  NEW_DEV_VERSION=`_tools/bin/semver -patch -build "dev" $(cat VERSION)` && echo "$NEW_DEV_VERSION" > VERSION
  git add VERSION
  git commit -m "Bump version - v$(cat VERSION)"
  git push --tag
  git push origin develop
  git push origin master
  ```
  - The `finish` option does it for you.

- Release version and metadata (if `VERSION` file in `develop` branch is `1.1.1+dev`) :
  - `make release start patch` start and publish the `release/1.1.1` git flow release branch
  - `make release start patch rc.1` start and publish the `release/1.1.1+rc.1` git flow release branch
  - `make release start minor` start and publish the `release/1.2.0` git flow release branch
  - `make release start minor foo.1` start and publish the `release/1.2.0+foo.1` git flow release branch
  - `make release start major` start and publish the `release/2.0.0` git flow release branch
  - `make release start major foo.1` start and publish the `release/2.0.0+foo.1` git flow release branch
  - `make release finish patch` make the `1.1.1` release and `VERSION` file in `develop` branch is `1.1.2+dev`
  - `make release finish patch rc.1` make the `1.1.1+rc.1` release and `VERSION` file in `develop` branch branch is `1.1.2+dev`
  - `make release finish minor` make the `1.2.0` release and `VERSION` file in `develop` branch is `1.2.1+dev`
  - `make release finish minor foo.1` make the `1.2.0+foo.1` release and `VERSION` file in `develop` branch is `1.2.0+dev`
  - `make release finish major` make the `2.0.0` release and `VERSION` file in `develop` branch is `2.0.1+dev`
  - `make release finish major foo.1` make the `2.0.0+foo.1` release and `VERSION` file in `develop` branch is `2.0.1+dev`

#### Manual steps

```bash
make tools
NEW_VERSION="x.y.z" && \
  git flow release start "v$NEW_VERSION" && \
  echo $NEW_VERSION > VERSION
git add VERSION
git commit -m "Bump version - v$(cat VERSION)"
awk \
  -v \
  log_title="## Unreleased\n\n- Nothing\n\n## $(cat VERSION) - $(date +%Y-%m-%d)" \
  '{gsub(/## Unreleased/,log_title)}1' \
  CHANGELOG.md > CHANGELOG.md.tmp && \
    mv CHANGELOG.md.tmp CHANGELOG.md
git add CHANGELOG.md
git commit -m "Improved CHANGELOG.md"
make package
git add _build/packaged/
git commit -m "Added v$(cat VERSION) packages"
git flow release publish "v$(cat VERSION)"
git flow release finish "v$(cat VERSION)"
# NB: _tools/bin/semver is compiled with "make tools"
NEW_DEV_VERSION=`_tools/bin/semver -patch -build "dev" $(cat VERSION)` && \
  echo $NEW_DEV_VERSION > VERSION
git add VERSION
git commit -m "Bump version - v$(cat VERSION)"
git push --tag
git push origin develop
git push origin master
```

## Use docker (no requirement)

- See gomeet/gomeet-builder docker image ([Docker Hub](https://hub.docker.com/r/gomeet/gomeet-builder/) - [Source](https://github.com/gomeet/gomeet-builder)).

## Working with gotools

If svc-profile repository is private and you use [Gogs](https://gogs.io/) has remote server.

To work with go tools (`go get` et `dep`) it's necesary to configure gogs, git and ssh.

1. Add your ssh key to your gogs user settings https://<GOGS_ADDRESS>/user/settings/ssh

2. In your local git config (`~/.gitconfig`) add these lines :

```
...
[url "ssh://<GOGS_SSH_USER>@<GOGS_ADDRESS>:<GOGS_SSH_PORT (default: 10022)>"]
	insteadOf = https://<GOGS_ADDRESS>
...
```

The SSH URL might require a trailing slash depending on the version of Git (observed on 2.9.3).

3. In your local ssh config (`~/.ssh/config`) add these lines :

```
...
Host <GOGS_ADDRESS>
  HostName <GOGS_ADDRESS>
  Port <GOGS_SSH_PORT (default: 10022)>
  User <GOGS_SSH_USER>
...
```

__WARNING__ : be sure that `ssh-agent` is running

## Uninstall

### Remove source

```bash
rm $GOPATH/bin/svc-profile
rm -rf $GOPATH/src/github.com/gomeet-examples
```

### Uninstall dependencies

- On Linux (Ubuntu Xenial)

```bash
sudo apt-get autoremove --purge golang-go protobuf-compiler git-flow
sudo add-apt-repository -r --purge ppa:longsleep/golang-backports
sudo add-apt-repository -r ppa:maarten-fonville/protobuf
sudo rm /etc/apt/sources.list.d/longsleep-ubuntu-golang-backports-xenial.list*
sudo rm /etc/apt/sources.list.d/maarten-fonville-ubuntu-protobuf-xenial.list*
sudo apt-get autoremove --purge build-essential git software-properties-common python-software-properties
sudo apt-get update

sed -i.bak ':a;N;$!ba;s/\nexport GOPATH=$(go env GOPATH)\nexport PATH=\${PATH}:\${GOPATH}\/bin//g' ~/.bashrc
unset GOPATH
source ~/.bashrc

# uninstall docker
sudo apt-get purge docker-ce
sudo rm -rf /var/lib/docker
sudo add-apt-repository -r "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-key del "Docker Release (CE deb) <docker@docker.com>"
```

- On MacOSX

```bash
brew rm go
brew rm git
brew rm protobuf
brew rm git-flow-avh

sed -i.bak '/export GOPATH=\$(go env GOPATH)/d' ~/.bashrc
unset GOPATH
source ~/.bashrc
```

- On Windows

```txt
TODO
```

## Some usual procedures

- To add a Gomeetexamples subservice as dependency see [this](../add_sub_service/README.md)
- To add a new gRPC service see [this](../add_grpc_service/README.md)
