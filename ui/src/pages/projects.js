import React, { Component } from 'react';
import { Button, Pane, Heading, majorScale, Table, Card } from 'evergreen-ui';
import axios from 'axios';

import config from '../config.js';
import TopHeader from '../components/TopHeader.js';

export default class Projects extends Component {
  state = {
    membershipsFull: []
  };

  componentDidMount() {
    axios
      .get(`${config.endpoint}/memberships?full`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          membershipsFull: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    return (
      <Pane
        display="flex"
        flexDirection="column"
        alignItems="center"
        background="tint1"
        flex={1}
        justifyContent="stretch"
      >
        <TopHeader
          user={this.props.user}
          heading="Select Project"
          showLogo={true}
          hideSwitchProjects={true}
          history={this.props.history}
        />
        <Card
          display="flex"
          flexDirection="column"
          margin={majorScale(2)}
          elevation={1}
          width="60%"
          background="white"
        >
          <Pane
            display="flex"
            flexDirection="row"
            justifyContent="space-between"
            alignItems="center"
          >
            <Heading paddingLeft={majorScale(2)}>Projects</Heading>
            <Button
              margin={majorScale(2)}
              appearance="primary"
              onClick={() => this.props.history.push(`/projects/create`)}
            >
              Create Project
            </Button>
          </Pane>
          {this.state.membershipsFull && this.state.membershipsFull.length > 0 && (
            <Table>
              <Table.Head background="tint2">
                <Table.TextHeaderCell>Project Name</Table.TextHeaderCell>
                <Table.TextHeaderCell>Devices</Table.TextHeaderCell>
                <Table.TextHeaderCell>Applications</Table.TextHeaderCell>
              </Table.Head>
              <Table.Body>
                {this.state.membershipsFull.map(membershipFull => (
                  <Table.Row
                    key={membershipFull.project.id}
                    isSelectable
                    onSelect={() =>
                      this.props.history.push(`/${membershipFull.project.name}`)
                    }
                  >
                    <Table.TextCell>
                      {membershipFull.project.name}
                    </Table.TextCell>
                    <Table.TextCell>
                      {membershipFull.project.deviceCounts.allCount}
                    </Table.TextCell>
                    <Table.TextCell>
                      {membershipFull.project.applicationCounts.allCount}
                    </Table.TextCell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          )}
        </Card>
      </Pane>
    );
  }
}
