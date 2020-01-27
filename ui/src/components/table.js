import React from 'react';
import { useTable, useSortBy } from 'react-table';
import styled from 'styled-components';
import { Icon } from 'evergreen-ui';

import { Column, Row, Text } from './core';

const Container = styled(Column)``;

Container.defaultProps = { borderRadius: 1, borderColor: 'white' };

const Cell = styled(Row)`
  flex: 1 0 0%;
  box-sizing: content-box;
  padding: 8px 16px;
`;

Cell.defaultProps = {
  overflow: 'hidden',
};

const TableRow = styled(Row)`
  align-items: center;
  border-bottom: 1px solid ${props => props.theme.colors.grays[1]};
  cursor: ${props => (props.selectable ? 'pointer' : 'default')};
  transition: ${props => props.theme.transitions[0]};

  &:hover {
    background-color: ${props =>
      props.selectable
        ? props.theme.colors.grays[4]
        : props.theme.colors.black};
  }
`;

const Header = styled(Row)`
  min-height: 50px;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
  text-transform: uppercase;
  align-items: center;

  & ${Cell}:hover svg {
    display: block !important;
  }
`;

Header.defaultProps = {
  fontSize: 1,
  fontWeight: 2,
  color: 'white',
  bg: 'grays.0',
};

const Table = ({ columns, data, onRowSelect, placeholder, editRow }) => {
  const selectable = !!onRowSelect;
  onRowSelect = onRowSelect || function() {};

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
  } = useTable(
    {
      columns,
      data,
    },
    useSortBy
  );

  const handleRowClick = index => () => {
    const selection = window.getSelection();
    // Only select row if user is not highlighting text
    if (selection.type !== 'Range') {
      onRowSelect(data[index]);
    }
  };

  return (
    <Container {...getTableProps()} overflowY="hidden">
      <Header flex={1}>
        {headerGroups.map(headerGroup => (
          <Row flex={1} {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map(column => (
              <Cell
                {...column.getHeaderProps(column.getSortByToggleProps())}
                style={{
                  ...column.style,
                  cursor: column.canSort ? 'pointer' : 'default',
                  alignSelf: 'center',
                }}
              >
                {column.render('Header')}
                <Row>
                  {column.isSorted ? (
                    <Icon
                      icon={column.isSortedDesc ? 'chevron-down' : 'chevron-up'}
                      size={14}
                      color="white"
                      marginLeft={8}
                    />
                  ) : column.canSort ? (
                    <Icon
                      size={14}
                      icon="expand-all"
                      marginLeft={8}
                      style={{ display: 'none' }}
                    />
                  ) : null}
                </Row>
              </Cell>
            ))}
          </Row>
        ))}
      </Header>
      <Column {...getTableBodyProps()} overflowY="auto">
        {rows.length === 0 && (
          <Row
            justifyContent="center"
            padding={4}
            borderBottom={0}
            borderColor="grays.1"
          >
            {placeholder}
          </Row>
        )}
        {rows.map((row, i) => {
          prepareRow(row);
          return (
            <TableRow
              {...row.getRowProps()}
              selectable={selectable}
              onClick={handleRowClick(row.index)}
              position="relative"
              style={{ transform: 'translate2d(0,0)' }}
            >
              {row.cells.map(cell => (
                <Cell
                  {...cell.getCellProps()}
                  style={{
                    justifyContent:
                      isNaN(cell.value) || cell.value === '-'
                        ? 'flex-start'
                        : 'flex-end',
                    ...cell.column.style,
                  }}
                  overflow={editRow ? 'visible' : 'hidden'}
                >
                  {cell.render('Cell')}
                </Cell>
              ))}
            </TableRow>
          );
        })}
      </Column>
    </Container>
  );
};

export default Table;
