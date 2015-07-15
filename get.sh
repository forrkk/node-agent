#!/bin/sh

agent_dir='/opt/wodby'
agent_version='0.0.1'
agent_url="http://someurl/agent/${agent_version}/agent.tar.gz"

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
${curl} ${agent_url} | tar xvz -C "${agent_dir}"
chmod +x "${agent_dir}"/bin/*

echo exec "${agent_dir}"/bin/wodby_agent "${@}"