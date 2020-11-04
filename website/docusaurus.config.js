const isProd = process.env.NODE_ENV === 'production';

module.exports = {
  title: 'fisherman',
  tagline: 'Small git hook management tool for developer.',
  url: 'https://fisherman.netlify.app',
  baseUrl: '/',
  onBrokenLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'evg4b',
  projectName: 'fisherman',
  themeConfig: {
    navbar: {
      title: 'fisherman',
      logo: { alt: 'fisherman logo', src: 'img/logo.png' },
      items: [
        { type: 'docsVersionDropdown', position: 'left' },
        { to: 'docs/', activeBasePath: 'docs', label: 'Docs', position: 'right' },
        { href: 'https://github.com/evg4b/fisherman', label: 'GitHub', position: 'right' },
      ],
    },
    footer: {
      style: 'dark',
      copyright: 'fisherman',
    },
    googleAnalytics: { trackingID: 'UA-128394725-2' },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          path: '../docs',
          includeCurrentVersion: true,
          showLastUpdateTime: true,
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        blog: false,
      },
    ],
  ],
};
