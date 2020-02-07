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
  flex: 0 0 auto;
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
    },
  })}
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
      padding={[0, 0, 6]}
      border={border ? 0 : undefined}
      borderColor="white"
      disabled={disabled}
      overflow="hidden"
      {...props}
    >
      {title && (
        <Column marginBottom={6} borderColor="white">
          {logo && (
            <Link
              href="https://deviceplane.com"
              marginX="auto"
              marginBottom={6}
            >
              <Logo size={50} />
            </Link>
          )}
          {top}
          <Row
            flexWrap="wrap"
            justifyContent="space-between"
            alignItems="flex-end"
            marginLeft={-4}
          >
            <Row alignItems="center" marginLeft={4}>
              <Text fontSize="28px" fontWeight={2} marginRight={4}>
                {title}
              </Text>
              {left}
            </Row>
            <Row justifyContent="center" marginTop={4} marginLeft={4}>
              {center}
            </Row>
            <Row marginLeft={4} marginTop={4}>
              {actions.map(
                (
                  {
                    href,
                    variant = 'primary',
                    title,
                    onClick,
                    disabled,
                    show = true,
                  },
                  index
                ) =>
                  show && (
                    <Button
                      key={title}
                      title={title}
                      href={href}
                      variant={variant}
                      onClick={onClick}
                      disabled={disabled}
                      marginLeft={index === 0 ? 0 : 4}
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
