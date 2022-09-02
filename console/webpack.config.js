/**
 * Created by yanggang on 2017/3/6.
 */
const path = require("path");
const webpack = require("webpack");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const UglifyJsPlugin = require("uglifyjs-webpack-plugin");
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin");

const env = process.env.NODE_ENV || 'development';

module.exports = {
    entry: {
        console: './main.js'
    },
    output: {
        path: path.join(__dirname, "../static/dist"),
        filename: "[name].js"
    },
    externals: {
        config: function () {
            return JSON.stringify(require('./config/' + env + '.json'));
        }()
    },
    optimization: {
        splitChunks: {
            chunks: 'all',
            automaticNameDelimiter: '-',
            minSize: 2000,
            cacheGroups: {
                react: {
                    test: (module, chunks) => /react/.test(module.context),
                    minChunks: 1,
                    name: "react",
                    priority: -10,
                },
                antd: {
                    test: (module, chunks) => /antd|rc-|async-validator/.test(module.context),
                    minChunks: 1,
                    name: "antd",
                    priority: -20,
                },
                vendors: {
                    test: /node_modules/,
                    minChunks: 1,
                    name: "vendors",
                    priority: -30,
                },
                default: {
                    minChunks: 1,
                    priority: -40,
                    chunks: 'all'
                }
            }
        },
        minimizer: [
            new UglifyJsPlugin({
                cache: true,
                parallel: true,
                sourceMap: true // set to true if you want JS source maps
            }),
            new OptimizeCSSAssetsPlugin({})
        ]
    },
    plugins: [
        new HtmlWebpackPlugin({
            title: "Riff Console",
            hash: true,
            inject: false,
            filename: "console.html",
            template: "index.html"
        }),
        new MiniCssExtractPlugin({
            filename: "[name].css",
            disable: false,
            allChunks: true
        }),
        new webpack.DefinePlugin({
            "process.env": { NODE_ENV: JSON.stringify(env) }
        })
    ],
    module: {
        rules: [{
            test: /\.flow/, loader: 'ignore-loader'
        }, {
            test: /\.js?$/,
            exclude: /node_modules/,
            use: {
                loader: "babel-loader",
                options: {
                    presets: ["env", "react", "stage-0"],
                    plugins: [
                        [
                            "import",
                            {
                                libraryName: "antd",
                                style: true
                            }
                        ]
                    ]
                }
            }
        }, {
            test: /.css$/,
            use: [
                MiniCssExtractPlugin.loader,
                "css-loader",
                {
                    loader: "postcss-loader",
                    options: {
                        ident: 'postcss',
                        plugins: (loader) => [
                            require('postcss-import')({ root: loader.resourcePath }),
                            require('postcss-cssnext')()
                        ]
                    }
                }
            ]
        }, {
            test: /.less$/,
            use: [
                MiniCssExtractPlugin.loader,
                "css-loader",
                {
                    loader: "postcss-loader",
                    options: {
                        ident: 'postcss',
                        plugins: (loader) => [
                            require('postcss-import')({ root: loader.resourcePath }),
                            require('postcss-cssnext')()
                        ]
                    }
                },
                {
                    loader: "less-loader"
                }
            ]
        }]
    }
};