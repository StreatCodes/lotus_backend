module.exports = {
	entry: {
		login: "./src/login.ts",
		pages: "./src/pages.ts"
	},
	module: {
		rules: [
			{
				test: /\.ts$/,
				exclude: /node_modules/,
				use: [
					{
						loader: "babel-loader",
						options: {
							presets: ["@babel/preset-typescript"],
							plugins: [
								["module:nanohtml",
								{
									useImport: true
								}]
							]
						}
					}
				]
			}
		]
	},
	resolve: {
		extensions: [".tsx", ".ts", ".js"]
	}
};
