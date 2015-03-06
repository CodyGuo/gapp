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
        url: 'https://auth.alipay.com/login/index.htm',
        inputText: 'Abc@1332',
        buttonType: false
    }

    $scope.position = {
         x: 846,
         y: 243
     };

    // 启动
    $scope.loadPage = function() {
        var url = $scope.option.url;
        //alert('加载页面: ' + url);
        //App_OpenWindow(url);
        cef.openWindow(url);
    };

    // 开始
    $scope.start = function() {
        //alert("start.");
        cef.callback("start")
    };

    // 模拟点击
    $scope.emuClick = function() {
        var url = $scope.option.url;
        var x = $scope.position.x;
        var y = $scope.position.y;
        var buttonType = $scope.option.buttonType;
        cef.callback("emuClick", url, x, y, buttonType);
    };

    // 模拟输入
    $scope.emuInput = function() {
        var url = $scope.option.url;
        var inputText = $scope.option.inputText;
        $scope.emuClick();
        cef.callback("emuInput", url, inputText);
    };
}
]);

app.controller("StrategyCtrl", ['$scope', '$rootScope', function($scope, $rootScope) {
    $scope.strategies = [
        {id:1, name:'aaa'}
    ];
}
]);
