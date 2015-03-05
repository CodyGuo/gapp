var app = angular.module( "Karl", ['ngRoute'] );

app.config(['$routeProvider', function($routeProvider){
    $routeProvider
        .when('/', {
            controller : 'RootCtrl',
            templateUrl : 'tpl/index.html'
        })
        .when('/conf', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/conf.html'
        })
        .when('/signup', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/signup.html'
        })
        .when('/login', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/login.html'
        })
        .otherwise({
            redirectTo : '/'
        });
}]);
