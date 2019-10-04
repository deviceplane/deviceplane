import React, { Component } from 'react';
import {
  Button,
  Pane,
  Heading,
  majorScale,
  Alert,
  toaster,
  Label,
  Textarea
} from 'evergreen-ui';
import axios from 'axios';

import config from '../../config.js';
import utils from '../../utils.js';
import InnerCard from '../../components/InnerCard.js';

export default class ApplicationScheduling extends Component {
  state = {
    schedulingRule: this.props.application.settings.schedulingRule,
    backendError: null
  };

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
