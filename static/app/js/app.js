var app = angular.module( "Karl", ['ngRoute'] );

app.config(['$routeProvider', function($routeProvider){
    $routeProvider
        .when('/', {
            controller : 'RootCtrl',
            templateUrl : 'tpl/index.html'
        })
        .when('/console', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/console.html'
        })
        .when('/strategy', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/strategy.html'
        })
        .when('/conf', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/conf.html'
        })
        .when('/feedback', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/feedback.html'
        })
        .when('/about', {
            controller : 'MainCtrl',
            templateUrl : 'tpl/about.html'
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
