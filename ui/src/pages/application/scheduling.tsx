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
  // @ts-ignore
} from 'evergreen-ui';
import axios from 'axios';

import config from '../../config';
import utils from '../../utils';
import InnerCard from '../../components/InnerCard';
import { DevicesFilterButtons } from '../../components/DevicesFilterButtons';
import { Query, Filter, DevicesFilter, LabelValueCondition } from '../../components/DevicesFilter';

interface Props {
  application: any,
  projectName: string,
  history: any,
}

interface State {
  schedulingRule: Query,
  backendError: any,
  showFilterDialog: boolean,
}

export default class ApplicationScheduling extends Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {
      schedulingRule: this.props.application.schedulingRule,
      backendError: null,
      showFilterDialog: false,
    };
  }

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
          schedulingRule: this.state.schedulingRule,
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

  filterDevices() {
    // TODO: fetch devices and show them
  }

  removeFilter = (index: number) => {
    this.setState({
      schedulingRule: this.state.schedulingRule.filter((_, i) => i !== index),
    }, this.filterDevices);
  };

  addFilter = (filter: Filter) => {
    this.setState({
      showFilterDialog: false,
      schedulingRule: [...this.state.schedulingRule, filter]
    }, this.filterDevices);
  };

  clearFilters = () => {
    this.setState({
      schedulingRule: []
    }, this.filterDevices);
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
            <Pane
              display={'flex'}
              justifyContent={'space-between'}
              alignItems={'center'}
              paddingBottom={majorScale(2)}
            >
              <Heading size={600}>
                Scheduling
              </Heading>
              <Button
                iconBefore="plus"
                onClick={() => this.setState({ showFilterDialog: true })}
              >
                Add Filter
              </Button>
            </Pane>
            <Pane
            backgroundColor={'#E4E7EB'}
            borderRadius={'5px'}
            minHeight={'60px'}
            >
              <DevicesFilterButtons
                query={this.state.schedulingRule}
                canRemoveFilter={true}
                removeFilter={this.removeFilter}
              />
            </Pane>
            <DevicesFilter
              whitelistedConditions={[LabelValueCondition]}
              show={this.state.showFilterDialog}
              onClose={() => this.setState({ showFilterDialog: false })}
              onSubmit={this.addFilter}
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
