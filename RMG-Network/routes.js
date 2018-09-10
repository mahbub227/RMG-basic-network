//SPDX-License-Identifier: Apache-2.0

var tuna = require('./controller.js');

module.exports = function(app){

  app.get('/get_product/:id', function(req, res){
    tuna.get_product(req, res);
  });
  app.get('/add_product/:product', function(req, res){
    tuna.add_product(req, res);
  });
  app.get('/get_all_product', function(req, res){
    tuna.get_all_product(req, res);
  });
  app.get('/change_owner/:product', function(req, res){
    tuna.change_owner(req, res);
  });
}
