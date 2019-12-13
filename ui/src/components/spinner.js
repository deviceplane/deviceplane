import React from 'react';
import SquareLoader from 'react-spinners/SquareLoader';
import styled from 'styled-components';
import { motion } from 'framer-motion';

import { Column } from './core';

const Container = styled(Column)`
  position: absolute;
  z-index: 99999;
  top: 0;
  left: 0;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.9);

  & > div {
    border-radius: 3px;
  }
`;

const Spinner = ({ fadeIn = 'full', name = 'cube-grid', show = false }) => {
  if (!show) {
    return null;
  }
  return (
    <motion.div
      initial={{
        opacity: 0,
      }}
      animate={{
        opacity: 1,
      }}
      transition={{ delay: 1 }}
    >
      <Container>
        <SquareLoader color="white" size={60} />
      </Container>
    </motion.div>
  );
};

export default Spinner;
