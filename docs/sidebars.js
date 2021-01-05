module.exports = {
  main: [
    'about',
    'getting-started',
    {
    	type: 'category',
    	label: 'Tutorial',
    	collapsed: false,
    	items: ['tutorial/running', 'tutorial/models'],
    },
    {
    	type: 'category',
    	label: 'Overview',
    	collapsed: false,
    	items: ['overview/schema', 'overview/api', 'overview/rpcs', 'overview/access', 'overview/identity', 'overview/records', 'overview/datastore'],
    },
  ],
};
