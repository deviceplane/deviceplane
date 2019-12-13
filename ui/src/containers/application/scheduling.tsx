// @ts-nocheck

import React, { useState } from 'react';
import {
  toaster,
  // @ts-ignore
} from 'evergreen-ui';
import { useNavigation } from 'react-navi';

import api from '../../api';
import utils from '../../utils';
import { DevicesFilterButtons } from '../../components/DevicesFilterButtons';
import {
  Query,
  Filter,
  DevicesFilter,
  LabelValueCondition,
} from '../../components/DevicesFilter';
import Card from '../../components/card';
import Alert from '../../components/alert';
import { Button, Row, Text } from '../../components/core';

interface Props {
  application: any;
  projectName: string;
  history: any;
}

interface State {
  schedulingRule: Query;
  backendError: any;
  showFilterDialog: boolean;
}

const Scheduling = ({
  route: {
    data: { params, application },
  },
}) => {
  const [schedulingRule, setSchedulingRule] = useState(
    application.schedulingRule
  );
  const [backendError, setBackendError] = useState();
  const [showFilterDialog, setShowFilterDialog] = useState();
  const navigation = useNavigation();

  const submit = async () => {
    setBackendError(null);

    try {
      await api.updateApplication({
        projectId: params.project,
        applicationId: application.name,
        data: {
          name: application.name,
          description: application.description,
          schedulingRule,
        },
      });

      toaster.success('Scheduling rule updated successfully.');
      navigation.navigate(
        `/${params.project}/applications/${application.name}`
      );
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Scheduling rule was not updated.');
      console.log(error);
    }
  };

  const filterDevices = () => {
    // TODO: fetch devices and show them
  };

  const removeFilter = (index: number) => {
    setSchedulingRule(schedulingRule.filter((_, i) => i !== index));
    filterDevices();
  };

  const addFilter = (filter: Filter) => {
    setShowFilterDialog(false);
    setSchedulingRule([...schedulingRule, filter]);
    filterDevices();
  };

  const clearFilters = () => {
    setSchedulingRule([]);
    filterDevices();
  };

  return (
    <Card
      title="Scheduling"
      size="xlarge"
      actions={[
        {
          title: 'Clear Filters',
          onClick: clearFilters,
          show: !!schedulingRule.length,
          variant: 'text',
        },
        {
          title: 'Add Filter',
          onClick: () => setShowFilterDialog(true),
          variant: 'secondary',
        },
      ]}
    >
      <Alert show={backendError} variant="error" description={backendError} />
      <Row bg="grays.0" borderRadius={1} minHeight={7} padding={2}>
        {schedulingRule.length ? (
          <DevicesFilterButtons
            canRemoveFilter
            query={schedulingRule}
            removeFilter={removeFilter}
          />
        ) : (
          <Row flex={1} justifyContent="center" alignItems="center">
            <Text>
              There are no <strong>Filters</strong>.
            </Text>
          </Row>
        )}
      </Row>
      <DevicesFilter
        whitelistedConditions={[LabelValueCondition]}
        show={showFilterDialog}
        onClose={() => setShowFilterDialog(false)}
        onSubmit={addFilter}
      />
      <Button
        title="Set Scheduling Rule"
        marginTop={6}
        onClick={submit}
        disabled={utils.deepEqual(application.schedulingRule, schedulingRule)}
      />
    </Card>
  );
};

export default Scheduling;
