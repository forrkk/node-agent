#!/bin/sh

agent_dir='/opt/wodby'
agent_bin_dir="${agent_dir}/bin"
agent_etc_dir="${agent_dir}/etc"
agent_version='0.0.1'
agent_url="https://github.com/Wodby/node-agent/releases/download/v${agent_version}/node-agent"

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

curl=''
if command_exists curl; then
	curl='curl -sSL'
elif command_exists wget; then
	curl='wget -qO-'
elif command_exists busybox && busybox --list-modules | grep -q wget; then
	curl='busybox wget -qO-'
fi
if [ -z "${curl}" ];then
 echo "error: curl or wget cannot be found, please install any"
 exit 1
fi

[ ! -d "${agent_dir}" ] && mkdir -p "${agent_dir}"
${curl} ${agent_url} > "${agent_bin_dir}"/node-agent
chmod +x "${agent_bin_dir}"/node-agent

exec "${agent_bin_dir}"/node-agent "${@}"