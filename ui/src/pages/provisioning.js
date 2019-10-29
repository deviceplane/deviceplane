import React, { Component, Fragment } from 'react';
import { Button, Pane, Table, Heading, majorScale, withTheme } from 'evergreen-ui';
import axios from 'axios';
import moment from 'moment';

import config from '../config.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';
import { buildLabelColorMap, renderLabels } from '../helpers/labels.js';

class Provisioning extends Component {
  constructor(props) {
    super(props);

    const palletteArray = Object.values(this.props.theme.palette);

    this.labelColors = [
      ...palletteArray.map(colors => colors.base),
      ...palletteArray.map(colors => colors.dark)
    ];

    this.state = {
      deviceRegistrationTokens: []
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        const labelColorMap = buildLabelColorMap({}, this.labelColors, response.data);
        this.setState({
          deviceRegistrationTokens: response.data,
          labelColorMap,
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
                    <Table.TextHeaderCell flexBasis={100}>Name</Table.TextHeaderCell>
                    <Table.TextHeaderCell flexBasis={50}>Created At</Table.TextHeaderCell>
                    <Table.TextHeaderCell flexBasis={50}>Devices Registered</Table.TextHeaderCell>
                    <Table.TextHeaderCell flexBasis={50}>Registration Limit</Table.TextHeaderCell>
                    <Table.TextHeaderCell flexBasis={150}>Labels</Table.TextHeaderCell>
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
                        flexGrow={1}
                        height="auto"
                        paddingY={majorScale(1)}
                        alignItems="flex-start"
                      >
                        <Table.TextCell flexBasis={100}>{token.name}</Table.TextCell>
                        <Table.TextCell flexBasis={50}>
                          {token.createdAt
                            ? moment(token.createdAt).fromNow()
                            : "-"}
                        </Table.TextCell>
                        <Table.TextCell flexBasis={50}>
                          {token.deviceCounts.allCount}
                        </Table.TextCell>
                        <Table.TextCell flexBasis={50}>
                          {typeof token.maxRegistrations === "number"
                            ? token.maxRegistrations
                            : "unlimited"}
                        </Table.TextCell>
                        <Table.TextCell flexBasis={150}>
                          {renderLabels(token.labels, this.state.labelColorMap)}
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

export default withTheme(Provisioning);