module.exports = {
  title: 'Aft',
  tagline: 'Free, open source, self-hosted backend as a service',
  url: 'https://your-docusaurus-test-site.com',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.svg',
  organizationName: 'awans.org', // Usually your GitHub org/user name.
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
          to: 'docs/',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
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
