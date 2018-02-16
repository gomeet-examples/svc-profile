# svc-profile bare-metal installation

To install svc-profile on bare-metal you just need to download the appropriate pre compiled binary.

- If you use [gogs](https://gogs.io/)

__TODO outdated__

```bash
mkdir svc-profile
cd svc-profile
GOGS_USER="<your gogs username>" && \
GOGS_PASSWORD="<your gogs password>" && \
VERSION="<version ex: 0.0.1>" && \
ARCH="<arch darwin-amd64|linux-amd64|linux-arm32|linux-arm64|windows-amd64>" && \
BIN_NAME=$(if [ "$ARCH" = "windows-amd64" ]; then echo "svc-profile.exe"; else echo "svc-profile"; fi) && \
curl -O "https://$GOGS_USER:$GOGS_PASSWORD@github.com/gomeet-examples/svc-profile/raw/v$VERSION/_build/packaged/$ARCH/$BIN_NAME" && \
curl -O "https://$GOGS_USER:$GOGS_PASSWORD@github.com/gomeet-examples/svc-profile/raw/v$VERSION/_build/packaged/$ARCH/SHA256SUM" && \
shasum -c SHA256SUM
chmod +x $BIN_NAME
```
