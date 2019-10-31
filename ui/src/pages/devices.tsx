import React, { Fragment, Component } from 'react';
import axios from 'axios';
import {
  Button,
  Pane,
  Table,
  Heading,
  Badge,
  majorScale,
  minorScale,
  Icon,
  Text,
  withTheme,
  Popover,
  // @ts-ignore
} from 'evergreen-ui';

import {
  Link
} from 'react-router-dom';

import config from '../config.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';
import { DevicesFilter, Query, Filter, Condition, DevicePropertyCondition, LabelExistenceCondition, LabelValueCondition, LabelValueConditionParams, LabelExistenceConditionParams, DevicePropertyConditionParams } from '../components/DevicesFilter';
import { DevicesFilterButtons } from '../components/DevicesFilterButtons';
import { buildLabelColorMap, renderLabels } from '../helpers/labels.js';

// Runtime type safety
import * as deviceTypes from '../components/DevicesFilter-ti';
import { createCheckers } from 'ts-interface-checker';
const typeCheckers = createCheckers(deviceTypes.default);

interface Props {
  theme: any,
  projectName: string,
  user: any,
  history: any,
}

interface State {
  labelColors: any[],
  labelColorMap: any,
  devices: any[],
  query: Query,
  showFilterDialog: boolean,
  popoverShown: boolean,
  defaultDeviceRegistrationTokenExists: boolean,
}

