const path = require('path');

const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = {
	mode: 'production',
	entry: {
		login: './src/login.js',
	},
	output: {
		path: path.resolve(__dirname, 'dist'),
		publicPath: '/',
	},
	module: {
		rules: [{
			test: /\.scss$/,
			use: [
				{ loader: MiniCssExtractPlugin.loader },
				{ loader: 'css-loader' },
				{ loader: 'postcss-loader',
					options: {
						postcssOptions: {
							plugins: ['autoprefixer'],
						}
					}
				}
			]
		}]
	},
	plugins: [
		new MiniCssExtractPlugin()
	],
	optimization: {
		minimize: true,
		minimizer: [`...`, new CssMinimizerPlugin()]
	}
}
