import React, { Component } from 'react';
import './../App.css';
import config from '../config.js';
import ChangePassword from './ChangePassword.js';
import CliDownload from './CliDownload.js';
import EditProfile from './EditProfile.js';
import UserAccessKeys from './UserAccessKeys.js';
import { Menu, Pane, majorScale, Heading, Text, SideSheet, Popover, Avatar, minorScale, Icon } from 'evergreen-ui';
import axios from 'axios';

export default class UserAvatar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      editProfileShown: false,
      cliDownloadShown: false,
      changePasswordShown: false,
      userAccessKeysShown: false
    };
  }

  handleLogout() {
    axios.post(`${config.endpoint}/logout`, null, {
      withCredentials: true
    })
    .then((response) => {
      window.location.reload();
    })
    .catch((error) => {
      console.log(error);
    });
  }

  renderInner() {
    const user = this.props.user;
    const name = `${user.firstName} ${user.lastName}`;
    return (
      <React.Fragment>
        <Menu>
          <Menu.Group>
            <Pane padding={majorScale(2)}>
              <Heading size={600}>{name}</Heading>
              <Text size={500}>{user.email}</Text>
            </Pane>
          </Menu.Group>
          <Menu.Divider />
          <Menu.Group>
            <Menu.Item onSelect={() => this.setState({ editProfileShown: true })}>Edit Profile</Menu.Item>
            <Menu.Item onSelect={() => this.setState({ changePasswordShown: true })}>Change Password</Menu.Item>
          </Menu.Group>
          <Menu.Divider />
          <Menu.Group>
            <Menu.Item onSelect={() => this.setState({ userAccessKeysShown: true })}>Manage Access Keys</Menu.Item>
            <Menu.Item onSelect={() => this.setState({ cliDownloadShown: true })}>Download CLI</Menu.Item>
          </Menu.Group>
          <Menu.Divider />
          <Menu.Group>
            {this.props.hideSwitchProjects ? (
              <React.Fragment></React.Fragment>
            ) : (
              <Menu.Item onSelect={() => this.props.history.push('/projects')}>Switch Project</Menu.Item>
            )}
            <Menu.Item onSelect={() => this.handleLogout()}>Logout</Menu.Item>
          </Menu.Group>
        </Menu>
      </React.Fragment>
    );
  }

  render() {
    const user = this.props.user;
    const name = `${user.firstName} ${user.lastName}`;
    return (
      <React.Fragment>
        <SideSheet
          isShown={this.state.cliDownloadShown}
          onCloseComplete={() => this.setState({ cliDownloadShown: false })}
          containerProps={{
            display: 'flex',
            flex: '1',
            flexDirection: 'column',
          }}
        >
          <CliDownload />
        </SideSheet>
        <SideSheet
          isShown={this.state.editProfileShown}
          onCloseComplete={() => this.setState({ editProfileShown: false })}
        >
          <EditProfile user={user}/>
        </SideSheet>
        <SideSheet
          isShown={this.state.userAccessKeysShown}
          onCloseComplete={() => this.setState({ userAccessKeysShown: false })}
        >
          <UserAccessKeys user={user} />
        </SideSheet>
        <SideSheet
          isShown={this.state.changePasswordShown}
          onCloseComplete={() => this.setState({ changePasswordShown: false })}
        >
          <ChangePassword user={user} />
        </SideSheet>
        <Popover content={this.renderInner()}>
          <Pane display="flex" alignItems="center">
            <Avatar isSolid name={name} size={36} />
            <Icon
              icon="caret-down"
              marginLeft={minorScale(1)}
              size={minorScale(4)}
            />
          </Pane>
        </Popover>
      </React.Fragment>
    );
  }
}