package filtering

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFiltering(t *testing.T) {
	var expectedFilteredMetrics = `# HELP cpu_guest_seconds_total Seconds the cpus spent in guests (VMs) for each mode.
# TYPE cpu_guest_seconds_total counter
cpu_guest_seconds_total{cpu="0",mode="nice"} 0
cpu_guest_seconds_total{cpu="0",mode="user"} 0
cpu_guest_seconds_total{cpu="1",mode="nice"} 0
cpu_guest_seconds_total{cpu="1",mode="user"} 0
cpu_guest_seconds_total{cpu="2",mode="nice"} 0
cpu_guest_seconds_total{cpu="2",mode="user"} 0
cpu_guest_seconds_total{cpu="3",mode="nice"} 0
cpu_guest_seconds_total{cpu="3",mode="user"} 0
# HELP cpu_seconds_total Seconds the cpus spent in each mode.
# TYPE cpu_seconds_total counter
cpu_seconds_total{cpu="0",mode="idle"} 540185.47
cpu_seconds_total{cpu="0",mode="iowait"} 235.69
cpu_seconds_total{cpu="0",mode="irq"} 0
cpu_seconds_total{cpu="0",mode="nice"} 0
cpu_seconds_total{cpu="0",mode="softirq"} 1606.73
cpu_seconds_total{cpu="0",mode="steal"} 0
cpu_seconds_total{cpu="0",mode="system"} 4613.14
cpu_seconds_total{cpu="0",mode="user"} 21853.84
cpu_seconds_total{cpu="1",mode="idle"} 543968.72
cpu_seconds_total{cpu="1",mode="iowait"} 160.87`
	efmArr := strings.Split(expectedFilteredMetrics, "\n")

	fm := FilterNodePrefix(rawUnfilteredHostMetrics)
	fmArr := strings.Split(fm, "\n")

	assert.Equal(
		t,
		strings.Join(fmArr[:len(efmArr)], "\n"),
		strings.Join(efmArr, "\n"),
	)
}
func TestEdgeCaseFiltering(t *testing.T) {
	var expectedFilteredMetrics = `# HELP cpu_guest_seconds_total Seconds the cpus spent node_cpu_stuff for each mode node node_ _node node_node.
# TYPE cpu_guest_seconds_total counter
cpu_guest_seconds_total{cpu="0",mode="node_nice_node"} 0
cpu_guest_seconds_total{cpu="0",mode="user"} 0
cpu_guest_seconds_total{cpu="1",mode="nice"} 0
cpu_guest_seconds_total{cpu="1",mode="user"} 0
cpu_guest_seconds_total{cpu="2",mode="nice"} 0
cpu_guest_seconds_total{cpu="2",mode="user"} 0
cpu_guest_seconds_total{cpu="3",mode="nice"} 0
cpu_guest_seconds_total{cpu="3",mode="user"} 0
# HELP cpu_seconds_total Seconds the cpus_node node_spent in each mode.
# TYPE cpu_seconds_total counter
cpu_seconds_total{cpu="0",mode="idle_node"} 540185.47
cpu_seconds_total{cpu="0",mode="iowait"} 235.69
cpu_seconds_total{cpu="0",mode="irq"} 0
cpu_seconds_total{cpu="0",mode="nice"} 0
cpu_seconds_total{cpu="0",mode="softirq"} 1606.73
cpu_seconds_total{cpu="0",mode="steal"} 0
cpu_seconds_total{cpu="0",mode="node_system"} 4613.14
cpu_seconds_total{cpu="0",mode="user"} 21853.84
cpu_seconds_total{cpu="1",mode="idle"} 543968.72
cpu_seconds_total{cpu="1",mode="iowait"} 160.87`

	var exampleUnfilteredMetrics = `# HELP node_cpu_guest_seconds_total Seconds the cpus spent node_cpu_stuff for each mode node node_ _node node_node.
# TYPE node_cpu_guest_seconds_total counter
node_cpu_guest_seconds_total{cpu="0",mode="node_nice_node"} 0
node_cpu_guest_seconds_total{cpu="0",mode="user"} 0
node_cpu_guest_seconds_total{cpu="1",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="1",mode="user"} 0
node_cpu_guest_seconds_total{cpu="2",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="2",mode="user"} 0
node_cpu_guest_seconds_total{cpu="3",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="3",mode="user"} 0
# HELP node_cpu_seconds_total Seconds the cpus_node node_spent in each mode.
# TYPE node_cpu_seconds_total counter
node_cpu_seconds_total{cpu="0",mode="idle_node"} 540185.47
node_cpu_seconds_total{cpu="0",mode="iowait"} 235.69
node_cpu_seconds_total{cpu="0",mode="irq"} 0
node_cpu_seconds_total{cpu="0",mode="nice"} 0
node_cpu_seconds_total{cpu="0",mode="softirq"} 1606.73
node_cpu_seconds_total{cpu="0",mode="steal"} 0
node_cpu_seconds_total{cpu="0",mode="node_system"} 4613.14
node_cpu_seconds_total{cpu="0",mode="user"} 21853.84
node_cpu_seconds_total{cpu="1",mode="idle"} 543968.72
node_cpu_seconds_total{cpu="1",mode="iowait"} 160.87`

	assert.Equal(
		t,
		expectedFilteredMetrics,
		FilterNodePrefix(exampleUnfilteredMetrics),
	)
}

