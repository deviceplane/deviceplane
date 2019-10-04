import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  Table,
  Dialog,
  majorScale,
  Button,
  Heading,
  Alert,
  toaster,
  TextInputField,
  Label,
  Textarea,
  Checkbox,
  Code
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import CustomSpinner from '../../components/CustomSpinner';
import InnerCard from '../../components/InnerCard';

export default class ServiceAccount extends Component {
  state = {
    serviceAccount: null,
    name: '',
    nameValidationMessage: null,
    description: '',
    allRoles: [],
    roleBindings: [],
    unchanged: true,
    showDeleteDialog: false
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.props.serviceAccountName}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          serviceAccount: response.data,
          name: response.data.name,
          description: response.data.description,
          roleBindings: this.createRoleBindings(
            response.data,
            this.state.allRoles
          )
        });
      })
      .catch(error => {
        console.log(error);
      });
    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/roles`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          allRoles: response.data,
          roleBindings: this.createRoleBindings(
            this.state.serviceAccount,
            response.data
          )
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  createRoleBindings = (serviceAccount, allRoles) => {
    var roleBindings = [];
    if (serviceAccount !== null) {
      for (var i = 0; i < allRoles.length; i++) {
        var hasRole = false;
        if (serviceAccount.roles && serviceAccount.roles.length > 0) {
          for (var j = 0; j < serviceAccount.roles.length; j++) {
            if (allRoles[i].id === serviceAccount.roles[j].id) {
              hasRole = true;
              break;
            }
          }
        }
        roleBindings.push({
          id: allRoles[i].id,
          name: allRoles[i].name,
          hasRoleBinding: hasRole
        });
      }
    }
    return roleBindings;
  };

  handleUpdateName = event => {
    this.setState({
      name: event.target.value,
      unchanged: false
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value,
      unchanged: false
    });
  };

  handleUpdateRoles = event => {
    const currentRoleBindings = this.state.roleBindings;
    var newRoleBindings = [];
    var hasRole = false;
    for (var i = 0; i < currentRoleBindings.length; i++) {
      hasRole = currentRoleBindings[i].hasRoleBinding;
      if (currentRoleBindings[i].id === event.target.id) {
        hasRole = event.target.checked;
      }
      newRoleBindings.push({
        id: currentRoleBindings[i].id,
        name: currentRoleBindings[i].name,
        hasRoleBinding: hasRole
      });
    }
    this.setState({
      roleBindings: newRoleBindings,
      unchanged: false
    });
  };

  handleUpdate() {
    var noError = true;
    var nameValidationMessage = utils.checkName(
      'service account',
      this.state.name
    );

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .put(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.state.serviceAccount.id}`,
        {
          name: this.state.name,
          description: this.state.description
        },
        {
          withCredentials: true
        }
      )
      .catch(error => {
        noError = false;
        console.log(error);
      });

    const updatedRoleBindings = this.state.roleBindings;
    const currentRoles = this.state.serviceAccount.roles;
    var addRoleBindings = [];
    var removeRoleBindings = [];

    for (var i = 0; i < updatedRoleBindings.length; i++) {
      var addRole = false;
      var removeRole = false;
      // add role binding to service account
      if (updatedRoleBindings[i].hasRoleBinding) {
        addRole = true;
      }
      //check if role binding already exists on service account
      if (currentRoles && currentRoles.length > 0) {
        for (var j = 0; j < currentRoles.length; j++) {
          if (updatedRoleBindings[i].id === currentRoles[j].id) {
            if (updatedRoleBindings[i].hasRoleBinding) {
              //if role binding already exists on service account, do not re-add role to service account
              addRole = false;
              break;
            } else {
              //if role binding already exists on service account, remove the role binding
              removeRole = true;
              break;
            }
          }
        }
      }
      if (addRole) {
        addRoleBindings.push(updatedRoleBindings[i]);
      }
      if (removeRole) {
        removeRoleBindings.push(updatedRoleBindings[i]);
      }
    }

    for (var k = 0; k < addRoleBindings.length; k++) {
      const roleId = addRoleBindings[k].id;
      if (noError) {
        noError = this.addRole(roleId);
      }
    }

    for (var l = 0; l < removeRoleBindings.length; l++) {
      const roleId = removeRoleBindings[l].id;
      if (noError) {
        noError = this.removeRole(roleId);
      }
    }

    if (noError) {
      toaster.success('Service account updated successfully.');
    } else {
      toaster.danger('Service account was not updated.');
    }
  }

  addRole = roleId => {
    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.state.serviceAccount.id}/roles/${roleId}/serviceaccountrolebindings`,
        {},
        {
          withCredentials: true
        }
      )
      .catch(error => {
        console.log(error);
        return false;
      });
    return true;
  };

  removeRole = roleId => {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.state.serviceAccount.id}/roles/${roleId}/serviceaccountrolebindings`,
        {
          withCredentials: true
        }
      )
      .catch(error => {
        console.log(error);
        return false;
      });
    return true;
  };

  handleDelete() {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.state.serviceAccount.id}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted service account.');
        this.props.history.push(
          `/${this.props.projectName}/iam/serviceaccounts`
        );
      })
      .catch(error => {
        this.setState({
          showDeleteDialog: false
        });
        console.log(error);
        toaster.danger('Service account was not deleted.');
      });
  }

  render() {
    const serviceAccount = this.state.serviceAccount;
    const roleBindings = this.state.roleBindings;
    return (
      <Pane width="70%">
        {serviceAccount ? (
          <InnerCard>
            <Pane padding={majorScale(4)}>
              <Heading paddingBottom={majorScale(2)} size={600}>
                Service Account / {serviceAccount.name}
              </Heading>
              <TextInputField
                label="Name"
                onChange={this.handleUpdateName}
                value={this.state.name}
                validationMessage={this.state.nameValidationMessage}
              />
              <Label
                htmlFor="description-textarea"
                marginBottom="4"
                display="block"
              >
                Description (optional)
              </Label>
              <Textarea
                id="description-textarea"
                height="100px"
                onChange={this.handleUpdateDescription}
                value={this.state.description}
              />
              <Heading size={500} paddingTop={majorScale(2)}>
                Choose Individual Roles
              </Heading>
              {roleBindings.map(role => (
                <Checkbox
                  key={role.id}
                  id={role.id}
                  label={role.name}
                  checked={role.hasRoleBinding}
                  onChange={event => this.handleUpdateRoles(event)}
                />
              ))}
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                disabled={this.state.unchanged}
                onClick={() => this.handleUpdate()}
              >
                Update Service Account
              </Button>
            </Pane>
            <ServiceAccountAccessKeys
              projectName={this.props.projectName}
              serviceAccount={serviceAccount}
              history={this.props.history}
            />
            <Pane
              borderTop="default"
              marginRight={majorScale(4)}
              marginLeft={majorScale(4)}
              marginBottom={majorScale(4)}
            >
              <Button
                marginTop={majorScale(4)}
                iconBefore="trash"
                intent="danger"
                onClick={() => this.setState({ showDeleteDialog: true })}
              >
                Delete Service Account...
              </Button>
            </Pane>
            <Pane>
              <Dialog
                isShown={this.state.showDeleteDialog}
                title="Delete Service Account"
                intent="danger"
                onCloseComplete={() =>
                  this.setState({ showDeleteDialog: false })
                }
                onConfirm={() => this.handleDelete()}
                confirmLabel="Delete Service Account"
              >
                You are about to delete the{' '}
                <strong>{this.props.serviceAccountName}</strong> service
                account.
              </Dialog>
            </Pane>
          </InnerCard>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}

class ServiceAccountAccessKeys extends Component {
  state = {
    accessKeys: null,
    newAccessKey: null,
    showAccessKeyCreated: false,
    backendError: null
  };

  componentDidMount() {
    this.loadAccessKeys();
  }

  loadAccessKeys() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.props.serviceAccount.id}/serviceaccountaccesskeys`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          accessKeys: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  createAccessKey() {
    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.props.serviceAccount.id}/serviceaccountaccesskeys`,
        {},
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showAccessKeyCreated: true,
          newAccessKey: response.data.value
        });
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Access key was not created successfully.');
          console.log(error);
        }
      });
  }

  deleteAccessKey = event => {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts/${this.props.serviceAccount.id}/serviceaccountaccesskeys/${event.target.id}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Successfully deleted access key.');
        this.loadAccessKeys();
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Access key was not deleted.');
          console.log(error);
        }
      });
  };

  closeAccessKeyDialog() {
    this.setState({
      showAccessKeyCreated: false
    });
    this.loadAccessKeys();
  }

  render() {
    const accessKeys = this.state.accessKeys;
    return (
      <React.Fragment>
        <Pane
          borderTop="default"
          marginRight={majorScale(4)}
          marginLeft={majorScale(4)}
          marginBottom={majorScale(4)}
        >
          <Pane
            display="flex"
            flexDirection="row"
            justifyContent="space-between"
            alignItems="center"
            paddingTop={majorScale(2)}
            paddingBottom={majorScale(2)}
          >
            <Heading>Access Keys</Heading>
            <Button appearance="primary" onClick={() => this.createAccessKey()}>
              Create Access Key
            </Button>
          </Pane>
          {this.state.backendError && (
            <Alert
              marginBottom={majorScale(2)}
              paddingTop={majorScale(2)}
              paddingBottom={majorScale(2)}
              intent="warning"
              title={this.state.backendError}
            />
          )}
          {accessKeys && accessKeys.length > 0 && (
            <Table>
              <Table.Head>
                <Table.TextHeaderCell>Access Key ID</Table.TextHeaderCell>
                <Table.TextHeaderCell>Created At</Table.TextHeaderCell>
                <Table.TextHeaderCell></Table.TextHeaderCell>
              </Table.Head>
              <Table.Body>
                {accessKeys.map(accessKey => (
                  <Table.Row key={accessKey.id}>
                    <Table.TextCell>{accessKey.id}</Table.TextCell>
                    <Table.TextCell>{accessKey.createdAt}</Table.TextCell>
                    <Table.TextCell>
                      <Button
                        iconBefore="trash"
                        intent="danger"
                        id={accessKey.id}
                        onClick={event => this.deleteAccessKey(event)}
                      >
                        Delete Access Key
                      </Button>
                    </Table.TextCell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          )}
        </Pane>
        <Pane>
          <Dialog
            isShown={this.state.showAccessKeyCreated}
            title="Access Key Created"
            onCloseComplete={() => this.closeAccessKeyDialog()}
            hasFooter={false}
          >
            <Pane display="flex" flexDirection="column">
              <Heading
                paddingTop={majorScale(2)}
                paddingBottom={majorScale(2)}
              >{`Access Key: `}</Heading>
              <Pane marginBottom={majorScale(4)}>
                <Code>{this.state.newAccessKey}</Code>
              </Pane>
            </Pane>
            <Alert
              intent="warning"
              title="Save the info above! This is the only time you'll be able to use it."
            >
              {`If you lose it, you'll need to create a new access key.`}
            </Alert>
            <Button
              marginTop={16}
              appearance="primary"
              onClick={() => this.closeAccessKeyDialog()}
            >
              Close
            </Button>
          </Dialog>
        </Pane>
      </React.Fragment>
    );
  }
}
