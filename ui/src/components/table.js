import React from 'react';
import { useTable, useSortBy } from 'react-table';
import styled from 'styled-components';
import { Icon } from 'evergreen-ui';

import { Column, Row } from './core';

const Container = styled(Column)``;

Container.defaultProps = { borderRadius: 1, borderColor: 'white' };

const Cell = styled(Row)`
  flex: 1 0 0%;
  overflow: hidden;
`;

Cell.defaultProps = {
  padding: 3,
};

const TableRow = styled(Row)`
  align-items: flex-start;
  border-bottom: 1px solid ${props => props.theme.colors.grays[2]};
  cursor: ${props => (props.selectable ? 'pointer' : 'default')};
  transition: background-color 150ms;

  &:hover {
    background-color: ${props =>
      props.selectable
        ? props.theme.colors.grays[1]
        : props.theme.colors.black};
  }
`;

const Header = styled(Row)`
  min-height: 50px;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
  text-transform: uppercase;
  align-items: center;
`;

Header.defaultProps = {
  fontSize: 0,
  fontWeight: 4,
  color: 'white',
  bg: 'grays.0',
};

const Table = ({ columns, data, onRowSelect, placeholder }) => {
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

  return (
    <Container {...getTableProps()}>
      <Header flex={1}>
        {headerGroups.map(headerGroup => (
          <Row flex={1} {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map(column => (
              <Cell
                {...column.getHeaderProps(column.getSortByToggleProps())}
                style={column.style}
              >
                {column.render('Header')}
                <Row>
                  {column.isSorted ? (
                    column.isSortedDesc ? (
                      <Icon
                        icon="chevron-down"
                        color="white"
                        size={14}
                        marginLeft={8}
                      />
                    ) : (
                      <Icon
                        icon="chevron-up"
                        color="white"
                        size={14}
                        marginLeft={8}
                      />
                    )
                  ) : (
                    ''
                  )}
                </Row>
              </Cell>
            ))}
          </Row>
        ))}
      </Header>
      <Column {...getTableBodyProps()}>
        {rows.length === 0 && (
          <Row
            justifyContent="center"
            padding={3}
            borderBottom={0}
            borderColor="grays.3"
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
              onClick={() => onRowSelect(data[row.index])}
            >
              {row.cells.map(cell => (
                <Cell {...cell.getCellProps()} style={cell.column.style || {}}>
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
