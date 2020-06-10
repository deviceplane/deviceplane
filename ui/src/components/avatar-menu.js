import React, { useState } from 'react';
import styled from 'styled-components';
import { useNavigation, useActive, useCurrentRoute } from 'react-navi';

import api from '../api';
import useToggle from '../hooks/useToggle';
import { Text, MenuItem } from './core';
import Popup from './popup';
import Popover from './popover';
import Avatar from './avatar';
import CliDownload from '../containers/cli-download';
import ChangePassword from '../containers/change-password';
import Profile from '../containers/profile';
import UserAccessKeys from '../containers/user-access-keys';

const Divider = styled.div`
  width: 100%;
  border-bottom: 1px solid ${props => props.theme.colors.grays[5]};
  margin: 8px 0;
`;

const AvatarMenu = () => {
  const {
    data: { context },
  } = useCurrentRoute();
  const [isCLI, toggleCLI] = useToggle();
  const [isUserProfile, toggleUserProfile] = useToggle();
  const [isUserAccessKeys, toggleUserAccessKeys] = useToggle();
  const [isChangePassword, toggleChangePassword] = useToggle();
  const navigation = useNavigation();
  const isProjectsRoute = useActive('/projects');
  const name = context.currentUser.name;

  return (
    <>
      <Popup show={isCLI} onClose={toggleCLI}>
        <CliDownload />
      </Popup>
      <Popup show={isUserProfile} onClose={toggleUserProfile}>
        <Profile user={context.currentUser} close={toggleUserProfile} />
      </Popup>
      <Popup show={isUserAccessKeys} onClose={toggleUserAccessKeys}>
        <UserAccessKeys user={context.currentUser} />
      </Popup>
      <Popup show={isChangePassword} onClose={toggleChangePassword}>
        <ChangePassword
          user={context.currentUser}
          close={toggleChangePassword}
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
            {!isProjectsRoute && (
              <>
                <Divider />
                <MenuItem onClick={() => navigation.navigate('/projects')}>
                  Projects
                </MenuItem>
              </>
            )}
            <Divider />
            <MenuItem
              onClick={() => {
                close();
                toggleUserProfile();
              }}
            >
              Profile
            </MenuItem>
            <MenuItem
              onClick={() => {
                close();
                toggleChangePassword();
              }}
            >
              Change Password
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={() => {
                close();
                toggleUserAccessKeys();
              }}
            >
              Access Keys
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={() => {
                close();
                toggleCLI();
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
