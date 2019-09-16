import React, { Component } from 'react';
import './../App.css';
import config from '../config.js';
import { Pane, majorScale, Heading, Button } from 'evergreen-ui';

export default class CliDownload extends Component {
  render() {
    return (
      <React.Fragment>
        <Pane zIndex={1} flexShrink={0} elevation={0} backgroundColor="white">
          <Pane padding={majorScale(2)}>
            <Heading size={600}>Download CLI</Heading>
          </Pane>
        </Pane>
        <Pane
          display="flex"
          flexDirection="column"
          margin={majorScale(4)}
        >
          <Button marginBottom={majorScale(3)} justifyContent="center" is="a" href={`${config.cliEndpoint}/latest/darwin/amd64/deviceplane`}>MacOS</Button>
          <Button marginBottom={majorScale(3)} justifyContent="center" is="a" href={`${config.cliEndpoint}/latest/windows/amd64/deviceplane.exe`}>Windows</Button>
          <Button marginBottom={majorScale(3)} justifyContent="center" is="a" href={`${config.cliEndpoint}/latest/linux/amd64/deviceplane`}>Linux AMD64</Button>
          <Button marginBottom={majorScale(3)} justifyContent="center" is="a" href={`${config.cliEndpoint}/latest/linux/arm/deviceplane`}>Linux ARM</Button>
          <Button marginBottom={majorScale(3)} justifyContent="center" is="a" href={`${config.cliEndpoint}/latest/linux/arm64/deviceplane`}>Linux ARM64</Button>
        </Pane>
      </React.Fragment>
    );
  }
}