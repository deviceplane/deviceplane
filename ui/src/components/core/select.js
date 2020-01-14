import React from 'react';
import ReactSelect, { components } from 'react-select';
import CreatableSelect from 'react-select/creatable';

import theme from '../../theme';

const styles = {
  container: () => ({
    display: 'flex',
    flex: 1,
    position: 'relative',
  }),
  option: (provided, { isFocused, isSelected, selectProps: { variant } }) => ({
    transition: 'background-color 200ms ease',
    backgroundColor: isSelected
      ? theme.colors.white
      : isFocused
      ? theme.colors.grays[3]
      : variant === 'black'
      ? theme.colors.black
      : theme.colors.grays[0],
    color: isSelected ? theme.colors.black : theme.colors.white,
    padding: '8px',
    cursor: 'pointer',
    margin: 0,
  }),
  menu: (provided, { selectProps: { variant } }) => ({
    ...provided,
    marginTop: '4px',
    backgroundColor:
      variant === 'black' ? theme.colors.black : theme.colors.grays[0],
    borderRadius: `${theme.radii[1]}px`,
    border: `1px solid ${theme.colors.white}`,
    boxShadow: 'none',
  }),
  menuList: provided => ({ ...provided, padding: 0 }),
  control: (_, { selectProps: { variant } }) => ({
    // none of react-select's styles are passed to <Control />
    display: 'flex',
    flex: 1,
    padding: 0,
    backgroundColor:
      variant === 'black' ? theme.colors.black : theme.colors.grays[0],
    borderRadius: `${theme.radii[1]}px`,
  }),
  input: () => ({
    padding: '4px',
    fontSize: theme.fontSizes[2],
    color: theme.colors.grays[12],
    fontWeight: theme.fontWeights[1],
  }),
  placeholder: () => ({
    fontSize: theme.fontSizes[2],
    color: theme.colors.grays[8],
  }),
  valueContainer: provided => ({ ...provided, padding: '0 8px' }),
  multiValue: () => ({
    display: 'flex',
    margin: '8px 8px 8px 0',
  }),
  multiValueLabel: () => ({
    padding: '4px 6px',
    backgroundColor: theme.colors.white,
    borderTopLeftRadius: '3px',
    borderBottomLeftRadius: '3px',
    color: theme.colors.black,
    fontSize: theme.fontSizes[1],
    fontWeight: theme.fontWeights[1],
  }),
  multiValueRemove: () => ({
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
    backgroundColor: theme.colors.black,
    color: theme.colors.red,
    borderTopRightRadius: '3px',
    borderBottomRightRadius: '3px',
    padding: '4px',
    fontSize: '18px',
    ':hover': {
      color: theme.colors.pureWhite,
    },
  }),
  singleValue: provided => {
    return { ...provided, color: theme.colors.grays[12] };
  },
  indicatorsContainer: provided => ({
    ...provided,
    cursor: 'pointer',
    color: theme.colors.white,
    ':hover': {
      color: theme.colors.black,
      backgroundColor: theme.colors.white,
    },
  }),
  clearIndicator: provided => ({
    ...provided,
    cursor: 'pointer',
    color: theme.colors.red,
    ':hover': {
      color: theme.colors.white,
    },
  }),
  indicatorSeparator: () => ({}),
};

const Option = props => {
  if (props.selectProps.optionComponent) {
    return (
      <components.Option {...props}>
        <props.selectProps.optionComponent {...props} {...props.data.props} />
      </components.Option>
    );
  }
  return <components.Option {...props} />;
};

const SingleValue = props => {
  if (props.selectProps.singleComponent) {
    return (
      <props.selectProps.singleComponent {...props} {...props.data.props} />
    );
  }
  return <components.SingleValue {...props} />;
};

const MultiValueLabel = props => {
  if (props.selectProps.multiComponent) {
    return (
      <props.selectProps.multiComponent {...props} {...props.data.props} />
    );
  }
  return <components.MultiValueLabel {...props} />;
};

const Select = ({ searchable, multi, disabled, creatable, ...props }) => {
  const SelectComponent = creatable ? CreatableSelect : ReactSelect;

  return (
    <SelectComponent
      styles={styles}
      isSearchable={searchable}
      isDisabled={disabled}
      isMulti={multi}
      components={{ Option, MultiValueLabel, SingleValue }}
      closeMenuOnSelect={multi ? false : true}
      menuPosition="fixed"
      {...props}
      isClearable={false}
    />
  );
};

export default Select;
