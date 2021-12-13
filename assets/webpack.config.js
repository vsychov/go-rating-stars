const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const HtmlMinimizerPlugin = require("html-minimizer-webpack-plugin");

module.exports = {
    mode: 'production',
    entry: './src/main.js',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: '[name].[contenthash].js',
        clean: true,
    },
    plugins: [new HtmlWebpackPlugin({
        publicPath: '/assets/assets/dist/'
    })],
    optimization: {
        minimize: true,
        minimizer: [
            `...`,
            new HtmlMinimizerPlugin(),
        ],
    },
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ["style-loader", "css-loader"],
            },
            {
                test: /\.svg/,
                type: 'asset/inline'
            }
        ]
    }
};