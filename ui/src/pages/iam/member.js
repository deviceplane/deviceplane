import React, { Component, Fragment } from 'react';
import axios from 'axios';
import {
  Pane,
  Dialog,
  majorScale,
  Button,
  Heading,
  Text,
  Checkbox,
  toaster
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import InnerCard from '../../components/InnerCard';
import CustomSpinner from '../../components/CustomSpinner';

export default class Member extends Component {
  state = {
    member: null,
    allRoles: [],
    roleBindings: [],
    unchanged: true,
    showRemoveDialog: false,
    isSelf: this.props.user.id === this.props.userId
  };

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/memberships/${this.props.userId}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          member: response.data,
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
            this.state.member,
            response.data
          )
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  createRoleBindings = (member, allRoles) => {
    var roleBindings = [];
    if (member !== null) {
      for (var i = 0; i < allRoles.length; i++) {
        var hasRole = false;
        if (member.roles && member.roles.length > 0) {
          for (var j = 0; j < member.roles.length; j++) {
            if (allRoles[i].id === member.roles[j].id) {
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

  handleUpdate = () => {
    const updatedRoleBindings = this.state.roleBindings;
    const currentRoles = this.state.member.roles;
    var addRoleBindings = [];
    var removeRoleBindings = [];

    for (var i = 0; i < updatedRoleBindings.length; i++) {
      var addRole = false;
      var removeRole = false;
      // add role binding to member
      if (updatedRoleBindings[i].hasRoleBinding) {
        addRole = true;
      }
      //check if role binding already exists on member
      if (currentRoles && currentRoles.length > 0) {
        for (var j = 0; j < currentRoles.length; j++) {
          if (updatedRoleBindings[i].id === currentRoles[j].id) {
            if (updatedRoleBindings[i].hasRoleBinding) {
              //if role binding already exists on member, do not re-add role to member
              addRole = false;
              break;
            } else {
              //if role binding already exists on member, remove the role binding
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

    var noError = true;

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
      toaster.success('Member updated successfully.');
    } else {
      toaster.danger('Member was not updated.');
    }
  };

  addRole = roleId => {
    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/memberships/${this.state.member.userId}/roles/${roleId}/membershiprolebindings`,
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
        `${config.endpoint}/projects/${this.props.projectName}/memberships/${this.state.member.userId}/roles/${roleId}/membershiprolebindings`,
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

  handleRemove = () => {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/memberships/${this.state.member.userId}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Successfully removed member.');
        this.props.history.push(`/${this.props.projectName}/iam/members`);
      })
      .catch(error => {
        this.setState({
          showRemoveDialog: false
        });
        toaster.danger('Member was not removed.');
        console.log(error);
      });
  };

  render() {
    const member = this.state.member;
    const roleBindings = this.state.roleBindings;
    return (
      <Pane width="70%">
        {member ? (
          <InnerCard>
            <Pane padding={majorScale(4)}>
              <Heading
                size={600}
              >{`${member.user.firstName} ${member.user.lastName}`}</Heading>
              <Text size={400}>{member.user.email}</Text>
              <Heading size={500} paddingTop={majorScale(4)}>
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
                Update Member
              </Button>
            </Pane>
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
                disabled={this.state.isSelf}
                onClick={() => this.setState({ showRemoveDialog: true })}
              >
                Remove Member...
              </Button>
            </Pane>
            <Pane>
              <Dialog
                isShown={this.state.showRemoveDialog}
                title="Remove Member"
                intent="danger"
                onCloseComplete={() =>
                  this.setState({ showRemoveDialog: false })
                }
                onConfirm={() => this.handleRemove()}
                confirmLabel="Remove Member"
              >
                You are about to remove the member (
                <strong>
                  {member.user.firstName} {member.user.lastName}
                </strong>
                ) from the project.
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
