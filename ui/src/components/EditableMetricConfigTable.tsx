import React, { Component, Fragment } from "react";
import axios from "axios";
import {
  Pane,
  Table,
  Dialog,
  IconButton,
  TextInputField,
  toaster
  // @ts-ignore
} from "evergreen-ui";

interface Props {
  setterEndpoint: string,
  configs: realConfig[],
}

interface State {
  configs: editableConfig[],
}

interface realConfig {
  metric: string,
  labels: string[],
  tags: string[],
}

interface editableConfig {
  values: configStringValues,
  editedValues?: configStringValues,
  mode: string,
}

interface configStringValues {
  metricStr: string,
  labelsStr: string,
  tagsStr: string,
}

const MODE = {
  NOEDIT: 'noedit',
  EDIT: 'edit',
}

export class EditableMetricConfigTable extends Component<Props, State> {
  constructor(props: Props) {
    super(props)

    this.state = {
      configs: this.props.configs.map(c => ({
        values: {
          metricStr: c.metric,
          labelsStr: c.labels.join(", "),
          tagsStr: c.tags.join(", "),
        },
        mode: MODE.NOEDIT,
      })),
    };
  }

  startConfigEdit = (i: number) => {
    var config = this.state.configs[i];
    config.mode = MODE.EDIT;
    config.editedValues = Object.assign({}, config.values);
    this.updateConfig(i, config);
  };

  editConfigField = (i: number, field: string, value: string) => {
    var config = this.state.configs[i];
    // @ts-ignore
    config.editedValues[field] = value;
    this.updateConfig(i, config);
  };

  updateConfig = (i: number, config: editableConfig, callback?: () => void) => {
    var configs = this.state.configs;
    configs[i] = config;
    this.setState({ configs }, callback);
  };

  // Final changes to config
  saveConfig = (i: number) => {
    var configs = this.state.configs;

    var config = configs[i];
    config.mode = MODE.NOEDIT;
    config.values = Object.assign({}, config.editedValues);
    config.values.labelsStr = reformatArrayString(config.values.labelsStr);
    config.values.tagsStr = reformatArrayString(config.values.tagsStr);
    config.editedValues = undefined;
    configs[i] = config;

    configs = configs.sort(configSortFn);
    this.setState({ configs }, this.syncUpdate);
  }

  revertConfig = (i: number) => {
    var config = this.state.configs[i];
    config.mode = MODE.NOEDIT;
    config.editedValues = undefined;
    this.updateConfig(i, config);
  }

  deleteConfig = (i: number) => {
    var configs = this.state.configs;
    configs.splice(i, 1);

    this.setState({ configs });
    this.syncUpdate();
  };

  addNewConfig = () => {
    var values = {
      metricStr: '',
      labelsStr: '',
      tagsStr: '',
    }

    var newConfig = {
      values: Object.assign({}, values),
      editedValues: Object.assign({}, values),
      mode: MODE.EDIT,
    };

    var configs = this.state.configs;
    configs.push(newConfig);
    this.setState({ configs });
  };

  syncUpdate = () => {
    axios
      .put(
        this.props.setterEndpoint,
        {
          configs: [{
            metrics: realConfigsFromEditableConfigs(this.state.configs)
          }],
        },
        {
          withCredentials: true
        }
      )
      .catch(error => {
        // if (utils.is4xx(error.response.status)) {
        //   this.setState({
        //     backendError: utils.convertErrorMessage(error.response.data)
        //   });
        // } else {
          console.log(error);
        // }
      });
  }

