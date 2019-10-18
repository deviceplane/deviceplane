import React, { Component, Fragment } from 'react';
import { Button, Pane, Table, Heading, majorScale } from 'evergreen-ui';
import axios from 'axios';
import moment from 'moment';

import config from '../config.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';

export default class Provisioning extends Component {
  state = {
    deviceRegistrationTokens: []
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          deviceRegistrationTokens: response.data
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
          heading="Provisioning"
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
              <Heading paddingLeft={majorScale(2)}>
                Device Registration Tokens
              </Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/provisioning/deviceregistrationtokens/create`
                  )
                }
              >
                Create Device Registration Token
              </Button>
            </Pane>
            {this.state.deviceRegistrationTokens &&
              this.state.deviceRegistrationTokens.length > 0 && (
                <Table>
                  <Table.Head>
                    <Table.TextHeaderCell>Name</Table.TextHeaderCell>
                    <Table.TextHeaderCell>Created At</Table.TextHeaderCell>
                    <Table.TextHeaderCell>Devices Registered</Table.TextHeaderCell>
                    <Table.TextHeaderCell>Registration Limit</Table.TextHeaderCell>
                  </Table.Head>
                  <Table.Body>
                    {this.state.deviceRegistrationTokens.map(token => (
                      <Table.Row
                        key={token.id}
                        isSelectable
                        onSelect={() =>
                          this.props.history.push(
                            `/${this.props.projectName}/provisioning/deviceregistrationtokens/${token.name}/overview`
                          )
                        }
                      >
                        <Table.TextCell>{token.name}</Table.TextCell>
                        <Table.TextCell>
                          {token.createdAt
                            ? moment(token.createdAt).fromNow()
                            : "-"}
                        </Table.TextCell>
                        <Table.TextCell>
                          {token.deviceCounts.allCount}
                        </Table.TextCell>
                        <Table.TextCell>
                          {typeof token.maxRegistrations === "number"
                            ? token.maxRegistrations
                            : "unlimited"}
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
