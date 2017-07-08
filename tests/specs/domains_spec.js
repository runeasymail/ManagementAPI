var apiClient = require('../apiClient');

apiClient.logged(function(frisby){
	frisby.create('Get domains')
		.get('domains')
		.expectStatus(200)
		.expectJSON({domains: null})
		.toss();
});

apiClient.unlogged(function(frisby){
	frisby.create('Get domains without token')
		.get('domains')
		.expectStatus(200)
		.expectBodyContains('Token is not correct')
		.toss();
});


