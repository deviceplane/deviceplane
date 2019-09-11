import React, { Component } from 'react';
import './../App.css';
import { Spinner, Pane } from 'evergreen-ui';

export default class CustomSpinner extends Component {
  render() {
    return (
      <Pane display="flex" alignItems="center" justifyContent="center" height={400}>	          
        <Spinner />
      </Pane>
    );
  }
}