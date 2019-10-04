import React, { Component } from 'react';
import {
  Textarea,
  Button,
  TextInputField,
  Pane,
  Tablist,
  SidebarTab,
  Tab,
  Table,
  Heading,
  Icon,
  Text,
  Code,
  Card,
  Label,
  Dialog,
  BackButton,
  IconButton,
  Badge,
  majorScale,
  minorScale,
  toaster,
  Link,
  Checkbox,
  Alert,
  SideSheet,
  TabNavigation
} from 'evergreen-ui';
import {
  BrowserRouter as Router,
  Route,
  Redirect,
  Switch
} from 'react-router-dom';
import axios from 'axios';
import moment from 'moment';

import './App.css';

import config from './config';
import segment from './segment';
import utils from './utils';

import Devices from './pages/devices';
import Applications from './pages/applications';
import Settings from './pages/settings';
import Projects from './pages/projects';
import Register from './pages/register';
import Login from './pages/login';

import CreateProject from './components/CreateProject';
import CustomSpinner from './components/CustomSpinner';
import DeviceSsh from './components/DeviceSsh';
import Editor from './components/Editor';
import InnerCard from './components/InnerCard';
import ResetPassword from './components/ResetPassword';
import TopHeader from './components/TopHeader';
import Logo from './components/logo';

class Iam extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tabs: ['members', 'serviceaccounts', 'roles'],
      tabLabels: ['Members', 'Service Accounts', 'Roles']
    };
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    var selectedIndex = 0;
    switch (match.params.iamTab) {
      case 'members':
        selectedIndex = 0;
        break;
      case 'serviceaccounts':
        selectedIndex = 1;
        break;
      case 'roles':
        selectedIndex = 2;
        break;
      default:
        this.props.history.push(`/${projectName}/iam`);
    }
    return (
      <Tablist border="default">
        {this.state.tabs.map((tab, index) => (
          <Tab
            key={tab}
            id={tab}
            onSelect={() =>
              this.props.history.push(`/${projectName}/iam/${tab}`)
            }
            isSelected={index === selectedIndex}
          >
            {this.state.tabLabels[index]}
          </Tab>
        ))}
      </Tablist>
    );
  };

  renderInner = match => {
    const user = this.props.user;
    const projectName = this.props.projectName;
    switch (match.params.iamTab) {
      case 'members':
        return (
          <MembersRouter
            user={user}
            projectName={projectName}
            match={match}
            history={this.props.history}
          />
        );
      case 'serviceaccounts':
        return (
          <ServiceAccountsRouter
            projectName={projectName}
            match={match}
            history={this.props.history}
          />
        );
      case 'roles':
        return (
          <RolesRouter
            projectName={projectName}
            match={match}
            history={this.props.history}
          />
        );
      default:
        return <Pane></Pane>;
    }
  };

  render() {
    const heading = 'IAM';
    return (
      <React.Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        <Pane
          display="flex"
          flexDirection="column"
          alignItems="center"
          background="white"
          width="100%"
          padding={majorScale(1)}
          borderBottom="default"
        >
          {this.renderTablist(this.props.match)}
        </Pane>
        {this.renderInner(this.props.match)}
      </React.Fragment>
    );
  }
}

class MembersRouter extends Component {
  render() {
    const projectName = this.props.projectName;
    const match = this.props.match;
    const user = this.props.user;
    return (
      <Switch>
        <Route
          path={`${match.path}/add`}
          render={route => (
            <AddMember projectName={projectName} history={route.history} />
          )}
        />
        <Route
          path={`${match.path}/:userId`}
          render={route => (
            <Member
              user={user}
              projectName={projectName}
              userId={route.match.params.userId}
              history={route.history}
            />
          )}
        />
        <Route
          exact
          path={match.path}
          render={route => (
            <Members
              projectName={projectName}
              match={match}
              history={route.history}
            />
          )}
        />
      </Switch>
    );
  }
}

