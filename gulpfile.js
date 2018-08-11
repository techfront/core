var autoprefixer = require('autoprefixer'),
    gulp = require('gulp'),
    mainBowerFiles = require('gulp-main-bower-files'),
    watch = require('gulp-watch'),
    sass = require('gulp-sass'),
    lost = require('lost'),
    pxtorem = require('postcss-pxtorem'),
    favicons = require("gulp-favicons"),
    filter = require('gulp-filter'),
    concat = require('gulp-concat'),
    postcss = require('gulp-postcss'),
    rename = require('gulp-rename'),
    stylelint = require('stylelint'),
    reporter = require('postcss-reporter'),
    ts = require('gulp-typescript'),
    tslint = require('gulp-tslint'),
    uglify = require('gulp-uglify'),
    stylefmt = require('gulp-stylefmt');

var plugins = [
    lost(),
    autoprefixer(),
    pxtorem({
        rootValue: 13,
        replace: true,
        propList: ['font', 'font-size', 'line-height', 'letter-spacing', 'margin', 'padding', 'height', 'width', 'border-radius'],
    })
];

var typings = 'typings/index.d.ts';

var source = 'src';
var output = 'tmp/src-bundle';
var publicOutput = 'public/assets';


gulp.task('build:styles-format', function () {
    return gulp.src(source + '/**/*.scss')
        .pipe(stylefmt())
        .pipe(gulp.dest(publicOutput + '/styles-fmt/'));
});

// Build Styles
gulp.task('build:styles', function () {
    return gulp.src([source + '/assets/**/init.scss', source + '/component/**/*.scss'])
        .pipe(sass().on('error', sass.logError))
        .pipe(postcss(plugins))
        .pipe(rename({dirname: ''}))
        .pipe(concat('app.css'))
        .pipe(gulp.dest(output + '/styles'));
});

// Build TS
gulp.task('build:ts', function () {
    var tsProject = ts.createProject('tsconfig.json');
    var tsResult = gulp.src([
        typings,
        source + '/assets/scripts/app.ts',
        source + '/assets/scripts/dropdown.ts',
        source + '/assets/scripts/dom.ts',
        source + '/assets/scripts/widget/*.ts',
        source + '/component/**/*.ts',
        source + '/assets/scripts/main.ts'
    ])
    .pipe(tslint({formatter: "verbose"}))
    .pipe(tslint.report())
    .pipe(tsProject());

    var srcFilter = filter(['**/*.js', '!**/dom.js'], { restore: true });
    var libFilter = filter('**/dom.js', { restore: true });

    return tsResult
        .pipe(srcFilter)
        .pipe(concat('app.js'))
        .pipe(gulp.dest(output + '/scripts'))
        .pipe(srcFilter.restore)

        .pipe(libFilter)
        .pipe(rename(function (path) {
            path.extname = ".min.js";
            path.dirname = "";
        }))
        .pipe(uglify())
        .pipe(gulp.dest(publicOutput + '/scripts'));
});

gulp.task('build:vendor-js', function () {
    return gulp.src('./bower.json')
        .pipe(mainBowerFiles())
        .pipe(filter('**/*.js', {restore: true}))
        .pipe(rename(function (path) {
            path.basename = path.dirname.split('/')[0];
            path.extname = ".min.js";
            path.dirname = "";
        }))
        .pipe(uglify())
        .pipe(gulp.dest(publicOutput + '/scripts'));
});

gulp.task('build:favicons', function () {
    return gulp.src('public/images/logos/favicon.png').pipe(favicons({
        appName: "TechFront",
        path: "/images/favicons/",
        url: "https://techfront.ru/",
        display: "standalone",
        orientation: "portrait",
        background: "#3498db",
        start_url: "/",
        version: 1.0,
        logging: false,
        online: false,
        html: "index.html",
        pipeHTML: true,
        replace: true
    })).pipe(gulp.dest('public/images/favicons'));
});

gulp.task('build:fonts', function () {
    return gulp.src(source + '/assets/fonts/**/*.{svg,ttf,woff,woff2,eot}')
        .pipe(rename({dirname: ''}))
        .pipe(gulp.dest(publicOutput + '/fonts'));
});

// Watch
gulp.task('watch', function () {
    gulp.watch(source + '/**/*.scss', ['build:styles']);
    gulp.watch(source + '/**/*.ts', ['build:scripts']);
});

gulp.task('build:scripts', ['build:ts', 'build:vendor-js']);
gulp.task('build', ['build:styles', 'build:scripts', 'build:fonts']);
gulp.task('develop', ['build', 'watch']);