export default {
  identify: () => {
    window.analytics.identify();
  },
  track: () => {
    window.analytics.track();
  },
  page: () => {
    window.analytics.page();
  },
};
