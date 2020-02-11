import React, { useState, useMemo, useEffect, useCallback } from 'react';

import api from '../api';
import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Row, Text, Input, Icon } from '../components/core';
import {
  DevicesFilter,
  OperatorIs,
  LabelValueCondition,
} from '../components/devices-filter';
import { DevicesFilterButtons } from '../components/devices-filter-buttons';
import DeviceStatus from '../components/device-status';
import { renderLabels } from '../helpers/labels';

import storage from '../storage';
import utils from '../utils';

const Params = {
  Filter: 'filter',
  OrderedColumn: 'order_by',
  OrderDirection: 'order',
};

const Devices = ({ route }) => {
  const [showFilterDialog, setShowFilterDialog] = useState();
  const [devices, setDevices] = useState(route.data.devices);
  const [deviceTotal, setDeviceTotal] = useState();
  const [filterQuery, setFilterQuery] = useState(
    storage.get('devicesFilter', route.data.params.project) || []
  );
  const [orderedColumn, setOrderedColumn] = useState();
  const [order, setOrder] = useState();
  const [filterToEdit, setFilterToEdit] = useState(null);
  const [searchInput, setSearchInput] = useState('');
  const [searchFocused, setSearchFocused] = useState();

  useEffect(() => {
    parseQueryString();
  }, []);

  useEffect(() => {
    queryDevices();
  }, [filterQuery, searchInput]);

  useEffect(() => {
    storage.set('devicesFilter', filterQuery, route.data.params.project);
  }, [filterQuery]);

  const addLabelFilter = useCallback(
    ({ key, value }) => {
      const labelFilter = [
        {
          type: LabelValueCondition,
          params: {
            key,
            operator: OperatorIs,
            value,
          },
        },
      ];
      if (!filterQuery.find(filter => utils.deepEqual(filter, labelFilter))) {
        setFilterQuery(filterQuery => [...filterQuery, labelFilter]);
      }
    },
    [filterQuery]
  );

  const columns = useMemo(
    () => [
      {
        Header: 'Status',
        accessor: 'status',
        Cell: ({ cell: { value } }) => <DeviceStatus status={value} />,
        maxWidth: '72px',
      },
      {
        Header: 'Name',
        accessor: 'name',
        minWidth: '200px',
      },
      {
        Header: 'IP Address',
        accessor: ({ info }) => info.ipAddress || '-',
        maxWidth: '140px',
        minWidth: '140px',
      },
      {
        Header: 'OS',
        accessor: ({ info }) => info.osRelease.prettyName || '-',
        minWidth: '200px',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value } }) =>
          value ? renderLabels(value, addLabelFilter) : null,
        minWidth: '300px',
      },
    ],
    [filterQuery]
  );
  const tableData = useMemo(() => [...devices, ...devices], [devices]);

  const fetchDevices = async queryString => {
    try {
      const { data, headers } = await api.devices({
        projectId: route.data.params.project,
        queryString,
      });
      setDeviceTotal(headers['total-device-count']);
      setDevices(data);
    } catch (error) {
      console.error(error);
    }
  };

  const queryDevices = () => {
    const query = [];

    if (orderedColumn) {
      query.push(`${Params.OrderedColumn}=${orderedColumn}`);
    }
    if (order) {
      query.push(`${Params.OrderDirection}=${order}`);
    }
    if (searchInput) {
      query.push(`search=${searchInput}`);
    }
    query.push(...buildFiltersQuery());
    const queryString = '?' + query.join('&');

    if (query.length) {
      window.history.replaceState(null, null, queryString);
    } else {
      window.history.replaceState(null, null, window.location.pathname);
    }

    fetchDevices(queryString);
  };

  const removeFilter = index => {
    setFilterQuery(filterQuery.filter((_, i) => i !== index));
  };

  const addFilter = filter => {
    setShowFilterDialog(false);
    if (filterToEdit !== null) {
      setFilterQuery(filterQuery =>
        filterQuery.map((f, index) => (index === filterToEdit ? filter : f))
      );
    } else {
      setFilterQuery(filterQuery => [...filterQuery, filter]);
    }
    setFilterToEdit(null);
  };

  const clearFilters = () => {
    setFilterQuery([]);
  };

  const buildFiltersQuery = () =>
    [...Array.from(new Set(filterQuery))].map(
      filter =>
        `${Params.Filter}=${encodeURIComponent(btoa(JSON.stringify(filter)))}`
    );

  const parseQueryString = () => {
    if (!window.location.search || window.location.search.length < 1) {
      return;
    }

    var builtQuery = [];
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
              );

              builtQuery.push(filter);
            } catch (e) {
              console.error('Error parsing filters:', e);
            }
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
    });
    setOrderedColumn(orderedColumn);
    setOrder(order);
    setFilterQuery(builtQuery);
  };

  return (
    <Layout>
      <Card
        title="Devices"
        size="full"
        left={
          deviceTotal && (
            <Text color="grays.8" fontWeight={2} fontSize={2}>
              (
              <Text color="white" fontWeight={1} fontSize={1} paddingX={1}>
                {filterQuery.length
                  ? `${devices.length} / ${deviceTotal}`
                  : deviceTotal}
              </Text>
              )
            </Text>
          )
        }
        center={
          <Row position="relative" alignItems="center">
            <Icon
              icon="search"
              size={16}
              color={searchFocused ? 'primary' : 'white'}
              style={{ position: 'absolute', left: 16 }}
            />
            <Input
              width="280px"
              bg="black"
              placeholder="Search by name or labels"
              paddingLeft={7}
              value={searchInput}
              onChange={e => setSearchInput(e.target.value)}
              onFocus={() => setSearchFocused(true)}
              onBlur={() => setSearchFocused(false)}
            />
          </Row>
        }
        actions={[
          ...(filterQuery.length
            ? [
                {
                  title: 'Clear Filters',
                  onClick: clearFilters,
                  variant: 'text',
                },
              ]
            : []),
          {
            title: 'Add Filter',
            variant: 'secondary',
            onClick: () => setShowFilterDialog(true),
          },
          {
            title: 'Register Device',
            href: `/${route.data.params.project}/devices/register`,
          },
        ]}
        maxHeight="100%"
      >
        {filterQuery.length > 0 && (
          <Row marginBottom={4}>
            <DevicesFilterButtons
              canRemoveFilter
              query={filterQuery}
              removeFilter={removeFilter}
              onEdit={index => {
                setFilterToEdit(index);
                setShowFilterDialog(true);
              }}
            />
          </Row>
        )}
        <Table
          columns={columns}
          data={tableData}
          rowHref={({ name }) =>
            `/${route.data.params.project}/devices/${name}`
          }
          placeholder={
            <Text>
              There are no <strong>Devices</strong>.
            </Text>
          }
        />
      </Card>

      {showFilterDialog && (
        <DevicesFilter
          filter={filterToEdit !== null && filterQuery[filterToEdit]}
          onClose={() => {
            setShowFilterDialog(false);
            setFilterToEdit(null);
          }}
          onSubmit={addFilter}
        />
      )}
    </Layout>
  );
};

export default Devices;
