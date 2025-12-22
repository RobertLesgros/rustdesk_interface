#!/usr/bin/env bash
source <(curl -fsSL https://raw.githubusercontent.com/community-scripts/ProxmoxVE/main/misc/build.func)
# Based on Community Scripts for ProxmoxVE
# Modified to install RustDesk Interface (custom)

APP="RustDesk Server (Interface)"
var_tags="${var_tags:-remote-desktop}"
var_cpu="${var_cpu:-1}"
var_ram="${var_ram:-512}"
var_disk="${var_disk:-2}"
var_os="${var_os:-debian}"
var_version="${var_version:-13}"
var_unprivileged="${var_unprivileged:-1}"

header_info "$APP"
variables
color
catch_errors

function update_script() {
  header_info
  check_container_storage
  check_container_resources

  if [[ ! -x /usr/bin/hbbr ]]; then
    msg_error "No ${APP} Installation Found!"
    exit
  fi

  RELEASE=$(curl -fsSL https://api.github.com/repos/rustdesk/rustdesk-server/releases/latest | grep "tag_name" | awk '{print substr($2, 2, length($2)-3) }')
  # Modified: Fetch from RobertLesgros/rustdesk_interface
  APIRELEASE=$(curl -fsSL https://api.github.com/repos/RobertLesgros/rustdesk_interface/releases/latest | grep "tag_name" | awk '{print substr($2, 3, length($2)-4) }')
  
  # Modified: Check .rustdesk-interface version file
  if [[ "${RELEASE}" != "$(cat ~/.rustdesk-hbbr)" ]] || [[ "${APIRELEASE}" != "$(cat ~/.rustdesk-interface)" ]] || [[ ! -f ~/.rustdesk-hbbr ]] || [[ ! -f ~/.rustdesk-interface ]]; then
    msg_info "Stopping Service"
    systemctl stop rustdesk-hbbr
    systemctl stop rustdesk-hbbs
    
    # Handle legacy service name if present
    if [[ -f /lib/systemd/system/rustdesk-api.service ]]; then
      systemctl stop rustdesk-api
      systemctl disable rustdesk-api
    fi
    # Handle new service name
    if [[ -f /lib/systemd/system/rustdesk-interface.service ]]; then
      systemctl stop rustdesk-interface
    fi
    msg_info "Stopped Service"

    fetch_and_deploy_gh_release "rustdesk-hbbr" "rustdesk/rustdesk-server" "binary" "latest" "/opt/rustdesk" "rustdesk-server-hbbr*amd64.deb"
    fetch_and_deploy_gh_release "rustdesk-hbbs" "rustdesk/rustdesk-server" "binary" "latest" "/opt/rustdesk" "rustdesk-server-hbbs*amd64.deb"
    fetch_and_deploy_gh_release "rustdesk-utils" "rustdesk/rustdesk-server" "binary" "latest" "/opt/rustdesk" "rustdesk-server-utils*amd64.deb"
    
    # Modified: Install rustdesk-interface
    fetch_and_deploy_gh_release "rustdesk-interface" "RobertLesgros/rustdesk_interface" "binary" "latest" "/opt/rustdesk" "rustdesk-interface-server*amd64.deb"

    msg_info "Starting services"
    systemctl start -q rustdesk-hbbr rustdesk-hbbs rustdesk-interface
    msg_ok "Services started"

    msg_ok "Updated successfully!"
  else
    msg_ok "No update required. ${APP} is already at v${RELEASE}"
  fi
  exit
}

start
build_container
description

msg_ok "Completed Successfully!\n"
echo -e "${CREATING}${GN}${APP} setup has been successfully initialized!${CL}"
echo -e "${INFO}${YW} Access it using the following URL:${CL}"
echo -e "${TAB}${GATEWAY}${BGN}${IP}:21114${CL}"
