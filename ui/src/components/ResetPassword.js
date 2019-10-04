import React, { Component } from 'react';
import {
  toaster,
  Pane,
  majorScale,
  minorScale,
  Heading,
  Alert,
  TextInputField,
  Button,
  Text,
  Link
} from 'evergreen-ui';
import axios from 'axios';

import config from '../config.js';
import utils from '../utils.js';
import Logo from './logo';

export default class ResetPassword extends Component {
  constructor(props) {
    super(props);
    this.state = {
      email: '',
      emailValidationMessage: null,
      showUserNotFound: false
    };
  }

  handleUpdateEmail = event => {
    this.setState({
      email: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();

    var emailValidationMessage = null;
    var showUserNotFound = false;

    if (!utils.emailRegex.test(this.state.email)) {
      emailValidationMessage = 'Please enter a valid email.';
    }

    this.setState({
      emailValidationMessage: emailValidationMessage,
      showUserNotFound: showUserNotFound
    });

    if (emailValidationMessage === null) {
      axios
        .post(`${config.endpoint}/recoverpassword`, {
          email: this.state.email
        })
        .then(response => {
          this.props.history.push(`/login`);
          toaster.success(
            'Password recovery email sent. Please check your email to reset your password.'
          );
        })
        .catch(error => {
          if (error.response.status === 404) {
            this.setState({
              showUserNotFound: true
            });
          } else {
            toaster.danger(
              'There was an error with your e-mail. Please contact us at info@deviceplane.com.'
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
        width={majorScale(50)}
        marginX="auto"
        marginY={200}
        padding={majorScale(6)}
        elevation={2}
      >
        <Pane
          display="flex"
          flexDirection="row"
          alignItems="center"
          marginX="auto"
        >
          <Pane paddingRight={minorScale(3)}>
            <Logo />
          </Pane>
          <Heading size={600}>Device Plane</Heading>
        </Pane>
        <React.Fragment>
          <Pane display="flex" justifyContent="center" padding={majorScale(2)}>
            <Heading size={600}>Reset your Password</Heading>
          </Pane>
          {this.state.showUserNotFound && (
            <Alert
              marginBottom={majorScale(2)}
              paddingTop={majorScale(2)}
              paddingBottom={majorScale(2)}
              intent="warning"
              title="User doesn't exist"
            >
              There is no user with that email address.
            </Alert>
          )}
          <Pane is="form">
            <TextInputField
              label="Email"
              onChange={this.handleUpdateEmail}
              value={this.state.email}
              isInvalid={this.state.emailValidationMessage !== null}
              validationMessage={this.state.emailValidationMessage}
            />
            <Pane>
              <Button
                width="100%"
                appearance="primary"
                justifyContent="center"
                onClick={this.handleSubmit}
              >
                Send Reset Password Email
              </Button>
            </Pane>
          </Pane>
        </React.Fragment>
        <Pane
          display="flex"
          flexDirection="row"
          justifyContent="center"
          paddingTop={majorScale(3)}
        >
          <Text color="muted">{`Already have an account?`}</Text>
          <Link paddingLeft={majorScale(1)} href="/login">
            Log in
          </Link>
        </Pane>
        <Pane display="flex" flexDirection="row" justifyContent="center">
          <Text color="muted">{`Don't have an account?`}</Text>
          <Link paddingLeft={majorScale(1)} href="/register">
            Sign up
          </Link>
        </Pane>
      </Pane>
    );
  }
}
