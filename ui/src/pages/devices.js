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
  withTheme
} from 'evergreen-ui';

import config from '../config.js';
import InnerCard from '../components/InnerCard.js';
import TopHeader from '../components/TopHeader.js';
import Filter from '../components/devices-filter';

class Devices extends Component {
  constructor(props) {
    super(props);

    const palletteArray = Object.values(this.props.theme.palette);

    this.labelColors = [
      ...palletteArray.map(colors => colors.base),
      ...palletteArray.map(colors => colors.dark)
    ];

    this.state = {
      devices: [],
      filters: [],
      showFilterDialog: false
    };
  }

  componentDidMount() {
    this.fetchDevices();
  }

  fetchDevices = () => {
    let queryParams = '';

    if (window.location.search && window.location.search.includes('filter=')) {
      queryParams = window.location.search.substr(1);
      const filters = queryParams.replace(/&/g, '').split('filter=');

      filters.forEach(encodedFilter => {
        if (encodedFilter) {
          try {
            const filter = JSON.parse(
              atob(decodeURIComponent(encodedFilter))
            ).filter(f => f.property && f.operator && f.value);

            if (filter.length) {
              this.setState(({ filters }) => ({
                filters: [...filters, filter]
              }));
            }
          } catch (e) {}
        }
      });
    }

    return axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices${
          queryParams ? `?${queryParams}` : ''
        }`,
        {
          withCredentials: true
        }
      )
      .then(({ data: devices }) => {
        devices = this.parseDevices(devices);
        var x = [];
        devices.forEach((device) => {
          Object.keys(device.labels).forEach((label) => {
            x.push(label);
          })
        });
        const labelKeys = [
          ...new Set(x)
        ].sort();
        const labelColorMap = labelKeys.reduce(
          (obj, key, i) => ({
            ...obj,
            [key]: this.labelColors[i % (this.labelColors.length - 1)]
          }),
          {}
        );
        this.setState({
          devices,
          labelColorMap
        });
      })
      .catch(console.log);
  };

  parseDevices = devices =>
    devices.map(device => device);

  buildFiltersQueryString = () =>
    [...new Set(this.state.filters)]
      .map(
        filter => `filter=${encodeURIComponent(btoa(JSON.stringify(filter)))}`
      )
      .join('&');

  filterDevices = () => {
    if (!this.state.filters.length) {
      window.history.pushState('', '', window.location.pathname);
      this.fetchDevices();
      return;
    }

    const filtersQueryString = this.buildFiltersQueryString();
    window.history.pushState('', '', `?${filtersQueryString}`);
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices?${filtersQueryString}`,
        {
          withCredentials: true
        }
      )
      .then(({ data: devices }) => {
        this.setState({
          devices: this.parseDevices(devices)
        });
      })
      .catch(console.log);
  };

  removeFilter = index => {
    this.setState(
      ({ filters }) => ({
        filters: filters.filter((_, i) => i !== index)
      }),
      this.filterDevices
    );
  };

  addFilter = filter => {
    this.setState(
      ({ filters }) => ({
        showFilterDialog: false,
        filters: [...filters, filter]
      }),
      this.filterDevices
    );
  };

  clearFilters = () => {
    this.setState({ filters: [] }, this.filterDevices);
  };

  renderFilter = ({ property, operator, key, value }) => {
    switch (property) {
      case 'label':
        return (
          <Fragment>
            <Text
              fontWeight={500}
              marginRight={minorScale(1)}
              style={{ textTransform: 'capitalize' }}
            >
              {property}
            </Text>
            <Text fontWeight={500} marginRight={minorScale(1)}>
              {operator}
            </Text>

            <Text fontWeight={700}>{key}</Text>
            {value && (
              <Fragment>
                <Text fontWeight={700} marginX={2}>
                  :
                </Text>
                <Text fontWeight={700}>{value}</Text>
              </Fragment>
            )}
          </Fragment>
        );

      case 'status':
        return (
          <Fragment>
            <Text
              fontWeight={500}
              marginRight={minorScale(1)}
              style={{ textTransform: 'capitalize' }}
            >
              {property}
            </Text>
            <Text fontWeight={500} marginRight={minorScale(1)}>
              {operator}
            </Text>
            <Text style={{ textTransform: 'capitalize' }} fontWeight={700}>
              {value}
            </Text>
          </Fragment>
        );
    }
  };

  renderLabels = device => (
    <Pane display="flex" flexWrap="wrap">
      {Object.keys(device.labels).map((key, i) => (
        <Pane
          display="flex"
          marginRight={minorScale(2)}
          marginY={minorScale(1)}
          overflow="hidden"
          key={key}
        >
          <Pane
            backgroundColor={this.state.labelColorMap[key]}
            paddingX={6}
            paddingY={2}
            color="white"
            borderTopLeftRadius={3}
            borderBottomLeftRadius={3}
            textOverflow="ellipsis"
            overflow="hidden"
            whiteSpace="nowrap"
          >
            {key}
          </Pane>
          <Pane
            backgroundColor="#E4E7EB"
            paddingX={6}
            paddingY={2}
            borderTopRightRadius={3}
            borderBottomRightRadius={3}
            textOverflow="ellipsis"
            overflow="hidden"
            whiteSpace="nowrap"
          >
            {device.labels[key]}
          </Pane>
        </Pane>
      ))}
    </Pane>
  );

  renderDeviceOs = device => {
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
          {this.renderLabels(device)}
        </Table.TextCell>
      </Table.Row>
    ));
  };

  render() {
    const { user, history, projectName } = this.props;
    const { showFilterDialog, filters } = this.state;

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
              <Pane display="flex" alignItems="center" padding={majorScale(2)}>
                <Pane marginLeft={majorScale(2)}>
                  {filters.length > 0 && (
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
                  <Button
                    appearance="primary"
                    onClick={() => history.push(`/${projectName}/devices/add`)}
                  >
                    Add Device
                  </Button>
                </Pane>
              </Pane>
            </Pane>
            {filters.length > 0 && (
              <Pane
                paddingX={majorScale(2)}
                paddingBottom={majorScale(2)}
                display="flex"
              >
                {filters.map((conditions, index) => (
                  <Pane
                    display="flex"
                    alignItems="center"
                    marginRight={minorScale(3)}
                    key={index}
                  >
                    <Pane
                      backgroundColor="#B7D4EF"
                      borderRadius={3}
                      paddingX={minorScale(2)}
                      paddingY={minorScale(1)}
                      display="flex"
                      alignItems="center"
                    >
                      {conditions.map((filter, i) => (
                        <Fragment key={i}>
                          {this.renderFilter(filter)}
                          {conditions.length > 1 &&
                            i !== conditions.length - 1 && (
                              <Text
                                fontSize={10}
                                fontWeight={700}
                                marginX={minorScale(3)}
                                color="white"
                              >
                                OR
                              </Text>
                            )}
                        </Fragment>
                      ))}
                      <Icon
                        marginLeft={minorScale(3)}
                        icon="cross"
                        appearance="minimal"
                        cursor="pointer"
                        color="white"
                        size={14}
                        onClick={() => this.removeFilter(index)}
                      />
                    </Pane>
                  </Pane>
                ))}
              </Pane>
            )}
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

        <Filter
          show={showFilterDialog}
          onClose={() => this.setState({ showFilterDialog: false })}
          onSubmit={this.addFilter}
        />
      </Fragment>
    );
  }
}

export default withTheme(Devices);
