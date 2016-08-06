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

  config.module = {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader"
      },

      {
        test: /\.css$/,
        loader: "style!css"
      }
    ]
  }

  config.resolve = {
    root: path.resolve(__dirname, "./app")
  };

  return config;
};
