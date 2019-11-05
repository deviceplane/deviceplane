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
  TextDropdownButton,
  Menu,
  Position,
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
  filterQuery: Query,
  showFilterDialog: boolean,
  popoverShown: boolean,
  defaultDeviceRegistrationTokenExists: boolean,

  orderedColumn?: string,
  order?: string,
  page: number,
}

const Params = {
  Filter: 'filter',
  Page: 'page',
  OrderedColumn: 'order_by',
  OrderDirection: 'order',
}

const Order = {
  ASC: 'asc',
  DESC: 'desc'
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
      filterQuery: [],
      showFilterDialog: false,
      popoverShown: false,
      defaultDeviceRegistrationTokenExists: false,

      orderedColumn: undefined,
      order: undefined,
      page: 0,
    };
  }

  async componentDidMount() {
    await this.parseQueryString()
    this.queryDevices();
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

  buildFiltersQuery = (): string[] =>
    [...Array.from(new Set(this.state.filterQuery))]
      .map(
        filter => `${Params.Filter}=${encodeURIComponent(btoa(JSON.stringify(filter)))}`
      );

  parseQueryString = () => {
    return new Promise((resolve) => {
      if (!window.location.search || window.location.search.length < 1) {
        resolve();
        return;
      }

      var builtQuery: Query = [];
      var page = 0;
      var orderedColumn = undefined;
      var order = undefined;

      let queryParams = window.location.search.substr(1).split('&');
      queryParams.forEach(queryParam => {
        const parts = queryParam.split('=');
        if (parts.length < 2) {
          return;
        }

        switch (parts[0]) {
          case Params.Filter: {
            let encodedFilter = parts[1];
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
            break;
          }
          case Params.Page: {
            let p = Number(parts[1]);
            if (!isNaN(p)) {
              page = p;
            }
            break;
          }
          case Params.OrderedColumn: {
            let p = parts[1];
            if (p) {
              orderedColumn = p;
            }
            break;
          }
          case Params.OrderDirection: {
            let p = parts[1];
            if (p) {
              order = p;
            }
            break;
          }
        }
      })
      this.setState({
        page: page,
        orderedColumn: orderedColumn,
        order: order,
        filterQuery: [...this.state.filterQuery, ...builtQuery]
      }, resolve);
    })
  }

  queryDevices = () => {
    var query: string[] = [];

    query.push(`${Params.Page}=${this.state.page}`)
    if (this.state.orderedColumn) {
      query.push(`${Params.OrderedColumn}=${this.state.orderedColumn}`)
    }
    if (this.state.order) {
      query.push(`${Params.OrderDirection}=${this.state.order}`)
    }
    query.push(...this.buildFiltersQuery())

    const queryString = '?' + query.join('&');
    window.history.pushState('', '', query.length ? queryString : window.location.pathname );
    this.fetchDevices(queryString);
};

  removeFilter = (index: number) => {
    this.setState({
      page: 0,
      filterQuery: this.state.filterQuery.filter((_, i) => i !== index),
    }, this.queryDevices);
  };

  addFilter = (filter: Filter) => {
    this.setState({
      page: 0,
      showFilterDialog: false,
      filterQuery: [...this.state.filterQuery, filter]
    }, this.queryDevices);
  };

  clearFilters = () => {
    this.setState({
      page: 0,
      filterQuery: []
    }, this.queryDevices);
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

  getIconForOrder = (order?: string) => {
    switch (order) {
      case Order.ASC:
        return 'arrow-up'
      case Order.DESC:
        return 'arrow-down'
      default:
        return 'caret-down'
    }
  }

  renderOrderedTableHeader = (title: string, jsonName: string) => {
    return (
      <Popover
        position={Position.BOTTOM_LEFT}
        content={({ close }: any) => (
          <Menu>
            <Menu.OptionsGroup
              title="Order"
              options={[
                { label: 'Ascending', value: Order.ASC },
                { label: 'Descending', value: Order.DESC }
              ]}
              selected={
                this.state.orderedColumn === jsonName ? this.state.order : null
              }
              onChange={(value: string) => {
                this.setState({
                  orderedColumn: jsonName,
                  order: value
                }, this.queryDevices)

                // Close the popover when you select a value.
                close()
              }}
            />
          </Menu>
        )}
      >
        <TextDropdownButton
          icon={
            this.state.orderedColumn === jsonName
              ? this.getIconForOrder(this.state.order)
              : 'caret-down'
          }
        >
          {title}
        </TextDropdownButton>
      </Popover>
    )
  }

  render() {
    const { user, history, projectName } = this.props;
    const { showFilterDialog, filterQuery } = this.state;

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
                  {filterQuery.length > 0 && (
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
            query={this.state.filterQuery}
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
                  {this.renderOrderedTableHeader('Status', 'status')}
                </Table.TextHeaderCell>
                <Table.TextHeaderCell>
                  {this.renderOrderedTableHeader('Name', 'name')}
                </Table.TextHeaderCell>
                <Table.TextHeaderCell
                  flexBasis={120}
                  flexShrink={0}
                  flexGrow={0}
                >
                  {this.renderOrderedTableHeader('IP Address', 'ipAddress')}
                </Table.TextHeaderCell>
                <Table.TextHeaderCell
                  flexBasis={120}
                  flexShrink={0}
                  flexGrow={0}
                >
                  OS
                  {/* In the future, we can add nesting and use osRelease.prettyName */}
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
