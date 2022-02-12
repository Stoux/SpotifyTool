const mix = require('laravel-mix');
const Dotenv = require('dotenv-webpack')
const fs = require('fs')

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

mix.copy('web/src/index.html', 'web/build/index.html')

// Configuration
mix.version();
mix.sourceMaps(false, 'inline-source-map');

if (MIX_BROWSERSYNC_PROXY) {
    mix.browserSync({
        notify: true,
        proxy: MIX_BROWSERSYNC_PROXY,
        files: [
            'web/build/',
        ]
    });
}

mix.webpackConfig({
    plugins: [
        new Dotenv({
            systemvars: true
        }),
    ]
});

// Replace the asset URLs in the index.html with their versioned URL
mix.then(() => {
    const manifestPath = 'web/build/mix-manifest.json';
    const htmlPath = 'web/build/index.html';

    if (fs.existsSync(manifestPath) && fs.existsSync(htmlPath)) {
        // Read the Manifest file & load the index.html
        const manifest = JSON.parse(fs.readFileSync(manifestPath, {encoding: 'utf8'}));
        let htmlContents = fs.readFileSync(htmlPath, { encoding: 'utf8'});
        Object.keys(manifest).forEach(key => {
            htmlContents = htmlContents.replace(`"${key}"`, `"${manifest[key]}"`);
        })
        fs.writeFileSync(htmlPath, htmlContents, { encoding: "utf8"})
    }
})