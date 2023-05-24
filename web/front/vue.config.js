module.exports = {
    devServer: {
        disableHostCheck: true
    },
    transpileDependencies: ['vuetify'],
    publicPath: './',
    outputDir: 'dist',
    chainWebpack: config => {
        config.plugin('html').tap(args => {
            args[0].title = 'LB小破站'
            return args
        })
    },
    productionSourceMap: false
}
