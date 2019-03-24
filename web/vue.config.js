module.exports = {
  runtimeCompiler: true,
  publicPath: '[{[ .StaticURL ]}]',
  pluginOptions: {
    i18n: {
      locale: 'en',
      fallbackLocale: 'en',
      localeDir: 'locales',
      enableInSFC: true
    }
  }
}
