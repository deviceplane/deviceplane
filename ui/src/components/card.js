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
  transition: all 250ms;
  width: 100%;

  ${variant({
    variants: {
      small: {
        maxWidth: 11,
      },
      medium: {
        maxWidth: 13,
      },
      large: {
        maxWidth: 14,
      },
      xlarge: {
        maxWidth: 15,
      },
      xxlarge: {
        maxWidth: 16,
      },
      full: {
        alignSelf: 'stretch',
      },
    },
  })}

  @media (max-width: 600px) {
    border-radius: 0;
  }
`;

const Card = ({
  size,
  title,
  subtitle,
  top = null,
  left = null,
  center = null,
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
      padding={7}
      border={border ? 1 : undefined}
      borderColor="white"
      disabled={disabled}
      {...props}
    >
      {logo && (
        <Link href="https://deviceplane.com" marginX="auto" marginBottom={8}>
          <Logo size={50} />
        </Link>
      )}
      {top}
      {title && (
        <Column marginBottom={6} borderColor="white">
          <Row
            alignItems="center"
            justifyContent="space-between"
            flexWrap="wrap"
          >
            <Row>
              <Text fontSize={5} fontWeight={2} marginRight={6}>
                {title}
              </Text>
              {left}
            </Row>
            {center && <Row marginX={6}>{center}</Row>}
            <Row>
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
                      marginLeft={4}
                    />
                  )
              )}
              {header}
            </Row>
          </Row>
          {subtitle && (
            <Row marginTop={1}>
              {typeof subtitle === 'string' ? (
                <Text fontSize={2} fontWeight={1} color="grays.8" marginTop={1}>
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
