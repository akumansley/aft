module.exports = {
  title: 'Aft',
  tagline: 'Free, open source, self-hosted backend as a service',
  url: '  https://aft.dev',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.svg',
  organizationName: 'awans', // Usually your GitHub org/user name.
  projectName: 'aft', // Usually your repo name.
  themeConfig: {
    colorMode: {
      defaultMode: 'dark',
      respectPrefersColorScheme: false,
      disableSwitch: true,
    },
    navbar: {
      title: 'Aft',
      logo: {
        alt: 'Site Logo',
        src: 'img/favicon.svg',
      },
      items: [
        {
          to: '/docs',
          label: 'Docs',
        },
        {
          href: 'https://github.com/awans/aft',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
          noFooter: true,
        },
      },
    ],
  ],
};
