import React, { Suspense, useEffect, useState } from 'react';
import ReactDOM from 'react-dom';
import { Router, View } from 'react-navi';
import { HelmetProvider } from 'react-helmet-async';
import { ThemeProvider } from 'styled-components';
import bugsnag from '@bugsnag/js';
import bugsnagReact from '@bugsnag/plugin-react';

import routes from './routes';
import api from './api';
import * as serviceWorker from './serviceWorker';
import theme from './theme';
import Page from './components/page';
import Spinner from './components/spinner';
import Intercom from './lib/intercom';

const App = () => {
  const [loaded, setLoaded] = useState();
  const [currentUser, setCurrentUser] = useState();

  const load = async () => {
    try {
      const { data: user } = await api.user();
      setCurrentUser(user);
      if (process.env.NODE_ENV !== 'development') {
        window.Intercom('boot', {
          app_id: process.env.REACT_APP_INTERCOM_ID,
          name: `${user.firstName} ${user.lastName}`,
          email: user.email,
        });
      }
    } catch (error) {
      console.log(error);
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
          </Page>
        </ThemeProvider>
      </Router>
    </HelmetProvider>
  );
};

if (process.env.NODE_ENV === 'development') {
  ReactDOM.render(<App />, document.getElementById('root'));
} else {
  const bugsnagClient = bugsnag(process.env.REACT_APP_BUGSNAG_KEY);
  bugsnagClient.use(bugsnagReact, React);
  const ErrorBoundary = bugsnagClient.getPlugin('react');

  ReactDOM.render(
    <ErrorBoundary>
      <App />
      <Intercom />
    </ErrorBoundary>,
    document.getElementById('root')
  );
}

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
