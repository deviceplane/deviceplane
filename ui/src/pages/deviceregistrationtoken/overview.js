import React, { Component, Fragment } from 'react';
import {
  Pane,
  majorScale,
  Heading,
  Card,
  Label
} from 'evergreen-ui';

import config from '../../config';
import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';
import { EditableLabelTable } from '../../components/EditableLabelTable';

export default class DeviceRegistrationTokenOverview extends Component {
  render() {
    var { deviceRegistrationToken } = this.props;
    var cards = [
      {
        title: "Name",
        value: deviceRegistrationToken.name
      },
      {
        title: "ID",
        value: deviceRegistrationToken.id
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
      },
      {
        title: "Labels",
        innerElement: (
          <EditableLabelTable
            getEndpoint={`${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/${this.props.deviceRegistrationToken.id}`}
            setEndpoint={`${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/${this.props.deviceRegistrationToken.id}/labels`}
            deleteEndpoint={`${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/${this.props.deviceRegistrationToken.id}/labels`}
          />
        )
      }
    ];
    const cardNodes = cards.map(card => (
      <InnerCard key={card.title}>
        <Heading paddingTop={majorScale(2)} paddingLeft={majorScale(2)}>
          {card.title}
        </Heading>
        <Card
          display="flex"
          flexDirection="column"
          alignItems="left"
          padding={majorScale(2)}
        >
          {card.innerElement || <Label>{card.value}</Label>}
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