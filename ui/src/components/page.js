import React from 'react';
import { NotFoundBoundary, useCurrentRoute } from 'react-navi';
import { createGlobalStyle } from 'styled-components';
import { Helmet } from 'react-helmet-async';

import { Box } from './core';
import NotFound from './not-found';

const GlobalStyle = createGlobalStyle`
  * {
    box-sizing: inherit;
  }

  ::selection {
    background-color: rgba(255,255,255,.99);
    color: #000;
  }

  html {
    box-sizing: border-box;
    font-family: ${props => props.theme.fonts.default};
    font-size: 16px;
    font-weight: 400;
    line-height: 1.2;
    background-color: ${props => props.theme.colors.pageBackground};
  }

  html, body {
    text-rendering: optimizeLegibility
    -webkit-font-smoothing: antialiased;
  }

  body {
    margin: 0;
  }

  html, body, #root, #root > div, main {
    height: 100%;
  }

  .ace_editor {
    background: ${props => props.theme.colors.grays[0]} !important;
    color: ${props => props.theme.colors.white} !important;
  }
  .ace_gutter {
    background: ${props => props.theme.colors.grays[1]} !important;
    color: ${props => props.theme.colors.white} !important;
  }
  .ace_gutter-active-line {
    background: ${props => props.theme.colors.grays[2]} !important;
  }
  .ace_active-line {
    background: ${props => props.theme.colors.grays[2]} !important;
  }

  svg[data-icon="caret-down"] {
    fill: ${props => props.theme.colors.white} !important;
  }

  div[data-evergreen-toaster-container] {
    position: relative;
    z-index: 99999999999;
  }
`;

const Page = ({ children }) => {
  const route = useCurrentRoute();
  return (
    <>
      <Helmet>
        {route.title && <title>{`${route.title} - Deviceplane`}</title>}
        <link
          href="https://fonts.googleapis.com/css?family=Rubik:300,400,500,700&display=swap"
          rel="stylesheet"
        />
      </Helmet>
      <GlobalStyle />
      <Box>
        <main>
          <NotFoundBoundary
            render={() => {
              if (route.data.context.currentUser) {
                return <NotFound />;
              } else {
                window.location.replace('/login');
                return null;
              }
            }}
          >
            {children}
          </NotFoundBoundary>
        </main>
      </Box>
    </>
  );
};

export default Page;
