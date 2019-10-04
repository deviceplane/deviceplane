import React from 'react';
import logo from '../assets/logo.png';

const Logo = ({ width = '44px', height = '40px' }) => (
  <img src={logo} width={width} height={height} alt="Logo" />
);

export default Logo;
