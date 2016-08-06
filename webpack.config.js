module.exports = function(options) {
  var path = require('path');
  var webpack = require('webpack');
  var config = {};

  config.context = __dirname;

  config.entry = {
    default: './app/entry.js'
  };

  config.output = {
    path: path.join(__dirname, "static"),
    filename: "bundle.js"
  }

  config.plugins = [
    new webpack.ProvidePlugin({
      'fetch': 'imports?this=>global!exports?global.fetch!whatwg-fetch'
    })
  ];

  config.module = {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader"
      },

      {
        test: /\.scss$/,
        loader: "style!css!postcss!sass"
      }
    ]
  }

  config.resolve = {
    root: path.resolve(__dirname, "./app")
  };

  return config;
};
