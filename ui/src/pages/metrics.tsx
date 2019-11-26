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
import { EditableMetricConfigTable } from '../components/EditableMetricConfigTable';
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

  stateMetricsConfig?: any[],
  hostMetricsConfig?: any[],
  serviceMetricsConfig?: any[],
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

      stateMetricsConfig: undefined,
      hostMetricsConfig: undefined,
      serviceMetricsConfig: undefined,
    };
  }

  componentDidMount() {
    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/metrictargetconfig/state`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          stateMetricsConfig: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });

    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/metrictargetconfig/host`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          hostMetricsConfig: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });

    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/metrictargetconfig/service`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          serviceMetricsConfig: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  renderMetricConfigTable = (params: any, metricRows: any[]) => {
    return (
      <Table>
        <Table.Head background="tint2">
          <Table.TextHeaderCell
            marginY={minorScale(1)}
          >
            Metric
          </Table.TextHeaderCell>

          <Table.TextHeaderCell
            marginY={minorScale(1)}
          >
            Labels
          </Table.TextHeaderCell>

          <Table.TextHeaderCell
            marginY={minorScale(1)}
          >
            Tags
          </Table.TextHeaderCell>
        </Table.Head>
        <Table.Body>{metricRows.map(this.renderMetricRow)}</Table.Body>
      </Table>
    )
  }

  renderMetricRow = (metric: any) => {
    return (
      <Table.Row
        key={metric.metric}
        isSelectable
        flexGrow={1}
        height="auto"
        paddingY={majorScale(1)}
        alignItems="flex-start"
      >
        <Table.TextCell
          marginY={minorScale(1)}
        >
          {metric.metric}
        </Table.TextCell>

        <Table.TextCell
          marginY={minorScale(1)}
        >
          {JSON.stringify(metric.labels)}
        </Table.TextCell>

        <Table.TextCell
          marginY={minorScale(1)}
        >
          {JSON.stringify(metric.tags)}
        </Table.TextCell>
      </Table.Row>
    );
  };

  render() {
    return (
      <Fragment>
        <TopHeader user={this.props.user} heading="Metrics" history={this.props.history} />
        <Pane width="70%">
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading padding={majorScale(2)}>State Metrics Config</Heading>
            </Pane>
            {this.state.stateMetricsConfig && <EditableMetricConfigTable
              setterEndpoint={
                `${config.endpoint}/projects/${this.props.projectName}/metrictargetconfig/state`
              }
              configs={
                this.state.stateMetricsConfig.length == 1 ?
                  this.state.stateMetricsConfig[0].metrics :
                  []
              }
            />}
          </InnerCard>
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading padding={majorScale(2)}>Host Metrics Config</Heading>
            </Pane>
            {this.state.hostMetricsConfig && <EditableMetricConfigTable
              setterEndpoint={
                `${config.endpoint}/projects/${this.props.projectName}/metrictargetconfig/host`
              }
              configs={
                this.state.hostMetricsConfig.length == 1 ?
                  this.state.hostMetricsConfig[0].metrics :
                  []
              }
            />}
          </InnerCard>
        </Pane>
      </Fragment>
    );
  }
}

export default withTheme(Devices);
