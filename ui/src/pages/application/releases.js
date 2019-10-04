import React, { Component, Fragment } from 'react';
import axios from 'axios';
import moment from 'moment';
import {
  Pane,
  Link,
  Table,
  SideSheet,
  Dialog,
  Label,
  majorScale,
  Button,
  Heading,
  Alert
} from 'evergreen-ui';

import utils from '../../utils';
import config from '../../config';
import Editor from '../../components/Editor';

export default class Releases extends Component {
  state = {
    releases: [],
    showRelease: false,
    selectedRelease: null
  };

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
      <Fragment>
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
      </Fragment>
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
