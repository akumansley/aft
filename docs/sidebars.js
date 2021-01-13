module.exports = {
  main: [
    'about',
    'getting-started',
    {
    	type: 'category',
    	label: 'Tutorial',
    	collapsed: false,
    	items: ['tutorial/running', 'tutorial/frontend-setup', 'tutorial/app-setup', 
                'tutorial/login', 'tutorial/user', 'tutorial/models', 'tutorial/creates',
                'tutorial/updates', 'tutorial/access', 'tutorial/review'],
    },
    {
    	type: 'category',
    	label: 'Overview',
    	collapsed: false,
    	items: ['overview/schema', 'overview/api', 'overview/rpcs', 'overview/access', 'overview/identity', 'overview/records', 'overview/internals'],
    },
  ],
};
