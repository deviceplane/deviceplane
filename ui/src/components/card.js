import React from 'react';
import styled from 'styled-components';
import { variant } from 'styled-system';

import Logo from './icons/logo';
import { Column, Row, Text, Button, Link } from './core';

const Container = styled(Column)`
  pointer-events: ${props => (props.disabled ? 'none' : 'auto')};
  opacity: ${props => (props.disabled ? 0.2 : 1)};
  overflow: ${props =>
    props.overflow ? `${props.overflow} !important` : 'initial'};

  ${variant({
    variants: {
      small: {
        width: 9,
      },
      medium: {
        width: 11,
      },
      large: {
        width: 12,
      },
      xlarge: {
        width: 13,
      },
      xxlarge: {
        width: 14,
      },
      full: {
        width: 'unset',
        alignSelf: 'stretch',
      },
    },
  })}

  @media (max-width: 600px) {
    width: 100%;
    border-radius: 0;
  }
`;

const Card = ({
  size,
  title,
  subtitle,
  top = null,
  border = false,
  logo,
  actions = [],
  header,
  children,
  disabled,
  ...props
}) => {
  return (
    <Container
      bg="black"
      color="white"
      variant={size}
      borderRadius={2}
      padding={6}
      border={border ? 1 : undefined}
      borderColor="white"
      disabled={disabled}
      {...props}
    >
      {logo && (
        <Link href="https://deviceplane.com" marginX="auto" marginBottom={6}>
          <Logo size={50} />
        </Link>
      )}
      {top}
      {title && (
        <Column marginBottom={5} borderColor="white">
          <Row alignItems="center" justifyContent="space-between">
            <Column>
              <Row>
                <Text fontSize={5} fontWeight={2}>
                  {title}
                </Text>
              </Row>
            </Column>
            <Row marginLeft={7}>
              {actions.map(
                ({
                  href,
                  variant = 'primary',
                  title,
                  onClick,
                  disabled,
                  show = true,
                }) =>
                  show && (
                    <Button
                      key={title}
                      title={title}
                      href={href}
                      variant={variant}
                      onClick={onClick}
                      disabled={disabled}
                      marginLeft={5}
                    />
                  )
              )}
              {header}
            </Row>
          </Row>
          {subtitle && (
            <Row marginTop={1}>
              {typeof subtitle === 'string' ? (
                <Text fontSize={1} fontWeight={1} color="grays.8" marginTop={1}>
                  {subtitle}
                </Text>
              ) : (
                subtitle
              )}
            </Row>
          )}
        </Column>
      )}
      {children}
    </Container>
  );
};

export default Card;