class Devices extends Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const palletteArray = Object.values(this.props.theme.palette);

    this.state = {
      labelColorMap: {},
      labelColors: [
        ...palletteArray.map((colors: any) => colors.base),
        ...palletteArray.map((colors: any) => colors.dark)
      ],
      devices: [],
      query: [],
      showFilterDialog: false,
      popoverShown: false,
      defaultDeviceRegistrationTokenExists: false,
    };
  }

  componentDidMount() {
    const queryString = this.parseFiltersQueryString()
    this.fetchDevices(queryString);
    this.checkDefaultDeviceRegistrationToken();
  }

  checkDefaultDeviceRegistrationToken = () => {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/default`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          defaultDeviceRegistrationTokenExists: true,
        });
      })
      .catch(error => {});
  }

  fetchDevices = (queryString: string) => {
    if (queryString.length) {
      queryString = '?' + queryString;
    }
    return axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices${queryString}`,
        {
          withCredentials: true
        }
      )
      .then(({ data: devices }) => {
        devices = this.parseDevices(devices);
        const labelColorMap = buildLabelColorMap(this.state.labelColorMap, this.state.labelColors, devices);
        this.setState({
          devices,
          labelColorMap
        });
      })
      .catch(console.log);
  };

  parseDevices = (devices: any[]) =>
    devices.map(device => device);

  buildFiltersQueryString = () =>
    [...Array.from(new Set(this.state.query))]
      .map(
        filter => `filter=${encodeURIComponent(btoa(JSON.stringify(filter)))}`
      )
      .join('&');

  parseFiltersQueryString = () => {
    let queryParams = '';

    if (window.location.search && window.location.search.includes('filter=')) {
      queryParams = window.location.search.substr(1);
      const filters = queryParams.replace(/&/g, '').split('filter=');

      var builtQuery: Query = [];

      filters.forEach(encodedFilter => {
        if (encodedFilter) {
          try {
            const filter = JSON.parse(
              atob(decodeURIComponent(encodedFilter))
            )

            const validFilter = filter.filter((c: Condition) => {
              return typeCheckers['Condition'].strictTest(c);
            });

            if (validFilter.length) {
              builtQuery.push(validFilter);
            }
          } catch (e) {
            console.log('Error parsing filters:', e)
          }
        }
      });

      this.setState({
        query: [...this.state.query, ...builtQuery]
      });
    }
    return queryParams;
  }

  filterDevices = () => {
    if (!this.state.query.length) {
      window.history.pushState('', '', window.location.pathname);
      this.fetchDevices('');
      return;
    }

    const filtersQueryString = this.buildFiltersQueryString();
    window.history.pushState('', '', `?${filtersQueryString}`);
    this.fetchDevices(filtersQueryString);
  };

  removeFilter = (index: number) => {
    this.setState({
      query: this.state.query.filter((_, i) => i !== index),
    }, this.filterDevices);
  };

  addFilter = (filter: Filter) => {
    this.setState({
      showFilterDialog: false,
      query: [...this.state.query, filter]
    }, this.filterDevices);
  };

  clearFilters = () => {
    this.setState({
      query: []
    }, this.filterDevices);
  };

  renderDeviceOs = (device: any) => {
    var innerText = '-';
    if (
      device.info.hasOwnProperty('osRelease') &&
      device.info.osRelease.hasOwnProperty('prettyName')
    ) {
      innerText = device.info.osRelease.prettyName;
    }
    return <Pane>{innerText}</Pane>;
  };

  renderTableBody = () => {
    const { history, projectName } = this.props;
    const { devices } = this.state;

    return devices.map(device => (
      <Table.Row
        key={device.id}
        isSelectable
        onSelect={() => history.push(`/${projectName}/devices/${device.name}`)}
        flexGrow={1}
        height="auto"
        paddingY={majorScale(1)}
        alignItems="flex-start"
      >
        <Table.TextCell
          flexBasis={90}
          flexShrink={0}
          flexGrow={0}
          alignItems="center"
          paddingRight="0"
          marginY={minorScale(1)}
        >
          {device.status === 'offline' ? (
            <Badge color="red">offline</Badge>
          ) : (
            <Badge color="green">online</Badge>
          )}
        </Table.TextCell>
        <Table.TextCell marginY={minorScale(1)}>{device.name}</Table.TextCell>
        <Table.TextCell
          flexBasis={120}
          flexShrink={0}
          flexGrow={0}
          marginY={minorScale(1)}
        >
          {device.info.hasOwnProperty('ipAddress') ? device.info.ipAddress : ''}
        </Table.TextCell>
        <Table.TextCell
          marginY={minorScale(1)}
          flexBasis={120}
          flexShrink={0}
          flexGrow={0}
        >
          {this.renderDeviceOs(device)}
        </Table.TextCell>
        <Table.TextCell flexGrow={2}>
          {renderLabels(device.labels, this.state.labelColorMap)}
        </Table.TextCell>
      </Table.Row>
    ));
  };

  render() {
    const { user, history, projectName } = this.props;
    const { showFilterDialog, query } = this.state;

    var addDeviceButtonHolder;
    if (this.state.defaultDeviceRegistrationTokenExists) {
      addDeviceButtonHolder = (
        <Button
          appearance="primary"
          onClick={() => history.push(`/${projectName}/devices/add`)}
        >
          Add Device
        </Button>
      );
    } else {
      addDeviceButtonHolder = (
        <Popover
          trigger="hover"
          isShown={this.state.popoverShown}
          content={
            <Pane
              display="flex"
              alignItems="center"
              justifyContent="center"
              flexDirection="column"
              width="250px"
              padding="20px"
              onMouseOver={() => {
                this.setState({ popoverShown: true });
              }}
              onMouseOut={() => {
                this.setState({ popoverShown: false });
              }}
            >
              <Text>
                There is no "default" device registration token, so adding
                devices from the UI is disabled.
              </Text>
              <Text paddingTop={minorScale(1)}>
                Device registration tokens can be created on the{" "}
                <Link style={{color: "blue"}} to={`/${projectName}/provisioning`}>
                  Provisioning
                </Link>{" "}
                page.
              </Text>
            </Pane>
          }
        >
          <Pane
            appearance="primary"
            onMouseOver={() => {
              this.setState({ popoverShown: true });
            }}
            onMouseOut={() => {
              this.setState({ popoverShown: false });
            }}
          >
            <Button disabled={true}>Add Device</Button>
          </Pane>
        </Popover>
      );
    }

    return (
      <Fragment>
        <TopHeader user={user} heading="Devices" history={history} />
        <Pane width="70%">
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Devices</Heading>
              <Pane alignItems="center" padding={majorScale(2)}>
                <Pane
                display="flex"
                marginLeft={majorScale(2)}>
                  {query.length > 0 && (
                    <Button
                      marginRight={majorScale(2)}
                      appearance="minimal"
                      intent="danger"
                      onClick={this.clearFilters}
                    >
                      Clear filters
                    </Button>
                  )}
                  <Button
                    iconBefore="plus"
                    marginRight={majorScale(2)}
                    onClick={() => this.setState({ showFilterDialog: true })}
                  >
                    Add Filter
                  </Button>
                  {addDeviceButtonHolder}
                </Pane>
              </Pane>
            </Pane>
            <DevicesFilterButtons
            query={this.state.query}
            canRemoveFilter={true}
            removeFilter={this.removeFilter}
            />
            <Table>
              <Table.Head background="tint2">
                <Table.TextHeaderCell
                  flexBasis={90}
                  flexShrink={0}
                  flexGrow={0}
                >
                  Status
                </Table.TextHeaderCell>
                <Table.TextHeaderCell>Name</Table.TextHeaderCell>
                <Table.TextHeaderCell
                  flexBasis={120}
                  flexShrink={0}
                  flexGrow={0}
                >
                  IP Address
                </Table.TextHeaderCell>
                <Table.TextHeaderCell
                  flexBasis={120}
                  flexShrink={0}
                  flexGrow={0}
                >
                  OS
                </Table.TextHeaderCell>
                <Table.TextHeaderCell flexGrow={2}>Labels</Table.TextHeaderCell>
              </Table.Head>
              <Table.Body>{this.renderTableBody()}</Table.Body>
            </Table>
          </InnerCard>
        </Pane>

        <DevicesFilter
          show={showFilterDialog}
          onClose={() => this.setState({ showFilterDialog: false })}
          onSubmit={this.addFilter}
        />
      </Fragment>
    );
  }
}

export default withTheme(Devices);
