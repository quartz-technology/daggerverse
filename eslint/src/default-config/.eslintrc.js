module.exports = {
  env: {
    node: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier",
    "plugin:prettier/recommended",
  ],
  parser: "@typescript-eslint/parser",
  plugins: ["@typescript-eslint"],
  root: true,
};
