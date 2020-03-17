import React from 'react';

import config from '../config';
import Card from '../components/card';
import Linux from '../components/icons/linux';
import Windows from '../components/icons/windows';
import Apple from '../components/icons/apple';
import { Row, Column, Button } from '../components/core';

const CliDownload = () => (
  <Card title="Download CLI" border size="large">
    <Row justifyContent="space-between">
      <Column width="80px" height="80px" alignItems="center">
        <Apple />
        <Row marginTop={4}>
          <Button
            variant="textPrimary"
            color="primary"
            title="Mac OS X"
            href={`${config.cliEndpoint}/latest/darwin/amd64/deviceplane`}
          />
        </Row>
      </Column>
      <Column width="80px" height="80px" alignItems="center">
        <Linux />
        <Row marginTop={4} justifyContent="center">
          <Button
            variant="textPrimary"
            color="primary"
            title="AMD64"
            marginRight={2}
            href={`${config.cliEndpoint}/latest/linux/amd64/deviceplane`}
          />
          <Button
            variant="textPrimary"
            color="primary"
            title="ARM"
            marginRight={2}
            href={`${config.cliEndpoint}/latest/linux/arm/deviceplane`}
          />
          <Button
            variant="textPrimary"
            color="primary"
            title="ARM64"
            href={`${config.cliEndpoint}/latest/linux/arm64/deviceplane`}
          />
        </Row>
      </Column>
      <Column width="80px" height="80px" alignItems="center">
        <Windows />
        <Row marginTop={4}>
          <Button
            variant="textPrimary"
            color="primary"
            title="Windows"
            href={`${config.cliEndpoint}/latest/windows/amd64/deviceplane.exe`}
          />
        </Row>
      </Column>
    </Row>
  </Card>
);

export default CliDownload;
