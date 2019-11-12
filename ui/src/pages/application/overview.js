import React from 'react';
import { Pane, majorScale, Heading, Card, Label, Button } from 'evergreen-ui';

import Releases from './releases';
import Editor from '../../components/Editor';
import InnerCard from '../../components/InnerCard';
import { DevicesFilterButtons } from '../../components/DevicesFilterButtons';
import { Link } from 'react-router-dom';

const ApplicationOverview = ({
  application: { latestRelease, name, schedulingRule },
  history,
  projectName
}) => {
  const currentConfig = latestRelease ? latestRelease.rawConfig : '';
  var schedulingRuleVisualized;
  if (schedulingRule.length) {
    schedulingRuleVisualized = (
      <DevicesFilterButtons query={schedulingRule} canRemoveFilter={false} />
    );
  } else {
    schedulingRuleVisualized = (<Label>No scheduling rule set. You can set one in the <Link to='./scheduling'>scheduling</Link> page.</Label>);
  }

  return (
    <Pane width="70%" paddingBottom={majorScale(4)}>
      <Heading
        paddingTop={majorScale(4)}
        paddingBottom={majorScale(1)}
        size={600}
      >
        {name}
      </Heading>
      <InnerCard>
        <Heading paddingTop={majorScale(2)} paddingLeft={majorScale(2)}>
          Scheduling Rule
        </Heading>
        <Card
          display="flex"
          flexDirection="column"
          alignItems="left"
          width="80%"
          padding={majorScale(2)}
        >
          {schedulingRuleVisualized}
        </Card>
      </InnerCard>
      <InnerCard>
        <Heading paddingTop={majorScale(2)} paddingLeft={majorScale(2)}>
          Current Config
        </Heading>
        <Card
          display="flex"
          flexDirection="column"
          alignItems="center"
          width="80%"
          padding={majorScale(2)}
        >
          <Editor width="100%" height="300px" value={currentConfig} readOnly />
        </Card>
      </InnerCard>
      <InnerCard>
        <Pane
          display="flex"
          flexDirection="row"
          alignItems="center"
          justifyContent="space-between"
        >
          <Heading padding={majorScale(2)}>Releases</Heading>
          <Button
            marginRight={majorScale(2)}
            appearance="primary"
            onClick={() =>
              history.push(`/${projectName}/applications/${name}/deploy`)
            }
          >
            Create New Release
          </Button>
        </Pane>
        <Releases
          projectName={projectName}
          applicationName={name}
          history={history}
        />
      </InnerCard>
    </Pane>
  );
};

export default ApplicationOverview;
