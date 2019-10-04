import React, { Component, Fragment } from 'react';
import { Switch, Route } from 'react-router-dom';
import { Pane, majorScale, Tablist, Tab } from 'evergreen-ui';

import TopHeader from '../../components/TopHeader';

import Roles from './roles';
import Role from './role';
import CreateRole from './create-role';

import Member from './member';
import Members from './members';
import AddMember from './add-member';

import ServiceAccount from './service-account';
import ServiceAccounts from './service-accounts';
import CreateServiceAccount from './create-service-account';

export default class Iam extends Component {
  state = {
    tabs: ['members', 'serviceaccounts', 'roles'],
    tabLabels: ['Members', 'Service Accounts', 'Roles']
  };

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
    return (
      <Fragment>
        <TopHeader
          user={this.props.user}
          heading="IAM"
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
      </Fragment>
    );
  }
}

const MembersRouter = ({ projectName, match, user }) => (
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

const ServiceAccountsRouter = ({ projectName, match, history }) => (
  <Switch>
    <Route
      path={`${match.path}/create`}
      render={route => (
        <CreateServiceAccount projectName={projectName} history={history} />
      )}
    />
    <Route
      path={`${match.path}/:serviceAccountName`}
      render={route => (
        <ServiceAccount
          projectName={projectName}
          serviceAccountName={route.match.params.serviceAccountName}
          history={history}
        />
      )}
    />
    <Route
      exact
      path={match.path}
      render={route => (
        <ServiceAccounts projectName={projectName} history={history} />
      )}
    />
  </Switch>
);

const RolesRouter = ({ projectName, match, history }) => (
  <Switch>
    <Route
      path={`${match.path}/create`}
      render={route => (
        <CreateRole projectName={projectName} history={history} />
      )}
    />
    <Route
      path={`${match.path}/:roleName`}
      render={route => (
        <Role
          projectName={projectName}
          roleName={route.match.params.roleName}
          history={history}
        />
      )}
    />
    <Route
      exact
      path={match.path}
      render={route => <Roles projectName={projectName} history={history} />}
    />
  </Switch>
);
