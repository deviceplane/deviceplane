import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useCurrentRoute } from 'react-navi';

import Popover from './popover';
import { Row, Text, Icon, Link } from './core';
import api from '../api';

const MenuItem = styled(Link)`
  text-transform: none;

  &:hover {
    color: ${props => props.theme.colors.pureWhite};
    background-color: ${props => props.theme.colors.grays[3]};
  }
`;
MenuItem.defaultProps = {
  paddingX: '6px',
  paddingY: 2,
  color: 'white',
  fontSize: 1,
  fontWeight: 1,
};

const Selector = styled(Row)`
  &:hover {
    border-color: ${props => props.theme.colors.primary};
  }
`;

const ProjectSelector = ({}) => {
  const { data } = useCurrentRoute();
  const [projects, setProjects] = useState([]);

  const loadProjects = async () => {
    try {
      const data = await api.projects();
      setProjects(data);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    loadProjects();
  }, []);

  if (!data.params.project) {
    return null;
  }

  return (
    <Popover
      button={({ show }) => (
        <Selector
          alignItems="center"
          padding={1}
          paddingX="6px"
          borderRadius={1}
          bg="black"
          border={0}
          borderColor={show ? 'primary' : 'white'}
        >
          <Text fontSize={1} fontWeight={1} color="white">
            {data.params.project}
          </Text>
          <Icon icon="caret-down" size={16} color="white" marginLeft={4} />
        </Selector>
      )}
      content={({ close }) =>
        projects.map(({ name }) => (
          <MenuItem href={`/${name}`} onClick={close}>
            {name}
          </MenuItem>
        ))
      }
      top="43px"
      width="180px"
    ></Popover>
  );
};

export default ProjectSelector;
