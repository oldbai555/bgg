module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-essential',
    'plugin:vue/vue3-strongly-recommended',
    'plugin:vue/vue3-recommended'
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    parser: '@typescript-eslint/parser',
    extraFileExtensions: ['.vue']
  },
  plugins: ['vue', '@typescript-eslint'],
  rules: {
    // Vue 相关规则
    'vue/multi-word-component-names': 'off', // 允许单单词组件名
    'vue/no-v-html': 'warn', // 警告使用 v-html
    'vue/require-default-prop': 'off', // 不要求 props 有默认值
    'vue/require-explicit-emits': 'warn', // 警告未显式声明的 emits
    'vue/html-self-closing': [
      'warn',
      {
        html: {
          void: 'always',
          normal: 'never',
          component: 'always'
        },
        svg: 'always',
        math: 'always'
      }
    ],
    'vue/max-attributes-per-line': [
      'warn',
      {
        singleline: 3,
        multiline: 1
      }
    ],
    'vue/singleline-html-element-content-newline': 'off',
    'vue/multiline-html-element-content-newline': 'off',
    
    // JavaScript/TypeScript 相关规则
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off', // 生产环境警告 console
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off', // 生产环境禁止 debugger
    'no-unused-vars': 'off', // 关闭，使用 TypeScript 的类型检查
    'no-undef': 'off', // 关闭，使用 TypeScript 的类型检查
    '@typescript-eslint/no-unused-vars': ['warn', {
      argsIgnorePattern: '^_',
      varsIgnorePattern: '^_'
    }], // TypeScript 未使用变量检查
    '@typescript-eslint/no-explicit-any': 'warn', // 警告使用 any
    '@typescript-eslint/explicit-module-boundary-types': 'off', // 不要求显式返回类型
    'prefer-const': 'warn', // 建议使用 const
    'no-var': 'error', // 禁止使用 var
    'eqeqeq': ['warn', 'always'], // 要求使用 === 和 !==
    'curly': ['warn', 'all'], // 要求所有控制语句使用大括号
    'no-eval': 'error', // 禁止使用 eval
    'no-implied-eval': 'error', // 禁止隐式 eval
    'no-new-func': 'error', // 禁止使用 new Function
    'no-script-url': 'error', // 禁止使用 javascript: URL
    
    // 代码风格
    'indent': ['warn', 2, { SwitchCase: 1 }], // 2 空格缩进
    'quotes': ['warn', 'single'], // 单引号
    'semi': ['warn', 'never'], // 不使用分号
    'comma-dangle': ['warn', 'never'], // 禁止尾随逗号
    'object-curly-spacing': ['warn', 'never'], // 对象大括号内无空格
    'array-bracket-spacing': ['warn', 'never'], // 数组方括号内无空格
    'space-before-function-paren': ['warn', {
      anonymous: 'always',
      named: 'never',
      asyncArrow: 'always'
    }],
    'keyword-spacing': ['warn', {
      before: true,
      after: true
    }],
    'space-infix-ops': 'warn', // 操作符周围有空格
    'space-unary-ops': ['warn', {
      words: true,
      nonwords: false
    }],
    'brace-style': ['warn', '1tbs'], // 大括号风格
    'comma-spacing': ['warn', {
      before: false,
      after: true
    }],
    'key-spacing': ['warn', {
      beforeColon: false,
      afterColon: true
    }],
    'no-trailing-spaces': 'warn', // 禁止尾随空格
    'no-multiple-empty-lines': ['warn', {
      max: 2,
      maxEOF: 1
    }], // 最多 2 个空行
    'eol-last': ['warn', 'always'] // 文件末尾换行
  },
  // 没有装 Prettier，indent 规则统一由上面的 ESLint 规则管，.vue 文件不例外；
  // 这个 override 本身仍然是必须的——`eslint .` 靠 overrides[].files 的 glob 才会把 *.vue 纳入检查范围
  overrides: [
    {
      files: ['*.vue']
    }
  ],
  ignorePatterns: [
    'node_modules/',
    'dist/',
    'build/',
    '*.min.js',
    '*.d.ts',
    'src/api/generated/'
  ]
}

