import React, { useState, useRef, useEffect } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';

import { Column } from './core';

const Container = styled(Column)`
  position: relative;
  z-index: 999;
`;

const Content = styled(Column)`
  position: absolute;

  top: 48px;
  right: 0;
  width: 240px;
`;

const Button = styled.button`
  appearance: none;
  background: none;
  padding: 0;
  margin: 0;
  border: none;
  outline: none;

  cursor: pointer;
`;

const Popover = ({ children, content = null }) => {
  const node = useRef();
  const [show, setShow] = useState();

  const handleClick = e => {
    if (!node.current.contains(e.target)) {
      if (show) {
        setShow(false);
      }
    }
  };

  useEffect(() => {
    if (show) {
      document.addEventListener('mousedown', handleClick);
    } else {
      document.removeEventListener('mousedown', handleClick);
    }

    return () => {
      document.removeEventListener('mousedown', handleClick);
    };
  }, [show]);

  return (
    <Container ref={node}>
      <Button
        onClick={() => {
          setShow(!show);
        }}
      >
        {children}
      </Button>
      <AnimatePresence>
        {show && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.15 }}
          >
            <Content bg="black" borderRadius={1} border={0} borderColor="white">
              {typeof content === 'function'
                ? content({ close: () => setShow(false) })
                : content}
            </Content>
          </motion.div>
        )}
      </AnimatePresence>
    </Container>
  );
};

export default Popover;