  renderMetricConfig(i: number) {
    var config = this.state.configs[i];
    switch (this.state.configs[i].mode) {
      case MODE.NOEDIT:
        return (
          <Fragment key={i}>
            <Table.Row>
              <Table.TextCell>{config.values.metricStr}</Table.TextCell>
              <Table.TextCell>{config.values.labelsStr}</Table.TextCell>
              <Table.TextCell>{config.values.tagsStr}</Table.TextCell>
              <Table.TextCell flexBasis={75} flexShrink={0} flexGrow={0}>
                <Pane display="flex">
                  <IconButton
                    icon="edit"
                    height={24}
                    appearance="minimal"
                    onClick={() => this.startConfigEdit(i)}
                  />
                  <IconButton
                    icon="trash"
                    height={24}
                    appearance="minimal"
                    // onClick={() => this.setShowRemoveDialog(i)}
                    onClick={() => this.deleteConfig(i)}
                  />
                </Pane>
              </Table.TextCell>
            </Table.Row>
            <Pane>
              <Dialog
                // isShown={config.showRemoveDialog}
                title="Remove config"
                intent="danger"
                // onCloseComplete={() => this.hideShowRemoveDialog(i)}
                onConfirm={() => this.deleteConfig(i)}
                confirmConfig="Remove config"
              >
                You are about to remove label <strong>{config.values.metricStr}</strong>
                .
              </Dialog>
            </Pane>
          </Fragment>
        );
      case MODE.EDIT:
        if (config.editedValues == undefined) {
          return <p>hi</p>;
        }
        return (
          <Table.Row key={i} height="auto">
            <Table.TextCell>
              <TextInputField
                label=""
                name={`edit-${config}${i}-metric`}
                value={config.editedValues.metricStr}
                onChange={(event: any) => {
                  this.editConfigField(i, 'metricStr', event.target.value);
                }}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`edit-${config}${i}-labels`}
                value={config.editedValues.labelsStr}
                onChange={(event: any) => {
                  this.editConfigField(i, 'labelsStr', event.target.value);
                }}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`edit-${config}${i}-tags`}
                value={config.editedValues.tagsStr}
                onChange={(event: any) => {
                  this.editConfigField(i, 'tagsStr', event.target.value);
                }}
                // isInvalid={config.valueValidationMessage !== null}
                // validationMessage={config.valueValidationMessage}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>

            <Table.TextCell flexBasis={75} flexShrink={0} flexGrow={0}>
              <Pane display="flex">
                <IconButton
                  icon="floppy-disk"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.saveConfig(i)}
                />
                <IconButton
                  icon="cross"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.revertConfig(i)}
                />
              </Pane>
            </Table.TextCell>
          </Table.Row>
        );
      default:
        return <Fragment />;
    }
  }

  render() {
    return (
      <Table>
        <Table.Head>
          <Table.TextHeaderCell>Metric</Table.TextHeaderCell>
          <Table.TextHeaderCell>Labels</Table.TextHeaderCell>
          <Table.TextHeaderCell>Tags</Table.TextHeaderCell>
          <Table.TextHeaderCell
            flexBasis={75}
            flexShrink={0}
            flexGrow={0}
          ></Table.TextHeaderCell>
        </Table.Head>
        <Table.Body>
          {this.state.configs.map((_, i) =>
            this.renderMetricConfig(i)
          )}
          <Table.Row key="add">
            <Table.TextCell>
              <IconButton
                icon="plus"
                height={24}
                appearance="minimal"
                onClick={() => this.addNewConfig()}
              />
            </Table.TextCell>
            <Table.TextCell></Table.TextCell>
            <Table.TextCell
              flexBasis={75}
              flexShrink={0}
              flexGrow={0}
            ></Table.TextCell>
          </Table.Row>
        </Table.Body>
      </Table>
    );
  }
}

function arrayFromArrayString(str: string) {
  return str.split(',').map(x => x.trim()).filter(x => x.length);
}

function reformatArrayString(str: string) {
  return arrayFromArrayString(str).join(", ");
}

function configSortFn(a: editableConfig, b: editableConfig) {
  var aMetric = a.values.metricStr;
  var bMetric = b.values.metricStr;
  if (a.editedValues) {
    aMetric = a.editedValues.metricStr;
  }
  if (b.editedValues) {
    bMetric = b.editedValues.metricStr;
  }

  // Make '*' sort to bottom, since it's lowest priority
  if (aMetric === '*') {
    return 1;
  }
  if (bMetric === '*') {
    return -1;
  }

  // Normal sort, otherwise
  if (aMetric < bMetric) {
    return -1;
  }
  if (aMetric === bMetric) {
    return 0;
  }
  return 1;
}

function realConfigsFromEditableConfigs(editableConfigs: editableConfig[]): realConfig[] {
  return editableConfigs.map(config => {
    return {
      metric: config.values.metricStr,
      labels: arrayFromArrayString(config.values.labelsStr),
      tags: arrayFromArrayString(config.values.tagsStr),
    }
  })
}