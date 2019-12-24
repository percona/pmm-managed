#!/usr/bin/python2

from __future__ import print_function, unicode_literals
import multiprocessing, os, subprocess, time


GO = 'https://dl.google.com/go/go1.12.10.linux-amd64.tar.gz'


def run_commands(commands):
    """Runs given shell commands and checks exit codes."""

    for cmd in commands:
        print(">", cmd)
        subprocess.check_call(cmd, shell=True)


def install_packages():
    """Installs required and useful RPM packages."""

    run_commands([
        # to install man pages
        "sed -i '/nodocs/d' /etc/yum.conf",

        # reinstall with man pages
        "yum reinstall -y yum rpm",

        "yum install -y gcc git make pkgconfig glibc-static \
            ansible-lint \
            mc tmux psmisc lsof which iproute \
            bash-completion bash-completion-extras \
            man man-pages",
    ])


def install_go():
    """Installs Go toolchain."""

    run_commands([
        "curl -sS {go} -o /tmp/golang.tar.gz".format(go=GO),
        "tar -C /usr/local -xzf /tmp/golang.tar.gz",
        "mkdir -p /root/go/bin",
        "update-alternatives --install '/usr/bin/go' 'go' '/usr/local/go/bin/go' 0",
        "update-alternatives --set go /usr/local/go/bin/go",
        "update-alternatives --install '/usr/bin/gofmt' 'gofmt' '/usr/local/go/bin/gofmt' 0",
        "update-alternatives --set gofmt /usr/local/go/bin/gofmt",
    ])

def make_install():
    """Runs make install."""

    run_commands([
        "make install",
    ])

def install_tools():
    """Installs Go developer tools."""

    run_commands([
        "curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh",
        "curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /root/go/bin",

        "rm -fr /tmp/tools && \
            mkdir -p /tmp/tools && \
            cd /tmp/tools && \
            go mod init tools && \
            env GOPROXY=https://proxy.golang.org go get -v \
                github.com/go-delve/delve/cmd/dlv \
                golang.org/x/tools/cmd/gopls"
    ])


def install_vendored_tools():
    """Installs pmm-managed-specific Go tools."""

    run_commands([
        "go install ./vendor/github.com/BurntSushi/go-sumtype",
        "go install ./vendor/github.com/vektra/mockery/cmd/mockery",
        "go install ./vendor/golang.org/x/tools/cmd/goimports",
        "go install ./vendor/gopkg.in/reform.v1/reform",
    ])


def setup():
    """Runs various setup commands."""

    run_commands([
        "supervisorctl stop pmm-managed",
        "psql --username=postgres --command='ALTER USER \"pmm-managed\" WITH SUPERUSER'",
    ])


def main():
    # install packages early as they will be required below
    install_packages_p = multiprocessing.Process(target=install_packages)
    install_packages_p.start()

    # install Go and wait for it
    install_go()

    # install tools (requires Go)
    install_tools_p = multiprocessing.Process(target=install_tools)
    install_tools_p.start()
    install_vendored_tools_p = multiprocessing.Process(target=install_vendored_tools)
    install_vendored_tools_p.start()

    # make install (requires make package)
    install_packages_p.join()
    make_install()

    # do basic setup
    setup()

    # wait for everything else to finish
    install_tools_p.join()
    install_vendored_tools_p.join()


MARKER = "/tmp/devcontainer-setup-done"
if os.path.exists(MARKER):
    print(MARKER, "exists, exiting.")
    exit(0)

start = time.time()
main()
print("Done in", time.time() - start)

open(MARKER, 'w').close()
