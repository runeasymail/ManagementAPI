var fs = require('fs'),
	frisby = require('frisby'),
	ini = require('ini');

var config = ini.parse(fs.readFileSync('../config.ini', 'utf-8'));
config.baseUrl = "http://localhost:8081/";

var logged = function(callback) {
	frisby.globalSetup({
		request: {
			baseUri: config.baseUrl
		}
	});

	frisby.create('Login')
		.post('auth', {
			username: config.auth.username,
			password: config.auth.password
		})
		.expectStatus(200)
		.expectJSON({result: true})
		.afterJSON(function(json) {
			frisby.globalSetup({
				request: {
					headers: {"Auth-token": json['token']},
					baseUri: config.baseUrl
				}
			});
			callback(frisby);
	}).toss();
};

var unlogged = function(callback) {
	frisby.globalSetup({
		request: {
			baseUri: config.baseUrl
		}
	});

	callback(frisby);
};


exports.logged = logged;
exports.unlogged = unlogged;

exports.BASE_URL = config.baseUrl;
exports.USERNAME = config.auth.username;
exports.PASSWORD = config.auth.password;