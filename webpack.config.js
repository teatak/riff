/**
 * Created by yanggang on 2017/3/6.
 */
var path = require("path");
var webpack = require('webpack');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var HtmlWebpackPlugin = require('html-webpack-plugin');
var autoprefixer = require('autoprefixer');
var precss = require('precss');

const env = process.env.NODE_ENV || 'development';

module.exports = {
    entry: {
        console: './console/main.js'
    },
    output: {
        path: path.join(__dirname, "static/dist"),
        filename: "[name].js"
    },
    externals: {
        'config': function(){
            return JSON.stringify(require('./console/config/' + env + '.json'));
        }()
    },
    plugins: function(){
        var options = [
            new HtmlWebpackPlugin({title:'Riff Console',hash:true,inject:false, filename: 'console.html', template: 'console/index.html'}),
            new ExtractTextPlugin({filename: "[name].css", disable: false, allChunks: true}),
            new webpack.DefinePlugin({'process.env': {NODE_ENV: JSON.stringify(env)}}),
        ];
        if(env === 'production' || env === 'testing') {
            options.push(new webpack.optimize.UglifyJsPlugin());
            options.push(new webpack.LoaderOptionsPlugin({minimize: true}));
        }
        return options;
    }(),
    module: {
        rules: [{
            test: /\.js?$/,
            exclude: /node_modules/,
            use: {
                loader: 'babel-loader',
                options: {
                    presets: ['env', 'react', 'stage-0'],
                    plugins: [["import", {
                        "libraryName": "antd",
                        "style": true
                    }]]
                }
            }
        },{
            test: /.css$/,
            use: ExtractTextPlugin.extract({
                fallback: "style-loader",
                use: [
                    "css-loader",
                    {
                        loader: 'postcss-loader',
                        options: {
                            plugins: function () {
                                return [precss,autoprefixer]
                            }
                        }
                    }
                ]
            })
        },{
            test: /.less$/,
            use: ExtractTextPlugin.extract({
                fallback: "style-loader",
                use: [
                    "css-loader",
                    {
                        loader: 'postcss-loader', options: {
                            plugins: function () {
                                return [precss,autoprefixer]
                            }
                        }
                    },
                    {
                        loader: 'less-loader', options: {
                            //modifyVars: themes()
                        }
                    }
                ]
            })
        },{
            test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
            loader: 'file-loader',
            options: {
                name: 'fonts/[name].[ext]',
            }
        },{
            test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: 'file-loader',
            options: {
                name: 'fonts/[name].[ext]',
            }
        },{
            test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: 'file-loader',
            options: {
                name: 'fonts/[name].[ext]',
            }
        },{
            test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, loader: 'file-loader',
            options: {
                name: 'fonts/[name].[ext]',
            }
        }
        ]
    },
};