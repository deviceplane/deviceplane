import React, { Component, Fragment } from 'react';
import {
  Textarea,
  Button,
  TextInputField,
  Pane,
  SidebarTab,
  Heading,
  Icon,
  Text,
  Code,
  Card,
  Label,
  BackButton,
  majorScale,
  toaster,
  Alert,
  TabNavigation
} from 'evergreen-ui';
import {
  BrowserRouter as Router,
  Route,
  Redirect,
  Switch
} from 'react-router-dom';
import axios from 'axios';

import './App.css';

import config from './config';
import segment from './segment';
import utils from './utils';

import Devices from './pages/devices';
import Device from './pages/device';
import Applications from './pages/applications';
import Application from './pages/application';
import ProjectSettings from './pages/project-settings';
import Projects from './pages/projects';
import Register from './pages/register';
import Login from './pages/login';
import Metrics from './pages/metrics';
import Iam from './pages/iam';
import PasswordRecovery from './pages/password-recovery';

import CreateProject from './components/CreateProject';
import CustomSpinner from './components/CustomSpinner';
import Editor from './components/Editor';
import InnerCard from './components/InnerCard';
import ResetPassword from './components/ResetPassword';
import TopHeader from './components/TopHeader';
import Logo from './components/logo';
import DeviceRegistrationToken from './pages/deviceregistrationtoken';
import Provisioning from './pages/provisioning';

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
      .get(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens/default`,
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

    var dockerCommands = [];
    dockerCommands.push([
      "curl https://install.deviceplane.com",
      "|",
      `VERSION=${config.agentVersion}`,
      `PROJECT=${this.state.project.id}`,
      `REGISTRATION_TOKEN=${this.state.deviceRegistrationToken.id}`,
      "bash",
    ].join(" "));

    if (window.location.hostname === "localhost") {
      dockerCommands.push([
        "go run cmd/agent/main.go",
        "--controller=http://localhost:8080/api",
        "--conf-dir=./cmd/agent/conf",
        "--state-dir=./cmd/agent/state",
        "--log-level=debug",
        `--project=${this.state.project.id}`,
        `--registration-token=${this.state.deviceRegistrationToken.id}`,
        "# note, this is the local version"
      ].join(" "));
    }

    return (
      <Fragment>
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
              </Pane>
              <Pane
                display="flex"
                flexDirection="row"
                alignItems="center"
                paddingBottom={majorScale(2)}
              >
                <Icon icon="info-sign" color="info" marginRight={8} />
                <Text>
                  Default device registration token with ID{' '}
                  <Code fontFamily="mono" background="#234361" color="white">
                    {this.state.deviceRegistrationToken.id}
                  </Code>{' '}
                  is being used.
                </Text>
              </Pane>
              <Card
                display="flex"
                flexDirection="column"
                padding={majorScale(2)}
                border="muted"
                background="tint2"
              >
                <Text>
                  Run the following command to register your device.
                </Text>
                {dockerCommands.map(cmd =>
                  <Card
                    marginTop={majorScale(1)}
                    padding={majorScale(1)}
                    background="#234361"
                  >
                    <Code fontFamily="mono" color="white">
                      {cmd}
                    </Code>
                  </Card>
                )}
              </Card>
            </Pane>
          </InnerCard>
        </Pane>
      </Fragment>
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
      <Fragment>
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
      </Fragment>
    );
  }
}

class CreateDeviceRegistrationToken extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      description: '',
      maxRegistrations: '',
      nameValidationMessage: null,
      maxRegistrationsValidationMessage: null,
      showDeleteDialog: false,
      backendError: null
    };
  }

  handleUpdateName = event => {
    this.setState({
      name: event.target.value,
    });
  };

  handleUpdateDescription = event => {
    this.setState({
      description: event.target.value,
    });
  };

  handleUpdateMaxRegistrations = event => {
    this.setState({
      maxRegistrations: event.target.value,
    });
  };

  handleSubmit = () => {
    var nameValidationMessage = utils.checkName('Device Registration Token', this.state.name);

    //always set validation message for name
    this.setState({
      nameValidationMessage: nameValidationMessage,
      backendError: null
    });

    if (nameValidationMessage !== null) {
      return;
    }

    // Convert max registrations to int or undefined
    var maxRegistrationsCleaned;
    if (this.state.maxRegistrations === '') {
      maxRegistrationsCleaned = null;
    } else {
      maxRegistrationsCleaned = Number(this.state.maxRegistrations);

      if (isNaN(maxRegistrationsCleaned)) {
        this.setState({
          maxRegistrationsValidationMessage: 'Max Registrations should either be a number or be left empty.'
        });
        return;
      }
    }

    axios
      .post(
        `${config.endpoint}/projects/${this.props.projectName}/deviceregistrationtokens`,
        {
          name: this.state.name,
          description: this.state.description,
          maxRegistrations: maxRegistrationsCleaned,
        },
        {
          withCredentials: true
        }
      )
      .then(response => {
        toaster.success('Device Registration Token created successfully.');
        this.props.history.push(
          // TODO: change this to an overview page once we get one
          `/${this.props.projectName}/provisioning`
        );
      })
      .catch(error => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Device Registration Token was not updated.');
          console.log(error);
        }
      });
  };


  handleCancel() {
    this.props.history.push(`/${this.props.projectName}/provisioning`);
  }

  render() {
    const heading = 'Create Device Registration Token';
    return (
      <Fragment>
        <TopHeader
          user={this.props.user}
          heading={heading}
          history={this.props.history}
        />
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
                Device Registration Token Settings
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
              <TextInputField
                paddingTop={majorScale(4)}
                label="Maximum Device Registrations"
                description="Limit the number of devices that can be registered using this token"
                hint="Leave empty to allow unlimited registrations"
                onChange={this.handleUpdateMaxRegistrations}
                value={this.state.maxRegistrations}
                validationMessage={this.state.maxRegistrationsValidationMessage}
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
      </Fragment>
    );
  }
}

class CreateRelease extends Component {
  constructor(props) {
    super(props);
    this.state = {
      externalData: [],
      rawConfig: '',
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
          rawConfig: response.data.rawConfig
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  handleSubmit() {
    var configError = utils.checkConfig('release', this.state.rawConfig);

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
          rawConfig: this.state.rawConfig
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
      <Fragment>
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
                value={this.state.rawConfig}
                onChange={value => this.setState({ rawConfig: value })}
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
      </Fragment>
    );
  }
}

class InnerOogie extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tabs: ['devices', 'provisioning', 'applications', 'metrics', 'iam'],
      tabLabels: ['Devices', 'Provisioning', 'Applications', 'Metrics', 'IAM'],
      icons: ['desktop', 'box', 'application', 'timeline-area-chart', 'user'],
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
      case 'provisioning':
        selectedIndex = 1;
        break;
      case 'applications':
        selectedIndex = 2;
        break;
      case 'metrics':
        selectedIndex = 3;
        break;
      case 'iam':
        selectedIndex = 4;
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
      case 'provisioning':
        return (
          <Switch>
            <Route
              path={`${match.path}/deviceregistrationtokens/create`}
              render={route => (
                <CreateDeviceRegistrationToken
                  user={user}
                  projectName={projectName}
                  history={this.props.history}
                />
              )}
            />
            <Route
              exact
              path={`${match.path}/deviceregistrationtokens/:tokenName`}
              render={route => (
                <Redirect
                  to={`${route.match.url}/overview`}
                />
              )}
            />
            <Route
              path={`${match.path}/deviceregistrationtokens/:tokenName/:appTab`}
              render={route => (
                <DeviceRegistrationToken
                  user={user}
                  projectName={projectName}
                  tokenName={route.match.params.tokenName}
                  applicationRoute={route.match}
                  history={this.props.history}
                />
              )}
            />
            <Route
              exact
              path={match.path}
              render={route => (
                <Provisioning
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
      case 'metrics':
        return (
          <Switch>
            <Route
              exact
              path={match.path}
              render={route => (
                <Metrics
                  user={user}
                  projectName={projectName}
                  match={route.match}
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
                <ProjectSettings
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
      <Fragment>
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
      </Fragment>
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
      <Fragment>
        {user ? (
          <OuterOogie user={user} history={this.props.history} />
        ) : (
          <CustomSpinner />
        )}
      </Fragment>
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
