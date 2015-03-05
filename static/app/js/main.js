function App_OpenWindow(url) {
    alert("App_OpenWindow");
    cef.openWindow(url);
}

app.controller("RootCtrl", ['$scope', '$rootScope', function($scope, $rootScope) {
 }
 ]);

app.controller("HeaderCtrl", ['$scope', '$rootScope', function($scope, $rootScope) {
  $rootScope.global = {
    Title: 'KarlBox',
    isPhone: false,
    isLogin: false
  };
}
]);

app.controller("MainCtrl", ['$scope', '$rootScope', function($scope, $rootScope) {
    $scope.option = {
        url: 'http://www.baidu.com/'
    };

    // 启动
    $scope.loadPage = function() {
        var url = $scope.option.url;
        //alert('加载页面: ' + url);
        //App_OpenWindow(url);
        cef.openWindow(url);
    };
}
]);

app.controller("StrategyCtrl", ['$scope', '$rootScope', function($scope, $rootScope) {
    $scope.strategies = [
        {id:1, name:'aaa'}
    ];
}
]);
