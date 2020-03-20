import React from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import api from '../../api';
import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Popup from '../../components/popup';
import Field from '../../components/field';
import { Text, Button, Form, Label, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  roles: yup.object(),
});

const Member = ({
  route: {
    data: { params, member, roles },
  },
}) => {
  const { control, handleSubmit, formState, errors } = useForm({
    validationSchema,
    defaultValues: {
      roles: roles.reduce(
        (obj, role) => ({
          ...obj,
          [role.name]: !!member.roles.find(({ name }) => name === role.name),
        }),
        {}
      ),
    },
  });
  const [isRemovePopup, toggleRemovePopup] = useToggle();
  const navigation = useNavigation();

  const removeMember = async () => {
    try {
      await api.removeMember({
        projectId: params.project,
        userId: member.userId,
      });
      toaster.success('Member removed.');
      navigation.navigate(`/${params.project}/iam/members`);
    } catch (error) {
      toaster.danger('Member was not removed.');
      console.error(error);
    }
  };

  const submit = async data => {
    let error = false;
    const roleArray = Object.keys(data.roles);
    for (let i = 0; i < roleArray.length; i++) {
      const role = roleArray[i];
      const choseRole = data.roles[role];
      const hasRole = member.roles.find(({ name }) => name === role);
      const roleId = roles.find(({ name }) => name === role).id;
      if (choseRole && !hasRole) {
        try {
          await api.addMembershipRoleBindings({
            projectId: params.project,
            userId: member.userId,
            roleId,
          });
        } catch (e) {
          error = true;
          console.error(e);
        }
      } else if (!choseRole && hasRole) {
        try {
          await api.removeMembershipRoleBindings({
            projectId: params.project,
            userId: member.userId,
            roleId,
          });
        } catch (e) {
          error = true;
          console.error(e);
        }
      }
    }

    if (error) {
      toaster.danger(
        'Roles for the member were not updated properly. Please check the roles of the member.'
      );
    } else {
      navigation.navigate(`/${params.project}/iam/members`);
      toaster.success('Member updated.');
    }
  };

  return (
    <>
      <Card
        title={member.user.name}
        subtitle={member.user.email}
        size="large"
        actions={[
          {
            title: 'Remove',
            onClick: toggleRemovePopup,
            variant: 'danger',
          },
        ]}
      >
        <Form onSubmit={handleSubmit(submit)}>
          <Label>Choose Individual Roles</Label>
          {roles.map(role => (
            <Field
              multi
              type="checkbox"
              key={role.id}
              name={`roles[${role.name}]`}
              label={role.name}
              control={control}
              errors={errors.roles && errors.roles[role.name]}
            />
          ))}
          <Button
            marginTop={3}
            title="Update"
            type="submit"
            disabled={!formState.dirty}
          />
        </Form>
      </Card>
      <Popup
        show={isRemovePopup}
        title="Remove Member"
        onClose={toggleRemovePopup}
      >
        <Card title="Remove Member" border size="large">
          <Text>
            You are about to remove the member (
            <strong>{member.user.name}</strong>) from the project.
          </Text>
          <Button
            marginTop={5}
            title="Remove"
            onClick={removeMember}
            variant="danger"
          />
        </Card>
      </Popup>
    </>
  );
};

export default Member;
