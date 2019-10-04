import React, { Component, Fragment } from 'react';
import axios from 'axios';
import { Pane, majorScale, Tablist, Tab } from 'evergreen-ui';

import config from '../../config';
import CustomSpinner from '../../components/CustomSpinner';
import TopHeader from '../../components/TopHeader';
import Overview from './overview';
import Settings from './settings';
import Ssh from './ssh';

export default class Device extends Component {
  state = {
    device: null,
    tabs: ['overview', 'ssh', 'settings'],
    tabLabels: ['Overview', 'SSH', 'Settings']
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.deviceName}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          device: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    const deviceName = match.params.deviceName;
    var selectedIndex = 0;
    switch (match.params.appTab) {
      case 'overview':
        selectedIndex = 0;
        break;
      case 'ssh':
        selectedIndex = 1;
        break;
      case 'settings':
        selectedIndex = 2;
        break;
      default:
        this.props.history.push(`/${projectName}/devices/${deviceName}`);
    }
    return (
      <Tablist border="default">
        {this.state.tabs.map((tab, index) => (
          <Tab
            key={tab}
            id={tab}
            onSelect={() =>
              this.props.history.push(
                `/${projectName}/devices/${deviceName}/${tab}`
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
            device={this.state.device}
            history={this.props.history}
          />
        );
      case 'ssh':
        return (
          <Ssh
            projectName={match.params.projectName}
            device={this.state.device}
            history={this.props.history}
          />
        );
      case 'settings':
        return (
          <Settings
            projectName={match.params.projectName}
            device={this.state.device}
            history={this.props.history}
          />
        );
      default:
        return <Pane></Pane>;
    }
  };

  render() {
    const device = this.state.device;
    const heading = 'Device / ' + this.props.deviceName;
    return (
      <Fragment>
        <TopHeader user={this.props.user} heading={heading} />
        {device ? (
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
              {this.renderTablist(this.props.deviceRoute)}
            </Pane>
            {this.renderInner(this.props.deviceRoute)}
          </Fragment>
        ) : (
          <CustomSpinner />
        )}
      </Fragment>
    );
  }
}
