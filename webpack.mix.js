const mix = require('laravel-mix');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const Dotenv = require('dotenv-webpack')

const {
    MIX_BROWSERSYNC_PROXY = 'http://localhost:8040',
} = process.env;



mix.setPublicPath("web/build")

mix.js('web/src/js/app.js', 'web/build')
	.vue({
		extractStyles: false,
		globalStyles: false,
	})

mix.sass('web/src/scss/styles.scss', 'web/build')
    .options({
        processCssUrls: false,
        autoprefixer: false,
    });



// Configuration
mix.version();
mix.sourceMaps(false, 'inline-source-map');

mix.browserSync({
    notify: true,
    proxy: MIX_BROWSERSYNC_PROXY,
    files: [
        'web/build/',
    ]
});

mix.webpackConfig({
    plugins: [
        new CopyWebpackPlugin({
            patterns: [
                {
                    from: 'src/index.html',
                    to: 'index.html',
                    context: "web",
                    globOptions: {
                        ignore: ['.DS_Store']
                    }
                },
            ],
        }),
        new Dotenv({
            systemvars: true
        }),
    ]
});
