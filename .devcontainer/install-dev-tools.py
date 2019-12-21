#!/usr/bin/python2

from __future__ import print_function, unicode_literals
import multiprocessing, subprocess, time


GO = 'https://dl.google.com/go/go1.12.10.linux-amd64.tar.gz'


def run_commands(commands):
    for cmd in commands:
        print(">", cmd)
        subprocess.check_call(cmd, shell=True)


def setup():
    run_commands([
        "supervisorctl stop pmm-managed",
        "psql --username=postgres --command='ALTER USER \"pmm-managed\" WITH SUPERUSER'",
    ])


def install_packages():
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
    run_commands([
        "make install",
    ])

def install_tools():
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
    run_commands([
        "go install ./vendor/github.com/BurntSushi/go-sumtype",
        "go install ./vendor/github.com/vektra/mockery/cmd/mockery",
        "go install ./vendor/golang.org/x/tools/cmd/goimports",
        "go install ./vendor/gopkg.in/reform.v1/reform",
    ])


def main():
    setup_p = multiprocessing.Process(target=setup)
    setup_p.start()

    install_packages_p = multiprocessing.Process(target=install_packages)
    install_packages_p.start()

    install_go()

    install_tools_p = multiprocessing.Process(target=install_tools)
    install_tools_p.start()

    install_vendored_tools_p = multiprocessing.Process(target=install_vendored_tools)
    install_vendored_tools_p.start()

    install_packages_p.join()
    make_install()

    setup_p.join()
    install_tools_p.join()
    install_vendored_tools_p.join()


start = time.time()
main()
print("Done in", time.time() - start)
