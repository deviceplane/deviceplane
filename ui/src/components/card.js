import React from 'react';
import styled from 'styled-components';
import { variant } from 'styled-system';

import Logo from './icons/logo';
import { Column, Row, Text, Button, Link } from './core';

const Container = styled(Column)`
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
      full: {
        width: 'unset',
        alignSelf: 'stretch',
      },
    },
  })}
`;

const Card = ({
  size = 'large',
  title,
  subtitle,
  top = null,
  border = false,
  logo,
  actions = [],
  children,
  ...props
}) => {
  return (
    <Container
      bg="black"
      color="white"
      variant={size}
      borderRadius={2}
      padding={6}
      border={border ? 0 : undefined}
      borderColor="white"
      boxShadow={1}
      overflow="hidden"
      {...props}
    >
      {logo && (
        <Link href="https://deviceplane.com" marginX="auto" marginBottom={6}>
          <Logo size={50} />
        </Link>
      )}
      {top}
      {title && (
        <Row
          justifyContent="space-between"
          alignItems="flex-end"
          marginBottom={5}
          borderColor="white"
        >
          <Column>
            <Text fontSize={5} fontWeight={4}>
              {title}
            </Text>
            {subtitle && (
              <Text fontSize={1} fontWeight={2} color="grays.8" marginTop={1}>
                {subtitle}
              </Text>
            )}
          </Column>
          <Row marginLeft={7}>
            {actions.map(
              ({ href, variant = 'primary', title, onClick, show = true }) =>
                show && (
                  <Button
                    key={title}
                    title={title}
                    href={href}
                    variant={variant}
                    onClick={onClick}
                    marginLeft={5}
                  />
                )
            )}
          </Row>
        </Row>
      )}
      {children}
    </Container>
  );
};

export default Card;
