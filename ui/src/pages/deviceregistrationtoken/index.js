import React, { Component, Fragment } from 'react';
import { Pane, majorScale, Tab, Tablist } from 'evergreen-ui';
import axios from 'axios';

import config from '../../config';
import TopHeader from '../../components/TopHeader';
import CustomSpinner from '../../components/CustomSpinner';
import Settings from './settings';
import Overview from './overview';

export default class DeviceRegistrationToken extends Component {
  state = {
    deviceRegistrationToken: null,
    tabs: ['overview', 'settings'],
    tabLabels: ['Overview', 'Settings']
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/${this.props.tokenName}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          deviceRegistrationToken: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    const tokenName = match.params.tokenName;
    var selectedIndex = 0;
    switch (match.params.appTab) {
      case 'overview':
        selectedIndex = 0;
        break;
      case 'settings':
        selectedIndex = 1;
        break;
      default:
        this.props.history.push(
          `/${projectName}/provisioning/deviceregistrationtokens/${tokenName}`
        );
    }
    return (
      <Tablist border="default">
        {this.state.tabs.map((tab, index) => (
          <Tab
            key={tab}
            id={tab}
            onSelect={() =>
              this.props.history.push(
                `/${projectName}/provisioning/deviceregistrationtokens/${tokenName}/${tab}`
              )
            }
            isSelected={index === selectedIndex}
          >
            {this.state.tabLabels[index]}
          </Tab>
        ))}
      </Tablist>
    );
  };

  renderInner = match => {
    switch (match.params.appTab) {
      case 'overview':
        return (
          <Overview
            projectName={match.params.projectName}
            deviceRegistrationToken={this.state.deviceRegistrationToken}
            history={this.props.history}
          />
        );
      case 'settings':
        return (
          <Settings
            projectName={match.params.projectName}
            deviceRegistrationToken={this.state.deviceRegistrationToken}
            history={this.props.history}
          />
        );
      default:
        return <Pane></Pane>;
      }
  };

  render() {
    const deviceRegistrationToken = this.state.deviceRegistrationToken;
    const heading = 'Device Registration Token / ' + this.props.tokenName;
    return (
      <Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        {deviceRegistrationToken ? (
          <Fragment>
            <Pane
              display="flex"
              flexDirection="column"
              alignItems="center"
              background="white"
              width="100%"
              padding={majorScale(1)}
              borderBottom="default"
            >
              {this.renderTablist(this.props.applicationRoute)}
            </Pane>
            {this.renderInner(this.props.applicationRoute)}
          </Fragment>
        ) : (
          <CustomSpinner />
        )}
      </Fragment>
    );
  }
}
