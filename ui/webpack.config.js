const webpack = require('webpack');
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CompressionPlugin = require('compression-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer')
  .BundleAnalyzerPlugin;

const config = {
  entry: './src/index.js',
  output: {
    publicPath: '/',
    path: path.resolve(__dirname, './dist'),
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: 'index.html',
      favicon: 'assets/favicon.ico',
    }),
    // new webpack.IgnorePlugin({
    //   resourceRegExp: /^\.\/locale$/,
    //   contextRegExp: /moment$/,
    // }),
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
  ],
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
        },
      },
      {
        test: /\.css$/,
        use: [
          'style-loader',
          { loader: 'css-loader', options: { importLoaders: 1 } },
          {
            loader: 'postcss-loader',
            options: {
              plugins: () => [require('autoprefixer')],
            },
          },
        ],
      },
      {
        test: /\.js\.map$/,
        use: {
          loader: 'file-loader',
        },
      },
    ],
  },
  node: {
    module: 'empty',
    dgram: 'empty',
    dns: 'mock',
    fs: 'empty',
    http2: 'empty',
    net: 'empty',
    tls: 'empty',
    child_process: 'empty',
  },
};

module.exports = () => {
  const production = process.env.NODE_ENV === 'production';

  if (production) {
    config.mode = 'production';
    config.devtool = 'source-map';
    config.output.filename = 'static/[name].[contenthash].js';
    config.optimization = {
      minimize: true,
      minimizer: [new TerserPlugin()],
      splitChunks: {
        chunks: 'all',
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendor',
            priority: -10,
            enforce: true,
          },
        },
      },
    };
    config.plugins = [
      ...config.plugins,
      new CleanWebpackPlugin(),
      new CompressionPlugin({
        algorithm: 'brotliCompress',
        test: /\.(js|css)$/,
      }),
      new webpack.optimize.LimitChunkCountPlugin({
        maxChunks: 3,
      }),
      //new BundleAnalyzerPlugin(),
    ];
  } else {
    config.mode = 'development';
    config.devtool = 'cheap-module-source-map';
    config.output.filename = 'bundle.js';
    config.devServer = {
      port: 3000,
      historyApiFallback: true,
      overlay: true,
      open: true,
      publicPath: '/',
      contentBase: './dist',
      hot: true,
      transportMode: 'ws',
      clientLogLevel: 'none',
      noInfo: true,
      proxy: {
        '/api': 'http://localhost:8080',
      },
    };
  }

  return config;
};
