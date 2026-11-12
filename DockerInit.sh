#!/bin/sh
set -eu

ARCH_INPUT="${1:-amd64}"
XRAY_VERSION="${2:-v26.11.9}"
case "$ARCH_INPUT" in
    amd64 | x86_64)
        XRAY_ARCH="64"
        OUTPUT_ARCH="amd64"
        ;;
    386 | i386)
        XRAY_ARCH="32"
        OUTPUT_ARCH="i386"
        ;;
    arm64 | armv8 | aarch64)
        XRAY_ARCH="arm64-v8a"
        OUTPUT_ARCH="arm64"
        ;;
    arm | armv7 | arm32)
        XRAY_ARCH="arm32-v7a"
        OUTPUT_ARCH="arm32"
        ;;
    armv6)
        XRAY_ARCH="arm32-v6"
        OUTPUT_ARCH="armv6"
        ;;
    armv5)
        XRAY_ARCH="arm32-v5"
        OUTPUT_ARCH="armv5"
        ;;
    s390x)
        XRAY_ARCH="s390x"
        OUTPUT_ARCH="s390x"
        ;;
    *)
        echo "Unsupported target architecture: $ARCH_INPUT" >&2
        exit 1
        ;;
esac

mkdir -p build/bin
cd build/bin
archive="Xray-linux-$XRAY_ARCH.zip"
url="https://github.com/XTLS/Xray-core/releases/download/$XRAY_VERSION/$archive"

wget -q -O "$archive" "$url"
unzip -q -o "$archive"
rm -f "$archive" geoip.dat geosite.dat
mv xray "xray-linux-$OUTPUT_ARCH"

# Fail the image build if any required rule database is missing.
wget -q -O geoip.dat https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget -q -O geosite.dat https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
wget -q -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -q -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -q -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -q -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat
