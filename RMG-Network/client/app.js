// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	
	$scope.queryAllProduct = function(){

		appFactory.queryAllProduct(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i]);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_product = array;
		});
	}

	$scope.queryProduct = function(){

		var id = $scope.product_id;

		appFactory.queryProduct(id, function(data){
			$scope.query_product = data;

			if ($scope.query_product == "Could not locate product"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordProduct = function(){

		appFactory.recordProduct($scope.product, function(data){
			$scope.create_product = data;
			$("#success_create").show();
		});
	}

	$scope.changeOwner = function(){

		appFactory.changeOwner($scope.product, function(data){
			$scope.change_owner = data;
			if ($scope.change_owner == "Error: no product found"){
				$("#error_holder").show();
				$("#success_holder").hide();
			} else{
				$("#success_holder").show();
				$("#error_holder").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllProduct = function(callback){

    	$http.get('/get_all_product/').success(function(output){
			callback(output)
		});
	}

	factory.queryProduct = function(id, callback){
    	$http.get('/get_product/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordProduct = function(data, callback){

		

		var product = data.id + "-" + data.name + "-" + data.description + "-" + data.price + "-" + data.owner;

    	$http.get('/add_product/'+product).success(function(output){
			callback(output)
		});
	}

	factory.changeOwner = function(data, callback){

		var product = data.id + "-" + data.owner;

    	$http.get('/change_owner/'+product).success(function(output){
			callback(output)
		});
	}

	return factory;
});


