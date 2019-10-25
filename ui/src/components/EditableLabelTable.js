import React, { Component, Fragment } from "react";
import axios from "axios";
import {
  Pane,
  Table,
  Dialog,
  IconButton,
  TextInputField,
  toaster
} from "evergreen-ui";

import utils from "../utils";
import config from "../config";

export class EditableLabelTable extends Component {
  constructor(props) {
    super(props)

    this.state = {
      labels: []
    };
  }

  componentDidMount() {
    axios
      .get(
        this.props.getEndpoint,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          labels: this.initializeLabels(response.data.labels)
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  initializeLabels = keyValues => {
    var labels = [];
    var labelKeys = Object.keys(keyValues);
    for (var i = 0; i < labelKeys.length; i++) {
      labels.push({
        key: labelKeys[i],
        value: keyValues[labelKeys[i]],
        mode: "default",
        keyValidationMessage: null,
        valueValidationMessage: null,
        showRemoveDialog: false
      });
    }
    return labels.sort(function(a, b) {
      if (a.key < b.key) {
        return -1;
      }
      if (a.key > b.key) {
        return 1;
      }
      return 0;
    });
  };

  createNewLabel = () => {
    var labels = this.state.labels;
    labels.push({
      key: "",
      value: "",
      mode: "new",
      keyValidationMessage: null,
      valueValidationMessage: null,
      showRemoveDialog: false
    });
    this.setState({
      labels: labels
    });
  };

  handleUpdate = (i, property) => {
    return event => {
      var labels = this.state.labels;
      labels[i][property] = event.target.value;
      this.setState({
        labels: labels
      });
    };
  };

  setEdit = i => {
    var editLabels = this.state.labels;
    editLabels[i]["mode"] = "edit";
    this.setState({
      labels: editLabels
    });
  };

  cancelEdit = i => {
    var editLabels = this.state.labels;
    editLabels[i]["mode"] = "default";
    this.setState({
      labels: editLabels
    });
  };

  setShowRemoveDialog = i => {
    var showRemoveDialogLabels = this.state.labels;
    showRemoveDialogLabels[i]["showRemoveDialog"] = true;
    this.setState({
      labels: showRemoveDialogLabels
    });
  };

  hideShowRemoveDialog = i => {
    var showRemoveDialogLabels = this.state.labels;
    showRemoveDialogLabels[i]["showRemoveDialog"] = false;
    this.setState({
      labels: showRemoveDialogLabels
    });
  };

  setLabel = (key, value, i) => {
    var updatedLabels = this.state.labels;
    var keyValidationMessage = utils.checkName("key", key);
    var valueValidationMessage = utils.checkName("value", value);

    if (keyValidationMessage === null) {
      for (var j = 0; j < updatedLabels.length; j++) {
        if (i !== j && key === updatedLabels[j]["key"]) {
          keyValidationMessage = "Key already exists.";
          break;
        }
      }
    }

    updatedLabels[i]["keyValidationMessage"] = keyValidationMessage;
    updatedLabels[i]["valueValidationMessage"] = valueValidationMessage;

    this.setState({
      labels: updatedLabels
    });

    if (
      keyValidationMessage === null &&
      valueValidationMessage === null &&
      key !== null &&
      value !== null
    ) {
      axios
        .put(
          this.props.setEndpoint,
          {
            key: key,
            value: value
          },
          {
            withCredentials: true
          }
        )
        .then(response => {
          var updatedLabels = this.state.labels;
          updatedLabels[i]["mode"] = "default";
          this.setState({
            labels: updatedLabels
          });
        })
        .catch(error => {
          console.log(error);
        });
    }
  };

  deleteLabel = (key, i) => {
    if (key !== "") {
      axios
        .delete(
          this.props.deleteEndpoint + `/${key}`,
          {
            withCredentials: true
          }
        )
        .then(response => {
          var removedLabels = this.state.labels;
          removedLabels.splice(i, 1);
          this.setState({
            labels: removedLabels
          });
        })
        .catch(error => {
          var hideRemoveDialogLabels = this.state.labels;
          hideRemoveDialogLabels[i]["showRemoveDialog"] = false;
          this.setState({
            labels: hideRemoveDialogLabels
          });
          toaster.danger("Label was not removed.");
          console.log(error);
        });
    } else {
      var removedLabels = this.state.labels;
      removedLabels.splice(i, 1);
      this.setState({
        labels: removedLabels
      });
    }
  };

  renderLabel(Label, i) {
    switch (Label.mode) {
      case "default":
        return (
          <Fragment key={Label.key}>
            <Table.Row>
              <Table.TextCell>{Label.key}</Table.TextCell>
              <Table.TextCell>{Label.value}</Table.TextCell>
              <Table.TextCell flexBasis={75} flexShrink={0} flexGrow={0}>
                <Pane display="flex">
                  <IconButton
                    icon="edit"
                    height={24}
                    appearance="minimal"
                    onClick={() => this.setEdit(i)}
                  />
                  <IconButton
                    icon="trash"
                    height={24}
                    appearance="minimal"
                    onClick={() => this.setShowRemoveDialog(i)}
                  />
                </Pane>
              </Table.TextCell>
            </Table.Row>
            <Pane>
              <Dialog
                isShown={Label.showRemoveDialog}
                title="Remove Label"
                intent="danger"
                onCloseComplete={() => this.hideShowRemoveDialog(i)}
                onConfirm={() => this.deleteLabel(Label.key, i)}
                confirmLabel="Remove Label"
              >
                You are about to remove label <strong>{Label.key}</strong>
                .
              </Dialog>
            </Pane>
          </Fragment>
        );
      case "edit":
        return (
          <Table.Row key={Label.key} height="auto">
            <Table.TextCell>{Label.key}</Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`edit-${Label.key}`}
                value={Label.value}
                onChange={event => this.handleUpdate(i, "value")(event)}
                isInvalid={Label.valueValidationMessage !== null}
                validationMessage={Label.valueValidationMessage}
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
                  onClick={() =>
                    this.setLabel(Label.key, Label.value, i)
                  }
                />
                <IconButton
                  icon="cross"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.cancelEdit(i)}
                />
              </Pane>
            </Table.TextCell>
          </Table.Row>
        );
      case "new":
        return (
          <Table.Row key={`new-${i}`} height="auto">
            <Table.TextCell>
              <TextInputField
                label=""
                name={`new-key-${i}`}
                value={Label.key}
                onChange={event => this.handleUpdate(i, "key")(event)}
                isInvalid={Label.keyValidationMessage !== null}
                validationMessage={Label.keyValidationMessage}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`new-value-${i}`}
                value={Label.value}
                onChange={event => this.handleUpdate(i, "value")(event)}
                isInvalid={Label.valueValidationMessage !== null}
                validationMessage={Label.valueValidationMessage}
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
                  onClick={() =>
                    this.setLabel(Label.key, Label.value, i)
                  }
                />
                <IconButton
                  icon="cross"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.deleteLabel(Label.key, i)}
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
          <Table.TextHeaderCell>Key</Table.TextHeaderCell>
          <Table.TextHeaderCell>Value</Table.TextHeaderCell>
          <Table.TextHeaderCell
            flexBasis={75}
            flexShrink={0}
            flexGrow={0}
          ></Table.TextHeaderCell>
        </Table.Head>
        <Table.Body>
          {this.state.labels.map((label, i) =>
            this.renderLabel(label, i)
          )}
          <Table.Row key="add">
            <Table.TextCell>
              <IconButton
                icon="plus"
                height={24}
                appearance="minimal"
                onClick={() => this.createNewLabel()}
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