class Members extends Component {
  constructor(props) {
    super(props);
    this.state = {
      members: []
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/memberships?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          members: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const members = this.state.members;
    return (
      <Pane width="70%">
        {members ? (
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Members</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/iam/members/add`
                  )
                }
              >
                Add Member
              </Button>
            </Pane>
            {members && members.length > 0 && (
              <Table>
                <Table.Head>
                  <Table.TextHeaderCell>Email</Table.TextHeaderCell>
                  <Table.TextHeaderCell>Name</Table.TextHeaderCell>
                  <Table.TextHeaderCell>Roles</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  {members.map(member => (
                    <Table.Row
                      key={member.userId}
                      isSelectable
                      onSelect={() =>
                        this.props.history.push(
                          `/${this.props.projectName}/iam/members/${member.userId}`
                        )
                      }
                    >
                      <Table.TextCell>{member.user.email}</Table.TextCell>
                      <Table.TextCell>{`${member.user.firstName} ${member.user.lastName}`}</Table.TextCell>
                      <Table.TextCell>
                        {member.roles.map(role => role.name).join(',')}
                      </Table.TextCell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
            )}
          </InnerCard>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}

class Member extends Component {
  constructor(props) {
    super(props);
    this.state = {
      member: null,
      allRoles: [],
      roleBindings: [],
      unchanged: true,
      showRemoveDialog: false,
      isSelf: this.props.user.id === this.props.userId
    };
  }

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

class AddMember extends Component {
  constructor(props) {
    super(props);
    this.state = {
      email: '',
      emailValidationMessage: null,
      backendError: null,
      roleBindings: []
    };
  }

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

class ServiceAccountsRouter extends Component {
  render() {
    const projectName = this.props.projectName;
    const match = this.props.match;
    return (
      <Switch>
        <Route
          path={`${match.path}/create`}
          render={route => (
            <CreateServiceAccount
              projectName={projectName}
              history={this.props.history}
            />
          )}
        />
        <Route
          path={`${match.path}/:serviceAccountName`}
          render={route => (
            <ServiceAccount
              projectName={projectName}
              serviceAccountName={route.match.params.serviceAccountName}
              history={this.props.history}
            />
          )}
        />
        <Route
          exact
          path={match.path}
          render={route => (
            <ServiceAccounts
              projectName={projectName}
              history={this.props.history}
            />
          )}
        />
      </Switch>
    );
  }
}

class ServiceAccounts extends Component {
  constructor(props) {
    super(props);
    this.state = {
      serviceAccounts: []
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          serviceAccounts: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const serviceAccounts = this.state.serviceAccounts;
    return (
      <Pane width="70%">
        {serviceAccounts ? (
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Service Accounts</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/iam/serviceaccounts/create`
                  )
                }
              >
                Create Service Account
              </Button>
            </Pane>
            {this.state.serviceAccounts &&
              this.state.serviceAccounts.length > 0 && (
                <Table>
                  <Table.Head>
                    <Table.TextHeaderCell>Service Account</Table.TextHeaderCell>
                  </Table.Head>
                  <Table.Body>
                    {serviceAccounts.map(serviceAccount => (
                      <Table.Row
                        key={serviceAccount.id}
                        isSelectable
                        onSelect={() =>
                          this.props.history.push(
                            `/${this.props.projectName}/iam/serviceaccounts/${serviceAccount.name}`
                          )
                        }
                      >
                        <Table.TextCell>{serviceAccount.name}</Table.TextCell>
                      </Table.Row>
                    ))}
                  </Table.Body>
                </Table>
              )}
          </InnerCard>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}

class ServiceAccount extends Component {
  constructor(props) {
    super(props);
    this.state = {
      serviceAccount: null,
      name: '',
      nameValidationMessage: null,
      description: '',
      allRoles: [],
      roleBindings: [],
      unchanged: true,
      showDeleteDialog: false
    };
  }

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
  constructor(props) {
    super(props);
    this.state = {
      accessKeys: null,
      newAccessKey: null,
      showAccessKeyCreated: false,
      backendError: null
    };
  }

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

class CreateServiceAccount extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      nameValidationMessage: null,
      description: '',
      backendError: null
    };
  }

  handleUpdateName = event => {
    this.setState({
      name: event.target.value
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();

    var nameValidationMessage = utils.checkName(
      'service account',
      this.state.name
    );

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/serviceaccounts`,
        {
          name: this.state.name,
          description: this.state.description
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        segment.track('Service Account Created');
        this.props.history.push(
          `/${this.props.projectName}/iam/serviceaccounts/`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Service Account was not created.');
          console.log(error);
        }
      });
  };

  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/iam/serviceaccounts/`);
  }

  render() {
    return (
      <Pane width="50%">
        <InnerCard>
          <Pane padding={majorScale(4)} is="form">
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
              Create Service Account
            </Heading>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
              isInvalid={this.state.nameValidationMessage !== null}
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
            />
            <Pane display="flex" flex="row">
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                onClick={this.handleSubmit}
              >
                Submit
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

class RolesRouter extends Component {
  render() {
    const projectName = this.props.projectName;
    const match = this.props.match;
    return (
      <Switch>
        <Route
          path={`${match.path}/create`}
          render={route => (
            <CreateRole
              projectName={projectName}
              history={this.props.history}
            />
          )}
        />
        <Route
          path={`${match.path}/:roleName`}
          render={route => (
            <Role
              projectName={projectName}
              roleName={route.match.params.roleName}
              history={this.props.history}
            />
          )}
        />
        <Route
          exact
          path={match.path}
          render={route => (
            <Roles projectName={projectName} history={this.props.history} />
          )}
        />
      </Switch>
    );
  }
}

class Roles extends Component {
  constructor(props) {
    super(props);
    this.state = {
      roles: []
    };
  }

  componentDidMount() {
    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}/roles`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          roles: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const roles = this.state.roles;
    return (
      <Pane width="70%">
        {roles ? (
          <InnerCard>
            <Pane
              display="flex"
              flexDirection="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Heading paddingLeft={majorScale(2)}>Roles</Heading>
              <Button
                margin={majorScale(2)}
                appearance="primary"
                onClick={() =>
                  this.props.history.push(
                    `/${this.props.projectName}/iam/roles/create`
                  )
                }
              >
                Create Role
              </Button>
            </Pane>
            {this.state.roles && this.state.roles.length > 0 && (
              <Table>
                <Table.Head>
                  <Table.TextHeaderCell>Role</Table.TextHeaderCell>
                </Table.Head>
                <Table.Body>
                  {roles.map(role => (
                    <Table.Row
                      key={role.id}
                      isSelectable
                      onSelect={() =>
                        this.props.history.push(
                          `/${this.props.projectName}/iam/roles/${role.name}`
                        )
                      }
                    >
                      <Table.TextCell>{role.name}</Table.TextCell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
            )}
          </InnerCard>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}

class Role extends Component {
  constructor(props) {
    super(props);
    this.state = {
      role: null,
      name: '',
      nameValidationMessage: null,
      description: '',
      config: '',
      unchanged: true,
      showDeleteDialog: false,
      backendError: null
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/roles/${this.props.roleName}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          role: response.data,
          name: response.data.name,
          description: response.data.description,
          config: response.data.config
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

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

  handleUpdate() {
    var nameValidationMessage = utils.checkName('role', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .put(
        `${config.endpoint}/projects/${this.props.projectName}/roles/${this.state.role.id}`,
        {
          name: this.state.name,
          description: this.state.description,
          config: this.state.config
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Successfully updated role.');
        this.props.history.push(`/${this.props.projectName}/iam/roles`);
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Role was not updated.');
          console.log(error);
        }
      });
  }

  handleDelete() {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/roles/${this.state.role.id}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted role.');
        this.props.history.push(`/${this.props.projectName}/iam/roles`);
      })
      .catch(error => {
        this.setState({
          showDeleteDialog: false
        });
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Role was not deleted.');
          console.log(error);
        }
      });
  }

  render() {
    const role = this.state.role;
    return (
      <Pane width="70%">
        {role ? (
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
                Role / {role.name}
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
              <Heading paddingTop={majorScale(2)} marginBottom={majorScale(2)}>
                Config
              </Heading>
              <Editor
                width="100%"
                height="300px"
                value={this.state.config}
                onChange={value =>
                  this.setState({ config: value, unchanged: false })
                }
              />
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                disabled={this.state.unchanged}
                onClick={() => this.handleUpdate()}
              >
                Update Role
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
                onClick={() => this.setState({ showDeleteDialog: true })}
              >
                Delete Role...
              </Button>
            </Pane>
            <Pane>
              <Dialog
                isShown={this.state.showDeleteDialog}
                title="Delete Role"
                intent="danger"
                onCloseComplete={() =>
                  this.setState({ showDeleteDialog: false })
                }
                onConfirm={() => this.handleDelete()}
                confirmLabel="Delete Role"
              >
                You are about to delete the{' '}
                <strong>{this.props.roleName}</strong> role.
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

class CreateRole extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      nameValidationMessage: null,
      description: '',
      config: '',
      backendError: null
    };
  }

  handleUpdateName = event => {
    this.setState({
      name: event.target.value
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value
    });
  };

  handleSubmit() {
    var nameValidationMessage = utils.checkName('role', this.state.name);
    var configError = utils.checkConfig('role config', this.state.config);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: configError
    });

    if (nameValidationMessage !== null || configError !== null) {
      return;
    }

    if (this.state.config === '') {
      this.setState({
        backendError: 'Please define the role config.'
      });
      return;
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/roles`,
        {
          name: this.state.name,
          description: this.state.description,
          config: this.state.config
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        segment.track('Role Created');
        this.props.history.push(`/${this.props.projectName}/iam/roles/`);
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Role was not created.');
          console.log(error);
        }
      });
  }

  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/iam/roles/`);
  }

  render() {
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
              Create Role
            </Heading>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
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
            />
            <Heading paddingTop={majorScale(2)} marginBottom={majorScale(2)}>
              Config
            </Heading>
            <Editor
              width="100%"
              height="300px"
              value={this.state.config}
              onChange={value => this.setState({ config: value })}
            />
            <Pane display="flex" flex="row">
              <Button
                marginTop={majorScale(2)}
                appearance="primary"
                onClick={() => this.handleSubmit()}
              >
                Submit
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

class Device extends Component {
  constructor(props) {
    super(props);
    this.state = {
      device: null,
      tabs: ['overview', 'ssh', 'settings'],
      tabLabels: ['Overview', 'SSH', 'Settings']
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.deviceName}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          device: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    const deviceName = match.params.deviceName;
    var selectedIndex = 0;
    switch (match.params.appTab) {
      case 'overview':
        selectedIndex = 0;
        break;
      case 'ssh':
        selectedIndex = 1;
        break;
      case 'settings':
        selectedIndex = 2;
        break;
      default:
        this.props.history.push(`/${projectName}/devices/${deviceName}`);
    }
    return (
      <Tablist border="default">
        {this.state.tabs.map((tab, index) => (
          <Tab
            key={tab}
            id={tab}
            onSelect={() =>
              this.props.history.push(
                `/${projectName}/devices/${deviceName}/${tab}`
              )
            }
            isSelected={index === selectedIndex}
          >
            {this.state.tabLabels[index]}
          </Tab>
        ))}
      </Tablist>
    );
  };

  renderInner = match => {
    switch (match.params.appTab) {
      case 'overview':
        return (
          <DeviceOverview
            projectName={match.params.projectName}
            device={this.state.device}
            history={this.props.history}
          />
        );
      case 'ssh':
        return (
          <DeviceSsh
            projectName={match.params.projectName}
            device={this.state.device}
            history={this.props.history}
          />
        );
      case 'settings':
        return (
          <DeviceSettings
            projectName={match.params.projectName}
            device={this.state.device}
            history={this.props.history}
          />
        );
      default:
        return <Pane></Pane>;
    }
  };

  render() {
    const device = this.state.device;
    const heading = 'Device / ' + this.props.deviceName;
    return (
      <React.Fragment>
        <TopHeader user={this.props.user} heading={heading} />
        {device ? (
          <React.Fragment>
            <Pane
              display="flex"
              flexDirection="column"
              alignItems="center"
              background="white"
              width="100%"
              padding={majorScale(1)}
              borderBottom="default"
            >
              {this.renderTablist(this.props.deviceRoute)}
            </Pane>
            {this.renderInner(this.props.deviceRoute)}
          </React.Fragment>
        ) : (
          <CustomSpinner />
        )}
      </React.Fragment>
    );
  }
}

class DeviceOverview extends Component {
  renderDeviceOs(device) {
    var innerText = '-';
    if (
      device.info.hasOwnProperty('osRelease') &&
      device.info.osRelease.hasOwnProperty('prettyName')
    ) {
      innerText = device.info.osRelease.prettyName;
    }
    return <Pane>{innerText}</Pane>;
  }

  render() {
    const device = this.props.device;
    return (
      <React.Fragment>
        <Pane width="70%" display="flex" flexDirection="column">
          {device ? (
            <React.Fragment>
              <InnerCard>
                <Heading padding={majorScale(2)}>Device Info</Heading>
                <Table>
                  <Table.Head>
                    <Table.TextHeaderCell
                      flexBasis={90}
                      flexShrink={0}
                      flexGrow={0}
                    >
                      Status
                    </Table.TextHeaderCell>
                    <Table.TextHeaderCell>IP Address</Table.TextHeaderCell>
                    <Table.TextHeaderCell>OS</Table.TextHeaderCell>
                  </Table.Head>
                  <Table.Body>
                    <Table.Row>
                      <Table.TextCell
                        flexBasis={90}
                        flexShrink={0}
                        flexGrow={0}
                        alignItems="center"
                        paddingRight="0"
                      >
                        {device.status === 'offline' ? (
                          <Badge color="red">offline</Badge>
                        ) : (
                          <Badge color="green">online</Badge>
                        )}
                      </Table.TextCell>
                      <Table.TextCell>
                        {device.info.hasOwnProperty('ipAddress')
                          ? device.info.ipAddress
                          : ''}
                      </Table.TextCell>
                      <Table.TextCell>
                        {this.renderDeviceOs(device)}
                      </Table.TextCell>
                    </Table.Row>
                  </Table.Body>
                </Table>
              </InnerCard>
              <InnerCard>
                <Heading padding={majorScale(2)}>Labels</Heading>
                <DeviceLabels
                  projectName={this.props.projectName}
                  device={device}
                />
              </InnerCard>
              <InnerCard>
                <Heading padding={majorScale(2)}>Services</Heading>
                <DeviceServices
                  projectName={this.props.projectName}
                  device={device}
                  history={this.props.history}
                />
              </InnerCard>
            </React.Fragment>
          ) : (
            <CustomSpinner />
          )}
        </Pane>
      </React.Fragment>
    );
  }
}

class DeviceServices extends Component {
  checkServices(applicationStatusInfo) {
    for (var i = 0; i < applicationStatusInfo.length; i++) {
      if (
        applicationStatusInfo[i].serviceStatuses &&
        applicationStatusInfo[i].serviceStatuses.length > 0
      ) {
        return true;
      }
    }
    return false;
  }

  render() {
    const applicationStatusInfo = this.props.device.applicationStatusInfo;
    return (
      <Pane>
        {this.checkServices(applicationStatusInfo) && (
          <Table>
            <Table.Head>
              <Table.TextHeaderCell>Service</Table.TextHeaderCell>
              <Table.TextHeaderCell>Current Release</Table.TextHeaderCell>
            </Table.Head>
            {applicationStatusInfo.map(applicationInfo => (
              <Pane key={applicationInfo.application.id}>
                <Table.Body>
                  {applicationInfo.serviceStatuses.map(
                    (serviceStatus, index) => (
                      <Table.Row key={index}>
                        <Table.TextCell>
                          <Button
                            color="#425A70"
                            appearance="minimal"
                            onClick={() =>
                              this.props.history.push(
                                `/${this.props.projectName}/applications/${applicationInfo.application.name}`
                              )
                            }
                          >
                            {applicationInfo.application.name} /{' '}
                            {serviceStatus.service}
                          </Button>
                        </Table.TextCell>
                        <Table.TextCell>
                          {serviceStatus.currentReleaseId}
                        </Table.TextCell>
                      </Table.Row>
                    )
                  )}
                </Table.Body>
              </Pane>
            ))}
          </Table>
        )}
      </Pane>
    );
  }
}

class DeviceLabels extends Component {
  constructor(props) {
    super(props);
    this.state = {
      labels: []
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/labels`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          labels: this.initializeLabels(response.data)
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  initializeLabels = keyValues => {
    var labels = [];
    for (var i = 0; i < keyValues.length; i++) {
      labels.push({
        key: keyValues[i].key,
        value: keyValues[i].value,
        mode: 'default',
        keyValidationMessage: null,
        valueValidationMessage: null,
        showRemoveDialog: false
      });
    }
    return labels;
  };

  createNewDeviceLabel = k => {
    var labels = this.state.labels;
    labels.push({
      key: '',
      value: '',
      mode: 'new',
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
    editLabels[i]['mode'] = 'edit';
    this.setState({
      labels: editLabels
    });
  };

  cancelEdit = i => {
    var editLabels = this.state.labels;
    editLabels[i]['mode'] = 'default';
    this.setState({
      labels: editLabels
    });
  };

  setShowRemoveDialog = i => {
    var showRemoveDialogLabels = this.state.labels;
    showRemoveDialogLabels[i]['showRemoveDialog'] = true;
    this.setState({
      labels: showRemoveDialogLabels
    });
  };

  hideShowRemoveDialog = i => {
    var showRemoveDialogLabels = this.state.labels;
    showRemoveDialogLabels[i]['showRemoveDialog'] = false;
    this.setState({
      labels: showRemoveDialogLabels
    });
  };

  setDeviceLabel = (key, value, i) => {
    var updatedLabels = this.state.labels;
    var keyValidationMessage = utils.checkName('key', key);
    var valueValidationMessage = utils.checkName('value', value);

    if (keyValidationMessage === null) {
      for (var j = 0; j < updatedLabels.length; j++) {
        if (i !== j && key === updatedLabels[j]['key']) {
          keyValidationMessage = 'Key already exists.';
          break;
        }
      }
    }

    updatedLabels[i]['keyValidationMessage'] = keyValidationMessage;
    updatedLabels[i]['valueValidationMessage'] = valueValidationMessage;

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
        .post(
          `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/labels`,
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
          updatedLabels[i]['mode'] = 'default';
          this.setState({
            labels: updatedLabels
          });
        })
        .catch(error => {
          console.log(error);
        });
    }
  };

  deleteDeviceLabel = (key, i) => {
    if (key !== '') {
      axios
        .delete(
          `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/labels/${key}`,
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
          hideRemoveDialogLabels[i]['showRemoveDialog'] = false;
          this.setState({
            labels: hideRemoveDialogLabels
          });
          toaster.danger('Device label was not removed.');
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

  renderDeviceLabel(deviceLabel, i) {
    switch (deviceLabel.mode) {
      case 'default':
        return (
          <React.Fragment key={deviceLabel.key}>
            <Table.Row>
              <Table.TextCell>{deviceLabel.key}</Table.TextCell>
              <Table.TextCell>{deviceLabel.value}</Table.TextCell>
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
                isShown={deviceLabel.showRemoveDialog}
                title="Remove Label"
                intent="danger"
                onCloseComplete={() => this.hideShowRemoveDialog(i)}
                onConfirm={() => this.deleteDeviceLabel(deviceLabel.key, i)}
                confirmLabel="Remove Label"
              >
                You are about to remove label <strong>{deviceLabel.key}</strong>
                .
              </Dialog>
            </Pane>
          </React.Fragment>
        );
      case 'edit':
        return (
          <Table.Row key={deviceLabel.key} height="auto">
            <Table.TextCell>{deviceLabel.key}</Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`edit-${deviceLabel.key}`}
                value={deviceLabel.value}
                onChange={event => this.handleUpdate(i, 'value')(event)}
                isInvalid={deviceLabel.valueValidationMessage !== null}
                validationMessage={deviceLabel.valueValidationMessage}
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
                    this.setDeviceLabel(deviceLabel.key, deviceLabel.value, i)
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
      case 'new':
        return (
          <Table.Row key={`new-${i}`} height="auto">
            <Table.TextCell>
              <TextInputField
                label=""
                name={`new-key-${i}`}
                value={deviceLabel.key}
                onChange={event => this.handleUpdate(i, 'key')(event)}
                isInvalid={deviceLabel.keyValidationMessage !== null}
                validationMessage={deviceLabel.keyValidationMessage}
                marginTop={8}
                marginBottom={8}
              />
            </Table.TextCell>
            <Table.TextCell>
              <TextInputField
                label=""
                name={`new-value-${i}`}
                value={deviceLabel.value}
                onChange={event => this.handleUpdate(i, 'value')(event)}
                isInvalid={deviceLabel.valueValidationMessage !== null}
                validationMessage={deviceLabel.valueValidationMessage}
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
                    this.setDeviceLabel(deviceLabel.key, deviceLabel.value, i)
                  }
                />
                <IconButton
                  icon="cross"
                  height={24}
                  appearance="minimal"
                  onClick={() => this.deleteDeviceLabel(deviceLabel.key, i)}
                />
              </Pane>
            </Table.TextCell>
          </Table.Row>
        );
      default:
        return <React.Fragment />;
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
          {this.state.labels.map((deviceLabel, i) =>
            this.renderDeviceLabel(deviceLabel, i)
          )}
          <Table.Row key="add">
            <Table.TextCell>
              <IconButton
                icon="plus"
                height={24}
                appearance="minimal"
                onClick={() => this.createNewDeviceLabel()}
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

class DeviceSettings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: props.device.name,
      nameValidationMessage: null,
      unchanged: true,
      showRemoveDialog: false,
      backendError: null
    };
  }

  handleUpdateName = event => {
    this.setState({
      name: event.target.value,
      unchanged: false
    });
  };

  handleUpdate = () => {
    var nameValidationMessage = utils.checkName('device', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .patch(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}`,
        {
          name: this.state.name
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Device updated successfully.');
        this.props.history.push(
          `/${this.props.projectName}/devices/${this.state.name}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device was not updated.');
          console.log(error);
        }
      });
  };

  handleRemove = () => {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Successfully deleted device.');
        this.props.history.push(`/${this.props.projectName}/devices`);
      })
      .catch(error => {
        this.setState({
          showRemoveDialog: false
        });
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device was not removed.');
          console.log(error);
        }
      });
  };

  render() {
    const device = this.props.device;
    return (
      <Pane width="50%">
        <InnerCard>
          <Pane display="flex" flexDirection="column" padding={majorScale(4)}>
            {this.state.backendError && (
              <Alert
                marginBottom={majorScale(2)}
                paddingTop={majorScale(2)}
                paddingBottom={majorScale(2)}
                intent="warning"
                title={this.state.backendError}
              />
            )}
            <Table paddingBottom={majorScale(2)}>
              <Table.Row>
                <Table.Cell paddingLeft="0">
                  <Heading size={600}>Device Settings</Heading>
                </Table.Cell>
                <Table.Cell flexBasis={90} flexShrink={0} flexGrow={0}>
                  {device.status === 'offline' ? (
                    <Badge color="red">offline</Badge>
                  ) : (
                    <Badge color="green">online</Badge>
                  )}
                </Table.Cell>
              </Table.Row>
            </Table>
            <Text>
              <strong>ID: </strong>
              {device.id}
            </Text>
            <TextInputField
              label="Name"
              onChange={this.handleUpdateName}
              value={this.state.name}
              validationMessage={this.state.nameValidationMessage}
              paddingTop={majorScale(2)}
            />
            <Pane marginTop={majorScale(2)}>
              <Button
                appearance="primary"
                disabled={this.state.unchanged}
                onClick={this.handleUpdate}
              >
                Update Settings
              </Button>
            </Pane>
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
              onClick={() => this.setState({ showRemoveDialog: true })}
            >
              Remove Device...
            </Button>
          </Pane>
          <Pane>
            <Dialog
              isShown={this.state.showRemoveDialog}
              title="Remove Device"
              intent="danger"
              onCloseComplete={() => this.setState({ showRemoveDialog: false })}
              onConfirm={() => this.handleRemove()}
              confirmLabel="Remove Device"
            >
              You are about to remove the <strong>{device.name}</strong> device.
            </Dialog>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}

class AddDevice extends Component {
  constructor(props) {
    super(props);
    this.state = {
      deviceRegistrationToken: null,
      project: null
    };
  }

  componentDidMount() {
    segment.page();

    this.getRegistrationToken();
    axios
      .get(`${config.endpoint}/projects/${this.props.projectName}`, {
        withCredentials: true
      })
      .then(response => {
        this.setState({
          project: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  getRegistrationToken = () => {
    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens`,
        null,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          deviceRegistrationToken: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  };

  handleAddNewDevice = () => {
    this.getRegistrationToken();
    toaster.success('New device token and command generated.');
  };

  render() {
    if (!this.state.deviceRegistrationToken || !this.state.project) {
      return <CustomSpinner />;
    }
    const heading = 'Add Device';
    return (
      <React.Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        <Pane width="70%">
          <InnerCard>
            <Pane
              paddingTop={majorScale(2)}
              paddingLeft={majorScale(4)}
              paddingRight={majorScale(4)}
              paddingBottom={majorScale(2)}
            >
              <Pane
                display="flex"
                flexDirection="row"
                alignItems="center"
                justifyContent="space-between"
                paddingBottom={majorScale(2)}
              >
                <BackButton
                  onClick={() =>
                    this.props.history.push(
                      `/${this.props.projectName}/devices`
                    )
                  }
                >
                  Devices
                </BackButton>
                <Button
                  appearance="primary"
                  onClick={() => this.handleAddNewDevice()}
                >
                  Add Another Device
                </Button>
              </Pane>
              <Pane
                display="flex"
                flexDirection="row"
                alignItems="center"
                paddingBottom={majorScale(2)}
              >
                <Icon icon="info-sign" color="info" marginRight={8} />
                <Text>
                  Device registration ID{' '}
                  <Code fontFamily="mono" background="#234361" color="white">
                    {this.state.deviceRegistrationToken.id}
                  </Code>{' '}
                  created.
                </Text>
              </Pane>
              <Card
                display="flex"
                flexDirection="column"
                padding={majorScale(2)}
                border="muted"
                background="tint2"
              >
                <Text>Run the following command to register your device.</Text>
                <Card
                  marginTop={majorScale(1)}
                  padding={majorScale(1)}
                  background="#234361"
                >
                  <Code fontFamily="mono" color="white">
                    docker run -d --restart=always --privileged --net=host
                    --pid=host -v /etc/deviceplane:/etc/deviceplane -v
                    /var/lib/deviceplane:/var/lib/deviceplane -v
                    /var/run/docker.sock:/var/run/docker.sock -v
                    /etc/os-release:/etc/os-release deviceplane/agent:
                    {config.agentVersion} --project={this.state.project.id}{' '}
                    --registration-token={this.state.deviceRegistrationToken.id}
                  </Code>
                </Card>
              </Card>
            </Pane>
          </InnerCard>
        </Pane>
      </React.Fragment>
    );
  }
}

class Application extends Component {
  constructor(props) {
    super(props);
    this.state = {
      application: null,
      tabs: ['overview', 'scheduling', 'settings'],
      tabLabels: ['Overview', 'Scheduling', 'Settings']
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.applicationName}?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          application: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    const applicationName = match.params.applicationName;
    var selectedIndex = 0;
    switch (match.params.appTab) {
      case 'overview':
        selectedIndex = 0;
        break;
      case 'scheduling':
        selectedIndex = 1;
        break;
      case 'settings':
        selectedIndex = 2;
        break;
      default:
        this.props.history.push(
          `/${projectName}/applicationName/${applicationName}`
        );
    }
    return (
      <Tablist border="default">
        {this.state.tabs.map((tab, index) => (
          <Tab
            key={tab}
            id={tab}
            onSelect={() =>
              this.props.history.push(
                `/${projectName}/applications/${applicationName}/${tab}`
              )
            }
            isSelected={index === selectedIndex}
          >
            {this.state.tabLabels[index]}
          </Tab>
        ))}
      </Tablist>
    );
  };

  renderInner = match => {
    switch (match.params.appTab) {
      case 'overview':
        return (
          <ApplicationOverview
            projectName={match.params.projectName}
            application={this.state.application}
            history={this.props.history}
          />
        );
      case 'scheduling':
        return (
          <ApplicationScheduling
            projectName={match.params.projectName}
            application={this.state.application}
            history={this.props.history}
          />
        );
      case 'settings':
        return (
          <ApplicationSettings
            projectName={match.params.projectName}
            application={this.state.application}
            history={this.props.history}
          />
        );
      default:
        return <Pane></Pane>;
    }
  };

  render() {
    const application = this.state.application;
    const heading = 'Application / ' + this.props.applicationName;
    return (
      <React.Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        {application ? (
          <React.Fragment>
            <Pane
              display="flex"
              flexDirection="column"
              alignItems="center"
              background="white"
              width="100%"
              padding={majorScale(1)}
              borderBottom="default"
            >
              {this.renderTablist(this.props.applicationRoute)}
            </Pane>
            {this.renderInner(this.props.applicationRoute)}
          </React.Fragment>
        ) : (
          <CustomSpinner />
        )}
      </React.Fragment>
    );
  }
}

class ApplicationOverview extends Component {
  render() {
    const currentConfig = this.props.application.latestRelease
      ? this.props.application.latestRelease.config
      : '';
    return (
      <Pane width="70%" paddingBottom={majorScale(4)}>
        <Heading
          paddingTop={majorScale(4)}
          paddingBottom={majorScale(1)}
          size={600}
        >
          {this.props.application.name}
        </Heading>
        <InnerCard>
          <Heading paddingTop={majorScale(2)} paddingLeft={majorScale(2)}>
            Scheduling Rule
          </Heading>
          <Card
            display="flex"
            flexDirection="column"
            alignItems="left"
            width="80%"
            padding={majorScale(2)}
          >
            <Label>{this.props.application.settings.schedulingRule}</Label>
          </Card>
        </InnerCard>
        <InnerCard>
          <Heading paddingTop={majorScale(2)} paddingLeft={majorScale(2)}>
            Current Config
          </Heading>
          <Card
            display="flex"
            flexDirection="column"
            alignItems="center"
            width="80%"
            padding={majorScale(2)}
          >
            <Editor
              width="100%"
              height="300px"
              value={currentConfig}
              readOnly
            />
          </Card>
        </InnerCard>
        <InnerCard>
          <Pane
            display="flex"
            flexDirection="row"
            alignItems="center"
            justifyContent="space-between"
          >
            <Heading padding={majorScale(2)}>Releases</Heading>
            <Button
              marginRight={majorScale(2)}
              appearance="primary"
              onClick={() =>
                this.props.history.push(
                  `/${this.props.projectName}/applications/${this.props.application.name}/deploy`
                )
              }
            >
              Create New Release
            </Button>
          </Pane>
          <Releases
            projectName={this.props.projectName}
            applicationName={this.props.application.name}
            history={this.props.history}
          />
        </InnerCard>
      </Pane>
    );
  }
}

class ApplicationScheduling extends Component {
  constructor(props) {
    super(props);
    this.state = {
      schedulingRule: props.application.settings.schedulingRule,
      backendError: null
    };
  }

  handleUpdateSchedulingRule = event => {
    this.setState({
      schedulingRule: event.target.value
    });
  };

  handleSubmit = () => {
    this.setState({
      backendError: null
    });

    axios
      .put(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.application.name}`,
        {
          name: this.props.application.name,
          description: this.props.application.description,
          settings: {
            schedulingRule: this.state.schedulingRule
          }
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Scheduling rule updated successfully.');
        this.props.history.push(
          `/${this.props.projectName}/applications/${this.props.application.name}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          console.log(error);
        }
      });
  };

  render() {
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
              Scheduling
            </Heading>
            <Label
              htmlFor="description-textarea"
              marginBottom="4"
              display="block"
            >
              Scheduling Rule
            </Label>
            <Textarea
              id="description-textarea"
              height="100px"
              onChange={this.handleUpdateSchedulingRule}
              value={this.state.schedulingRule}
            />
            <Button
              marginTop={majorScale(2)}
              appearance="primary"
              onClick={() => this.handleSubmit()}
            >
              Submit
            </Button>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}

class ApplicationSettings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: props.application.name,
      nameValidationMessage: null,
      description: props.application.description,
      unchanged: true,
      showDeleteDialog: false,
      backendError: null
    };
  }

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

  handleUpdate = () => {
    var nameValidationMessage = utils.checkName('application', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .put(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.application.name}`,
        {
          name: this.state.name,
          description: this.state.description,
          settings: this.props.application.settings
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Application updated successfully.');
        this.props.history.push(
          `/${this.props.projectName}/applications/${this.state.name}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Application was not updated.');
          console.log(error);
        }
      });
  };

  handleDelete() {
    axios
      .delete(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.application.name}`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showDeleteDialog: false
        });
        toaster.success('Successfully deleted application.');
        this.props.history.push(`/${this.props.projectName}/applications`);
      })
      .catch(error => {
        this.setState({
          showDeleteDialog: false
        });
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Application was not deleted.');
          console.log(error);
        }
      });
  }

  render() {
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
              Application Settings
            </Heading>
            <TextInputField
              label="Application Name"
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
            <Button
              marginTop={majorScale(2)}
              appearance="primary"
              disabled={this.state.unchanged}
              onClick={this.handleUpdate}
            >
              Update Settings
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
              onClick={() => this.setState({ showDeleteDialog: true })}
            >
              Delete Application...
            </Button>
          </Pane>
          <Pane>
            <Dialog
              isShown={this.state.showDeleteDialog}
              title="Delete Application"
              intent="danger"
              onCloseComplete={() => this.setState({ showDeleteDialog: false })}
              onConfirm={() => this.handleDelete()}
              confirmLabel="Delete Application"
            >
              You are about to delete the{' '}
              <strong>{this.props.application.name}</strong> application.
            </Dialog>
          </Pane>
        </InnerCard>
      </Pane>
    );
  }
}

class CreateApplication extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      nameValidationMessage: null,
      description: '',
      backendError: null
    };
  }

  handleUpdateName = event => {
    this.setState({
      name: event.target.value
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();

    var nameValidationMessage = utils.checkName('application', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/applications`,
        {
          name: this.state.name,
          description: this.state.description
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        segment.track('Application Created');
        this.props.history.push(
          `/${this.props.projectName}/applications/${this.state.name}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Application was not created.');
          console.log(error);
        }
      });
  };

  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/applications`);
  }

  render() {
    const heading = 'Create Application';
    return (
      <React.Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        <Pane width="50%">
          <InnerCard>
            <Pane padding={majorScale(4)} is="form">
              {this.state.backendError && (
                <Alert
                  marginBottom={majorScale(2)}
                  paddingTop={majorScale(2)}
                  paddingBottom={majorScale(2)}
                  intent="warning"
                  title={this.state.backendError}
                />
              )}
              <TextInputField
                label="Name"
                onChange={this.handleUpdateName}
                value={this.state.name}
                isInvalid={this.state.nameValidationMessage !== null}
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
              />
              <Pane display="flex" flex="row">
                <Button
                  marginTop={majorScale(2)}
                  appearance="primary"
                  onClick={this.handleSubmit}
                >
                  Submit
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
      </React.Fragment>
    );
  }
}

class CreateRelease extends Component {
  constructor(props) {
    super(props);
    this.state = {
      externalData: [],
      config: '',
      backendError: null
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.applicationName}/releases/latest`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          externalData: response.data,
          config: response.data.config
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  handleSubmit() {
    var configError = utils.checkConfig('release', this.state.config);

    this.setState({
      backendError: configError
    });

    if (configError !== null) {
      return;
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.applicationName}/releases`,
        {
          config: this.state.config
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        segment.track('Release Created');
        this.props.history.push(
          `/${this.props.projectName}/applications/${this.props.applicationName}`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          console.log(error);
        }
      });
  }

  handleCancel() {
    this.props.history.push(
      `/${this.props.projectName}/applications/${this.props.applicationName}`
    );
  }

  render() {
    const heading = 'Create Release';
    return (
      <React.Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
        <Pane width="70%">
          <InnerCard>
            <Pane padding={majorScale(2)}>
              {this.state.backendError && (
                <Alert
                  marginBottom={majorScale(2)}
                  paddingTop={majorScale(2)}
                  paddingBottom={majorScale(2)}
                  intent="warning"
                  title={this.state.backendError}
                />
              )}
              <Heading htmlFor="textarea-2" marginBottom={majorScale(2)}>
                Config
              </Heading>
              <Editor
                width="100%"
                height="300px"
                value={this.state.config}
                onChange={value => this.setState({ config: value })}
              />
              <Pane display="flex" flex="row">
                <Button
                  marginTop={majorScale(2)}
                  appearance="primary"
                  onClick={() => this.handleSubmit()}
                >
                  Submit
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
      </React.Fragment>
    );
  }
}

class Releases extends Component {
  constructor(props) {
    super(props);
    this.state = {
      releases: [],
      showRelease: false,
      selectedRelease: null
    };
  }

  componentDidMount() {
    axios
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.applicationName}/releases?full`,
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          releases: response.data
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  getReleasedBy = release => {
    if (release) {
      if (release.createdByUser) {
        const memberUrl = '../../iam/members/' + release.createdByUser.id;
        return (
          <Link color="neutral" href={memberUrl}>
            {release.createdByUser.firstName} {release.createdByUser.lastName}
          </Link>
        );
      } else if (release.createdByServiceAccount) {
        const serviceAccountUrl =
          '../../iam/serviceaccounts/' + release.createdByServiceAccount.name;
        return (
          <Link color="neutral" href={serviceAccountUrl}>
            {release.createdByServiceAccount.name}
          </Link>
        );
      }
    }
    return '-';
  };

  showSelectedRelease = release => {
    this.setState({
      showRelease: true,
      selectedRelease: release
    });
  };

  render() {
    return (
      <React.Fragment>
        <Pane>
          {this.state.releases && this.state.releases.length > 0 && (
            <Table>
              <Table.Head>
                <Table.TextHeaderCell flexGrow={3} flexShrink={3}>
                  Release
                </Table.TextHeaderCell>
                <Table.TextHeaderCell flexGrow={2} flexShrink={2}>
                  Released By
                </Table.TextHeaderCell>
                <Table.TextHeaderCell>Started</Table.TextHeaderCell>
                <Table.TextHeaderCell>Device Count</Table.TextHeaderCell>
              </Table.Head>
              <Table.Body>
                {this.state.releases.map(release => (
                  <Table.Row
                    key={release.id}
                    isSelectable
                    onSelect={() => this.showSelectedRelease(release)}
                  >
                    <Table.TextCell flexGrow={3} flexShrink={3}>
                      {release.id}
                    </Table.TextCell>
                    <Table.TextCell flexGrow={2} flexShrink={2}>
                      {this.getReleasedBy(release)}
                    </Table.TextCell>
                    <Table.TextCell>
                      {moment(release.createdAt).fromNow()}
                    </Table.TextCell>
                    <Table.TextCell>
                      {release.deviceCounts.allCount}
                    </Table.TextCell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          )}
        </Pane>
        <SideSheet
          isShown={this.state.showRelease}
          onCloseComplete={() => this.setState({ showRelease: false })}
        >
          <Release
            release={this.state.selectedRelease}
            projectName={this.props.projectName}
            applicationName={this.props.applicationName}
            history={this.props.history}
          ></Release>
        </SideSheet>
      </React.Fragment>
    );
  }
}

class Release extends Component {
  constructor(props) {
    super(props);
    this.state = {
      backendError: null,
      showConfirmDialog: false
    };
  }

  getReleasedBy = release => {
    if (release) {
      if (release.createdByUser) {
        const memberUrl = '../../iam/members/' + release.createdByUser.id;
        return (
          <Link color="neutral" href={memberUrl}>
            {release.createdByUser.firstName} {release.createdByUser.lastName}
          </Link>
        );
      } else if (release.createdByServiceAccount) {
        const serviceAccountUrl =
          '../../iam/serviceaccounts/' + release.createdByServiceAccount.name;
        return (
          <Link color="neutral" href={serviceAccountUrl}>
            {release.createdByServiceAccount.name}
          </Link>
        );
      }
    }
    return '-';
  };

  revertRelease = () => {
    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/applications/${this.props.applicationName}/releases`,
        {
          config: this.props.release.config
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        this.setState({
          showConfirmDialog: false
        });
        // segment.track('Release Created');
        this.props.history.push(
          `/${this.props.projectName}/applications/${this.props.applicationName}`
        );
      })
      .catch(error => {
        this.setState({
          showConfirmDialog: false
        });
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          console.log(error);
        }
      });
  };

  render() {
    const release = this.props.release;
    return (
      <React.Fragment>
        <Pane zIndex={1} flexShrink={0} elevation={0} backgroundColor="white">
          <Pane padding={majorScale(2)}>
            <Heading size={600}>Release / {release.id}</Heading>
          </Pane>
        </Pane>
        <Pane display="flex" flexDirection="column" margin={majorScale(2)}>
          {this.state.backendError && (
            <Alert
              marginBottom={majorScale(2)}
              paddingTop={majorScale(2)}
              paddingBottom={majorScale(2)}
              intent="warning"
              title={this.state.backendError}
            />
          )}
          <Pane
            display="flex"
            flexDirection="column"
            marginBottom={majorScale(2)}
          >
            <Label>
              <strong>Released By:</strong> {this.getReleasedBy(release)}
            </Label>
            <Label>
              <strong>Started:</strong> {moment(release.createdAt).fromNow()}
            </Label>
          </Pane>
          <Pane
            display="flex"
            justifyContent="space-between"
            alignItems="center"
            marginBottom={majorScale(2)}
          >
            <Heading size={600}>Config</Heading>
            <Button
              appearance="primary"
              justifyContent="center"
              onClick={() => this.setState({ showConfirmDialog: true })}
            >
              Revert to this Release
            </Button>
          </Pane>
          <Editor width="100%" height="300px" value={release.config} readOnly />
        </Pane>
        <Dialog
          isShown={this.state.showConfirmDialog}
          title="Revert Release"
          onCloseComplete={() => this.setState({ showConfirmDialog: false })}
          onConfirm={() => this.revertRelease()}
          confirmLabel="Revert Release"
        >
          This will create a new release to application{' '}
          <strong>{this.props.applicationName}</strong> using the config from
          release <strong>{release.id}</strong>.
        </Dialog>
      </React.Fragment>
    );
  }
}

class InnerOogie extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tabs: ['devices', 'applications', 'iam'],
      tabLabels: ['Devices', 'Applications', 'IAM'],
      icons: ['desktop', 'application', 'user'],
      footerTabs: ['settings'],
      footerTabLabels: ['Settings'],
      footerIcons: ['settings']
    };
  }

  renderTablist = match => {
    const projectName = match.params.projectName;
    var selectedIndex = 100;
    var footerSelectedIndex = 100;
    switch (match.params.tab) {
      case 'devices':
        selectedIndex = 0;
        break;
      case 'applications':
        selectedIndex = 1;
        break;
      case 'iam':
        selectedIndex = 2;
        break;
      case 'settings':
        footerSelectedIndex = 0;
        break;
      default:
        this.props.history.push(`/${projectName}`);
    }
    return (
      <Pane
        display="flex"
        flexDirection="column"
        alignItems="center"
        position="sticky"
        top="0"
        left="0"
        overflow="auto"
        borderRight="default"
        height="100vh"
      >
        <Pane
          display="flex"
          flexDirection="column"
          alignItems="center"
          padding={majorScale(2)}
          width="100%"
        >
          <a href="/">
            <Logo />
          </a>
        </Pane>
        <TabNavigation>
          {this.state.tabs.map((tab, index) => (
            <SidebarTab
              key={tab}
              id={tab}
              is="a"
              onSelect={() => this.props.history.push(`/${projectName}/${tab}`)}
              isSelected={index === selectedIndex}
              aria-controls={`panel-${tab}`}
              padding="0"
              margin="0"
              height="auto"
            >
              <Pane
                display="flex"
                flexDirection="column"
                alignItems="center"
                padding={majorScale(2)}
              >
                <Icon icon={this.state.icons[index]} />
                <Text paddingTop={majorScale(1)}>
                  {this.state.tabLabels[index]}
                </Text>
              </Pane>
            </SidebarTab>
          ))}
        </TabNavigation>
        <Pane
          display="flex"
          flexDirection="column"
          alignItems="center"
          marginTop="auto"
          width="100%"
        >
          {this.state.footerTabs.map((tab, index) => (
            <SidebarTab
              key={tab}
              id={tab}
              onSelect={() => this.props.history.push(`/${projectName}/${tab}`)}
              isSelected={index === footerSelectedIndex}
              aria-controls={`panel-${tab}`}
              padding="0"
              margin="0"
              height="auto"
              width="100%"
            >
              <Pane
                display="flex"
                flexDirection="column"
                alignItems="center"
                padding={majorScale(2)}
              >
                <Icon icon={this.state.footerIcons[index]} />
                <Text paddingTop={majorScale(1)}>
                  {this.state.footerTabLabels[index]}
                </Text>
              </Pane>
            </SidebarTab>
          ))}
        </Pane>
      </Pane>
    );
  };

  renderInner = match => {
    const projectName = match.params.projectName;
    const user = this.props.user;
    switch (match.params.tab) {
      case 'devices':
        return (
          <Switch>
            <Route
              path={`${match.path}/add`}
              render={route => (
                <AddDevice
                  user={user}
                  projectName={projectName}
                  history={this.props.history}
                />
              )}
            />
            <Route
              exact
              path={`${match.path}/:deviceName`}
              render={route => <Redirect to={`${route.match.url}/overview`} />}
            />
            <Route
              path={`${match.path}/:deviceName/:appTab`}
              render={route => (
                <Device
                  user={user}
                  projectName={projectName}
                  deviceName={route.match.params.deviceName}
                  deviceRoute={route.match}
                  history={this.props.history}
                />
              )}
            />
            <Route
              exact
              path={match.path}
              render={route => (
                <Devices
                  user={user}
                  projectName={projectName}
                  history={this.props.history}
                />
              )}
            />
          </Switch>
        );
      case 'applications':
        return (
          <Switch>
            <Route
              path={`${match.path}/create`}
              render={route => (
                <CreateApplication
                  user={user}
                  projectName={projectName}
                  history={this.props.history}
                />
              )}
            />
            <Route
              path={`${match.path}/:applicationName/deploy`}
              render={route => (
                <CreateRelease
                  user={user}
                  projectName={projectName}
                  applicationName={route.match.params.applicationName}
                  history={this.props.history}
                />
              )}
            />
            <Route
              exact
              path={`${match.path}/:applicationName`}
              render={route => <Redirect to={`${route.match.url}/overview`} />}
            />
            <Route
              path={`${match.path}/:applicationName/:appTab`}
              render={route => (
                <Application
                  user={user}
                  projectName={projectName}
                  applicationName={route.match.params.applicationName}
                  applicationRoute={route.match}
                  history={this.props.history}
                />
              )}
            />
            <Route
              exact
              path={match.path}
              render={route => (
                <Applications
                  user={user}
                  projectName={projectName}
                  history={this.props.history}
                />
              )}
            />
          </Switch>
        );
      case 'iam':
        return (
          <Switch>
            <Route
              path={`${match.path}/:iamTab`}
              render={route => (
                <Iam
                  user={user}
                  projectName={projectName}
                  match={route.match}
                  history={this.props.history}
                />
              )}
            />
            <Route
              render={route => <Redirect to={`${route.match.url}/members`} />}
            />
          </Switch>
        );
      case 'settings':
        return (
          <Switch>
            <Route
              exact
              path={match.path}
              render={route => (
                <Settings
                  user={user}
                  projectName={projectName}
                  history={this.props.history}
                />
              )}
            />
          </Switch>
        );
      default:
        return <Pane></Pane>;
    }
  };

  render() {
    return (
      <React.Fragment>
        <TabNavigation width="120px">
          {this.renderTablist(this.props.match)}
        </TabNavigation>
        <Pane
          display="flex"
          flexDirection="column"
          alignItems="center"
          background="tint1"
          width="100%"
        >
          {this.renderInner(this.props.match)}
        </Pane>
      </React.Fragment>
    );
  }
}

class Confirm extends Component {
  componentDidMount() {
    axios
      .post(
        `${config.endpoint}/completeregistration`,
        {
          registrationTokenValue: this.props.token
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        window.location.reload();
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    return <Heading>Confirming registration...</Heading>;
  }
}

class PasswordRecovery extends Component {
  constructor(props) {
    super(props);
    this.state = {
      invalidToken: false,
      password: '',
      passwordConfirmation: '',
      passwordValidationMessage: null,
      passwordConfirmationValidationMessage: null,
      passwordRecoveryError: null
    };
  }

  componentDidMount() {
    axios
      .get(`${config.endpoint}/passwordrecoverytokens/${this.props.token}`)
      .then(response => {
        this.setState({
          invalidToken: this.isTokenExpired(response.data.expiresAt)
        });
      })
      .catch(error => {
        this.setState({
          invalidToken: true
        });
        console.log(error);
      });
  }

  isTokenExpired(expiration) {
    const expiratonTime = new Date(expiration).getTime();
    const currentTime = new Date().getTime();
    return currentTime > expiratonTime;
  }

  handleUpdatePassword = event => {
    this.setState({
      password: event.target.value
    });
  };

  handleUpdatePasswordConfirmation = event => {
    this.setState({
      passwordConfirmation: event.target.value
    });
  };

  handleSubmit = () => {
    var passwordValidationMessage = null;
    var passwordConfirmationValidationMessage = null;

    if (!utils.passwordRegex.test(this.state.password)) {
      passwordValidationMessage = 'Please enter a valid password.';
    }

    if (
      this.state.passwordConfirmation === '' ||
      this.state.password !== this.state.passwordConfirmation
    ) {
      passwordConfirmationValidationMessage = 'Passwords must match.';
    }

    this.setState({
      passwordValidationMessage: passwordValidationMessage,
      passwordConfirmationValidationMessage: passwordConfirmationValidationMessage,
      passwordRecoveryError: null
    });

    if (
      passwordValidationMessage === null &&
      passwordConfirmationValidationMessage === null
    ) {
      axios
        .post(`${config.endpoint}/changepassword`, {
          passwordRecoveryTokenValue: this.props.token,
          password: this.state.password
        })
        .then(response => {
          toaster.success('Password successfully changed.');
          this.props.history.push(`/login`);
        })
        .catch(error => {
          if (utils.is4xx(error.response.status)) {
            this.setState({
              passwordRecoveryError: utils.convertErrorMessage(
                error.response.data
              )
            });
          } else {
            toaster.danger(
              'Something went wrong with changing your password. Please contact us at support@deviceplane.com.'
            );
            console.log(error);
          }
        });
    }
  };

  render() {
    return (
      <Pane
        display="flex"
        flexDirection="column"
        alignItems="center"
        padding={majorScale(6)}
        width={majorScale(50)}
        marginX="auto"
        marginY={200}
        elevation={2}
      >
        <Pane
          display="flex"
          alignItems="center"
          paddingBottom={majorScale(2)}
          marginX="auto"
        >
          <Pane paddingRight={minorScale(3)}>
            <Logo />
          </Pane>
          <Heading size={600}>Device Plane</Heading>
        </Pane>
        {this.state.passwordRecoveryError && (
          <Alert
            marginBottom={majorScale(2)}
            paddingTop={majorScale(2)}
            paddingBottom={majorScale(2)}
            intent="warning"
            title={this.state.passwordRecoveryError}
          />
        )}
        {this.state.invalidToken ? (
          <React.Fragment>
            <Alert
              marginBottom={majorScale(2)}
              paddingTop={majorScale(2)}
              paddingBottom={majorScale(2)}
              intent="warning"
              title="Your recovery token has expired. Please reset your password again."
            />
            <Link href="/forgot">Reset your password</Link>
          </React.Fragment>
        ) : (
          <React.Fragment>
            <Pane>
              <TextInputField
                type="password"
                label="New Password"
                hint="Password must be at least 8 characters, contain at least one lower case letter, one upper case letter, and no spaces."
                required
                onChange={this.handleUpdatePassword}
                value={this.state.password}
                isInvalid={this.state.passwordValidationMessage !== null}
                validationMessage={this.state.passwordValidationMessage}
              />
              <TextInputField
                type="password"
                label="Re-enter Password"
                required
                onChange={this.handleUpdatePasswordConfirmation}
                value={this.state.passwordConfirmation}
                isInvalid={
                  this.state.passwordConfirmationValidationMessage !== null
                }
                validationMessage={
                  this.state.passwordConfirmationValidationMessage
                }
              />
            </Pane>
            <Button
              appearance="primary"
              justifyContent="center"
              onClick={this.handleSubmit}
            >
              Submit
            </Button>
          </React.Fragment>
        )}
      </Pane>
    );
  }
}

class OuterOogie extends Component {
  componentDidMount() {
    const user = this.props.user;
    window.Intercom('boot', {
      app_id: 'vm7fcuub',
      name: `${user.firstName} ${user.lastName}`,
      email: user.email
    });
  }

  render() {
    const user = this.props.user;
    return (
      <Pane display="flex" flexGrow={1} minHeight="100vh">
        <Switch>
          <Route exact path="/" render={route => <Redirect to="/projects" />} />
          <Route
            exact
            path="/projects"
            render={route => (
              <Projects
                user={user}
                match={route.match}
                history={route.history}
              />
            )}
          />
          <Route
            exact
            path="/projects/create"
            render={route => (
              <CreateProject
                user={user}
                match={route.match}
                history={route.history}
              />
            )}
          />
          <Route
            exact
            path="/:projectName"
            render={route => <Redirect to={`${route.match.url}/devices`} />}
          />
          <Route
            path="/:projectName/:tab"
            render={route => (
              <InnerOogie
                user={user}
                match={route.match}
                history={route.history}
              />
            )}
          />
        </Switch>
      </Pane>
    );
  }
}

class Authenticated extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: null
    };
  }

  componentDidMount() {
    axios
      .get(`${config.endpoint}/me`, {
        withCredentials: true
      })
      .then(response => {
        const user = response.data;
        segment.identify(user.id, {
          firstName: user.firstName,
          lastName: user.lastName,
          email: user.email
        });
        this.setState({
          user: user
        });
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.props.history.push('/login');
        } else {
          console.log(error);
        }
      });
  }

  render() {
    const user = this.state.user;
    return (
      <React.Fragment>
        {user ? (
          <OuterOogie user={user} history={this.props.history} />
        ) : (
          <CustomSpinner />
        )}
      </React.Fragment>
    );
  }
}

class Unauthenticated extends Component {
  constructor(props) {
    super(props);
    this.state = {
      authenticationCheckCompleted: false
    };
  }

  componentDidMount() {
    axios
      .get(`${config.endpoint}/me`, {
        withCredentials: true
      })
      .then(response => {
        this.props.history.push('/');
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            authenticationCheckCompleted: true
          });
        } else {
          console.log(error);
        }
      });
  }

  render() {
    return (
      <Pane>
        {this.state.authenticationCheckCompleted ? (
          <Switch>
            <Route exact path="/forgot" component={ResetPassword} />
            <Route exact path="/login" component={Login} />
            <Route exact path="/register" component={Register} />
            <Route
              path="/confirm/:token"
              render={route => <Confirm token={route.match.params.token} />}
            />
            <Route
              path="/recover/:token"
              render={route => (
                <PasswordRecovery
                  token={route.match.params.token}
                  history={route.history}
                />
              )}
            />
          </Switch>
        ) : (
          <CustomSpinner />
        )}
      </Pane>
    );
  }
}

class App extends Component {
  render() {
    return (
      <Router>
        <Switch>
          <Redirect from="/index.html" to="/" />
          <Route exact path="/forgot" component={Unauthenticated} />
          <Route exact path="/login" component={Unauthenticated} />
          <Route exact path="/register" component={Unauthenticated} />
          <Route path="/confirm/:token" component={Unauthenticated} />
          <Route path="/recover/:token" component={Unauthenticated} />
          <Route component={Authenticated} />
        </Switch>
      </Router>
    );
  }
}

export default App;
