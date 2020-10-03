const isProd = process.env.NODE_ENV === 'production';

module.exports = {
  title: 'fisherman',
  tagline: 'Small git hooks tool for developer.',
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
      links: [
        {
          title: 'Docs',
          items: [
            { label: 'Style Guide', to: 'docs/' },
            { label: 'Second Doc', to: 'docs/doc2/' },
          ],
        },
        {
          title: 'Community',
          items: [
            { label: 'Stack Overflow', href: 'https://stackoverflow.com/questions/tagged/docusaurus' },
            { label: 'Discord', href: 'https://discordapp.com/invite/docusaurus' },
            { label: 'Twitter', href: 'https://twitter.com/docusaurus' },
          ],
        },
        {
          title: 'More',
          items: [
            { label: 'Blog', to: 'blog' },
            { label: 'GitHub', href: 'https://github.com/facebook/docusaurus' },
          ],
        },
      ],
    },
    googleAnalytics: { trackingID: 'UA-128394725-2' },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          includeCurrentVersion: !isProd,
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
