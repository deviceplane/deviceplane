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

const bugsnagClient = bugsnag('2a74913ba2df5151edd8d25004bbc820');
bugsnagClient.use(bugsnagReact, React);
const ErrorBoundary = bugsnagClient.getPlugin('react');

const App = () => {
  const [loaded, setLoaded] = useState();
  const [currentUser, setCurrentUser] = useState();

  const load = async () => {
    try {
      const response = await api.user();
      setCurrentUser(response.data);
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

ReactDOM.render(
  <ErrorBoundary>
    <App />
  </ErrorBoundary>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
