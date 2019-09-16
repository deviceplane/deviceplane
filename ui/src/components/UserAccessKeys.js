import React, { Component } from 'react';
import './../App.css';
import config from '../config.js';
import utils from '../utils.js';
import { toaster, Pane, majorScale, Heading, Alert, Button, Table, Dialog, Code } from 'evergreen-ui';
import axios from 'axios';

export default class UserAccessKeys extends Component {
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
    axios.get(`${config.endpoint}/useraccesskeys`, {
      withCredentials: true
    })
      .then((response) => {
        this.setState({
          accessKeys: response.data
        });
      })
      .catch((error) => {
        console.log(error);
      });
  }

  createAccessKey() {
    this.setState({
      backendError: null
    });

    axios.post(`${config.endpoint}/useraccesskeys`, {
    }, {
      withCredentials: true
    })
      .then((response) => {
        this.setState({
          showAccessKeyCreated: true,
          newAccessKey: response.data.value
        });
      })
      .catch((error) => {
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

  deleteAccessKey = (event) => {
    this.setState({
      backendError: null
    });

    axios.delete(`${config.endpoint}/useraccesskeys/${event.target.id}`, {
      withCredentials: true
    })
      .then((response) => {
        toaster.success('Successfully deleted access key.');
        this.loadAccessKeys();
      })
      .catch((error) => {
        if (utils.is4xx(error.response.status)) {
          this.setState({
            backendError: utils.convertErrorMessage(error.response.data)
          });
        } else {
          toaster.danger('Access key was not deleted.');
          console.log(error);
        }
      });
  }

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
        <Pane zIndex={1} flexShrink={0} elevation={0} backgroundColor="white">
          <Pane padding={majorScale(2)}>
            <Heading size={600}>User Access Keys</Heading>
          </Pane>
        </Pane>
        <Pane
          display="flex"
          flexDirection="column"
          margin={majorScale(4)}
        >
          {this.state.backendError && (
            <Alert marginBottom={majorScale(2)} paddingTop={majorScale(2)} paddingBottom={majorScale(2)} intent="warning" title={this.state.backendError} />
          )}
          <Pane><Button
            marginBottom={majorScale(2)}
            appearance="primary"
            onClick={() => this.createAccessKey()}
            justifyContent="center"
          >
            Create Access Key
          </Button>
          </Pane>
          {this.state.accessKeys && this.state.accessKeys.length > 0 && (
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
                      <Button iconBefore="trash" intent="danger" id={accessKey.id} onClick={(event) => this.deleteAccessKey(event)}>Delete Access Key</Button>
                    </Table.TextCell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          )}
          <Dialog
            isShown={this.state.showAccessKeyCreated}
            title="Access Key Created"
            onCloseComplete={() => this.closeAccessKeyDialog()}
            hasFooter={false}
          >
            <Pane
              display="flex"
              flexDirection="column"
            >
              <Heading paddingTop={majorScale(2)} paddingBottom={majorScale(2)}>{`Access Key: `}</Heading>
              <Pane marginBottom={majorScale(4)}>
                <Code>{this.state.newAccessKey}</Code>
              </Pane>
            </Pane>
            <Alert intent="warning" title="Save the info above! This is the only time you'll be able to use it.">
              {`If you lose it, you'll need to create a new access key.`}
            </Alert>
            <Button marginTop={16} appearance="primary" onClick={() => this.closeAccessKeyDialog()}>
              Close
            </Button>
          </Dialog>
        </Pane>
      </React.Fragment>
    );
  }
}