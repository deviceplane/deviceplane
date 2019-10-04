import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  Table,
  Dialog,
  majorScale,
  Button,
  Heading,
  Badge,
  IconButton,
  TextInputField,
  toaster
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';

export default class Members extends Component {
  state = {
    members: []
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/memberships?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          members: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const members = this.state.members;
    return (
      <Pane width="70%">
        {members ? (
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Members</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/iam/members/add`
                  )
                }
              >
                Add Member
              </Button>
            </Pane>
            {members && members.length > 0 && (
              <Table>
                <Table.Head>
                  <Table.TextHeaderCell>Email</Table.TextHeaderCell>
                  <Table.TextHeaderCell>Name</Table.TextHeaderCell>
                  <Table.TextHeaderCell>Roles</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  {members.map(member => (
                    <Table.Row
                      key={member.userId}
                      isSelectable
                      onSelect={() =>
                        this.props.history.push(
                          `/${this.props.projectName}/iam/members/${member.userId}`
                        )
                      }
                    >
                      <Table.TextCell>{member.user.email}</Table.TextCell>
                      <Table.TextCell>{`${member.user.firstName} ${member.user.lastName}`}</Table.TextCell>
                      <Table.TextCell>
                        {member.roles.map(role => role.name).join(',')}
                      </Table.TextCell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
            )}
          </InnerCard>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}
