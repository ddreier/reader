const webpack = require('webpack');

module.exports = {
    entry: {
        main: "./javascripts/application.js",
        vendor: [
            "./javascripts/jquery-min.js",
            "./javascripts/bootstrap-min.js",
            "./javascripts/bootstrap.file-input.js",
            "./javascripts/mousetrap-min.js",
            "./javascripts/jquery-visible-min.js",
            "./javascripts/underscore-min.js",
            "./javascripts/backbone-min.js",
        ]
    },
    plugins: [
        new webpack.IgnorePlugin({resourceRegExp: /underscore/})
    ]
}