const rawUnfilteredHostMetrics = `# HELP node_cpu_guest_seconds_total Seconds the cpus spent in guests (VMs) for each mode.
# TYPE node_cpu_guest_seconds_total counter
node_cpu_guest_seconds_total{cpu="0",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="0",mode="user"} 0
node_cpu_guest_seconds_total{cpu="1",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="1",mode="user"} 0
node_cpu_guest_seconds_total{cpu="2",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="2",mode="user"} 0
node_cpu_guest_seconds_total{cpu="3",mode="nice"} 0
node_cpu_guest_seconds_total{cpu="3",mode="user"} 0
# HELP node_cpu_seconds_total Seconds the cpus spent in each mode.
# TYPE node_cpu_seconds_total counter
node_cpu_seconds_total{cpu="0",mode="idle"} 540185.47
node_cpu_seconds_total{cpu="0",mode="iowait"} 235.69
node_cpu_seconds_total{cpu="0",mode="irq"} 0
node_cpu_seconds_total{cpu="0",mode="nice"} 0
node_cpu_seconds_total{cpu="0",mode="softirq"} 1606.73
node_cpu_seconds_total{cpu="0",mode="steal"} 0
node_cpu_seconds_total{cpu="0",mode="system"} 4613.14
node_cpu_seconds_total{cpu="0",mode="user"} 21853.84
node_cpu_seconds_total{cpu="1",mode="idle"} 543968.72
node_cpu_seconds_total{cpu="1",mode="iowait"} 160.87
node_cpu_seconds_total{cpu="1",mode="irq"} 0
node_cpu_seconds_total{cpu="1",mode="nice"} 0
node_cpu_seconds_total{cpu="1",mode="softirq"} 290.09
node_cpu_seconds_total{cpu="1",mode="steal"} 0
node_cpu_seconds_total{cpu="1",mode="system"} 5419.9
node_cpu_seconds_total{cpu="1",mode="user"} 24590.12
node_cpu_seconds_total{cpu="2",mode="idle"} 546289.31
node_cpu_seconds_total{cpu="2",mode="iowait"} 185.16
node_cpu_seconds_total{cpu="2",mode="irq"} 0
node_cpu_seconds_total{cpu="2",mode="nice"} 0
node_cpu_seconds_total{cpu="2",mode="softirq"} 203.27
node_cpu_seconds_total{cpu="2",mode="steal"} 0
node_cpu_seconds_total{cpu="2",mode="system"} 4115.2
node_cpu_seconds_total{cpu="2",mode="user"} 23729.88
node_cpu_seconds_total{cpu="3",mode="idle"} 547303.49
node_cpu_seconds_total{cpu="3",mode="iowait"} 183.2
node_cpu_seconds_total{cpu="3",mode="irq"} 0
node_cpu_seconds_total{cpu="3",mode="nice"} 0
node_cpu_seconds_total{cpu="3",mode="softirq"} 208.09
node_cpu_seconds_total{cpu="3",mode="steal"} 0
node_cpu_seconds_total{cpu="3",mode="system"} 4125.51
node_cpu_seconds_total{cpu="3",mode="user"} 23413.07
# HELP node_disk_io_now The number of I/Os currently in progress.
# TYPE node_disk_io_now gauge
node_disk_io_now{device="mmcblk0"} 0
node_disk_io_now{device="mmcblk0p1"} 0
node_disk_io_now{device="mmcblk0p2"} 0
# HELP node_disk_io_time_seconds_total Total seconds spent doing I/Os.
# TYPE node_disk_io_time_seconds_total counter
node_disk_io_time_seconds_total{device="mmcblk0"} 1128.96
node_disk_io_time_seconds_total{device="mmcblk0p1"} 0.44
node_disk_io_time_seconds_total{device="mmcblk0p2"} 1128.8600000000001
# HELP node_disk_io_time_weighted_seconds_total The weighted # of seconds spent doing I/Os.
# TYPE node_disk_io_time_weighted_seconds_total counter
node_disk_io_time_weighted_seconds_total{device="mmcblk0"} 2842.02
node_disk_io_time_weighted_seconds_total{device="mmcblk0p1"} 0.47000000000000003
node_disk_io_time_weighted_seconds_total{device="mmcblk0p2"} 2841.4900000000002
# HELP node_disk_read_bytes_total The total number of bytes read successfully.
# TYPE node_disk_read_bytes_total counter
node_disk_read_bytes_total{device="mmcblk0"} 1.937740288e+09
node_disk_read_bytes_total{device="mmcblk0p1"} 2.645504e+06
node_disk_read_bytes_total{device="mmcblk0p2"} 1.934558208e+09
# HELP node_disk_read_time_seconds_total The total number of seconds spent by all reads.
# TYPE node_disk_read_time_seconds_total counter
node_disk_read_time_seconds_total{device="mmcblk0"} 191.82
node_disk_read_time_seconds_total{device="mmcblk0p1"} 0.42
node_disk_read_time_seconds_total{device="mmcblk0p2"} 191.34
# HELP node_disk_reads_completed_total The total number of reads completed successfully.
# TYPE node_disk_reads_completed_total counter
node_disk_reads_completed_total{device="mmcblk0"} 19544
node_disk_reads_completed_total{device="mmcblk0p1"} 103
node_disk_reads_completed_total{device="mmcblk0p2"} 19428
# HELP node_disk_reads_merged_total The total number of reads merged.
# TYPE node_disk_reads_merged_total counter
node_disk_reads_merged_total{device="mmcblk0"} 6578
node_disk_reads_merged_total{device="mmcblk0p1"} 656
node_disk_reads_merged_total{device="mmcblk0p2"} 5922
# HELP node_disk_write_time_seconds_total This is the total number of seconds spent by all writes.
# TYPE node_disk_write_time_seconds_total counter
node_disk_write_time_seconds_total{device="mmcblk0"} 2651.11
node_disk_write_time_seconds_total{device="mmcblk0p1"} 0.05
node_disk_write_time_seconds_total{device="mmcblk0p2"} 2651.06
# HELP node_disk_writes_completed_total The total number of writes completed successfully.
# TYPE node_disk_writes_completed_total counter
node_disk_writes_completed_total{device="mmcblk0"} 458980
node_disk_writes_completed_total{device="mmcblk0p1"} 2
node_disk_writes_completed_total{device="mmcblk0p2"} 458978
# HELP node_disk_writes_merged_total The number of writes merged.
# TYPE node_disk_writes_merged_total counter
node_disk_writes_merged_total{device="mmcblk0"} 1.168384e+06
node_disk_writes_merged_total{device="mmcblk0p1"} 0
node_disk_writes_merged_total{device="mmcblk0p2"} 1.168384e+06
# HELP node_disk_written_bytes_total The total number of bytes written successfully.
# TYPE node_disk_written_bytes_total counter
node_disk_written_bytes_total{device="mmcblk0"} 8.06265344e+09
node_disk_written_bytes_total{device="mmcblk0p1"} 1024
node_disk_written_bytes_total{device="mmcblk0p2"} 8.062652416e+09
# HELP node_filesystem_avail_bytes Filesystem space available to non-root users in bytes.
# TYPE node_filesystem_avail_bytes gauge
node_filesystem_avail_bytes{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 2.1707776e+07
node_filesystem_avail_bytes{device="/dev/root",fstype="ext4",mountpoint="/"} 4.951965696e+09
node_filesystem_avail_bytes{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 4.3700224e+08
node_filesystem_avail_bytes{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 5.238784e+06
# HELP node_filesystem_device_error Whether an error occurred while getting statistics for the given device.
# TYPE node_filesystem_device_error gauge
node_filesystem_device_error{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 0
node_filesystem_device_error{device="/dev/root",fstype="ext4",mountpoint="/"} 0
node_filesystem_device_error{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 0
node_filesystem_device_error{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 0
# HELP node_filesystem_files Filesystem total file nodes.
# TYPE node_filesystem_files gauge
node_filesystem_files{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 0
node_filesystem_files{device="/dev/root",fstype="ext4",mountpoint="/"} 887168
node_filesystem_files{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 118681
node_filesystem_files{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 118681
# HELP node_filesystem_files_free Filesystem total free file nodes.
# TYPE node_filesystem_files_free gauge
node_filesystem_files_free{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 0
node_filesystem_files_free{device="/dev/root",fstype="ext4",mountpoint="/"} 524393
node_filesystem_files_free{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 118245
node_filesystem_files_free{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 118678
# HELP node_filesystem_free_bytes Filesystem free space in bytes.
# TYPE node_filesystem_free_bytes gauge
node_filesystem_free_bytes{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 2.1707776e+07
node_filesystem_free_bytes{device="/dev/root",fstype="ext4",mountpoint="/"} 5.599297536e+09
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 4.3700224e+08
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 5.238784e+06
# HELP node_filesystem_readonly Filesystem read-only status.
# TYPE node_filesystem_readonly gauge
node_filesystem_readonly{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 0
node_filesystem_readonly{device="/dev/root",fstype="ext4",mountpoint="/"} 0
node_filesystem_readonly{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 0
node_filesystem_readonly{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 0
# HELP node_filesystem_size_bytes Filesystem size in bytes.
# TYPE node_filesystem_size_bytes gauge
node_filesystem_size_bytes{device="/dev/mmcblk0p1",fstype="vfat",mountpoint="/boot"} 4.4271104e+07
node_filesystem_size_bytes{device="/dev/root",fstype="ext4",mountpoint="/"} 1.5190114304e+10
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",mountpoint="/run"} 4.86117376e+08
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",mountpoint="/run/lock"} 5.24288e+06
# HELP node_load1 1m load average.
# TYPE node_load1 gauge
node_load1 0.19
# HELP node_load15 15m load average.
# TYPE node_load15 gauge
node_load15 0.31
# HELP node_load5 5m load average.
# TYPE node_load5 gauge
node_load5 0.33
# HELP node_memory_Active_anon_bytes Memory information field Active_anon_bytes.
# TYPE node_memory_Active_anon_bytes gauge
node_memory_Active_anon_bytes 9.3691904e+07
# HELP node_memory_Active_bytes Memory information field Active_bytes.
# TYPE node_memory_Active_bytes gauge
node_memory_Active_bytes 2.46403072e+08
# HELP node_memory_Active_file_bytes Memory information field Active_file_bytes.
# TYPE node_memory_Active_file_bytes gauge
node_memory_Active_file_bytes 1.52711168e+08
# HELP node_memory_AnonPages_bytes Memory information field AnonPages_bytes.
# TYPE node_memory_AnonPages_bytes gauge
node_memory_AnonPages_bytes 1.75386624e+08
# HELP node_memory_Bounce_bytes Memory information field Bounce_bytes.
# TYPE node_memory_Bounce_bytes gauge
node_memory_Bounce_bytes 0
# HELP node_memory_Buffers_bytes Memory information field Buffers_bytes.
# TYPE node_memory_Buffers_bytes gauge
node_memory_Buffers_bytes 5.4116352e+07
# HELP node_memory_Cached_bytes Memory information field Cached_bytes.
# TYPE node_memory_Cached_bytes gauge
node_memory_Cached_bytes 5.75582208e+08
# HELP node_memory_CmaFree_bytes Memory information field CmaFree_bytes.
# TYPE node_memory_CmaFree_bytes gauge
node_memory_CmaFree_bytes 6.90176e+06
# HELP node_memory_CmaTotal_bytes Memory information field CmaTotal_bytes.
# TYPE node_memory_CmaTotal_bytes gauge
node_memory_CmaTotal_bytes 8.388608e+06
# HELP node_memory_CommitLimit_bytes Memory information field CommitLimit_bytes.
# TYPE node_memory_CommitLimit_bytes gauge
node_memory_CommitLimit_bytes 5.9097088e+08
# HELP node_memory_Committed_AS_bytes Memory information field Committed_AS_bytes.
# TYPE node_memory_Committed_AS_bytes gauge
node_memory_Committed_AS_bytes 8.0744448e+08
# HELP node_memory_Dirty_bytes Memory information field Dirty_bytes.
# TYPE node_memory_Dirty_bytes gauge
node_memory_Dirty_bytes 77824
# HELP node_memory_Inactive_anon_bytes Memory information field Inactive_anon_bytes.
# TYPE node_memory_Inactive_anon_bytes gauge
node_memory_Inactive_anon_bytes 1.19570432e+08
# HELP node_memory_Inactive_bytes Memory information field Inactive_bytes.
# TYPE node_memory_Inactive_bytes gauge
node_memory_Inactive_bytes 5.58768128e+08
# HELP node_memory_Inactive_file_bytes Memory information field Inactive_file_bytes.
# TYPE node_memory_Inactive_file_bytes gauge
node_memory_Inactive_file_bytes 4.39197696e+08
# HELP node_memory_KernelStack_bytes Memory information field KernelStack_bytes.
# TYPE node_memory_KernelStack_bytes gauge
node_memory_KernelStack_bytes 1.777664e+06
# HELP node_memory_Mapped_bytes Memory information field Mapped_bytes.
# TYPE node_memory_Mapped_bytes gauge
node_memory_Mapped_bytes 1.09707264e+08
# HELP node_memory_MemAvailable_bytes Memory information field MemAvailable_bytes.
# TYPE node_memory_MemAvailable_bytes gauge
node_memory_MemAvailable_bytes 6.58018304e+08
# HELP node_memory_MemFree_bytes Memory information field MemFree_bytes.
# TYPE node_memory_MemFree_bytes gauge
node_memory_MemFree_bytes 9.6243712e+07
# HELP node_memory_MemTotal_bytes Memory information field MemTotal_bytes.
# TYPE node_memory_MemTotal_bytes gauge
node_memory_MemTotal_bytes 9.72234752e+08
# HELP node_memory_Mlocked_bytes Memory information field Mlocked_bytes.
# TYPE node_memory_Mlocked_bytes gauge
node_memory_Mlocked_bytes 12288
# HELP node_memory_NFS_Unstable_bytes Memory information field NFS_Unstable_bytes.
# TYPE node_memory_NFS_Unstable_bytes gauge
node_memory_NFS_Unstable_bytes 0
# HELP node_memory_PageTables_bytes Memory information field PageTables_bytes.
# TYPE node_memory_PageTables_bytes gauge
node_memory_PageTables_bytes 1.7408e+06
# HELP node_memory_SReclaimable_bytes Memory information field SReclaimable_bytes.
# TYPE node_memory_SReclaimable_bytes gauge
node_memory_SReclaimable_bytes 3.1768576e+07
# HELP node_memory_SUnreclaim_bytes Memory information field SUnreclaim_bytes.
# TYPE node_memory_SUnreclaim_bytes gauge
node_memory_SUnreclaim_bytes 2.4133632e+07
# HELP node_memory_Shmem_bytes Memory information field Shmem_bytes.
# TYPE node_memory_Shmem_bytes gauge
node_memory_Shmem_bytes 4.898816e+07
# HELP node_memory_Slab_bytes Memory information field Slab_bytes.
# TYPE node_memory_Slab_bytes gauge
node_memory_Slab_bytes 5.5902208e+07
# HELP node_memory_SwapCached_bytes Memory information field SwapCached_bytes.
# TYPE node_memory_SwapCached_bytes gauge
node_memory_SwapCached_bytes 98304
# HELP node_memory_SwapFree_bytes Memory information field SwapFree_bytes.
# TYPE node_memory_SwapFree_bytes gauge
node_memory_SwapFree_bytes 1.006592e+08
# HELP node_memory_SwapTotal_bytes Memory information field SwapTotal_bytes.
# TYPE node_memory_SwapTotal_bytes gauge
node_memory_SwapTotal_bytes 1.04853504e+08
# HELP node_memory_Unevictable_bytes Memory information field Unevictable_bytes.
# TYPE node_memory_Unevictable_bytes gauge
node_memory_Unevictable_bytes 12288
# HELP node_memory_VmallocChunk_bytes Memory information field VmallocChunk_bytes.
# TYPE node_memory_VmallocChunk_bytes gauge
node_memory_VmallocChunk_bytes 0
# HELP node_memory_VmallocTotal_bytes Memory information field VmallocTotal_bytes.
# TYPE node_memory_VmallocTotal_bytes gauge
node_memory_VmallocTotal_bytes 1.140850688e+09
# HELP node_memory_VmallocUsed_bytes Memory information field VmallocUsed_bytes.
# TYPE node_memory_VmallocUsed_bytes gauge
node_memory_VmallocUsed_bytes 0
# HELP node_memory_WritebackTmp_bytes Memory information field WritebackTmp_bytes.
# TYPE node_memory_WritebackTmp_bytes gauge
node_memory_WritebackTmp_bytes 0
# HELP node_memory_Writeback_bytes Memory information field Writeback_bytes.
# TYPE node_memory_Writeback_bytes gauge
node_memory_Writeback_bytes 0
# HELP node_network_receive_bytes_total Network device statistic receive_bytes.
# TYPE node_network_receive_bytes_total counter
node_network_receive_bytes_total{device="docker0"} 4509
node_network_receive_bytes_total{device="eth0"} 0
node_network_receive_bytes_total{device="lo"} 1545
node_network_receive_bytes_total{device="veth110d75c"} 2172
node_network_receive_bytes_total{device="vethd7cda3f"} 2771
node_network_receive_bytes_total{device="wlan0"} 4.080929985e+09
# HELP node_network_receive_compressed_total Network device statistic receive_compressed.
# TYPE node_network_receive_compressed_total counter
node_network_receive_compressed_total{device="docker0"} 0
node_network_receive_compressed_total{device="eth0"} 0
node_network_receive_compressed_total{device="lo"} 0
node_network_receive_compressed_total{device="veth110d75c"} 0
node_network_receive_compressed_total{device="vethd7cda3f"} 0
node_network_receive_compressed_total{device="wlan0"} 0
# HELP node_network_receive_drop_total Network device statistic receive_drop.
# TYPE node_network_receive_drop_total counter
node_network_receive_drop_total{device="docker0"} 0
node_network_receive_drop_total{device="eth0"} 0
node_network_receive_drop_total{device="lo"} 0
node_network_receive_drop_total{device="veth110d75c"} 0
node_network_receive_drop_total{device="vethd7cda3f"} 0
node_network_receive_drop_total{device="wlan0"} 0
# HELP node_network_receive_errs_total Network device statistic receive_errs.
# TYPE node_network_receive_errs_total counter
node_network_receive_errs_total{device="docker0"} 0
node_network_receive_errs_total{device="eth0"} 0
node_network_receive_errs_total{device="lo"} 0
node_network_receive_errs_total{device="veth110d75c"} 0
node_network_receive_errs_total{device="vethd7cda3f"} 0
node_network_receive_errs_total{device="wlan0"} 0
# HELP node_network_receive_fifo_total Network device statistic receive_fifo.
# TYPE node_network_receive_fifo_total counter
node_network_receive_fifo_total{device="docker0"} 0
node_network_receive_fifo_total{device="eth0"} 0
node_network_receive_fifo_total{device="lo"} 0
node_network_receive_fifo_total{device="veth110d75c"} 0
node_network_receive_fifo_total{device="vethd7cda3f"} 0
node_network_receive_fifo_total{device="wlan0"} 0
# HELP node_network_receive_frame_total Network device statistic receive_frame.
# TYPE node_network_receive_frame_total counter
node_network_receive_frame_total{device="docker0"} 0
node_network_receive_frame_total{device="eth0"} 0
node_network_receive_frame_total{device="lo"} 0
node_network_receive_frame_total{device="veth110d75c"} 0
node_network_receive_frame_total{device="vethd7cda3f"} 0
node_network_receive_frame_total{device="wlan0"} 0
# HELP node_network_receive_multicast_total Network device statistic receive_multicast.
# TYPE node_network_receive_multicast_total counter
node_network_receive_multicast_total{device="docker0"} 0
node_network_receive_multicast_total{device="eth0"} 0
node_network_receive_multicast_total{device="lo"} 0
node_network_receive_multicast_total{device="veth110d75c"} 0
node_network_receive_multicast_total{device="vethd7cda3f"} 0
node_network_receive_multicast_total{device="wlan0"} 853334
# HELP node_network_receive_packets_total Network device statistic receive_packets.
# TYPE node_network_receive_packets_total counter
node_network_receive_packets_total{device="docker0"} 31
node_network_receive_packets_total{device="eth0"} 0
node_network_receive_packets_total{device="lo"} 11
node_network_receive_packets_total{device="veth110d75c"} 14
node_network_receive_packets_total{device="vethd7cda3f"} 17
node_network_receive_packets_total{device="wlan0"} 1.9557287e+07
# HELP node_network_transmit_bytes_total Network device statistic transmit_bytes.
# TYPE node_network_transmit_bytes_total counter
node_network_transmit_bytes_total{device="docker0"} 3.481729e+06
node_network_transmit_bytes_total{device="eth0"} 0
node_network_transmit_bytes_total{device="lo"} 1545
node_network_transmit_bytes_total{device="veth110d75c"} 6.965802e+06
node_network_transmit_bytes_total{device="vethd7cda3f"} 6.965646e+06
node_network_transmit_bytes_total{device="wlan0"} 2.549999149e+09
# HELP node_network_transmit_carrier_total Network device statistic transmit_carrier.
# TYPE node_network_transmit_carrier_total counter
node_network_transmit_carrier_total{device="docker0"} 0
node_network_transmit_carrier_total{device="eth0"} 0
node_network_transmit_carrier_total{device="lo"} 0
node_network_transmit_carrier_total{device="veth110d75c"} 0
node_network_transmit_carrier_total{device="vethd7cda3f"} 0
node_network_transmit_carrier_total{device="wlan0"} 0
# HELP node_network_transmit_colls_total Network device statistic transmit_colls.
# TYPE node_network_transmit_colls_total counter
node_network_transmit_colls_total{device="docker0"} 0
node_network_transmit_colls_total{device="eth0"} 0
node_network_transmit_colls_total{device="lo"} 0
node_network_transmit_colls_total{device="veth110d75c"} 0
node_network_transmit_colls_total{device="vethd7cda3f"} 0
node_network_transmit_colls_total{device="wlan0"} 0
# HELP node_network_transmit_compressed_total Network device statistic transmit_compressed.
# TYPE node_network_transmit_compressed_total counter
node_network_transmit_compressed_total{device="docker0"} 0
node_network_transmit_compressed_total{device="eth0"} 0
node_network_transmit_compressed_total{device="lo"} 0
node_network_transmit_compressed_total{device="veth110d75c"} 0
node_network_transmit_compressed_total{device="vethd7cda3f"} 0
node_network_transmit_compressed_total{device="wlan0"} 0
# HELP node_network_transmit_drop_total Network device statistic transmit_drop.
# TYPE node_network_transmit_drop_total counter
node_network_transmit_drop_total{device="docker0"} 0
node_network_transmit_drop_total{device="eth0"} 0
node_network_transmit_drop_total{device="lo"} 0
node_network_transmit_drop_total{device="veth110d75c"} 0
node_network_transmit_drop_total{device="vethd7cda3f"} 0
node_network_transmit_drop_total{device="wlan0"} 0
# HELP node_network_transmit_errs_total Network device statistic transmit_errs.
# TYPE node_network_transmit_errs_total counter
node_network_transmit_errs_total{device="docker0"} 0
node_network_transmit_errs_total{device="eth0"} 0
node_network_transmit_errs_total{device="lo"} 0
node_network_transmit_errs_total{device="veth110d75c"} 0
node_network_transmit_errs_total{device="vethd7cda3f"} 0
node_network_transmit_errs_total{device="wlan0"} 0
# HELP node_network_transmit_fifo_total Network device statistic transmit_fifo.
# TYPE node_network_transmit_fifo_total counter
node_network_transmit_fifo_total{device="docker0"} 0
node_network_transmit_fifo_total{device="eth0"} 0
node_network_transmit_fifo_total{device="lo"} 0
node_network_transmit_fifo_total{device="veth110d75c"} 0
node_network_transmit_fifo_total{device="vethd7cda3f"} 0
node_network_transmit_fifo_total{device="wlan0"} 0
# HELP node_network_transmit_packets_total Network device statistic transmit_packets.
# TYPE node_network_transmit_packets_total counter
node_network_transmit_packets_total{device="docker0"} 9071
node_network_transmit_packets_total{device="eth0"} 0
node_network_transmit_packets_total{device="lo"} 11
node_network_transmit_packets_total{device="veth110d75c"} 18149
node_network_transmit_packets_total{device="vethd7cda3f"} 18153
node_network_transmit_packets_total{device="wlan0"} 1.9709704e+07
# HELP node_scrape_collector_duration_seconds node_exporter: Duration of a collector scrape.
# TYPE node_scrape_collector_duration_seconds gauge
node_scrape_collector_duration_seconds{collector="cpu"} 0.00428738
node_scrape_collector_duration_seconds{collector="diskstats"} 0.004817796
node_scrape_collector_duration_seconds{collector="filesystem"} 0.0028265
node_scrape_collector_duration_seconds{collector="loadavg"} 0.000554425
node_scrape_collector_duration_seconds{collector="meminfo"} 0.002967177
node_scrape_collector_duration_seconds{collector="netdev"} 0.007897889
node_scrape_collector_duration_seconds{collector="textfile"} 0.000191509
node_scrape_collector_duration_seconds{collector="time"} 5.099e-05
# HELP node_scrape_collector_success node_exporter: Whether a collector succeeded.
# TYPE node_scrape_collector_success gauge
node_scrape_collector_success{collector="cpu"} 1
node_scrape_collector_success{collector="diskstats"} 1
node_scrape_collector_success{collector="filesystem"} 1
node_scrape_collector_success{collector="loadavg"} 1
node_scrape_collector_success{collector="meminfo"} 1
node_scrape_collector_success{collector="netdev"} 1
node_scrape_collector_success{collector="textfile"} 1
node_scrape_collector_success{collector="time"} 1
# HELP node_textfile_scrape_error 1 if there was an error opening or reading a file, 0 otherwise
# TYPE node_textfile_scrape_error gauge
node_textfile_scrape_error 0
# HELP node_time_seconds System time in seconds since epoch (1970).
# TYPE node_time_seconds gauge
node_time_seconds 1.577749661874304e+09
# HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
# TYPE promhttp_metric_handler_errors_total counter
promhttp_metric_handler_errors_total{cause="encoding"} 0
promhttp_metric_handler_errors_total{cause="gathering"} 0
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 2
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0

`
