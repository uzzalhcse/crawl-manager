module.exports = {
    apps: [
        {
            name: 'crawl-manager-frontend',
            port: '3000',
            exec_mode: 'cluster',
            instances: 'max',
            script: './.output/server/index.mjs',
            env: {
                HOST: '0.0.0.0',
                PORT: 3000
            }
        }
    ]
}
