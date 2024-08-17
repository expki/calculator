module.exports = {
    mode: 'production',
    entry: './public/wasm_exec.js',
    output: {
        filename: 'wasm_exec.js',
        path: require('path').resolve(__dirname, 'dist')
    },
};
