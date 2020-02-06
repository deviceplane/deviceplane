import React from 'react';
import { useTable, useSortBy } from 'react-table';
import styled from 'styled-components';
import { useLinkProps } from 'react-navi';

import { Column, Row, Icon } from './core';

const Container = styled(Column)``;

Container.defaultProps = { borderRadius: 1, borderColor: 'white' };

const Cell = styled(Row)`
  flex: 1 0 0%;
  min-width: 50px;
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

const A = styled.a`
  text-decoration: none;
  color: unset;
`;

const LinkRow = ({ children, href, ...rest }) => {
  const linkProps = useLinkProps({ href });
  return (
    <A {...linkProps} {...rest}>
      {children}
    </A>
  );
};

const Header = styled(Row)`
  min-height: 50px;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
  text-transform: uppercase;
  align-items: center;

  & ${Cell} svg {
    transition: fill 200ms;
  }

  & ${Cell}:hover svg {
    fill: ${props => props.theme.colors.white} !important;
  }
`;

Header.defaultProps = {
  fontSize: 1,
  fontWeight: 2,
  color: 'white',
  bg: 'grays.0',
};

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
    <Container {...getTableProps()} overflowY="hidden">
      <Header flex={1}>
        {headerGroups.map(headerGroup => (
          <Row flex={1} {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map(column => (
              <Cell
                {...column.getHeaderProps(column.getSortByToggleProps())}
                title=""
                style={{
                  ...column.style,
                  cursor: column.canSort ? 'pointer' : 'default',
                  alignSelf: 'center',
                  justifyContent: 'space-between',
                }}
              >
                {column.render('Header')}
                <Row marginLeft={2} alignItems="center">
                  {column.isSorted ? (
                    <Icon
                      icon={column.isSortedDesc ? 'chevron-down' : 'chevron-up'}
                      size={14}
                      color="white"
                    />
                  ) : column.canSort ? (
                    <Icon size={12} icon="expand-all" color="grays.5" />
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
          const cells = row.cells.map(cell => (
            <Cell
              {...cell.getCellProps()}
              style={{
                justifyContent:
                  isNaN(cell.value) || cell.value === '-'
                    ? 'flex-start'
                    : 'flex-end',
                ...cell.column.style,
                ...cell.column.cellStyle,
              }}
              overflow={editRow ? 'visible' : 'hidden'}
            >
              {cell.render('Cell')}
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
          if (rowHref) {
            return <LinkRow href={rowHref(row.original)}>{tableRow}</LinkRow>;
          }
          return tableRow;
        })}
      </Column>
    </Container>
  );
};

export default Table;
