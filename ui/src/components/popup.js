import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';
import { Icon } from 'evergreen-ui';
import { motion, AnimatePresence } from 'framer-motion';

import { Column } from './core';

const Overlay = styled(Column)`
  position: fixed;
  z-index: 9999999;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background-color: ${props => props.theme.colors.overlay};
`;

const Container = styled(Column)`
  position: relative;
  z-index: 9999999;
`;

const Content = styled(Column)`
  max-height: 80vh;
  overflow: ${props => props.overflow};
  & > div {
    overflow: scroll;
  }
`;

const CloseButton = styled.button`
  display: flex;
  appearance: none;
  border: none;
  outline: none;
  margin: 0;
  padding: 0;
  position: absolute;
  top: 0px;
  left: -64px;
  padding: 4px;
  border-radius: 999px;
  z-index: 9999999;
  cursor: pointer;
  border: 2px solid ${props => props.theme.colors.white};

  transition: background-color 200ms;
  background-color: ${props => props.theme.colors.white};

  &:hover {
    background-color: ${props => props.theme.colors.black};
  }

  & svg {
    transition: fill 200ms;
  }

  &:hover svg {
    fill: ${props => props.theme.colors.white} !important;
  }
`;

const Popup = ({ children, show, onClose, overflow = 'hidden' }) => {
  const node = useRef();

  const handleClick = e => {
    if (!node.current.contains(e.target)) {
      onClose();
    }
  };

  const handleKeyDown = e => {
    if (e.key === 'Escape') {
      onClose();
    }
  };

  useEffect(() => {
    if (show) {
      document.addEventListener('keydown', handleKeyDown);
      document.addEventListener('mousedown', handleClick);
      document.body.style.overflow = 'hidden';
    } else {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('mousedown', handleClick);
      document.body.style.overflow = 'initial';
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('mousedown', handleClick);
      document.body.style.overflow = 'initial';
    };
  }, [show]);

  return (
    <AnimatePresence>
      {show && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.1 }}
        >
          <Overlay>
            <Container ref={node}>
              <CloseButton onClick={onClose}>
                <Icon icon="cross" size={20} color="black" />
              </CloseButton>
              <motion.div
                initial={{ opacity: 0, y: 150, scale: 0.75 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                exit={{ opacity: 0, y: -150, scale: 0.5 }}
                transition={{ duration: 0.2, delay: 0.1 }}
              >
                <Content overflow={overflow}>{children}</Content>
              </motion.div>
            </Container>
          </Overlay>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default Popup;
