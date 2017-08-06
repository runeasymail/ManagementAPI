var apiClient = require('../apiClient');
var utils = require('../utils');

apiClient.logged(function(frisby){
	frisby.create('Get domains1')
		.get('domains')
		.expectStatus(200)
		.toss();

	frisby.create('Create domain')
		.post('domains', {
			domain: utils.randomString(8) + ".com"
		})
		.expectStatus(200)
		.toss();

	// Bug?
	frisby.create('Create domain negative')
		.post('domains', {
			domain: utils.randomString(6) + "com"
		})
		.expectStatus(200)
		.toss();

	frisby.create('Create domain as url')
		.post('domains', {
			domain: "https://" + utils.randomString(8) +".com"
		})
		.expectStatus(200)
		.toss();

	frisby.create('Create domain and add email account')
		.post('domains', {
			domain: utils.randomString(8) + ".com",
			username: utils.randomString(5) + "@" + utils.randomString(4) + "." + utils.randomString(3),
			password: utils.randomString(10)
		})
		.expectStatus(200)
		.toss();

	//Bug?
	frisby.create('Create domain and add email account negative')
		.post('domains', {
			domain: utils.randomString(8) + "." + utils.randomString(3),
			username: utils.randomString(5),
			password: utils.randomString(10)
		})
		.expectStatus(200)
		.toss();

	frisby.create('Get domains2')
		.get('domains')
		.expectStatus(200)
		.afterJSON(function(json){
			frisby.create('Get accounts')
		      	.get('users/' + json.domains[0].id)
		      	.expectStatus(200)
		    	.toss()
		})
		.toss();

	frisby.create('Get domain and create account')
		.get('domains')
		.expectStatus(200)
		.afterJSON(function(json){
			frisby.create('Get accounts')
		      	.post('users/' + json.domains[0].id, {
		      		username: utils.randomString(5) + "@" + utils.randomString(3) + "." + utils.randomString(3)
		      	})
		      	.expectStatus(200)
		    	.toss()
		})
		.toss();

	//It's not valid email
	frisby.create('Get domain and create account negative')
		.get('domains')
		.expectStatus(200)
		.afterJSON(function(json){
			frisby.create('Get accounts')
		      	.post('users/' + json.domains[0].id, {
		      		username: utils.randomString(5)
		      	})
		      	.expectStatus(200)
		    	.toss()
		})
		.toss();

	frisby.create('Get domain and change user password')
		.get('domains')
		.expectStatus(200)
		.afterJSON(function(json){
			var domainId = json.domains[0].id;
			frisby.create('Get user accounts')
		      	.get('users/' + domainId)
		      	.expectStatus(200)
		      	.afterJSON(function(json){
		      		var accountId = json.users[0].id
		      		frisby.create('Change password')
				      	.post('user-change-password', {
				      		id: accountId,
				      		password: utils.randomString(10),
				      		domain_id: domainId
				      	})
				      	.expectStatus(200)
				    	.toss()
		      	})
		    	.toss()
		})
		.toss();

	// Bug?
	frisby.create('Get domain and change user password negative')
		.get('domains')
		.expectStatus(200)
		.afterJSON(function(json){
			var domainId = json.domains[0].id;
			frisby.create('Get user accounts')
		      	.get('users/' + domainId)
		      	.expectStatus(200)
		      	.afterJSON(function(json){
		      		var accountId = json.users[0].id
		      		frisby.create('Change password')
				      	.post('user-change-password', {
				      		id: utils.randomNum(5),
				      		password: utils.randomString(10),
				      		domain_id: utils.randomNum(2)
				      	})
				      	.expectStatus(200)
				    	.toss()
		      	})
		    	.toss()
		})
		.toss();
});

apiClient.unlogged(function(frisby){
	frisby.create('Get domains without token')
		.get('domains')
		.expectStatus(200)
		.expectBodyContains('Token is not correct')
		.toss();
});


