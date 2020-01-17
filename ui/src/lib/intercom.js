import React, { useEffect } from 'react';

export const LoadIntercom = () => {
  useEffect(() => {
    if (!process.env.INTERCOM_ID) {
      return;
    }

    (function() {
      var w = window;
      var ic = w.Intercom;
      if (typeof ic === 'function') {
        ic('reattach_activator');
        ic('update', w.intercomSettings);
      } else {
        var d = document;
        var i = function() {
          i.c(arguments);
        };
        i.q = [];
        i.c = function(args) {
          i.q.push(args);
        };
        w.Intercom = i;
        var l = function() {
          var s = d.createElement('script');
          s.type = 'text/javascript';
          s.async = true;
          s.src =
            'https://widget.intercom.io/widget/' + '${process.env.INTERCOM_ID}';
          var x = d.getElementsByTagName('script')[0];
          x.parentNode.insertBefore(s, x);
        };
        if (w.attachEvent) {
          w.attachEvent('onload', l);
        } else {
          w.addEventListener('load', l, false);
        }
      }
    })();
    window.intercomSettings = {
      app_id: process.env.INTERCOM_ID,
    };
  }, []);

  return null;
};

export const bootIntercom = user => {
  if (window.Intercom) {
    window.Intercom('boot', {
      app_id: process.env.INTERCOM_ID,
      name: `${user.firstName} ${user.lastName}`,
      email: user.email,
    });
  }
};
