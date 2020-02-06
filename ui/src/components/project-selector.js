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
    background-color: ${props => props.theme.colors.grays[0]};
  }
`;
MenuItem.defaultProps = {
  paddingX: 2,
  paddingY: 2,
  color: 'white',
  fontSize: 2,
  fontWeight: 1,
};

const ProjectSelector = ({}) => {
  const { data } = useCurrentRoute();
  const [projects, setProjects] = useState([]);

  const loadProjects = async () => {
    try {
      const data = await api.projects();
      setProjects(data);
    } catch (error) {
      console.log(error);
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
      content={({ close }) =>
        projects.map(({ name }) => (
          <MenuItem href={`/${name}`} onClick={close}>
            {name}
          </MenuItem>
        ))
      }
      top="40px"
      width="180px"
    >
      <Row
        alignItems="center"
        padding={2}
        borderRadius={1}
        bg="black"
        border={0}
        borderColor="white"
      >
        <Text fontSize={2} fontWeight={2} color="white">
          {data.params.project}
        </Text>
        <Icon icon="caret-down" size={16} color="white" marginLeft={4} />
      </Row>
    </Popover>
  );
};

export default ProjectSelector;
