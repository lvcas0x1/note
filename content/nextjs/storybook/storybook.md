---
date: 2026-09-22
---
### Installation

npx storybook@latest init

### Run without telemetry

npm run storybook -- --disable-telemetry

```javascript
// .storybook/main.ts

// Replace your-framework with the framework you are using (e.g., react-webpack5, vue3-vite)
import type { StorybookConfig } from "@storybook/your-framework";

const config: StorybookConfig = {
  framework: "@storybook/your-framework",
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|ts|tsx)"],
  core: {
    disableTelemetry: true, // 👈 Disables telemetry
  },
};

export default config;
```

### Stop Crash reports (disabled by default)

```javascript
// .storybook/main.ts

// Replace your-framework with the framework you are using (e.g., react-webpack5, vue3-vite)
import type { StorybookConfig } from "@storybook/your-framework";

const config: StorybookConfig = {
  framework: "@storybook/your-framework",
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|ts|tsx)"],
  core: {
    enableCrashReports: true, // 👈 Appends the crash reports to the telemetry events
  },
};

export default config;
```

### Storybook with TailwindCss

- Install addon
  yarn add -D @storybook/addon-styling

- Regist storybook

```javascript
module.exports = {
  stories: ["../stories/**/*.mdx", "../stories/**/*.stories.@(js|jsx|ts|tsx)"],
  addons: ["@storybook/addon-essentials", "@storybook/addon-styling"],
};
```

- Auto Config

```shell
# Run the postinstall script from the root of your project
node node_modules/@storybook/addon-styling/bin/postinstall.js
```

### Opt-out Telemetry

npm run storybook -- --disable-telemetry
