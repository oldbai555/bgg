module.exports = {
  transpileDependencies: ['vuetify'],
  assetsDir: 'static',
  chainWebpack: config => {
    config.plugin('html').tap(args => {
      args[0].title = '白与慧的歌'
      return args
    })
  },
  productionSourceMap: false
}
