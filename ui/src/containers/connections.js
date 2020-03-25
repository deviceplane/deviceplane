import React, { useMemo, useState } from 'react';
import { useTable, useSortBy, useRowSelect } from 'react-table';
import { useNavigation } from 'react-navi';

import api from '../api';
import Layout from '../components/layout';
import Card from '../components/card';
import Popup from '../components/popup';
import Table, { SelectColumn } from '../components/table';
import { Row, Button, Text, toaster } from '../components/core';

const Connections = ({
  route: {
    data: { params, connections },
  },
}) => {
  const navigation = useNavigation();
  const [showDeletePopup, setShowDeletePopup] = useState();

  const columns = useMemo(
    () => [
      SelectColumn,
      {
        Header: 'Name',
        accessor: 'name',
      },
      {
        Header: 'Port',
        accessor: 'port',
        maxWidth: '100px',
        minWidth: '100px',
      },
      {
        Header: 'Protocol',
        accessor: ({ protocol }) => protocol.toUpperCase(),
        maxWidth: '120px',
        minWidth: '120px',
      },
    ],
    []
  );
  const tableData = useMemo(() => connections, [connections]);

  const { selectedFlatRows, ...tableProps } = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy,
    useRowSelect
  );

  const submitDelete = async () => {
    setShowDeletePopup(false);

    for (let i = 0; i < selectedFlatRows.length; i++) {
      const { name } = selectedFlatRows[i].original;
      try {
        await api.deleteConnection({
          projectId: params.project,
          connectionId: name,
        });
      } catch (error) {
        console.error(error);
        toaster.danger('Connection deletion failed.');
        break;
      }
    }

    navigation.refresh();
  };

  return (
    <Layout alignItems="center">
      <Card
        title="Connections"
        size="xlarge"
        maxHeight="100%"
        actions={[
          {
            title: 'Create Connection',
            href: 'create',
          },
        ]}
      >
        <Row marginBottom={3}>
          <Button
            marginRight={4}
            title="Delete"
            variant="tertiaryDanger"
            disabled={selectedFlatRows.length === 0}
            onClick={() => setShowDeletePopup(true)}
          />
        </Row>
        <Table
          {...tableProps}
          rowHref={({ name }) => `/${params.project}/connections/${name}`}
          placeholder={
            <Text>
              There are no <strong>Connections</strong>.
            </Text>
          }
        />
      </Card>
      <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
        <Card
          border
          title={`Delete Connection${selectedFlatRows.length > 1 ? 's' : ''}`}
          size="large"
        >
          <Text>
            You are about to delete the{' '}
            <strong>
              {selectedFlatRows
                .map(({ original: { name } }) => name)
                .join(', ')}
            </strong>{' '}
            connection{selectedFlatRows.length > 1 ? 's' : ''}.
          </Text>
          <Button
            marginTop={5}
            title="Delete"
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </Layout>
  );
};

export default Connections;
