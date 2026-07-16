/**
 * 与 styles/variables.scss 的 $screen-sm 保持同一数值，JS/TS 侧无法直接读取 SCSS 变量，
 * 此常量是断点在脚本侧的唯一来源，改动需同步 variables.scss
 */
export const MOBILE_BREAKPOINT = 768
