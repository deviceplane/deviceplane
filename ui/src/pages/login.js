import React, { Component } from 'react';
import {
  Button,
  Pane,
  Heading,
  majorScale,
  Text,
  Link,
  Alert,
  minorScale,
  TextInputField
} from 'evergreen-ui';
import axios from 'axios';

import config from '../config';
import utils from '../utils';
import Logo from '../components/logo';

export default class Login extends Component {
  state = {
    email: '',
    password: '',
    emailValidationMessage: null,
    passwordValidationMessage: null,
    submitError: null
  };

  handleUpdateEmail = event => {
    this.setState({
      email: event.target.value
    });
  };

  handleUpdatePassword = event => {
    this.setState({
      password: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();
    var emailValidationMessage = null;
    var passwordValidationMessage = null;
    var submitError = null;

    if (!utils.emailRegex.test(this.state.email)) {
      emailValidationMessage = 'Please enter a valid email.';
    }

    if (this.state.password === '') {
      passwordValidationMessage = 'Please enter a password.';
    }

    this.setState({
      emailValidationMessage: emailValidationMessage,
      passwordValidationMessage: passwordValidationMessage,
      submitError: submitError
    });

    if (emailValidationMessage === null && passwordValidationMessage === null) {
      axios
        .post(
          `${config.endpoint}/login`,
          {
            email: this.state.email,
            password: this.state.password
          },
          {
            withCredentials: true
          }
        )
        .then(response => {
          window.location.reload();
        })
        .catch(error => {
          this.setState({
            submitError: 'Invalid Username/Password'
          });
          console.log(error);
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
          marginX="auto"
          alignItems="center"
          paddingBottom={majorScale(2)}
        >
          <Pane paddingRight={minorScale(3)}>
            <Logo />
          </Pane>
          <Heading size={600}>Device Plane</Heading>
        </Pane>
        {this.state.submitError && (
          <Alert
            marginBottom={majorScale(2)}
            paddingTop={majorScale(2)}
            paddingBottom={majorScale(2)}
            intent="warning"
            title={this.state.submitError}
          />
        )}
        <Pane is="form">
          <TextInputField
            label="Email"
            onChange={this.handleUpdateEmail}
            value={this.state.email}
            isInvalid={this.state.emailValidationMessage !== null}
            validationMessage={this.state.emailValidationMessage}
          />
          <TextInputField
            type="password"
            label="Password"
            onChange={this.handleUpdatePassword}
            value={this.state.password}
            isInvalid={this.state.passwordValidationMessage !== null}
            validationMessage={this.state.passwordValidationMessage}
          />
          <Button
            width="100%"
            appearance="primary"
            justifyContent="center"
            onClick={this.handleSubmit}
          >
            Submit
          </Button>
        </Pane>
        <Pane
          display="flex"
          flexDirection="row"
          justifyContent="center"
          paddingTop={majorScale(3)}
        >
          <Text color="muted">{`Forgot your Password?`}</Text>
          <Link paddingLeft={majorScale(1)} href="/forgot">
            Reset your password
          </Link>
        </Pane>
        <Pane
          display="flex"
          flexDirection="row"
          justifyContent="center"
          paddingTop={majorScale(1)}
        >
          <Text color="muted">{`Don't have an account?`}</Text>
          <Link paddingLeft={majorScale(1)} href="/register">
            Sign up
          </Link>
        </Pane>
      </Pane>
    );
  }
}
