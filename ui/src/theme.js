const breakpoints = ['600px', '800px', '1000px'];

const space = [
  0,
  4,
  8,
  12,
  16,
  24,
  32,
  48,
  64,
  96,
  128,
  256,
  384,
  448,
  512,
  768,
  1024,
];

const colors = {
  primary: '#6fccff',

  black: '#000',
  white: '#E6E6E6',
  green: '#A2FAA3',
  red: '#F93943',

  pureWhite: '#fff',

  transparent: 'transparent',

  overlay: 'rgba(24, 24, 24, .8)',

  pageBackground: '#181818',

  grays: [
    '#121212',
    '#141414',
    '#181818',
    '#222222',
    '#363636',
    '#484848',
    '#5A5A5A',
    '#6C6C6C',
    '#7E7E7E',
    '#909090',
    '#A2A2A2',
    '#B4B4B4',
    '#C6C6C6',
    '#D8D8D8',
    '#EAEAEA',
  ],
};

export default {
  colors,

  space,

  sizes: space,

  fonts: {
    default: `Rubik,Roboto,sans-serif`,
    code: 'source-code-pro,Menlo,Monaco,Consolas,monospace',
  },

  fontSizes: [12, 14, 16, 18, 24, 36, 48, 64, 72],

  fontWeights: [300, 400, 500, 700],

  radii: [0, 4, 6, 8, 12, 16, 9999, '100%'],

  breakpoints,

  borders: [`1px solid ${colors.primary}`, `3px solid ${colors.primary}`],

  mediaQueries: {
    small: `@media screen and (min-width: ${breakpoints[0]})`,
    medium: `@media screen and (min-width: ${breakpoints[1]})`,
    large: `@media screen and (min-width: ${breakpoints[2]})`,
  },

  transitions: ['all 200ms ease'],
};

export const labelColors = [
  '#75c800',
  '#599900',
  '#c5b500',
  '#978b00',
  '#1fcf0f',
  '#0d9f00',
  '#e4a679',
  '#d36d24',
  '#00c9b8',
  '#009a8d',
  '#ec9ea5',
  '#df5f6a',
  '#eb9ac9',
  '#dd57a5',
  '#7dbae5',
  '#2d8fd5',
  '#e498ea',
  '#d153dd',
  '#c4a6ed',
  '#a172e3',
  '#a9aeee',
  '#7780e4',
];
