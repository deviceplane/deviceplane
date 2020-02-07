import React from 'react';
import { useTable, useSortBy } from 'react-table';
import styled from 'styled-components';
import { useLinkProps } from 'react-navi';

import { Column, Row, Grid, Icon } from './core';

const A = styled.a`
  text-decoration: none;
  color: unset;
  flex: 1;
  margin: -8px -12px;
  padding: 8px 12px;
`;

const LinkCell = ({ children, href, ...rest }) => {
  const linkProps = useLinkProps({ href });
  return (
    <A {...linkProps} {...rest}>
      {children}
    </A>
  );
};

const StyledTable = styled(Grid).attrs({ as: 'table' })`
  width: auto;
  border-collapse: collapse;
`;

const TableHead = styled.thead`
  display: contents;
`;

const TableBody = styled.tbody`
  display: contents;
`;

const TableRow = styled.tr`
  border-bottom: 1px solid ${props => props.theme.colors.grays[1]};
  cursor: ${props => (props.selectable ? 'pointer' : 'default')};
  transition: ${props => props.theme.transitions[0]};
  display: contents;

  &:hover td {
    background-color: ${props =>
      props.selectable
        ? props.theme.colors.grays[4]
        : props.theme.colors.black};
  }
`;

const HeaderCell = styled.th`
  position: sticky;
  top: 0;
  text-transform: uppercase;
  font-size: 14px;
  font-weight: 500;
  padding: 16px 12px;
  text-align: left;
  color: ${props => props.theme.colors.white};
  background-color: ${props => props.theme.colors.grays[0]};

  & ${Cell} svg {
    transition: fill 200ms;
  }

  & ${Cell}:hover svg {
    fill: ${props => props.theme.colors.white} !important;
  }
`;

const Cell = styled.td`
  display: flex;
  padding: 8px 12px;
`;

const Table = ({
  columns,
  data,
  onRowSelect,
  placeholder,
  editRow,
  rowHref,
}) => {
  const selectable = onRowSelect || rowHref;
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
    <>
      <StyledTable
        {...getTableProps()}
        overflowY="hidden"
        gridTemplateColumns={columns
          .map(
            ({ minWidth = 'min-content', maxWidth = '1fr' }) =>
              `minmax(${minWidth}, ${maxWidth})`
          )
          .join(' ')}
      >
        <TableHead>
          <TableRow>
            {headerGroups.map(headerGroup =>
              headerGroup.headers.map(column => (
                <HeaderCell
                  {...column.getHeaderProps(column.getSortByToggleProps())}
                  title=""
                  style={{
                    cursor: column.canSort ? 'pointer' : 'default',
                  }}
                >
                  <Row justifyContent="space-between">
                    {column.render('Header')}
                    {column.isSorted ? (
                      <Icon
                        icon={
                          column.isSortedDesc ? 'chevron-down' : 'chevron-up'
                        }
                        size={14}
                        color="white"
                        marginLeft={2}
                      />
                    ) : column.canSort ? (
                      <Icon
                        size={12}
                        icon="expand-all"
                        color="grays.5"
                        marginLeft={2}
                      />
                    ) : null}
                  </Row>
                </HeaderCell>
              ))
            )}
          </TableRow>
        </TableHead>
        <TableBody {...getTableBodyProps()} overflowY="auto">
          {rows.map((row, i) => {
            prepareRow(row);
            const cells = row.cells.map(cell => (
              <Cell
                {...cell.getCellProps()}
                style={{
                  textAlign:
                    isNaN(cell.value) || cell.value === '-' ? 'left' : 'right',
                  ...cell.column.cellStyle,
                }}
                selectable={selectable}
                overflow={editRow ? 'visible' : 'hidden'}
              >
                {rowHref ? (
                  <LinkCell href={rowHref(row.original)}>
                    {cell.render('Cell')}
                  </LinkCell>
                ) : (
                  cell.render('Cell')
                )}
              </Cell>
            ));
            const tableRow = (
              <TableRow
                {...row.getRowProps()}
                selectable={selectable}
                onClick={handleRowClick(row.index)}
                position="relative"
                style={{ transform: 'translate2d(0,0)' }}
              >
                {cells}
              </TableRow>
            );
            return tableRow;
          })}
        </TableBody>
      </StyledTable>
      {rows.length === 0 && (
        <Row
          flex={1}
          justifyContent="center"
          padding={4}
          borderBottom={0}
          borderColor="grays.1"
        >
          {placeholder}
        </Row>
      )}
    </>
  );
};

export default Table;
