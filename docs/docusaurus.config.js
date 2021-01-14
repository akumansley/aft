module.exports = {
  title: 'Aft',
  tagline: 'Free, open source, self-hosted backend as a service',
  url: '  https://awans.github.io',
  baseUrl: '/aft/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.svg',
  organizationName: 'awans', // Usually your GitHub org/user name.
  projectName: 'aft', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Aft',
      logo: {
        alt: 'Site Logo',
        src: 'img/favicon.svg',
      },
      items: [
        {
          href: 'https://github.com/awans/aft',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {},
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          routeBasePath: '/',
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
