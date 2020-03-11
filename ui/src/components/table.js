import React from 'react';
import styled from 'styled-components';
import { useLinkProps } from 'react-navi';

import { labelColor } from '../helpers/labels';
import { DeviceLabelKey } from './device-label';
import { Checkbox, Button, Row, Grid, Icon, Link } from './core';

export const SelectColumn = {
  id: 'select',
  Header: ({ toggleAllRowsSelected, isAllRowsSelected }) => (
    <Checkbox checked={isAllRowsSelected} onChange={toggleAllRowsSelected} />
  ),
  Cell: ({ row: { isSelected, toggleRowSelected } }) => (
    <Row onClick={e => e.stopPropagation()} alignSelf="flex-start">
      <Checkbox checked={isSelected} onChange={toggleRowSelected} />
    </Row>
  ),
  minWidth: '40px',
  maxWidth: '40px',
  cellStyle: {
    backgroundColor: 'grays.5',
  },
};

export const DeviceLabelKeyColumn = {
  Header: 'Labels',
  accessor: 'labels',
  Cell: ({ cell: { value } }) => (
    <Row marginY={-2}>
      {value.map(label => (
        <DeviceLabelKey key={label} label={label} color={labelColor(label)} />
      ))}
    </Row>
  ),
};

export const SaveOrCancelCell = ({ onSave, onCancel }) => (
  <Row>
    <Button
      title={<Icon icon="floppy-disk" size={16} color="primary" />}
      variant="icon"
      onClick={onSave}
    />
    <Button
      title={<Icon icon="cross" size={16} color="pureWhite" />}
      variant="iconSecondary"
      onClick={onCancel}
      marginLeft={3}
    />
  </Row>
);

const A = styled.a`
  overflow-x: hidden;
  overflow-y: inherit;
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
  overflow: auto;
`;

const TableHead = styled.thead`
  display: contents;
`;

const TableBody = styled.tbody`
  display: contents;
`;

const TableRow = styled.tr`
  cursor: ${props => (props.clickable ? 'pointer' : 'default')};
  transition: ${props => props.theme.transitions[0]};
  display: contents;

  &:hover td {
    background-color: ${props =>
      props.clickable ? props.theme.colors.grays[3] : props.theme.colors.black};
  }
`;

const Cell = styled.td`
  display: flex;
  padding: 8px 12px;
  border-bottom: 1px solid ${props => props.theme.colors.grays[3]};
`;

const HeaderCell = styled.th`
  position: sticky;
  top: 0;
  z-index: 1;
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
    fill: ${props => props.theme.colors.white};
  }
`;

const Table = ({
  columns,
  rows,
  prepareRow,
  placeholder,
  maxHeight,
  headers,
  getTableBodyProps,
  getTableProps,
  onRowClick,
  rowHref,
}) => {
  const clickable = onRowClick || rowHref;

  return (
    <>
      <StyledTable
        {...getTableProps()}
        maxHeight={maxHeight}
        gridTemplateColumns={columns
          .map(col => {
            const minWidth = col.minWidth || 'min-content';
            const maxWidth =
              col.maxWidth === Number.MAX_SAFE_INTEGER ? '1fr' : col.maxWidth;
            return `minmax(${minWidth}, ${maxWidth})`;
          })
          .join(' ')}
      >
        <TableHead>
          <TableRow>
            {headers.map(column => (
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
                      icon={column.isSortedDesc ? 'chevron-down' : 'chevron-up'}
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
            ))}
          </TableRow>
        </TableHead>
        <TableBody {...getTableBodyProps()} overflowY="auto">
          {rows.map(row => {
            prepareRow(row);
            const cells = row.cells.map(cell => (
              <Cell
                {...cell.getCellProps()}
                style={{
                  textAlign:
                    isNaN(cell.value) || cell.value === '-' ? 'left' : 'right',
                  ...cell.column.cellStyle,
                }}
                clickable={clickable}
              >
                {rowHref && cell.column.id !== 'select' ? (
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
                clickable={clickable}
                onClick={() => onRowClick && onRowClick(row.original)}
                position="relative"
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
          padding={3}
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
