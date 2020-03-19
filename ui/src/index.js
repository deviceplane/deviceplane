import React, { Suspense, useEffect, useState } from 'react';
import ReactDOM from 'react-dom';
import { Router, View } from 'react-navi';
import { HelmetProvider } from 'react-helmet-async';
import { ThemeProvider } from 'styled-components';
import bugsnag from '@bugsnag/js';
import bugsnagReact from '@bugsnag/plugin-react';

import routes from './routes';
import api from './api';
import theme from './theme';
import Page from './components/page';
import Spinner from './components/spinner';
import { ToastManager } from './components/core';
import { LoadIntercom, bootIntercom } from './lib/intercom';
import segment, { LoadSegment } from './lib/segment';

const App = () => {
  const [loaded, setLoaded] = useState();
  const [currentUser, setCurrentUser] = useState();

  const load = async () => {
    try {
      const { data: user } = await api.user();
      segment.identify(user.id, {
        name: user.name,
        email: user.email,
      });
      setCurrentUser(user);
      bootIntercom(user);
    } catch (error) {
      console.error(error);
    }
    setLoaded(true);
  };

  useEffect(() => {
    load();
  }, []);

  if (!loaded) {
    return null;
  }

  return (
    <HelmetProvider>
      <Router routes={routes} context={{ currentUser, setCurrentUser }}>
        <ThemeProvider theme={theme}>
          <Page>
            <Suspense fallback={<Spinner />}>
              <View />
            </Suspense>
            <ToastManager />
          </Page>
        </ThemeProvider>
      </Router>
    </HelmetProvider>
  );
};

if (process.env.BUGSNAG_KEY) {
  const bugsnagClient = bugsnag(process.env.BUGSNAG_KEY);
  bugsnagClient.use(bugsnagReact, React);
  const ErrorBoundary = bugsnagClient.getPlugin('react');

  ReactDOM.render(
    <ErrorBoundary>
      <LoadSegment />
      <LoadIntercom />
      <App />
    </ErrorBoundary>,
    document.querySelector('#root')
  );
} else {
  ReactDOM.render(
    <>
      <LoadSegment />
      <LoadIntercom />
      <App />
    </>,
    document.querySelector('#root')
  );
}
