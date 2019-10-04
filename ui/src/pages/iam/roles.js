import React, { Component } from 'react';
import axios from 'axios';
import { Pane, Table, majorScale, Button, Heading } from 'evergreen-ui';

import config from '../../config';
import CustomSpinner from '../../components/CustomSpinner';
import InnerCard from '../../components/InnerCard';

export default class Roles extends Component {
  state = {
    roles: []
  };

  componentDidMount() {
    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/roles`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          roles: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const roles = this.state.roles;
    return (
      <Pane width="70%">
        {roles ? (
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Roles</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/iam/roles/create`
                  )
                }
              >
                Create Role
              </Button>
            </Pane>
            {this.state.roles && this.state.roles.length > 0 && (
              <Table>
                <Table.Head>
                  <Table.TextHeaderCell>Role</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  {roles.map(role => (
                    <Table.Row
                      key={role.id}
                      isSelectable
                      onSelect={() =>
                        this.props.history.push(
                          `/${this.props.projectName}/iam/roles/${role.name}`
                        )
                      }
                    >
                      <Table.TextCell>{role.name}</Table.TextCell>
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
