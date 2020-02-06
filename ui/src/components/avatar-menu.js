import React, { useState } from 'react';
import styled from 'styled-components';
import { useNavigation, useActive, useCurrentRoute } from 'react-navi';

import api from '../api';
import { Row, Text, Icon, MenuItem } from './core';
import Popup from './popup';
import Popover from './popover';
import Avatar from './avatar';
import CliDownload from '../containers/cli-download';
import ChangePassword from '../containers/change-password';
import Profile from '../containers/profile';
import UserAccessKeys from '../containers/user-access-keys';
import SSHKeys from '../containers/ssh-keys';

const Divider = styled.div`
  width: 100%;
  border-bottom: 1px solid ${props => props.theme.colors.grays[5]};
  margin: 8px 0;
`;

const AvatarMenu = () => {
  const {
    data: { context },
  } = useCurrentRoute();
  const [showCLI, setShowCLI] = useState();
  const [showUserProfile, setShowUserProfile] = useState();
  const [showUserAccessKeys, setShowUserAccessKeys] = useState();
  const [showChangePassword, setShowChangePassword] = useState();
  const [showSSHKeys, setShowSSHKeys] = useState();
  const navigation = useNavigation();
  const isProjectsRoute = useActive('/projects');
  const name = `${context.currentUser.firstName} ${context.currentUser.lastName}`;

  return (
    <>
      <Popup show={showCLI} onClose={() => setShowCLI(false)}>
        <CliDownload />
      </Popup>
      <Popup show={showUserProfile} onClose={() => setShowUserProfile(false)}>
        <Profile
          user={context.currentUser}
          close={() => setShowUserProfile(false)}
        />
      </Popup>
      <Popup
        show={showUserAccessKeys}
        onClose={() => setShowUserAccessKeys(false)}
      >
        <UserAccessKeys user={context.currentUser} />
      </Popup>
      <Popup show={showSSHKeys} onClose={() => setShowSSHKeys(false)}>
        <SSHKeys user={context.currentUser} />
      </Popup>
      <Popup
        show={showChangePassword}
        onClose={() => setShowChangePassword(false)}
      >
        <ChangePassword
          user={context.currentUser}
          close={() => setShowChangePassword(false)}
        />
      </Popup>
      <Popover
        top="46px"
        right={0}
        width="240px"
        button={({ show }) => (
          <Avatar
            name={name}
            color={show ? 'primary' : 'white'}
            borderColor={show ? 'primary' : 'white'}
          />
        )}
        content={({ close }) => (
          <>
            <Text
              fontSize={3}
              fontWeight={2}
              paddingX={3}
              marginX={1}
              paddingTop={2}
            >
              {name}
            </Text>
            <Text
              fontSize={1}
              marginBottom={1}
              paddingX={3}
              marginX={1}
              color="grays.8"
            >
              {context.currentUser.email}
            </Text>
            <Divider />
            {!isProjectsRoute && (
              <MenuItem onClick={() => navigation.navigate('/projects')}>
                Projects
              </MenuItem>
            )}
            <Divider />
            <MenuItem
              onClick={() => {
                close();
                setShowUserProfile(true);
              }}
            >
              Profile
            </MenuItem>
            <MenuItem
              onClick={() => {
                close();
                setShowChangePassword(true);
              }}
            >
              Change Password
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={() => {
                close();
                setShowUserAccessKeys(true);
              }}
            >
              Access Keys
            </MenuItem>
            <MenuItem
              onClick={() => {
                close();
                setShowSSHKeys(true);
              }}
            >
              SSH Keys
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={() => {
                close();
                setShowCLI(true);
              }}
            >
              Download CLI
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={async () => {
                context.setCurrentUser(null);
                await api.logout();
                navigation.navigate('/login');
              }}
              paddingBottom={3}
              marginBottom={2}
            >
              Logout
            </MenuItem>
          </>
        )}
      />
    </>
  );
};

export default AvatarMenu;
