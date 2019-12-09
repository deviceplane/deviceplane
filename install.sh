#!/bin/sh
set -e

# --- helper functions for logs ---
info()
{
    echo '[INFO] ' "$@"
}
fatal()
{
    echo '[ERROR] ' "$@" >&2
    exit 1
}

# --- fatal if no systemd ---
verify_system() {
    if [ -d /run/systemd ]; then
        return
    fi
    fatal 'Cannot find systemd as the process supervisor'
}

# --- escape most punctuation characters, except quotes, forward slash, and space ---
escape() {
    printf '%s' "$@" | sed -e 's/\([][!#$%&()*;<=>?\_`{|}]\)/\\\1/g;'
}

# --- escape double quotes ---
escape_dq() {
    printf '%s' "$@" | sed -e 's/"/\\"/g'
}

# --- define needed environment variables ---
setup_env() {
    # --- system name ---
    SYSTEM_NAME=deviceplane-agent

    # --- set related files from system name ---
    SERVICE=${SYSTEM_NAME}.service

    # --- use sudo if we are not already root ---
    SUDO=sudo
    if [ $(id -u) -eq 0 ]; then
        SUDO=
    fi

    # --- use binary install directory if defined or create default ---
    if [ -n "${INSTALL_BIN_DIR}" ]; then
        BIN_DIR=${INSTALL_BIN_DIR}
    else
        BIN_DIR=/usr/local/bin
    fi

    # --- use systemd directory if defined or create default ---
    if [ -n "${INSTALL_SYSTEMD_DIR}" ]; then
        SYSTEMD_DIR="${INSTALL_SYSTEMD_DIR}"
    else
        SYSTEMD_DIR=/etc/systemd/system
    fi

    # --- service file location ---
    FILE_SERVICE=${SYSTEMD_DIR}/${SERVICE}
}

# --- set arch and suffix, fatal if architecture not supported ---
setup_verify_arch() {
    if [ -z "$ARCH" ]; then
        ARCH=$(uname -m)
    fi
    case $ARCH in
        amd64)
            ARCH=amd64
            ;;
        x86_64)
            ARCH=amd64
            ;;
        arm64)
            ARCH=arm64
            ;;
        aarch64)
            ARCH=arm64
            ;;
        arm*)
            ARCH=arm
            ;;
        *)
            fatal "Unsupported architecture $ARCH"
    esac
}

# --- verify existence of network downloader executable ---
verify_downloader() {
    # Return failure if it doesn't exist or is no executable
    [ -x "$(which $1)" ] || return 1

    # Set verified executable as our downloader program and return success
    DOWNLOADER=$1
    return 0
}

# --- create tempory directory and cleanup when done ---
setup_tmp() {
    TMP_DIR=$(mktemp -d -t deviceplane.XXXXXXXXXX)
    TMP_BIN=${TMP_DIR}/deviceplane-agent
    cleanup() {
        code=$?
        set +e
        trap - EXIT
        rm -rf ${TMP_DIR}
        exit $code
    }
    trap cleanup INT EXIT
}

# --- download from github url ---
download() {
    [ $# -eq 2 ] || fatal 'download needs exactly 2 arguments'

    case $DOWNLOADER in
        curl)
            curl -o $1 -sfL $2
            ;;
        wget)
            wget -qO $1 $2
            ;;
        *)
            fatal "Incorrect executable '$DOWNLOADER'"
            ;;
    esac

    # Abort if download command failed
    [ $? -eq 0 ] || fatal 'Download failed'
}

# --- download binary from github url ---
download_binary() {
    BIN_URL=https://agent.deviceplane.com/${VERSION}/linux/${ARCH}/deviceplane-agent
    info "Downloading binary ${BIN_URL}"
    download ${TMP_BIN} ${BIN_URL}
}

# --- setup permissions and move binary to system directory ---
setup_binary() {
    chmod 755 ${TMP_BIN}
    info "Installing deviceplane-agent to ${BIN_DIR}/deviceplane-agent"
    $SUDO chown root:root ${TMP_BIN}
    $SUDO mv -f ${TMP_BIN} ${BIN_DIR}/deviceplane-agent
}

# --- download and verify ---
download_and_verify() {
    setup_verify_arch
    verify_downloader curl || verify_downloader wget || fatal 'Can not find curl or wget for downloading files'
    setup_tmp
    download_binary
    setup_binary
}

# --- write systemd service file ---
create_systemd_service_file() {
    info "systemd: Creating service file ${FILE_SERVICE}"
    $SUDO tee ${FILE_SERVICE} >/dev/null << EOF
[Unit]
Description=Deviceplane agent
Documentation=https://deviceplane.com/docs/
Wants=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
KillMode=process
Delegate=yes
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=${BIN_DIR}/deviceplane-agent \\
    --project=${PROJECT} --registration-token=${REGISTRATION_TOKEN}

EOF
}

# --- enable and start systemd service ---
systemd_enable() {
    info "systemd: Enabling ${SYSTEM_NAME} unit"
    $SUDO systemctl enable ${FILE_SERVICE} >/dev/null
    $SUDO systemctl daemon-reload >/dev/null
}

systemd_start() {
    info "systemd: Starting ${SYSTEM_NAME}"
    $SUDO systemctl restart ${SYSTEM_NAME}
}

# --- run the install process --
{
    verify_system
    setup_env "$@"
    download_and_verify
    create_systemd_service_file
    systemd_enable
    systemd_start
}
