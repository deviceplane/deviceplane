import React, { Component, Fragment } from 'react';
import { Pane, majorScale, Tab, Tablist } from 'evergreen-ui';
import axios from 'axios';

import config from '../../config';
import TopHeader from '../../components/TopHeader';
import CustomSpinner from '../../components/CustomSpinner';
import Overview from './overview';
import Scheduling from './scheduling';
import Settings from './settings';

export default class Application extends Component {
  state = {
    application: null,
    tabs: ['overview', 'scheduling', 'settings'],
    tabLabels: ['Overview', 'Scheduling', 'Settings']
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.applicationName}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          application: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    const applicationName = match.params.applicationName;
    var selectedIndex = 0;
    switch (match.params.appTab) {
      case 'overview':
        selectedIndex = 0;
        break;
      case 'scheduling':
        selectedIndex = 1;
        break;
      case 'settings':
        selectedIndex = 2;
        break;
      default:
        this.props.history.push(
          `/${projectName}/applicationName/${applicationName}`
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
                `/${projectName}/applications/${applicationName}/${tab}`
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
            application={this.state.application}
            history={this.props.history}
          />
        );
      case 'scheduling':
        return (
          <Scheduling
            projectName={match.params.projectName}
            application={this.state.application}
            history={this.props.history}
          />
        );
      case 'settings':
        return (
          <Settings
            projectName={match.params.projectName}
            application={this.state.application}
            history={this.props.history}
          />
        );
      default:
        return <Pane></Pane>;
    }
  };

  render() {
    const application = this.state.application;
    const heading = 'Application / ' + this.props.applicationName;
    return (
      <Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        {application ? (
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
