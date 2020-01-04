import utils from '../utils';

export const getMetricLabel = metric =>
  metric
    .split('_')
    .filter(s => s !== 'promhttp')
    .map(utils.capitalize)
    .map(s => (s === 'Io' ? 'IO' : s === 'Cpu' ? 'CPU' : s))
    .join(' ');
