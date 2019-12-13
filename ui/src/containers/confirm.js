import React, { useEffect } from 'react';
import { useNavigation } from 'react-navi';
import { motion } from 'framer-motion';

import api from '../api';
import { Column, Heading } from '../components/core';

const Confirm = ({
  route: {
    data: { params },
  },
}) => {
  const navigation = useNavigation();
  useEffect(() => {
    api
      .completeRegistration({ registrationTokenValue: params.token })
      .then(() => navigation.navigate('/login'))
      .catch(console.log);
  }, []);

  return (
    <Column flex={1} alignItems="center" paddingTop={9}>
      <motion.div
        initial={{
          opacity: 0,
        }}
        animate={{
          opacity: 1,
        }}
        transition={{ delay: 1 }}
      >
        <Heading>Confirming registration</Heading>
      </motion.div>
    </Column>
  );
};

export default Confirm;
