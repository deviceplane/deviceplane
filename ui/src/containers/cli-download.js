import React from 'react';

import config from '../config';
import Card from '../components/card';
import { Button } from '../components/core';

const CliDownload = () => (
  <Card title="Download CLI" border size="large">
    <Button
      title="MacOS"
      marginBottom={4}
      href={`${config.cliEndpoint}/latest/darwin/amd64/deviceplane`}
    />
    <Button
      title="Windows"
      marginBottom={4}
      href={`${config.cliEndpoint}/latest/windows/amd64/deviceplane.exe`}
    />
    <Button
      title="Linux AMD64"
      marginBottom={4}
      href={`${config.cliEndpoint}/latest/linux/amd64/deviceplane`}
    />
    <Button
      title="Linux ARM"
      marginBottom={4}
      href={`${config.cliEndpoint}/latest/linux/arm/deviceplane`}
    />
    <Button
      title="Linux ARM64"
      href={`${config.cliEndpoint}/latest/linux/arm64/deviceplane`}
    />
  </Card>
);

export default CliDownload;
