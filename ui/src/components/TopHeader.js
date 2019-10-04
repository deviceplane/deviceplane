import React, { Component } from 'react';
import './../App.css';
import UserAvatar from './UserAvatar.js';
import { Pane, majorScale, Heading } from 'evergreen-ui';

import logo from '../assets/logo.png';

export default class TopHeader extends Component {
  render() {
    return (
      <Pane
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        width="100%"
        background="white"
        borderBottom="default"
        padding={majorScale(1)}
      >
        <Pane flex={1}>
          {this.props.showLogo && (
            <img src={logo} alt="Logo" height="30px" width="35px" />
          )}
        </Pane>
        <Pane display="flex" justifyContent="center" flex={1}>
          <Heading size={500}>{this.props.heading}</Heading>
        </Pane>
        <Pane
          display="flex"
          justifyContent="flex-end"
          role="button"
          cursor="pointer"
          flex={1}
        >
          <UserAvatar
            user={this.props.user}
            history={this.props.history}
            hideSwitchProjects={this.props.hideSwitchProjects}
          />
        </Pane>
      </Pane>
    );
  }
}
