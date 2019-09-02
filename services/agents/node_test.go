// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package agents

import (
	"testing"

	"github.com/percona/pmm/api/agentpb"
	"github.com/stretchr/testify/assert"

	"github.com/percona/pmm-managed/models"
)

func TestNodeExporterConfig(t *testing.T) {
	t.Run("Linux", func(t *testing.T) {
		node := &models.Node{}
		exporter := &models.Agent{}
		actual := nodeExporterConfig(node, exporter)
		expected := &agentpb.SetStateRequest_AgentProcess{
			Type:               agentpb.Type_NODE_EXPORTER,
			TemplateLeftDelim:  "{{",
			TemplateRightDelim: "}}",
			Args: []string{
				"--collector.arp=false",
				"--collector.bcache=false",
				"--collector.bonding=true",
				"--collector.buddyinfo=false",
				"--collector.conntrack=false",
				"--collector.cpu=false",
				"--collector.diskstats=true",
				"--collector.drbd=false",
				"--collector.edac=false",
				"--collector.entropy=true",
				"--collector.filefd=true",
				"--collector.filesystem=true",
				"--collector.hwmon=false",
				"--collector.infiniband=false",
				"--collector.interrupts=false",
				"--collector.ipvs=false",
				"--collector.ksmd=false",
				"--collector.loadavg=true",
				"--collector.logind=false",
				"--collector.mdadm=false",
				"--collector.meminfo=true",
				"--collector.meminfo_numa=true",
				"--collector.mountstats=false",
				"--collector.netclass=false",
				"--collector.netdev=true",
				"--collector.netstat.fields=^(.*_(InErrors|InErrs|InCsumErrors)" +
					"|Tcp_(ActiveOpens|PassiveOpens|RetransSegs|CurrEstab|AttemptFails|OutSegs|InSegs|EstabResets|OutRsts|OutSegs)|Tcp_Rto(Algorithm|Min|Max)" +
					"|Udp_(RcvbufErrors|SndbufErrors)|Udp(6?|Lite6?)_(InDatagrams|OutDatagrams|RcvbufErrors|SndbufErrors|NoPorts)" +
					"|Icmp6?_(OutEchoReps|OutEchos|InEchos|InEchoReps|InAddrMaskReps|InAddrMasks|OutAddrMaskReps|OutAddrMasks|InTimestampReps|InTimestamps" +
					"|OutTimestampReps|OutTimestamps|OutErrors|InDestUnreachs|OutDestUnreachs|InTimeExcds|InRedirects|OutRedirects|InMsgs|OutMsgs)" +
					"|IcmpMsg_(InType3|OutType3)|Ip(6|Ext)_(InOctets|OutOctets)|Ip_Forwarding|TcpExt_(Listen.*|Syncookies.*|TCPTimeouts))$",
				"--collector.netstat=true",
				"--collector.nfs=false",
				"--collector.nfsd=false",
				"--collector.ntp=false",
				"--collector.processes=false",
				"--collector.qdisc=false",
				"--collector.runit=false",
				"--collector.sockstat=false",
				"--collector.standard.go=true",
				"--collector.standard.process=true",
				"--collector.stat=true",
				"--collector.supervisord=false",
				"--collector.systemd=false",
				"--collector.tcpstat=false",
				"--collector.textfile.directory.hr=/usr/local/percona/pmm2/collectors/textfile-collector/high-resolution",
				"--collector.textfile.directory.lr=/usr/local/percona/pmm2/collectors/textfile-collector/low-resolution",
				"--collector.textfile.directory.mr=/usr/local/percona/pmm2/collectors/textfile-collector/medium-resolution",
				"--collector.textfile.hr=true",
				"--collector.textfile.lr=true",
				"--collector.textfile.mr=true",
				"--collector.time=true",
				"--collector.timex=false",
				"--collector.uname=true",
				"--collector.vmstat.fields=^(pg(steal_(kswapd|direct)|refill|alloc)_(movable|normal|dma3?2?)" +
					"|nr_(dirty.*|slab.*|vmscan.*|isolated.*|free.*|shmem.*|i?n?active.*|anon_transparent_.*|writeback.*|unstable" +
					"|unevictable|mlock|mapped|bounce|page_table_pages|kernel_stack)|drop_slab|slabs_scanned|pgd?e?activate" +
					"|pgpg(in|out)|pswp(in|out)|pgm?a?j?fault)$",
				"--collector.vmstat=true",
				"--collector.wifi=false",
				"--collector.xfs=false",
				"--collector.zfs=false",
				"--web.disable-exporter-metrics",
				"--web.listen-address=:{{ .listen_port }}",
			},
		}
		assertNoDuplicateFlags(t, actual.Args)
		assert.Equal(t, expected.Args, actual.Args)
		assert.Equal(t, expected.Env, actual.Env)
		assert.Equal(t, expected, actual)
	})

	t.Run("MacOS", func(t *testing.T) {
		node := &models.Node{
			Distro: "darwin",
		}
		exporter := &models.Agent{}
		actual := nodeExporterConfig(node, exporter)
		expected := &agentpb.SetStateRequest_AgentProcess{
			Type:               agentpb.Type_NODE_EXPORTER,
			TemplateLeftDelim:  "{{",
			TemplateRightDelim: "}}",
			Args: []string{
				"--collector.textfile.directory.hr=/usr/local/percona/pmm2/collectors/textfile-collector/high-resolution",
				"--collector.textfile.directory.lr=/usr/local/percona/pmm2/collectors/textfile-collector/low-resolution",
				"--collector.textfile.directory.mr=/usr/local/percona/pmm2/collectors/textfile-collector/medium-resolution",
				"--web.disable-exporter-metrics",
				"--web.listen-address=:{{ .listen_port }}",
			},
		}
		assertNoDuplicateFlags(t, actual.Args)
		assert.Equal(t, expected.Args, actual.Args)
		assert.Equal(t, expected.Env, actual.Env)
		assert.Equal(t, expected, actual)
	})
}
