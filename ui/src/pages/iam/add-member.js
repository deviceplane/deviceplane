import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  majorScale,
  Button,
  Heading,
  Alert,
  toaster,
  Checkbox,
  TextInputField
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import segment from '../../segment';
import InnerCard from '../../components/InnerCard';

export default class AddMember extends Component {
  state = {
    email: '',
    emailValidationMessage: null,
    backendError: null,
    roleBindings: []
  };

  componentDidMount() {
    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/roles`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          roleBindings: this.createRoleBindings(response.data)
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  createRoleBindings = allRoles => {
    var roleBindings = [];
    for (var i = 0; i < allRoles.length; i++) {
      roleBindings.push({
        id: allRoles[i].id,
        name: allRoles[i].name,
        hasRoleBinding: false
      });
    }
    return roleBindings;
  };

  handleUpdateEmail = event => {
    this.setState({
      email: event.target.value
    });
  };

  handleAddRoleBindings = event => {
    const roleBindings = this.state.roleBindings;
    var updatedRoleBindings = [];
    var hasRole = false;
    for (var i = 0; i < roleBindings.length; i++) {
      hasRole = roleBindings[i].hasRoleBinding;
      if (roleBindings[i].id === event.target.id) {
        hasRole = event.target.checked;
      }
      updatedRoleBindings.push({
        id: roleBindings[i].id,
        name: roleBindings[i].name,
        hasRoleBinding: hasRole
      });
    }
    this.setState({
      roleBindings: updatedRoleBindings
    });
  };

  handleSubmit = () => {
    this.setState({
      emailValidationMessage: null,
      backendError: null
    });

    if (!utils.emailRegex.test(this.state.email)) {
      this.setState({
        emailValidationMessage: 'Please enter a valid email.'
      });
      return;
    }

    if (this.state.emailValidationMessage !== null) {
      this.setState({
        emailValidationMessage: null
      });
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/memberships`,
        {
          email: this.state.email
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        segment.track('Member Added');

        const userId = response.data.userId;
        const roleBindings = this.state.roleBindings;
        var noError = true;

        for (var i = 0; i < roleBindings.length; i++) {
          const roleId = roleBindings[i].id;
          if (roleBindings[i].hasRoleBinding && noError) {
            noError = this.addRole(userId, roleId);
          }
        }

        if (noError) {
          this.props.history.push(`/${this.props.projectName}/iam/members/`);
          toaster.success('Member was added successfully.');
        } else {
          toaster.warning(
            'Member was added successfully, but role bindings for the member were not updated properly. Please check the roles of the member.'
          );
        }
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          console.log(error);
          toaster.danger('Member was not added.');
        }
      });
  };

  addRole = (userId, roleId) => {
    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/memberships/${userId}/roles/${roleId}/membershiprolebindings`,
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

  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/iam/members/`);
  }

  render() {
    const roleBindings = this.state.roleBindings;
    return (
      <Pane width="50%">
        <InnerCard>
          <Pane padding={majorScale(4)}>
            {this.state.backendError && (
              <Alert
                marginBottom={majorScale(2)}
                paddingTop={majorScale(2)}
                paddingBottom={majorScale(2)}
                intent="warning"
                title={this.state.backendError}
              />
            )}
            <Heading paddingBottom={majorScale(4)} size={600}>
              Add Member
            </Heading>
            <TextInputField
              label="Email"
              onChange={this.handleUpdateEmail}
              value={this.state.email}
              isInvalid={this.state.emailValidationMessage !== null}
              validationMessage={this.state.emailValidationMessage}
            />
            <Heading size={500}>Choose Individual Roles</Heading>
            {roleBindings.map(role => (
              <Checkbox
                key={role.id}
                id={role.id}
                label={role.name}
                checked={role.hasRoleBinding}
                onChange={event => this.handleAddRoleBindings(event)}
              />
            ))}
            <Pane display="flex" flex="row">
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                onClick={() => this.handleSubmit()}
              >
                Add Member
              </Button>
              <Button
                marginTop={majorScale(2)}
                marginLeft={majorScale(2)}
                onClick={() => this.handleCancel()}
              >
                Cancel
              </Button>
            </Pane>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}
