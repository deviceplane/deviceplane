import React, { useState, useEffect } from 'react';
import styled, { keyframes } from 'styled-components';

import utils from '../../utils';
import { Box, Row, Column } from './box';
import Button from './button';
import Text from './text';
import Icon from './icon';

export const ToastTypes = ['INFO', 'SUCCESS', 'DANGER'].reduce(
  (obj, type) => ({ ...obj, [type]: type }),
  {}
);

export const ToastEvent = 'TOAST';

const getColor = type => {
  switch (type) {
    case ToastTypes.DANGER:
      return 'red';
    case ToastTypes.SUCCESS:
      return 'green';
    case ToastTypes.INFO:
    default:
      return 'primary';
  }
};

const getIcon = type => {
  switch (type) {
    case ToastTypes.DANGER:
      return 'warning-sign';
    case ToastTypes.SUCCESS:
      return 'tick-circle';
    case ToastTypes.INFO:
    default:
      return 'info-sign';
  }
};

const barAnim = keyframes`
from {
    transform: scaleX(1);
} to {
    transform: scaleX(0);
}
`;

const toastAnim = keyframes`
0% {
    opacity: 0;
    transform: scale(.8) translateY(-140%);
} 
10%,90% {
    opacity: 1;
    transform: scale(1) translateY(0);
}
100% {
    opacity: 0;
    transform: scale(.8) translateY(-140%);
}
`;

const ProgressBar = styled(Box)`
  height: 3px;
  opacity: 0.5;
  width: 100%;

  transform-origin: 0 50%;

  animation: 3s ${barAnim} linear forwards 500ms;
`;

const Container = styled(Column)`
  animation: 4s ${toastAnim} ease forwards;
`;

export const Toast = ({ content, type = ToastTypes.INFO, close }) => {
  const color = getColor(type);

  useEffect(() => {
    setTimeout(close, 4000);
  }, []);

  return (
    <Container
      borderRadius={1}
      bg={color}
      minWidth="200px"
      marginBottom={3}
      zIndex={999}
      style={{ pointerEvents: 'all' }}
    >
      <Row
        paddingX={4}
        paddingY={2}
        alignItems="center"
        justifyContent="space-between"
      >
        <Row>
          <Icon marginRight={4} icon={getIcon(type)} size={14} color="black" />
          <Text color="black" fontSize={1} fontWeight={1}>
            {content}
          </Text>
        </Row>

        <Button
          variant="text"
          title={<Icon marginLeft={4} size={12} icon="cross" color="black" />}
          onClick={close}
        />
      </Row>
      <ProgressBar bg="black" />
    </Container>
  );
};

export const ToastManager = () => {
  const [toasts, setToasts] = useState([]);

  const addToast = ({ detail: toast }) => {
    setToasts(toasts => [...toasts, { ...toast }]);
  };

  const removeToast = id => () =>
    setToasts(toasts => toasts.filter(toast => toast.id !== id));

  useEffect(() => {
    window.addEventListener(ToastEvent, addToast);

    return () => window.removeEventListener(ToastEvent, addToast);
  }, []);

  return (
    <Column
      position="absolute"
      right={0}
      top={4}
      alignItems="center"
      width="100%"
      style={{ pointerEvents: 'none' }}
    >
      {toasts.map(({ id, ...toast }) => (
        <Toast key={id} {...toast} close={removeToast(id)} />
      ))}
    </Column>
  );
};

export const toaster = {
  success: (content, id = utils.id()) => {
    window.dispatchEvent(
      new CustomEvent(ToastEvent, {
        bubbles: true,
        detail: { content, type: ToastTypes.SUCCESS, id },
      })
    );
  },

  danger: (content, id = utils.id()) => {
    window.dispatchEvent(
      new CustomEvent(ToastEvent, {
        bubbles: true,
        detail: { content, type: ToastTypes.DANGER, id },
      })
    );
  },

  info: (content, id = utils.id()) => {
    window.dispatchEvent(
      new CustomEvent(ToastEvent, {
        bubbles: true,
        detail: { content, type: ToastTypes.INFO, id },
      })
    );
  },
};
