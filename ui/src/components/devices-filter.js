import React, { Fragment, Component } from 'react';
import {
  Pane,
  Dialog,
  Select,
  IconButton,
  Button,
  majorScale,
  Strong,
  TextInput,
  minorScale
} from 'evergreen-ui';

const initialState = {
  filters: [{ property: 'status', operator: 'is', value: 'online' }]
};

export default class DevicesFilter extends Component {
  state = initialState;

  renderOperatorField = (filter, index) => {
    switch (filter.property) {
      case 'label':
        return (
          <Select
            value={filter.operator}
            onChange={event => {
              const { value: operator } = event.target;
              this.setState(({ filters }) => ({
                filters: filters.map((f, i) =>
                  i === index ? { ...filter, operator } : f
                )
              }));
            }}
            marginRight={majorScale(1)}
          >
            <option value="is">is</option>
            <option value="is not">is not</option>
            <option value="has key">has key</option>
            <option value="does not have key">does not have key</option>
          </Select>
        );
      case 'status':
        return (
          <Select
            value={filter.operator}
            onChange={event => {
              const { value: operator } = event.target;
              this.setState(({ filters }) => ({
                filters: filters.map((f, i) =>
                  i === index ? { ...filter, operator } : f
                )
              }));
            }}
            marginRight={majorScale(1)}
          >
            <option value="is">is</option>
            <option value="is not">is not</option>
          </Select>
        );
    }
  };

  renderValueField = (filter, index) => {
    switch (filter.property) {
      case 'label':
        switch (filter.operator) {
          case 'is':
          case 'is not':
            return (
              <Pane
                display="flex"
                flexDirection="column"
                flex="1"
                marginRight={majorScale(1)}
              >
                <TextInput
                  width="auto"
                  placeholder="Key"
                  marginBottom={minorScale(2)}
                  onChange={event => {
                    const { value: key } = event.target;
                    this.setState(({ filters }) => ({
                      filters: filters.map((f, i) =>
                        i === index ? { ...filter, key } : f
                      )
                    }));
                  }}
                />
                <TextInput
                  width="auto"
                  placeholder="Value"
                  onChange={event => {
                    const { value } = event.target;
                    this.setState(({ filters }) => ({
                      filters: filters.map((f, i) =>
                        i === index ? { ...filter, value } : f
                      )
                    }));
                  }}
                />
              </Pane>
            );
          case 'has key':
          case 'does not have key':
            return (
              <Pane display="flex" flex="1" marginRight={majorScale(1)}>
                <TextInput
                  width="auto"
                  placeholder="Key"
                  onChange={event => {
                    const { value: key } = event.target;
                    this.setState(({ filters }) => ({
                      filters: filters.map((f, i) =>
                        i === index ? { ...filter, key } : f
                      )
                    }));
                  }}
                />
              </Pane>
            );
        }
      case 'status':
        return (
          <Select
            value={filter.value}
            onChange={event => {
              const { value } = event.target;
              this.setState(({ filters }) => ({
                filters: filters.map((f, i) =>
                  i === index ? { ...filter, value } : f
                )
              }));
            }}
            marginRight={majorScale(1)}
          >
            <option value="online">Online</option>
            <option value="offline">Offline</option>
          </Select>
        );
    }
  };

  render() {
    const { show, onClose, onSubmit } = this.props;
    const { filters } = this.state;

    return (
      <Pane>
        <Dialog
          isShown={show}
          title="Filter Devices"
          onCloseComplete={onClose}
          onConfirm={() => {
            onSubmit(filters);
            this.setState(initialState);
          }}
          confirmLabel="Filter"
          hasCancel={false}
          style={{ maxHeight: majorScale(12), overflowY: 'auto' }}
        >
          <Pane display="flex" flexDirection="column">
            {filters.map((filter, index) => (
              <Fragment key={index}>
                <Pane display="flex" justifyContent="space-around">
                  <Select
                    value={filter.property}
                    onChange={event => {
                      const { value: property } = event.target;
                      this.setState(({ filters }) => ({
                        filters: filters.map((f, i) =>
                          i === index
                            ? {
                                ...filter,
                                property,
                                value:
                                  property === 'label'
                                    ? undefined
                                    : filter.value
                              }
                            : f
                        )
                      }));
                    }}
                    marginRight={majorScale(1)}
                  >
                    <option value="status">Status</option>
                    <option value="label">Label</option>
                  </Select>

                  {this.renderOperatorField(filter, index)}

                  {this.renderValueField(filter, index)}

                  {index > 0 ? (
                    <IconButton
                      icon="cross"
                      intent="danger"
                      appearance="minimal"
                      onClick={() =>
                        this.setState(({ filters }) => ({
                          filters: filters.filter((_, i) => i !== index)
                        }))
                      }
                    />
                  ) : (
                    <Pane width={majorScale(4)} />
                  )}
                </Pane>
                {index < filters.length - 1 && (
                  <Pane marginY={majorScale(2)}>
                    <Strong
                      size={300}
                      paddingX={majorScale(2)}
                      paddingY={majorScale(1)}
                      backgroundColor="#E4E7EB"
                      borderRadius={3}
                    >
                      OR
                    </Strong>
                  </Pane>
                )}
              </Fragment>
            ))}
          </Pane>
          <Pane display="flex" flexDirection="column" marginTop={majorScale(4)}>
            <Pane>
              <Button
                intent="none"
                iconBefore="plus"
                onClick={() => {
                  this.setState(({ filters }) => ({
                    filters: [
                      ...filters,
                      { property: 'status', operator: 'is', value: 'online' }
                    ]
                  }));
                }}
              >
                Add Condition
              </Button>
            </Pane>
          </Pane>
        </Dialog>
      </Pane>
    );
  }
}
