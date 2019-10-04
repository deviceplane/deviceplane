import React, { Component, Fragment } from 'react';
import { Button, Pane, Table, Heading, majorScale } from 'evergreen-ui';
import axios from 'axios';
import moment from 'moment';

import config from '../config.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';

export default class Applications extends Component {
  state = {
    applications: []
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/applications?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          applications: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    return (
      <Fragment>
        <TopHeader
          user={this.props.user}
          heading="Applications"
          history={this.props.history}
        />
        <Pane width="70%">
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Applications</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/applications/create`
                  )
                }
              >
                Create Application
              </Button>
            </Pane>
            {this.state.applications && this.state.applications.length > 0 && (
              <Table>
                <Table.Head>
                  <Table.TextHeaderCell>Name</Table.TextHeaderCell>
                  <Table.TextHeaderCell>Last Release</Table.TextHeaderCell>
                  <Table.TextHeaderCell>Device Count</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  {this.state.applications.map(application => (
                    <Table.Row
                      key={application.id}
                      isSelectable
                      onSelect={() =>
                        this.props.history.push(
                          `/${this.props.projectName}/applications/${application.name}`
                        )
                      }
                    >
                      <Table.TextCell>{application.name}</Table.TextCell>
                      <Table.TextCell>
                        {application.latestRelease
                          ? moment(
                              application.latestRelease.createdAt
                            ).fromNow()
                          : '-'}
                      </Table.TextCell>
                      <Table.TextCell>
                        {application.deviceCounts.allCount}
                      </Table.TextCell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
            )}
          </InnerCard>
        </Pane>
      </Fragment>
    );
  }
}
