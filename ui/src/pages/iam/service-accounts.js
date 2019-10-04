import React, { Component, Fragment } from 'react';
import axios from 'axios';
import { Pane, Table, majorScale, Button, Heading } from 'evergreen-ui';

import config from '../../config';
import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';

export default class ServiceAccounts extends Component {
  state = {
    serviceAccounts: []
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          serviceAccounts: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const serviceAccounts = this.state.serviceAccounts;
    return (
      <Pane width="70%">
        {serviceAccounts ? (
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Service Accounts</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/iam/serviceaccounts/create`
                  )
                }
              >
                Create Service Account
              </Button>
            </Pane>
            {this.state.serviceAccounts &&
              this.state.serviceAccounts.length > 0 && (
                <Table>
                  <Table.Head>
                    <Table.TextHeaderCell>Service Account</Table.TextHeaderCell>
                  </Table.Head>
                  <Table.Body>
                    {serviceAccounts.map(serviceAccount => (
                      <Table.Row
                        key={serviceAccount.id}
                        isSelectable
                        onSelect={() =>
                          this.props.history.push(
                            `/${this.props.projectName}/iam/serviceaccounts/${serviceAccount.name}`
                          )
                        }
                      >
                        <Table.TextCell>{serviceAccount.name}</Table.TextCell>
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
