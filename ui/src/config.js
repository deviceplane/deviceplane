const development = {
  endpoint: 'http://localhost:8080/api',
  wsEndpoint: 'ws://localhost:8080/api',
};

const endpointBase = window.location.port
  ? `${window.location.hostname}:${window.location.port}/api`
  : `${window.location.hostname}/api`;

const production = {
  endpoint: `${window.location.protocol}//${endpointBase}`,
  wsEndpoint:
    window.location.protocol === 'http:'
      ? `ws://${endpointBase}`
      : `wss://${endpointBase}`,
};

const metricProperties = {
  device: {
    id: 'device',
    label: 'Device',
    description: 'Includes a Datadog tag with the device name.',
  },
};

const supportedDeviceMetricProperties = [metricProperties.device];

const supportedServiceMetricProperties = [metricProperties.device];

const supportedDeviceMetrics = [
  'cpu_guest_seconds_total',
  'cpu_seconds_total',
  'disk_io_now',
  'disk_io_time_seconds_total',
  'disk_io_time_weighted_seconds_total',
  'disk_read_bytes_total',
  'disk_read_time_seconds_total',
  'disk_reads_completed_total',
  'disk_reads_merged_total',
  'disk_write_time_seconds_total',
  'disk_writes_completed_total',
  'disk_writes_merged_total',
  'disk_written_bytes_total',
  'filesystem_avail_bytes',
  'filesystem_device_error',
  'filesystem_files',
  'filesystem_files_free',
  'filesystem_free_bytes',
  'filesystem_readonly',
  'filesystem_size_bytes',
  'load',
  'memory_Active_anon_bytes',
  'memory_Active_bytes',
  'memory_Active_file_bytes',
  'memory_AnonPages_bytes',
  'memory_Bounce_bytes',
  'memory_Buffers_bytes',
  'memory_Cached_bytes',
  'memory_CmaFree_bytes',
  'memory_CmaTotal_bytes',
  'memory_CommitLimit_bytes',
  'memory_Committed_AS_bytes',
  'memory_Dirty_bytes',
  'memory_Inactive_anon_bytes',
  'memory_Inactive_bytes',
  'memory_Inactive_file_bytes',
  'memory_KernelStack_bytes',
  'memory_Mapped_bytes',
  'memory_MemAvailable_bytes',
  'memory_MemFree_bytes',
  'memory_MemTotal_bytes',
  'memory_Mlocked_bytes',
  'memory_NFS_Unstable_bytes',
  'memory_PageTables_bytes',
  'memory_SReclaimable_bytes',
  'memory_SUnreclaim_bytes',
  'memory_Shmem_bytes',
  'memory_Slab_bytes',
  'memory_SwapCached_bytes',
  'memory_SwapFree_bytes',
  'memory_SwapTotal_bytes',
  'memory_Unevictable_bytes',
  'memory_VmallocChunk_bytes',
  'memory_VmallocTotal_bytes',
  'memory_VmallocUsed_bytes',
  'memory_WritebackTmp_bytes',
  'memory_Writeback_bytes',
  'network_receive_bytes_total',
  'network_receive_compressed_total',
  'network_receive_drop_total',
  'network_receive_errs_total',
  'network_receive_fifo_total',
  'network_receive_frame_total',
  'network_receive_multicast_total',
  'network_receive_packets_total',
  'network_transmit_bytes_total',
  'network_transmit_carrier_total',
  'network_transmit_colls_total',
  'network_transmit_compressed_total',
  'network_transmit_drop_total',
  'network_transmit_errs_total',
  'network_transmit_fifo_total',
  'network_transmit_packets_total',
  'scrape_collector_duration_seconds',
  'scrape_collector_success',
  'textfile_scrape_error',
  'time_seconds',
  'promhttp_metric_handler_errors_total',
  'promhttp_metric_handler_requests_in_flight ',
  'promhttp_metric_handler_requests_total',
];

const config =
  process.env.REACT_APP_ENVIRONMENT === 'development'
    ? development
    : production;

export default {
  agentVersion: '1.9.1',
  cliEndpoint: 'https://cli.deviceplane.com',
  supportedDeviceMetrics,
  supportedDeviceMetricProperties,
  supportedServiceMetricProperties,
  ...config,
};
