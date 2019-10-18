import React, { Component, Fragment } from 'react';
import {
  Pane,
  majorScale,
  Heading,
  Card,
  Label
} from 'evergreen-ui';

import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';

export default class DeviceRegistrationTokenOverview extends Component {
  render() {
    var { deviceRegistrationToken } = this.props;
    var cards = [
      {
        title: "Name",
        value: deviceRegistrationToken.name
      },
      {
        title: "Description",
        value: deviceRegistrationToken.description
      },
      {
        title: "Devices Registered",
        value: deviceRegistrationToken.deviceCounts.allCount
      },
      {
        title: "Maximum Device Registrations",
        value: deviceRegistrationToken.maxRegistrations
      }
    ];
    const cardNodes = cards.map(card => (
      <InnerCard>
        <Heading paddingTop={majorScale(2)} paddingLeft={majorScale(2)}>
          {card.title}
        </Heading>
        <Card
          display="flex"
          flexDirection="column"
          alignItems="left"
          width="80%"
          padding={majorScale(2)}
        >
          <Label>{card.value}</Label>
        </Card>
      </InnerCard>
    ));

    return (
      <Pane width="70%" display="flex" flexDirection="column">
        {deviceRegistrationToken ? (
          <Fragment>
            {cardNodes}
          </Fragment>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}